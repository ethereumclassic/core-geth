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
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/vars"
)

// TestDifficultyETCMinimum verifies that difficulty never falls below the minimum (131072).
func TestDifficultyETCMinimum(t *testing.T) {
	config := params.ClassicChainConfig
	parent := &types.Header{
		Number:     big.NewInt(10_500_000),
		Difficulty: vars.MinimumDifficulty,
		Time:       1000,
	}
	// With a very large time gap, adjustment goes negative but is clamped.
	diff := CalcDifficulty(config, 100000, parent)
	if diff.Cmp(vars.MinimumDifficulty) < 0 {
		t.Errorf("difficulty %v below minimum %v", diff, vars.MinimumDifficulty)
	}
}

// TestDifficultyETCAdjustmentDirection verifies that shorter block times increase
// difficulty and longer block times decrease it (EIP-100B algorithm).
func TestDifficultyETCAdjustmentDirection(t *testing.T) {
	config := params.ClassicChainConfig
	parentDiff := big.NewInt(1_000_000_000)
	parentTime := uint64(1_600_000_000)
	parentNumber := big.NewInt(10_500_000) // Phoenix fork active (EIP-100B)

	parent := &types.Header{
		Number:     parentNumber,
		Difficulty: parentDiff,
		Time:       parentTime,
		UncleHash:  types.EmptyUncleHash,
	}

	// Fast block (1 second gap) — difficulty should increase
	fastDiff := CalcDifficulty(config, parentTime+1, parent)
	if fastDiff.Cmp(parentDiff) <= 0 {
		t.Errorf("fast block: difficulty should increase, got %v <= parent %v", fastDiff, parentDiff)
	}

	// Slow block (30 second gap) — difficulty should decrease
	slowDiff := CalcDifficulty(config, parentTime+30, parent)
	if slowDiff.Cmp(parentDiff) >= 0 {
		t.Errorf("slow block: difficulty should decrease, got %v >= parent %v", slowDiff, parentDiff)
	}
}

// TestDifficultyECIP1041BombRemoval verifies that after ECIP-1041 (DisposalBlock),
// the difficulty bomb component is completely removed from the calculation.
// On ETC Classic, DisposalBlock = 5,900,000.
func TestDifficultyECIP1041BombRemoval(t *testing.T) {
	config := params.ClassicChainConfig
	parentDiff := big.NewInt(500_000_000_000)
	parentTime := uint64(1_500_000_000)

	// Test at a very high block number (post-ECIP1041).
	// Without bomb removal, difficulty would contain a huge exponential term.
	parent := &types.Header{
		Number:     big.NewInt(15_000_000),
		Difficulty: parentDiff,
		Time:       parentTime,
		UncleHash:  types.EmptyUncleHash,
	}

	diff := CalcDifficulty(config, parentTime+13, parent)

	// The adjustment for a 13-second block on a 1T difficulty parent should
	// be a small delta. If the bomb were active, it would add 2^((15M/100000)-2)
	// which is astronomically large.
	// With bomb removed, diff should be close to parentDiff (within ~10%).
	ratio := new(big.Int).Mul(diff, big.NewInt(100))
	ratio.Div(ratio, parentDiff)
	if ratio.Int64() < 90 || ratio.Int64() > 110 {
		t.Errorf("post-ECIP1041 difficulty diverged too far from parent: got %v, parent %v (ratio %d%%)",
			diff, parentDiff, ratio.Int64())
	}
}

