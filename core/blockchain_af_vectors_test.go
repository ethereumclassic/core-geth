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

package core

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ECBP-1100 (MESS) conformance vectors.
//
// The vectors are computed from ECIP-1100's pseudocode and not from this client, so these
// tests measure conformance to the specification rather than agreement with the current
// implementation. Provenance and the validation performed before adoption are recorded in
// core/testdata/README.md.
//
// These exist because MESS's REJECTING branch had no coverage at all: `errReorgFinality`
// was asserted nowhere in the tree, and the test that once covered it,
// TestBlockChain_AF_ECBP1100, is skipped as disused since the sinusoidal-to-cubic change.
// Two full genesis syncs, across two networks, evaluated 169 reorganisations and permitted
// every one -- live running cannot reach the rejecting branch, because the reorgs a healthy
// network produces are exactly what the curve is built to permit.
const messVectorsFile = "testdata/ecbp1100_mess_vectors.json"

type messVectors struct {
	MessArtificialFinality struct {
		CurveDenominator string `json:"curveDenominator"`
		CurveVectors     []struct {
			TimeDeltaSeconds string `json:"timeDeltaSeconds"`
			CurveNumerator   string `json:"curveNumerator"`
		} `json:"curveVectors"`
		DecisionVectors []struct {
			CommonAncestorTimestamp string `json:"commonAncestorTimestamp"`
			LocalHeadTimestamp      string `json:"localHeadTimestamp"`
			ProposedHeadTimestamp   string `json:"proposedHeadTimestamp"`
			LocalSubchainTd         string `json:"localSubchainTd"`
			ProposedSubchainTd      string `json:"proposedSubchainTd"`
			Rejected                bool   `json:"rejected"`
			Note                    string `json:"note"`
		} `json:"decisionVectors"`
		SubchainVectors []struct {
			Name            string      `json:"name"`
			CommonAncestor  messBlock   `json:"commonAncestor"`
			LocalSegment    []messBlock `json:"localSegment"`
			ProposedSegment []messBlock `json:"proposedSegment"`
			Rejected        bool        `json:"rejected"`
			Note            string      `json:"note"`
		} `json:"subchainVectors"`
	} `json:"messArtificialFinality"`
}

type messBlock struct {
	Number          string `json:"number"`
	Timestamp       string `json:"timestamp"`
	Difficulty      string `json:"difficulty"`
	TotalDifficulty string `json:"totalDifficulty"`
}

func (b messBlock) num() uint64  { return mustBig(b.Number).Uint64() }
func (b messBlock) time() uint64 { return mustBig(b.Timestamp).Uint64() }

func mustBig(s string) *big.Int {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("bad decimal in MESS vectors: " + s)
	}
	return v
}

func loadMessVectors(t *testing.T) *messVectors {
	t.Helper()
	blob, err := os.ReadFile(messVectorsFile)
	if err != nil {
		t.Fatalf("reading %s: %v", messVectorsFile, err)
	}
	v := new(messVectors)
	if err := json.Unmarshal(blob, v); err != nil {
		t.Fatalf("parsing %s: %v", messVectorsFile, err)
	}
	// A vector file that loaded but carries nothing would pass every assertion below
	// while testing none of them.
	m := v.MessArtificialFinality
	if len(m.CurveVectors) == 0 || len(m.DecisionVectors) == 0 || len(m.SubchainVectors) == 0 {
		t.Fatalf("vector file is empty: curve=%d decision=%d subchain=%d",
			len(m.CurveVectors), len(m.DecisionVectors), len(m.SubchainVectors))
	}
	return v
}

// TestECBP1100Vectors_Curve asserts the polynomial against ECIP-1100's own values.
func TestECBP1100Vectors_Curve(t *testing.T) {
	v := loadMessVectors(t)
	for _, c := range v.MessArtificialFinality.CurveVectors {
		x, want := mustBig(c.TimeDeltaSeconds), mustBig(c.CurveNumerator)
		if got := ecbp1100PolynomialV(x); got.Cmp(want) != 0 {
			t.Errorf("curve at x=%s: got %s, want %s", x, got, want)
		}
	}
}

