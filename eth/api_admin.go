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
	"sync/atomic"

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

	// messDeactivationGiven reports whether the MESS deactivation block was set rather than
	// bundled: at startup, by a MESS flag or config file line, or since, by MessDeactivate.
	// MessActivate and MessNoDisable keep a deactivation that was set, as --mess.activate keeps
	// one given with --mess.deactivate, and move a bundled one out of reach.
	messDeactivationGiven atomic.Bool
}

// NewAdminAPI creates a new instance of AdminAPI.
func NewAdminAPI(eth *Ethereum) *AdminAPI {
	api := &AdminAPI{eth: eth}
	api.messDeactivationGiven.Store(eth.config.OverrideECBP1100Deactivate != nil)
	return api
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

// MessActivate sets the ECBP-1100 (MESS) activation block and reports whether MESS applies
// afterwards. The block is a height, or "latest" or "pending", which both mean the current head;
// "finalized" and "safe" are refused.
//
// It is the runtime counterpart of --mess.activate, and does what the flag does: MESS applies
// from the block with no end, unless a deactivation block was given, with --mess.deactivate, a
// config file line or MessDeactivate. Then the later of the two blocks decides: an activation
// after the deactivation block turns MESS on, and one at or before it leaves MESS off once the
// head reaches the deactivation block, which the node then logs. Mess turns MESS on or off
// whatever the two blocks are.
//
// This mutates chain configuration, and does not persist. To read the current state without
// changing it, use MessStatus.
func (api *AdminAPI) MessActivate(blockNr rpc.BlockNumber) (bool, error) {
	head := api.eth.blockchain.CurrentBlock().Number
	i, err := ecbp1100ActivationBlock(blockNr, head.Uint64())
	if err != nil {
		return false, err
	}
	if err := applyMESSActivation(api.eth.blockchain.Config(), i, api.messDeactivationGiven.Load()); err != nil {
		return false, err
	}
	status := api.MessStatus()
	if !status.Enabled && status.NodeSwitch {
		// The activation was set as asked and MESS still does not apply. Say which block keeps
		// it off rather than leaving a bare false to be interpreted.
		if d := api.eth.blockchain.Config().GetECBP1100DeactivateTransition(); messKeptOffBy(i, d, head.Uint64()) {
			log.Warn("ECBP1100 (MESS) stays off: the deactivation block is at or after the activation block set, and the chain has reached it",
				"activation", i, "deactivation", *d, "head", head,
				"hint", "admin_mess(true) turns MESS on")
		}
	}
	return status.Enabled, nil
}

// messKeptOffBy reports whether deactivation is the block keeping MESS off at head once the
// activation is set to activation. The later block decides, so that is a deactivation at or
// after the activation which the head has reached; an activation after the deactivation turns
// MESS on from its own block instead.
func messKeptOffBy(activation uint64, deactivation *uint64, head uint64) bool {
	return deactivation != nil && activation <= *deactivation && *deactivation <= head
}

// applyMESSActivation sets the activation block and does with the deactivation what
// --mess.activate does: one that was given stays, and the later of the two blocks decides; a
// bundled one moves out of reach, so MESS applies from the activation with no end. It is separate
// from MessActivate so that it can be tested against a chain configuration without standing up a
// node.
func applyMESSActivation(config ctypes.ChainConfigurator, activation uint64, deactivationGiven bool) error {
	if err := config.SetECBP1100Transition(&activation); err != nil {
		return err
	}
	if deactivationGiven {
		return nil
	}
	never := core.ECBP1100Unreachable
	return config.SetECBP1100DeactivateTransition(&never)
}

// MessDeactivate sets the ECBP-1100 (MESS) deactivation block and reports whether MESS applies
// afterwards. The block is read as for MessActivate.
//
// It is the runtime counterpart of --mess.deactivate. Of the activation and the deactivation
// block, the later one decides: a deactivation at or after the activation turns MESS off from
// that block, and one before the activation does not turn it off. MessActivate and MessNoDisable
// keep a deactivation set here, as --mess.activate keeps one given with --mess.deactivate.
//
// This mutates chain configuration, and does not persist. To read the current state without
// changing it, use MessStatus.
func (api *AdminAPI) MessDeactivate(blockNr rpc.BlockNumber) (bool, error) {
	head := api.eth.blockchain.CurrentBlock().Number
	i, err := ecbp1100DeactivationBlock(blockNr, head.Uint64())
	if err != nil {
		return false, err
	}
	if err := api.eth.blockchain.Config().SetECBP1100DeactivateTransition(&i); err != nil {
		return false, err
	}
	api.messDeactivationGiven.Store(true)
	return api.MessStatus().Enabled, nil
}

// Mess turns the ECBP-1100 (MESS) chain-selection defense on or off in the running node and
// reports the state afterwards. It is the runtime counterpart of --mess and --mess=false, and
// it works whatever the two blocks are, so the caller does not have to know which one keeps
// MESS off.
//
// On pushes the deactivation out of reach and, if the activation is unset or sits beyond the
// head, as --mess=false leaves it, moves it back to the network's own activation block, as
// --mess does, or to block 0 where the network ships none or the head has not reached it. Off
// pushes the activation out of reach, so MESS does not apply whatever the deactivation block is.
//
// It does not touch the node-level switch that the low-peer-count and stale-head safeguards
// turn off, so Enabled can still be false with the window open. NodeSwitch reports that switch,
// and MessNoDisable or --mess.nodisable keeps it on.
//
// This mutates chain configuration, and does not persist: the flags and the config file
// decide again at the next start.
func (api *AdminAPI) Mess(enable bool) (ECBP1100Status, error) {
	head := api.eth.blockchain.CurrentBlock().Number
	if err := applyMESSSwitch(api.eth.blockchain.Config(), head.Uint64(), api.eth.messBundledActivation, enable); err != nil {
		return ECBP1100Status{}, err
	}
	return api.MessStatus(), nil
}

// applyMESSSwitch moves whichever end of the ECBP-1100 window is holding it shut, which is
// what makes Mess a switch rather than a block setter. bundled is the activation block the
// network ships, or nil. It is separate from Mess so that the decision can be tested against a
// chain configuration without standing up a node.
func applyMESSSwitch(config ctypes.ChainConfigurator, head uint64, bundled *uint64, enable bool) error {
	never := core.ECBP1100Unreachable
	if !enable {
		return config.SetECBP1100Transition(&never)
	}
	if err := resetMESSActivation(config, head, bundled); err != nil {
		return err
	}
	return config.SetECBP1100DeactivateTransition(&never)
}

// resetMESSActivation moves an activation that keeps MESS off however the deactivation is set,
// one unset or beyond the head, back to the network's own activation block, bundled, as --mess
// does when it clears an off switch. A network that ships none, or whose block the head has not
// reached, gets block 0.
//
// Not to the head, because the head is not a floor: a reorg puts the local head below where it
// was, and every site that consults this window consults it against the local head. An
// activation pinned at the head a switch was thrown at would therefore switch MESS off for the
// depth of any reorg, which is the event it exists for. The network's block, like block 0, is
// fixed and below the head, rather than tied to the moment of the call.
//
// An activation at or below the head is a block the operator chose, and is left alone.
func resetMESSActivation(config ctypes.ChainConfigurator, head uint64, bundled *uint64) error {
	if a := config.GetECBP1100Transition(); a != nil && *a <= head {
		return nil
	}
	from := uint64(0)
	if bundled != nil && *bundled <= head {
		from = *bundled
	}
	return config.SetECBP1100Transition(&from)
}

// MessNoDisable is the runtime counterpart of --mess.nodisable, and reports the state afterwards.
//
// True does what the flag does. It turns MESS on: an activation that keeps MESS off moves as Mess
// moves it, and the deactivation moves out of reach unless one was given. And it keeps MESS on
// once it is on: the low-peer-count and stale-head safeguards no longer switch it off. Like the
// flag, it does not switch MESS on while those safeguards hold it off; MESS comes on at the next
// sync, and then stays on.
//
// False lets the safeguards switch MESS off again, as --mess.nodisable=false does, and leaves the
// two blocks alone.
//
// This mutates chain configuration, and does not persist: the flags and the config file
// decide again at the next start.
func (api *AdminAPI) MessNoDisable(enable bool) (ECBP1100Status, error) {
	if !enable {
		api.eth.blockchain.ArtificialFinalityNoDisable(0)
		return api.MessStatus(), nil
	}
	head := api.eth.blockchain.CurrentBlock().Number.Uint64()
	if err := applyMESSNoDisable(api.eth.blockchain.Config(), head, api.eth.messBundledActivation, api.messDeactivationGiven.Load()); err != nil {
		return ECBP1100Status{}, err
	}
	api.eth.blockchain.ArtificialFinalityNoDisable(1)
	return api.MessStatus(), nil
}

// applyMESSNoDisable does to the window what --mess.nodisable does: MESS on, an activation that
// keeps it off moving as resetMESSActivation moves it, and the deactivation moving out of reach
// unless one was given. It is separate from MessNoDisable so that it can be tested against a
// chain configuration without standing up a node.
func applyMESSNoDisable(config ctypes.ChainConfigurator, head uint64, bundled *uint64, deactivationGiven bool) error {
	if err := resetMESSActivation(config, head, bundled); err != nil {
		return err
	}
	if deactivationGiven {
		return nil
	}
	never := core.ECBP1100Unreachable
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
		return 0, fmt.Errorf("%v does not name a block to %s MESS at", blockNr, verb)
	}
}

