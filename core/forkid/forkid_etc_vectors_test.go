// Copyright 2026 The core-geth Authors
// This file is part of the core-geth library.
//
// The core-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The core-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the core-geth library. If not, see <http://www.gnu.org/licenses/>.

package forkid

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

// EIP-2124 fork-identifier conformance vectors for Ethereum Classic mainnet.
//
// The values are computed from EIP-2124 rather than read back from this client, so these
// assert conformance to the specification. Provenance is in testdata/README.md.
//
// Fork identifiers gate peering: two nodes whose identifiers disagree do not connect. A
// silent change here is a network partition, not a test failure, which is why it is worth
// pinning against an external oracle.
type etcForkIDVectors struct {
	ForkIdentifierMainnet struct {
		GenesisHash string   `json:"genesisHash"`
		ForkBlocks  []string `json:"forkBlocks"`
		Vectors     []struct {
			Head      string `json:"head"`
			ForkHash  string `json:"forkHash"`
			ForkNext  string `json:"forkNext"`
			ForkIDRlp string `json:"forkIdRlp"`
		} `json:"vectors"`
	} `json:"forkIdentifierMainnet"`
}

func TestETCForkIDVectors(t *testing.T) {
	blob, err := os.ReadFile("testdata/etc_fork_identifiers.json")
	if err != nil {
		t.Fatalf("reading vectors: %v", err)
	}
	v := new(etcForkIDVectors)
	if err := json.Unmarshal(blob, v); err != nil {
		t.Fatalf("parsing vectors: %v", err)
	}
	m := v.ForkIdentifierMainnet
	// An empty vector set would pass every assertion below while testing none.
	if len(m.Vectors) == 0 {
		t.Fatal("vector file carries no vectors")
	}

	genesis := core.GenesisToBlock(params.DefaultClassicGenesisBlock(), nil)

	// The vectors are only meaningful against the genesis they were computed for.
	if got, want := genesis.Hash(), common.HexToHash(m.GenesisHash); got != want {
		t.Fatalf("genesis hash: got %s, want %s -- vectors do not describe this chain", got.Hex(), want.Hex())
	}

	for _, c := range m.Vectors {
		head := mustUint(t, c.Head)
		id := NewID(params.ClassicChainConfig, genesis, head, 0)

		var wantHash [4]byte
		hb, err := hexutil.Decode(c.ForkHash)
		if err != nil || len(hb) != 4 {
			t.Fatalf("head %d: bad forkHash %q in vectors", head, c.ForkHash)
		}
		copy(wantHash[:], hb)

		if id.Hash != wantHash {
			t.Errorf("head %d: forkHash got %#x, want %#x", head, id.Hash, wantHash)
		}
		if want := mustUint(t, c.ForkNext); id.Next != want {
			t.Errorf("head %d: forkNext got %d, want %d", head, id.Next, want)
		}

		// The RLP encoding is what actually travels in the handshake, so assert the
		// bytes on the wire rather than only the struct fields.
		enc, err := rlp.EncodeToBytes(id)
		if err != nil {
			t.Fatalf("head %d: encoding fork id: %v", head, err)
		}
		if got, want := hexutil.Encode(enc), c.ForkIDRlp; got != want {
			t.Errorf("head %d: forkIdRlp got %s, want %s", head, got, want)
		}
	}
	t.Logf("%d EIP-2124 vectors asserted against ClassicChainConfig", len(m.Vectors))
	_ = binary.BigEndian
	_ = big.NewInt
}

func mustUint(t *testing.T, s string) uint64 {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("bad decimal %q in vectors", s)
	}
	return v.Uint64()
}
