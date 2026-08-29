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

package params

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// Required-block-hash vectors, read from an archive node rather than from a client's
// configuration. Provenance is in testdata/README.md.
//
// The hash required at block 1,920,000 is what makes this chain Ethereum Classic. The
// rejected alternatives in the vectors include the forked chain's block at the same height,
// so naming the wrong side of the DAO fork fails here rather than surfacing later as a node
// quietly syncing the wrong chain.
type etcRequiredHashVectors struct {
	RequiredBlockHeaders struct {
		RequiredHashes []struct {
			Block string `json:"block"`
			Hash  string `json:"hash"`
			Role  string `json:"role"`
		} `json:"requiredHashes"`
		Vectors []struct {
			Block         string `json:"block"`
			Hash          string `json:"hash"`
			Accepted      bool   `json:"accepted"`
			HashTakenFrom string `json:"hashTakenFrom"`
			Note          string `json:"note"`
		} `json:"vectors"`
	} `json:"requiredBlockHeaders"`
}

func TestETCRequiredBlockHashes(t *testing.T) {
	blob, err := os.ReadFile("testdata/etc_required_block_headers.json")
	if err != nil {
		t.Fatalf("reading vectors: %v", err)
	}
	v := new(etcRequiredHashVectors)
	if err := json.Unmarshal(blob, v); err != nil {
		t.Fatalf("parsing vectors: %v", err)
	}
	r := v.RequiredBlockHeaders
	if len(r.RequiredHashes) == 0 || len(r.Vectors) == 0 {
		t.Fatalf("vector file is empty: required=%d vectors=%d", len(r.RequiredHashes), len(r.Vectors))
	}

	cfg := ClassicChainConfig.RequireBlockHashes
	if len(cfg) == 0 {
		t.Fatal("ClassicChainConfig requires no block hashes -- the DAO-fork pin is absent")
	}

	// Every hash the vectors say is required must be configured, at the right height.
	for _, want := range r.RequiredHashes {
		n := mustDec(t, want.Block)
		got, ok := cfg[n]
		if !ok {
			t.Errorf("block %d (%s): no required hash configured, want %s", n, want.Role, want.Hash)
			continue
		}
		if got != common.HexToHash(want.Hash) {
			t.Errorf("block %d (%s): required hash is %s, want %s", n, want.Role, got.Hex(), want.Hash)
		}
	}

	// And each accept/reject vector must agree with what is configured. The rejected
	// entries are the load-bearing half: they include the forked chain's hash.
	var accepted, rejected int
	for _, c := range r.Vectors {
		n := mustDec(t, c.Block)
		configured, constrained := cfg[n]
		// A height with no configured hash is unconstrained: any hash is accepted.
		// The vectors exercise both neighbors of a required height precisely to pin
		// that the requirement does not leak onto adjacent blocks.
		matches := !constrained || configured == common.HexToHash(c.Hash)
		if matches != c.Accepted {
			t.Errorf("block %d hash %s: constrained=%v accepted=%v, want accepted=%v (%s; %s)",
				n, c.Hash, constrained, matches, c.Accepted, c.HashTakenFrom, c.Note)
		}
		if c.Accepted {
			accepted++
		} else {
			rejected++
		}
	}
	if rejected == 0 {
		t.Error("no rejection vectors ran -- a configuration matching everything would pass")
	}
	t.Logf("%d required hashes, %d accept and %d reject vectors asserted", len(r.RequiredHashes), accepted, rejected)
}

func mustDec(t *testing.T, s string) uint64 {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("bad decimal %q in vectors", s)
	}
	return v.Uint64()
}