// ecbp1100ActivationBlock resolves the block MessActivate was given.
func ecbp1100ActivationBlock(blockNr rpc.BlockNumber, head uint64) (uint64, error) {
	return ecbp1100Block(blockNr, head, "activate")
}

// ecbp1100DeactivationBlock resolves the block MessDeactivate was given.
func ecbp1100DeactivationBlock(blockNr rpc.BlockNumber, head uint64) (uint64, error) {
	return ecbp1100Block(blockNr, head, "deactivate")
}

// ECBP1100Status reports whether the ECBP-1100 (MESS) chain-selection defense is
// enabled at the current head.
//
// The field names follow the two Best Practice documents. ECBP-1100 added MESS,
// activated it by default, and added the flags that let an operator turn it on or
// off. ECBP-1110 adjusts that default to off; it did not remove MESS from the
// client, and any operator can still turn it on. So the deactivation end is
// reported as DefaultDisabledAtBlock: past it MESS is off by default, not removed.
//
// Block numbers are hex quantities, as elsewhere in the JSON-RPC API. As JSON numbers they
// lost precision in JavaScript: the console showed --mess=false's activation block,
// 18446744073709551614, as 18446744073709552000.
type ECBP1100Status struct {
	// Enabled reports whether MESS is presently applied to reorganization
	// decisions. It requires that the node switch is on and that the head has
	// reached the activation block, and has not reached the deactivation block
	// unless the activation comes after it.
	Enabled bool `json:"enabled"`
	// NodeSwitch reports the node-level setting alone, which the low-peer-count and
	// stale-head safeguards turn off, independent of any height. --mess=false does not
	// touch it: it moves the activation block out of reach instead.
	NodeSwitch bool `json:"nodeSwitch"`
	// NoDisable reports whether the low-peer-count and stale-head safeguards are kept from
	// switching MESS off once it is on, as --mess.nodisable and MessNoDisable set.
	NoDisable bool `json:"noDisable"`
	// ActivatedAtBlock is the height from which the mechanism is available,
	// per ECBP-1100. Nil when unset.
	ActivatedAtBlock *hexutil.Uint64 `json:"activatedAtBlock"`
	// DefaultDisabledAtBlock is the height from which MESS stops applying: the
	// deactivation block of ECBP-1110's default, or one the operator set. It
	// disables MESS at that height without removing it, and an activation after it
	// turns MESS on again. Nil means nothing stops MESS once activated, as with
	// ECBP-1100's default; --mess, --mess.activate and --mess.nodisable, and Mess,
	// MessActivate and MessNoDisable on a running node, set 0xfffffffffffffffe instead, a
	// block no chain reaches.
	DefaultDisabledAtBlock *hexutil.Uint64 `json:"defaultDisabledAtBlock"`
	// Head is the block number the heights above were evaluated against.
	Head hexutil.Uint64 `json:"head"`
}

