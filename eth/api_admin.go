// Copyright 2023 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

// AdminAPI is the collection of Ethereum full node related APIs for node
// administration.
type AdminAPI struct {
	eth *Ethereum
}

// NewAdminAPI creates a new instance of AdminAPI.
func NewAdminAPI(eth *Ethereum) *AdminAPI {
	return &AdminAPI{eth: eth}
}

// ExportChain exports the current blockchain into a local file,
// or a range of blocks if first and last are non-nil.
func (api *AdminAPI) ExportChain(file string, first *uint64, last *uint64) (bool, error) {
	if first == nil && last != nil {
		return false, errors.New("last cannot be specified without first")
	}
	if first != nil && last == nil {
		head := api.eth.BlockChain().CurrentHeader().Number.Uint64()
		last = &head
	}
	if _, err := os.Stat(file); err == nil {
		// File already exists. Allowing overwrite could be a DoS vector,
		// since the 'file' may point to arbitrary paths on the drive.
		return false, errors.New("location would overwrite an existing file")
	}
	// Make sure we can create the file to export into
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return false, err
	}
	defer out.Close()

	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	// Export the blockchain
	if first != nil {
		if err := api.eth.BlockChain().ExportN(writer, *first, *last); err != nil {
			return false, err
		}
	} else if err := api.eth.BlockChain().Export(writer); err != nil {
		return false, err
	}
	return true, nil
}

func hasAllBlocks(chain *core.BlockChain, bs []*types.Block) bool {
	for _, b := range bs {
		if !chain.HasBlock(b.Hash(), b.NumberU64()) {
			return false
		}
	}

	return true
}

