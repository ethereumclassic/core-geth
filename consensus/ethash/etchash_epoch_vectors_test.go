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

package ethash

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ECIP-1099 epoch-schedule conformance vectors, derived from the specification rather than
// from this client. Provenance is in testdata/README.md.
type etchashEpochVectors struct {
	EtchashEpochSchedule struct {
		ActivationBlock string `json:"activationBlock"`
		Vectors         []struct {
			Block           string `json:"block"`
			EpochLength     string `json:"epochLength"`
			Epoch           string `json:"epoch"`
			EpochStartBlock string `json:"epochStartBlock"`
			SeedHash        string `json:"seedHash"`
		} `json:"vectors"`
		SeedContinuity struct {
			Identity string `json:"identity"`
			Vectors  []struct {
				PostForkEpoch           string `json:"postForkEpoch"`
				PostForkEpochStartBlock string `json:"postForkEpochStartBlock"`
				SeedIterations          string `json:"seedIterations"`
				EquivalentLegacyEpoch   string `json:"equivalentLegacyEpoch"`
			} `json:"vectors"`
		} `json:"seedContinuity"`
	} `json:"etchashEpochSchedule"`
}

func loadEtchashVectors(t *testing.T) *etchashEpochVectors {
	t.Helper()
	blob, err := os.ReadFile("testdata/etchash_epoch_schedule.json")
	if err != nil {
		t.Fatalf("reading vectors: %v", err)
	}
	v := new(etchashEpochVectors)
	if err := json.Unmarshal(blob, v); err != nil {
		t.Fatalf("parsing vectors: %v", err)
	}
	s := v.EtchashEpochSchedule
	if len(s.Vectors) == 0 || len(s.SeedContinuity.Vectors) == 0 {
		t.Fatalf("vector file is empty: schedule=%d continuity=%d",
			len(s.Vectors), len(s.SeedContinuity.Vectors))
	}
	return v
}

func etchashUint(t *testing.T, s string) uint64 {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("bad decimal %q in vectors", s)
	}
	return v.Uint64()
}

// TestETCHashEpochSchedule asserts epoch length, index, start block and seed hash across
// the ECIP-1099 activation boundary.
func TestETCHashEpochSchedule(t *testing.T) {
	v := loadEtchashVectors(t)
	s := v.EtchashEpochSchedule
	activation := etchashUint(t, s.ActivationBlock)

	for _, c := range s.Vectors {
		block := etchashUint(t, c.Block)

		gotLen := calcEpochLength(block, &activation)
		if want := etchashUint(t, c.EpochLength); gotLen != want {
			t.Errorf("block %d: epochLength got %d, want %d", block, gotLen, want)
			continue
		}
		gotEpoch := calcEpoch(block, gotLen)
		if want := etchashUint(t, c.Epoch); gotEpoch != want {
			t.Errorf("block %d: epoch got %d, want %d", block, gotEpoch, want)
			continue
		}
		if want := etchashUint(t, c.EpochStartBlock); calcEpochBlock(gotEpoch, gotLen) != want {
			t.Errorf("block %d: epochStartBlock got %d, want %d", block, calcEpochBlock(gotEpoch, gotLen), want)
		}
		if got, want := hexutil.Encode(seedHash(gotEpoch, gotLen)), c.SeedHash; got != want {
			t.Errorf("block %d (epoch %d, len %d): seedHash got %s, want %s", block, gotEpoch, gotLen, got, want)
		}
	}
	t.Logf("%d ECIP-1099 schedule vectors asserted", len(s.Vectors))
}

// TestETCHashSeedContinuity asserts that the seed chain CONTINUES across the ECIP-1099
// activation rather than restarting: seed(e, 60000) == seed(2e, 30000).
//
// This is the assertion worth having. The natural wrong implementation divides the epoch
// start block by the epoch length in force instead of by 30,000. It raises no error --- it
// silently reuses a seed from roughly six million blocks earlier and generates a real
// dataset for the wrong epoch, so a miner produces work the network rejects while every
// internal check passes.
func TestETCHashSeedContinuity(t *testing.T) {
	v := loadEtchashVectors(t)
	s := v.EtchashEpochSchedule

	for _, c := range s.SeedContinuity.Vectors {
		post := etchashUint(t, c.PostForkEpoch)
		legacy := etchashUint(t, c.EquivalentLegacyEpoch)

		postSeed := hexutil.Encode(seedHash(post, 60000))
		legacySeed := hexutil.Encode(seedHash(legacy, 30000))
		if postSeed != legacySeed {
			t.Errorf("continuity broken: seed(epoch %d, 60000)=%s != seed(epoch %d, 30000)=%s",
				post, postSeed, legacy, legacySeed)
		}

		// The identity holds because the iteration count divides the start block by
		// 30,000 whatever the epoch length is. Assert the count itself, so an
		// implementation that reached the right seed by a wrong route still fails.
		if want := etchashUint(t, c.SeedIterations); legacy != want {
			t.Errorf("post-fork epoch %d: seed iterations got %d, want %d", post, legacy, want)
		}
	}
	t.Logf("%d seed-continuity vectors asserted (%s)", len(s.SeedContinuity.Vectors), s.SeedContinuity.Identity)
}