// TestDifficultyECIP1010BombPause verifies that ECIP-1010 pauses the difficulty
// bomb at block 3,000,000 and continues at 5,000,000 (3M + 2M length) on Classic.
func TestDifficultyECIP1010BombPause(t *testing.T) {
	// Create a config that has ECIP-1010 but NOT ECIP-1041 (bomb removal).
	// This lets us test the pause/continue logic in isolation.
	config := &coregeth.CoreGethChainConfig{
		NetworkID: 1,
		ChainID:   big.NewInt(61),
		Ethash:    new(ctypes.EthashConfig),

		EIP2FBlock: big.NewInt(1150000),
		EIP7FBlock: big.NewInt(1150000),

		EIP155Block:  big.NewInt(3000000),
		EIP160FBlock: big.NewInt(3000000),

		// ECIP-1010: pause at 3M, length 2M (continue at 5M)
		ECIP1010PauseBlock: big.NewInt(3_000_000),
		ECIP1010Length:     big.NewInt(2_000_000),

		// Byzantium-eq (EIP-100B) at block 8,772,000 — NOT active for our test blocks
		EIP100FBlock: big.NewInt(8_772_000),

		// NO DisposalBlock — bomb is NOT removed
	}

	parentDiff := big.NewInt(100_000_000_000)
	parentTime := uint64(1_500_000_000)

	// Block 3,500,000 is during the pause. The explosion reference point
	// should be frozen at block 3,000,000.
	parentDuring := &types.Header{
		Number:     big.NewInt(3_499_999),
		Difficulty: parentDiff,
		Time:       parentTime,
	}
	diffDuring := CalcDifficulty(config, parentTime+13, parentDuring)

	// Block 4,000,000 is still during the pause. Explosion reference
	// should still be 3,000,000.
	parentLater := &types.Header{
		Number:     big.NewInt(3_999_999),
		Difficulty: parentDiff,
		Time:       parentTime,
	}
	diffLater := CalcDifficulty(config, parentTime+13, parentLater)

	// During the pause, the bomb component should be the same regardless of
	// block number (since the reference point is frozen at 3M).
	if diffDuring.Cmp(diffLater) != 0 {
		t.Errorf("bomb should be frozen during ECIP-1010 pause:\n  block 3.5M difficulty: %v\n  block 4.0M difficulty: %v",
			diffDuring, diffLater)
	}

	// After the continue point (5M), the bomb should resume with an offset.
	// Block 5,500,000: reference = 5.5M - (5M - 3M) = 3.5M
	// This should produce a LARGER bomb than during the pause (ref=3M).
	parentAfter := &types.Header{
		Number:     big.NewInt(5_499_999),
		Difficulty: parentDiff,
		Time:       parentTime,
	}
	diffAfter := CalcDifficulty(config, parentTime+13, parentAfter)

	// After continue, difficulty should be >= during pause (bomb resumes)
	if diffAfter.Cmp(diffDuring) < 0 {
		t.Errorf("bomb should resume after ECIP-1010 continue:\n  during pause: %v\n  after continue: %v",
			diffDuring, diffAfter)
	}
}

// TestDifficultyECIP1099EpochLength verifies that ECIP-1099 doubles the DAG
// epoch length from 30,000 to 60,000 blocks.
func TestDifficultyECIP1099EpochLength(t *testing.T) {
	cases := []struct {
		name        string
		block       uint64
		ecip1099    *uint64
		wantEpochLen uint64
	}{
		{
			name:        "pre-ECIP1099 default epoch",
			block:       1_000_000,
			ecip1099:    nil,
			wantEpochLen: 30000,
		},
		{
			name:        "Classic pre-ECIP1099 (block < 11.7M)",
			block:       11_699_999,
			ecip1099:    uint64Ptr(11_700_000),
			wantEpochLen: 30000,
		},
		{
			name:        "Classic at ECIP1099 (block = 11.7M)",
			block:       11_700_000,
			ecip1099:    uint64Ptr(11_700_000),
			wantEpochLen: 60000,
		},
		{
			name:        "Classic post-ECIP1099",
			block:       15_000_000,
			ecip1099:    uint64Ptr(11_700_000),
			wantEpochLen: 60000,
		},
		{
			name:        "Mordor at ECIP1099 (block = 2.52M)",
			block:       2_520_000,
			ecip1099:    uint64Ptr(2_520_000),
			wantEpochLen: 60000,
		},
		{
			name:        "Mordor pre-ECIP1099",
			block:       2_519_999,
			ecip1099:    uint64Ptr(2_520_000),
			wantEpochLen: 30000,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calcEpochLength(tc.block, tc.ecip1099)
			if got != tc.wantEpochLen {
				t.Errorf("calcEpochLength(%d) = %d, want %d", tc.block, got, tc.wantEpochLen)
			}
		})
	}
}