// TestECBP1100Vectors_Decision drives the real ecbp1100() through synthetic headers whose
// total difficulties reproduce each vector's two subchain sums, and asserts the verdict.
//
// Rejection is signalled by errReorgFinality, which before this file was asserted nowhere.
func TestECBP1100Vectors_Decision(t *testing.T) {
	v := loadMessVectors(t)
	for i, c := range v.MessArtificialFinality.DecisionVectors {
		ancestor := &types.Header{Number: big.NewInt(1000), Time: mustBig(c.CommonAncestorTimestamp).Uint64(), Difficulty: big.NewInt(1)}
		local := &types.Header{Number: big.NewInt(1001), Time: mustBig(c.LocalHeadTimestamp).Uint64(), Difficulty: big.NewInt(1)}
		proposedParent := &types.Header{Number: big.NewInt(1000), Time: mustBig(c.CommonAncestorTimestamp).Uint64(), Difficulty: big.NewInt(2)}
		proposed := &types.Header{
			Number:     big.NewInt(1001),
			Time:       mustBig(c.ProposedHeadTimestamp).Uint64(),
			Difficulty: big.NewInt(0),
			ParentHash: proposedParent.Hash(),
		}
		// Base TD is arbitrary: ecbp1100 subtracts the ancestor's from both sides.
		const base = 1_000_000
		td := map[common.Hash]*big.Int{
			ancestor.Hash():       big.NewInt(base),
			local.Hash():          new(big.Int).Add(big.NewInt(base), mustBig(c.LocalSubchainTd)),
			proposedParent.Hash(): new(big.Int).Add(big.NewInt(base), mustBig(c.ProposedSubchainTd)),
		}
		getTD := func(h common.Hash, _ uint64) *big.Int {
			if v, ok := td[h]; ok {
				return v
			}
			return big.NewInt(0)
		}

		err := ecbp1100(ancestor, local, proposed, getTD)
		if gotRejected := err != nil; gotRejected != c.Rejected {
			t.Errorf("decision[%d] (local=%s proposed=%s x=%d): rejected=%v want %v (err=%v)\n  note: %s",
				i, c.LocalSubchainTd, c.ProposedSubchainTd,
				local.Time-ancestor.Time, gotRejected, c.Rejected, err, c.Note)
		}
	}
}

// TestECBP1100Vectors_Subchain asserts that the client DERIVES the comparison's inputs
// correctly from two competing segments -- summing difficulty rather than counting blocks,
// measuring from the common ancestor, and taking the curve's age off the LOCAL head.
//
// That last point is the one an implementer gets wrong: the comment above ecbp1100() says
// proposed.Time while the code correctly uses current.Time.
func TestECBP1100Vectors_Subchain(t *testing.T) {
	v := loadMessVectors(t)
	for _, c := range v.MessArtificialFinality.SubchainVectors {
		t.Run(c.Name, func(t *testing.T) {
			ancestorTD := mustBig(c.CommonAncestor.TotalDifficulty)
			ancestor := &types.Header{
				Number:     new(big.Int).SetUint64(c.CommonAncestor.num()),
				Time:       c.CommonAncestor.time(),
				Difficulty: big.NewInt(1),
			}
			td := map[common.Hash]*big.Int{ancestor.Hash(): ancestorTD}

			// Build a segment, accumulating total difficulty from the ancestor.
			build := func(seg []messBlock, salt int64) (head *types.Header, parent *types.Header) {
				cur := new(big.Int).Set(ancestorTD)
				prev := ancestor
				for _, b := range seg {
					cur = new(big.Int).Add(cur, mustBig(b.Difficulty))
					h := &types.Header{
						Number:     new(big.Int).SetUint64(b.num()),
						Time:       b.time(),
						Difficulty: mustBig(b.Difficulty),
						ParentHash: prev.Hash(),
						Extra:      big.NewInt(salt).Bytes(), // keep the two segments distinct
					}
					td[h.Hash()] = new(big.Int).Set(cur)
					parent = prev
					prev = h
				}
				return prev, parent
			}
			local, _ := build(c.LocalSegment, 1)
			proposed, _ := build(c.ProposedSegment, 2)

			getTD := func(h common.Hash, _ uint64) *big.Int {
				if v, ok := td[h]; ok {
					return v
				}
				return big.NewInt(0)
			}

			err := ecbp1100(ancestor, local, proposed, getTD)
			if gotRejected := err != nil; gotRejected != c.Rejected {
				t.Errorf("rejected=%v want %v (err=%v)\n  note: %s", gotRejected, c.Rejected, err, c.Note)
			}
		})
	}
}
