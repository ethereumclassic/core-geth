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

package mutations

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

// ECIP-1017 emission-schedule conformance vectors, computed from the specification rather
// than read back from this client. Provenance is in testdata/README.md.
type etcEraVectors struct {
	EraEmissionSchedule struct {
		EraLength           string   `json:"eraLength"`
		ValidOmmerDistances []string `json:"validOmmerDistances"`
		Vectors             []struct {
			Block                      string            `json:"block"`
			Era                        string            `json:"era"`
			WinnerBaseReward           string            `json:"winnerBaseReward"`
			IncluderBonusPerUncle      string            `json:"includerBonusPerUncle"`
			UncleMinerRewardByDistance map[string]string `json:"uncleMinerRewardByDistance"`
		} `json:"vectors"`
	} `json:"eraEmissionSchedule"`
}

func mustU256(t *testing.T, s string) *uint256.Int {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("bad decimal %q in vectors", s)
	}
	return uint256.MustFromBig(v)
}

func mustBigDec(t *testing.T, s string) *big.Int {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("bad decimal %q in vectors", s)
	}
	return v
}

// TestETCEraEmissionVectors asserts the ECIP-1017 issuance schedule: era index, winner
// reward, includer bonus per ommer, and the distance-dependent ommer reward.
func TestETCEraEmissionVectors(t *testing.T) {
	blob, err := os.ReadFile("testdata/etc_era_emission_schedule.json")
	if err != nil {
		t.Fatalf("reading vectors: %v", err)
	}
	v := new(etcEraVectors)
	if err := json.Unmarshal(blob, v); err != nil {
		t.Fatalf("parsing vectors: %v", err)
	}
	s := v.EraEmissionSchedule
	if len(s.Vectors) == 0 {
		t.Fatal("vector file carries no vectors")
	}
	eraLength := mustBigDec(t, s.EraLength)

	// The base reward ECIP-1017 reduces from. Era 0 pays it in full.
	baseReward := uint256.MustFromBig(new(big.Int).SetUint64(5e18))

	for _, c := range s.Vectors {
		block := mustBigDec(t, c.Block)

		era := GetBlockEra(block, eraLength)
		if want := mustBigDec(t, c.Era); era.Cmp(want) != 0 {
			// The boundary is the subtle part: block 5,000,000 is era 0, not era 1.
			t.Errorf("block %s: era got %s, want %s", block, era, want)
			continue
		}

		if got, want := GetBlockWinnerRewardByEra(era, baseReward), mustU256(t, c.WinnerBaseReward); got.Cmp(want) != 0 {
			t.Errorf("block %s (era %s): winner reward got %s, want %s", block, era, got, want)
		}

		// The includer bonus is 1/32 of the era's reward, per ommer included.
		oneUncle := []*types.Header{{Number: new(big.Int).Sub(block, big.NewInt(1))}}
		if got, want := GetBlockWinnerRewardForUnclesByEra(era, oneUncle, baseReward), mustU256(t, c.IncluderBonusPerUncle); got.Cmp(want) != 0 {
			t.Errorf("block %s (era %s): includer bonus got %s, want %s", block, era, got, want)
		}

		// Era 0 pays an ommer by how far back it sits; later eras pay a flat rate.
		for distStr, wantStr := range c.UncleMinerRewardByDistance {
			dist := mustBigDec(t, distStr)
			header := &types.Header{Number: block}
			uncle := &types.Header{Number: new(big.Int).Sub(block, dist)}
			if got, want := GetBlockUncleRewardByEra(era, header, uncle, baseReward), mustU256(t, wantStr); got.Cmp(want) != 0 {
				t.Errorf("block %s (era %s) ommer at distance %s: got %s, want %s",
					block, era, distStr, got, want)
			}
		}
	}
	t.Logf("%d ECIP-1017 emission vectors asserted across %d ommer distances",
		len(s.Vectors), len(s.ValidOmmerDistances))
}