// TestDifficultyECIP1099EpochCalculation verifies epoch numbering across the
// ECIP-1099 transition boundary.
func TestDifficultyECIP1099EpochCalculation(t *testing.T) {
	cases := []struct {
		name      string
		block     uint64
		epochLen  uint64
		wantEpoch uint64
	}{
		// Pre-ECIP1099: 30,000 block epochs
		{"epoch 0 (block 0)", 0, 30000, 0},
		{"epoch 0 (block 29999)", 29999, 30000, 0},
		{"epoch 1 (block 30000)", 30000, 30000, 1},
		{"epoch 389 (block 11.7M)", 11_700_000, 30000, 390},

		// Post-ECIP1099: 60,000 block epochs
		{"epoch 195 (block 11.7M, doubled)", 11_700_000, 60000, 195},
		{"epoch 195 (block 11.759M)", 11_759_999, 60000, 195},
		{"epoch 196 (block 11.76M)", 11_760_000, 60000, 196},

		// Mordor ECIP-1099 at 2,520,000
		{"Mordor epoch 42 (block 2.52M, doubled)", 2_520_000, 60000, 42},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calcEpoch(tc.block, tc.epochLen)
			if got != tc.wantEpoch {
				t.Errorf("calcEpoch(%d, %d) = %d, want %d",
					tc.block, tc.epochLen, got, tc.wantEpoch)
			}
		})
	}
}

// TestDifficultyClassicVsMordorConfigs verifies that both Classic and Mordor
// configs produce valid difficulty calculations at their respective fork boundaries.
func TestDifficultyClassicVsMordorConfigs(t *testing.T) {
	parentDiff := big.NewInt(200_000_000_000)
	parentTime := uint64(1_600_000_000)

	tests := []struct {
		name   string
		config ctypes.ChainConfigurator
		blocks []int64
	}{
		{
			name:   "Classic",
			config: params.ClassicChainConfig,
			blocks: []int64{
				1,              // Pre-Homestead
				1_150_000,      // Homestead
				3_000_000,      // ECIP-1010 pause
				5_900_000,      // ECIP-1041 (bomb removal)
				8_772_000,      // Atlantis (Byzantium-eq, EIP-100B)
				10_500_839,     // Phoenix (Istanbul-eq)
				13_189_133,     // Magneto (Berlin-eq)
				19_250_000,     // Spiral (Shanghai-eq)
			},
		},
		{
			name:   "Mordor",
			config: params.MordorChainConfig,
			blocks: []int64{
				1,
				301_243,    // Atlantis
				999_983,    // Agharta
				2_520_000,  // ECIP-1099 (Etchash)
				3_985_893,  // Magneto
				5_520_000,  // Mystique
				9_957_000,  // Spiral
			},
		},
	}

	for _, tt := range tests {
		for _, bn := range tt.blocks {
			parent := &types.Header{
				Number:     big.NewInt(bn - 1),
				Difficulty: parentDiff,
				Time:       parentTime,
				UncleHash:  types.EmptyUncleHash,
			}
			diff := CalcDifficulty(tt.config, parentTime+13, parent)
			if diff.Sign() <= 0 {
				t.Errorf("%s block %d: difficulty must be positive, got %v", tt.name, bn, diff)
			}
			if diff.Cmp(vars.MinimumDifficulty) < 0 {
				t.Errorf("%s block %d: difficulty %v below minimum %v", tt.name, bn, diff, vars.MinimumDifficulty)
			}
		}
	}
}

func uint64Ptr(v uint64) *uint64 {
	return &v
}