// MessStatus reports the ECBP-1100 (MESS) state at the current head and changes nothing.
//
// MessActivate answers a similar question by first assigning the activation block it is
// passed, so it cannot be used to observe a running node. This exists so that state can be read
// without altering it.
func (api *AdminAPI) MessStatus() ECBP1100Status {
	bc := api.eth.blockchain
	return ecbp1100Status(bc.Config(), bc.CurrentBlock().Number, bc.IsArtificialFinalityEnabled(),
		bc.IsArtificialFinalityNoDisable())
}

func ecbp1100Status(config ctypes.ChainConfigurator, head *big.Int, nodeSwitch, noDisable bool) ECBP1100Status {
	return ECBP1100Status{
		Enabled:                nodeSwitch && config.IsEnabled(config.GetECBP1100Transition, head),
		NodeSwitch:             nodeSwitch,
		NoDisable:              noDisable,
		ActivatedAtBlock:       (*hexutil.Uint64)(config.GetECBP1100Transition()),
		DefaultDisabledAtBlock: (*hexutil.Uint64)(config.GetECBP1100DeactivateTransition()),
		Head:                   hexutil.Uint64(head.Uint64()),
	}
}

// The MESS calls under their ECBP-1100 names, as the MESS flags keep theirs: --ecbp1100 for
// --mess.activate, --override.ecbp1100.deactivate for --mess.deactivate and --ecbp1100.nodisable
// for --mess.nodisable. Each does exactly what its MESS name does.

// Ecbp1100 is MessActivate under its ECBP-1100 name, which it had first.
func (api *AdminAPI) Ecbp1100(blockNr rpc.BlockNumber) (bool, error) {
	return api.MessActivate(blockNr)
}

// Ecbp1100Deactivate is MessDeactivate under its ECBP-1100 name.
func (api *AdminAPI) Ecbp1100Deactivate(blockNr rpc.BlockNumber) (bool, error) {
	return api.MessDeactivate(blockNr)
}

// Ecbp1100NoDisable is MessNoDisable under its ECBP-1100 name.
func (api *AdminAPI) Ecbp1100NoDisable(enable bool) (ECBP1100Status, error) {
	return api.MessNoDisable(enable)
}

// Ecbp1100Status is MessStatus under its ECBP-1100 name, which it had first.
func (api *AdminAPI) Ecbp1100Status() ECBP1100Status {
	return api.MessStatus()
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