// ImportChain imports a blockchain from a local file.
func (api *AdminAPI) ImportChain(file string) (bool, error) {
	// Make sure the can access the file to import
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()

	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}

	// Run actual the import in pre-configured batches
	stream := rlp.NewStream(reader, 0)

	blocks, index := make([]*types.Block, 0, 2500), 0
	for batch := 0; ; batch++ {
		// Load a batch of blocks from the input file
		for len(blocks) < cap(blocks) {
			block := new(types.Block)
			if err := stream.Decode(block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			// ignore the genesis block when importing blocks
			if block.NumberU64() == 0 {
				continue
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}

		if hasAllBlocks(api.eth.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		// Import the batch and reset the buffer
		if _, err := api.eth.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}

// Ecbp1100 sets the ECBP-1100 (MESS) activation block and reports whether the
// mechanism is active afterwards. The block is a height, or "latest" or "pending",
// which both mean the current head; "finalized" and "safe" are refused.
//
// It is the runtime counterpart of --mess.activate and sets that end of the window only,
// so it can report false with the activation set exactly as asked: MESS is a window, and
// a deactivation block at or below the head closes it whatever the activation says. Mess
// is the switch that reaches both ends; Ecbp1100Status shows which end is responsible.
//
// This mutates chain configuration. To read the current state without changing
// it, use Ecbp1100Status.
func (api *AdminAPI) Ecbp1100(blockNr rpc.BlockNumber) (bool, error) {
	head := api.eth.blockchain.CurrentBlock().Number
	i, err := ecbp1100ActivationBlock(blockNr, head.Uint64())
	if err != nil {
		return false, err
	}
	if err := api.eth.blockchain.Config().SetECBP1100Transition(&i); err != nil {
		return false, err
	}
	status := api.Ecbp1100Status()
	if !status.Enabled && status.NodeSwitch {
		// The activation was set as asked and MESS still does not apply. Say which end is
		// responsible rather than leaving a bare false to be interpreted.
		if d := api.eth.blockchain.Config().GetECBP1100DeactivateTransition(); d != nil && *d <= head.Uint64() {
			log.Warn("ECBP1100 (MESS) activation set, but the window is closed at the head",
				"activation", i, "deactivation", *d, "head", head,
				"hint", "admin_mess(true) reaches both ends")
		}
	}
	return status.Enabled, nil
}

// Ecbp1100Deactivate sets the ECBP-1100 (MESS) deactivation block and reports whether the
// mechanism is active afterwards. It is the runtime counterpart of --mess.deactivate, and
// the other half of the pair Ecbp1100 begins: MESS applies from the activation block up to
// this one.
//
// This mutates chain configuration. To read the current state without changing
// it, use Ecbp1100Status.
func (api *AdminAPI) Ecbp1100Deactivate(blockNr rpc.BlockNumber) (bool, error) {
	head := api.eth.blockchain.CurrentBlock().Number
	i, err := ecbp1100DeactivationBlock(blockNr, head.Uint64())
	if err != nil {
		return false, err
	}
	if err := api.eth.blockchain.Config().SetECBP1100DeactivateTransition(&i); err != nil {
		return false, err
	}
	return api.Ecbp1100Status().Enabled, nil
}

// Mess turns the ECBP-1100 (MESS) chain-selection defense on or off in the running node and
// reports the state afterwards. It is the runtime counterpart of --mess, and the call an
// operator wants: MESS is a window with two ends, and this reaches whichever end is holding
// it shut, so neither the caller nor the documentation has to explain the encoding.
//
// On pushes the deactivation out of reach, and drops the activation to block 0 if it sits
// beyond the head, which is the state --mess=false leaves. Off pushes the activation out of
// reach, the method ECIP-1110 itself documents for disabling MESS.
//
// It does not touch the node-level switch that the low-peer-count and stale-head safeguards
// turn off, so Enabled can still be false with the window open; NodeSwitch reports that end
// and --mess.nodisable is what pins it.
//
// This mutates chain configuration, and does not persist: the flags and the config file
// decide again at the next start.
func (api *AdminAPI) Mess(enable bool) (ECBP1100Status, error) {
	head := api.eth.blockchain.CurrentBlock().Number
	if err := applyMESSSwitch(api.eth.blockchain.Config(), head.Uint64(), enable); err != nil {
		return ECBP1100Status{}, err
	}
	return api.Ecbp1100Status(), nil
}

// applyMESSSwitch moves whichever end of the ECBP-1100 window is holding it shut, which is
// what makes Mess a switch rather than a block setter. It is separate from Mess so that the
// decision can be tested against a chain configuration without standing up a node.
func applyMESSSwitch(config ctypes.ChainConfigurator, head uint64, enable bool) error {
	never := messNever
	if !enable {
		return config.SetECBP1100Transition(&never)
	}
	// An activation beyond the head keeps MESS off however the deactivation is set, so that
	// end has to move. It goes to block 0 rather than to the head, because the head is not a
	// floor: a reorg puts the local head below where it was, and every site that consults
	// this window consults it against the local head. An activation pinned at the head a
	// switch was thrown at would therefore switch MESS off for the depth of any reorg, which
	// is the event it exists for. Block 0 cannot be dipped below.
	//
	// An activation at or below the head is a block the operator chose, and is left alone.
	if a := config.GetECBP1100Transition(); a == nil || *a > head {
		from := uint64(0)
		if err := config.SetECBP1100Transition(&from); err != nil {
			return err
		}
	}
	return config.SetECBP1100DeactivateTransition(&never)
}

// ecbp1100Block resolves the block an admin ECBP-1100 call was given. A tag such as
// "latest" arrives as a negative rpc.BlockNumber, and read as a height it lands beyond any
// chain, which would quietly mean "never". The tags that name the chain head resolve to it;
// the tags this chain has no height for are refused. Both ends of the window resolve here,
// so the trap is handled in one place rather than once per setter.
func ecbp1100Block(blockNr rpc.BlockNumber, head uint64, verb string) (uint64, error) {
	switch {
	case blockNr >= 0:
		return uint64(blockNr), nil
	case blockNr == rpc.LatestBlockNumber, blockNr == rpc.PendingBlockNumber:
		return head, nil
	default:
		return 0, fmt.Errorf("ecbp1100: %v does not name a block to %s at", blockNr, verb)
	}
}

// ecbp1100ActivationBlock resolves the block Ecbp1100 was given.
func ecbp1100ActivationBlock(blockNr rpc.BlockNumber, head uint64) (uint64, error) {
	return ecbp1100Block(blockNr, head, "activate")
}

// ecbp1100DeactivationBlock resolves the block Ecbp1100Deactivate was given.
func ecbp1100DeactivationBlock(blockNr rpc.BlockNumber, head uint64) (uint64, error) {
	return ecbp1100Block(blockNr, head, "deactivate")
}

// ECBP1100Status reports whether the ECBP-1100 (MESS) chain-selection defense is
// enabled at the current head.
//
// The vocabulary matters and the code has historically got it wrong. ECBP-1100
// added MESS and shipped it enabled by default. ECBP-1110 turned the bundled
// default off at Spiral; it did not remove the mechanism, which remains in the
// client and available to any operator who enables it. So a height configured
// under ECBP-1110 disables MESS by default at that block rather than
// deactivating it, and the two are not the same claim.
//
// Block numbers are hex quantities, as elsewhere in the JSON-RPC API. As JSON numbers they
// lost precision in JavaScript: the console showed --mess=false's activation block,
// 18446744073709551614, as 18446744073709552000.
type ECBP1100Status struct {
	// Enabled reports whether MESS is presently applied to reorganization
	// decisions. It requires both that the node switch is on and that the head
	// has reached the activation block.
	Enabled bool `json:"enabled"`
	// NodeSwitch reports the node-level setting alone, which the low-peer-count and
	// stale-head safeguards turn off, independent of any height. --mess=false does not
	// touch it: it moves the activation block out of reach instead.
	NodeSwitch bool `json:"nodeSwitch"`
	// ActivatedAtBlock is the height from which the mechanism is available,
	// per ECBP-1100. Nil when unset.
	ActivatedAtBlock *hexutil.Uint64 `json:"activatedAtBlock"`
	// DefaultDisabledAtBlock is the height from which the bundled default is
	// off, per ECBP-1110, and is what this client ships for Ethereum Classic and
	// Mordor. A non-nil value disables MESS by default at that height without
	// removing it; nil means nothing stops it once activated, which is the state
	// --mess leaves and the one v1.13.0 shipped.
	DefaultDisabledAtBlock *hexutil.Uint64 `json:"defaultDisabledAtBlock"`
	// Head is the block number the heights above were evaluated against.
	Head hexutil.Uint64 `json:"head"`
}

// Ecbp1100Status reports the ECBP-1100 (MESS) state at the current head and
// changes nothing.
//
// Ecbp1100 above answers a similar question by first assigning the activation
// block it is passed, so it cannot be used to observe a running node. This
// exists so that state can be read without altering it.
func (api *AdminAPI) Ecbp1100Status() ECBP1100Status {
	return ecbp1100Status(api.eth.blockchain.Config(), api.eth.blockchain.CurrentBlock().Number,
		api.eth.blockchain.IsArtificialFinalityEnabled())
}

func ecbp1100Status(config ctypes.ChainConfigurator, head *big.Int, nodeSwitch bool) ECBP1100Status {
	return ECBP1100Status{
		Enabled:                nodeSwitch && config.IsEnabled(config.GetECBP1100Transition, head),
		NodeSwitch:             nodeSwitch,
		ActivatedAtBlock:       (*hexutil.Uint64)(config.GetECBP1100Transition()),
		DefaultDisabledAtBlock: (*hexutil.Uint64)(config.GetECBP1100DeactivateTransition()),
		Head:                   hexutil.Uint64(head.Uint64()),
	}
}

// MaxPeers sets the maximum peer limit for the protocol manager and the p2p server.
func (api *AdminAPI) MaxPeers(n int) (bool, error) {
	api.eth.handler.maxPeers = n
	api.eth.p2pServer.MaxPeers = n

	for i := api.eth.handler.peers.len(); i > n; i = api.eth.handler.peers.len() {
		p := api.eth.handler.peers.WorstPeer()
		if p == nil {
			break
		}
		api.eth.handler.removePeer(p.ID())
	}
	return true, nil
}
