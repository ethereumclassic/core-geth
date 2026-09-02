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
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
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
// mechanism is active afterwards.
//
// This mutates chain configuration. To read the current state without changing
// it, use Ecbp1100Status.
func (api *AdminAPI) Ecbp1100(blockNr rpc.BlockNumber) (bool, error) {
	i := uint64(blockNr.Int64())
	err := api.eth.blockchain.Config().SetECBP1100Transition(&i)
	return api.eth.blockchain.IsArtificialFinalityEnabled() &&
		api.eth.blockchain.Config().IsEnabled(
			api.eth.blockchain.Config().GetECBP1100Transition,
			api.eth.blockchain.CurrentBlock().Number), err
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
type ECBP1100Status struct {
	// Enabled reports whether MESS is presently applied to reorganization
	// decisions. It requires both that the node switch is on and that the head
	// has reached the activation block.
	Enabled bool `json:"enabled"`
	// NodeSwitch reports the node-level setting alone, which --mess=false and
	// the low-peer-count safeguard both turn off, independent of any height.
	NodeSwitch bool `json:"nodeSwitch"`
	// ActivatedAtBlock is the height from which the mechanism is available,
	// per ECBP-1100. Nil when unset.
	ActivatedAtBlock *uint64 `json:"activatedAtBlock"`
	// DefaultDisabledAtBlock is the height from which the bundled default is
	// off, per ECBP-1110. Nil means the shipped default stays on, which is this
	// client's configuration; a non-nil value disables MESS by default at that
	// height without removing it.
	DefaultDisabledAtBlock *uint64 `json:"defaultDisabledAtBlock"`
	// Head is the block number the heights above were evaluated against.
	Head uint64 `json:"head"`
}

// Ecbp1100Status reports the ECBP-1100 (MESS) state at the current head and
// changes nothing.
//
// Ecbp1100 above answers a similar question by first assigning the activation
// block it is passed, so it cannot be used to observe a running node. This
// exists so that state can be read without altering it.
func (api *AdminAPI) Ecbp1100Status() ECBP1100Status {
	config := api.eth.blockchain.Config()
	head := api.eth.blockchain.CurrentBlock().Number
	nodeSwitch := api.eth.blockchain.IsArtificialFinalityEnabled()
	return ECBP1100Status{
		Enabled:                nodeSwitch && config.IsEnabled(config.GetECBP1100Transition, head),
		NodeSwitch:             nodeSwitch,
		ActivatedAtBlock:       config.GetECBP1100Transition(),
		DefaultDisabledAtBlock: config.GetECBP1100DeactivateTransition(),
		Head:                   head.Uint64(),
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
