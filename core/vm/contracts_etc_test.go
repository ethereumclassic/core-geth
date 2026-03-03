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

package vm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/vars"
)

// addr returns a precompile address from a single byte.
func addr(b byte) common.Address {
	return common.BytesToAddress([]byte{b})
}

// TestETCPrecompilesPerFork verifies the correct set of precompiles is active
// at each ETC hard fork boundary. This catches accidental activation or
// deactivation of precompiles during chain config changes.
func TestETCPrecompilesPerFork(t *testing.T) {
	config := params.ClassicChainConfig
	var zero uint64

	tests := []struct {
		name     string
		block    int64
		expected []byte // addresses that MUST be present
		absent   []byte // addresses that MUST NOT be present
	}{
		{
			name:     "pre-Atlantis (block 8,771,999)",
			block:    8_771_999,
			expected: []byte{1, 2, 3, 4},
			absent:   []byte{5, 6, 7, 8, 9},
		},
		{
			name:     "Atlantis (block 8,772,000) — adds modexp + bn256",
			block:    8_772_000,
			expected: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			absent:   []byte{9, 10, 11},
		},
		{
			name:     "Agharta (block 9,573,000) — no precompile changes",
			block:    9_573_000,
			expected: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			absent:   []byte{9},
		},
		{
			name:     "Phoenix (block 10,500,839) — adds blake2F",
			block:    10_500_839,
			expected: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
			absent:   []byte{10, 11},
		},
		{
			name:     "Magneto (block 13,189,133) — no new precompile addresses",
			block:    13_189_133,
			expected: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
			absent:   []byte{10, 11},
		},
		{
			name:     "Spiral (block 19,250,000) — no new precompile addresses",
			block:    19_250_000,
			expected: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
			absent:   []byte{10, 11},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			precompiles := PrecompiledContractsForConfig(config, big.NewInt(tt.block), &zero)
			for _, b := range tt.expected {
				if precompiles[addr(b)] == nil {
					t.Errorf("missing expected precompile at address 0x%02x", b)
				}
			}
			for _, b := range tt.absent {
				if precompiles[addr(b)] != nil {
					t.Errorf("unexpected precompile present at address 0x%02x", b)
				}
			}
		})
	}
}

// TestETCBn256GasRepricing verifies that bn256 precompiles switch from
// Byzantium gas pricing to Istanbul gas pricing at the Phoenix fork.
func TestETCBn256GasRepricing(t *testing.T) {
	config := params.ClassicChainConfig
	var zero uint64

	// At Atlantis: Byzantium gas pricing
	atlantis := PrecompiledContractsForConfig(config, big.NewInt(8_772_000), &zero)
	addGasAtlantis := atlantis[addr(6)].RequiredGas(nil)
	mulGasAtlantis := atlantis[addr(7)].RequiredGas(nil)

	if addGasAtlantis != vars.Bn256AddGasByzantium {
		t.Errorf("Atlantis bn256Add gas = %d, want %d (Byzantium)", addGasAtlantis, vars.Bn256AddGasByzantium)
	}
	if mulGasAtlantis != vars.Bn256ScalarMulGasByzantium {
		t.Errorf("Atlantis bn256ScalarMul gas = %d, want %d (Byzantium)", mulGasAtlantis, vars.Bn256ScalarMulGasByzantium)
	}

	// At Phoenix: Istanbul gas pricing (cheaper)
	phoenix := PrecompiledContractsForConfig(config, big.NewInt(10_500_839), &zero)
	addGasPhoenix := phoenix[addr(6)].RequiredGas(nil)
	mulGasPhoenix := phoenix[addr(7)].RequiredGas(nil)

	if addGasPhoenix != vars.Bn256AddGasIstanbul {
		t.Errorf("Phoenix bn256Add gas = %d, want %d (Istanbul)", addGasPhoenix, vars.Bn256AddGasIstanbul)
	}
	if mulGasPhoenix != vars.Bn256ScalarMulGasIstanbul {
		t.Errorf("Phoenix bn256ScalarMul gas = %d, want %d (Istanbul)", mulGasPhoenix, vars.Bn256ScalarMulGasIstanbul)
	}

	// Sanity: Istanbul should be cheaper than Byzantium
	if addGasPhoenix >= addGasAtlantis {
		t.Errorf("Istanbul bn256Add (%d) should be cheaper than Byzantium (%d)", addGasPhoenix, addGasAtlantis)
	}
}

// TestETCModExpEIP2565Repricing verifies that the modexp precompile switches
// to EIP-2565 gas repricing at the Magneto fork. EIP-2565 introduces:
// - ceil(x/8)^2 complexity formula (simpler than EIP-198)
// - divisor of 3 (vs 20 in EIP-198)
// - minimum gas of 200
func TestETCModExpEIP2565Repricing(t *testing.T) {
	config := params.ClassicChainConfig
	var zero uint64

	// Use a 32-byte base, 32-byte exponent, 32-byte modulus — large enough for
	// both EIP-198 and EIP-2565 to produce meaningful (non-zero) gas costs.
	input := make([]byte, 96+96) // 32+32+32 header + 32+32+32 data
	input[31] = 32               // baseLen = 32
	input[63] = 32               // expLen = 32
	input[95] = 32               // modLen = 32
	// Fill data bytes with non-zero values
	for i := 96; i < 192; i++ {
		input[i] = 0xFF
	}

	// Before Magneto (Atlantis): original EIP-198 pricing
	preMagneto := PrecompiledContractsForConfig(config, big.NewInt(8_772_000), &zero)
	gasOld := preMagneto[addr(5)].RequiredGas(input)

	// At Magneto: EIP-2565 pricing
	magneto := PrecompiledContractsForConfig(config, big.NewInt(13_189_133), &zero)
	gasNew := magneto[addr(5)].RequiredGas(input)

	// Both should produce non-zero gas costs
	if gasOld == 0 {
		t.Errorf("pre-Magneto modexp gas should be non-zero")
	}
	if gasNew == 0 {
		t.Errorf("Magneto modexp gas should be non-zero")
	}

	// EIP-2565 changes the gas formula — costs differ for this input size
	if gasOld == gasNew {
		t.Errorf("modexp gas should differ between EIP-198 (%d) and EIP-2565 (%d)", gasOld, gasNew)
	}

	// EIP-2565 minimum is 200
	if gasNew < 200 {
		t.Errorf("EIP-2565 modexp gas = %d, should be >= 200 (minimum)", gasNew)
	}

	t.Logf("modexp gas: pre-Magneto (EIP-198) = %d, Magneto (EIP-2565) = %d", gasOld, gasNew)
}

// TestETCMordorPrecompiles verifies the Mordor testnet has the same precompile
// set as Classic mainnet at equivalent fork levels.
func TestETCMordorPrecompiles(t *testing.T) {
	mordor := params.MordorChainConfig
	classic := params.ClassicChainConfig
	var zero uint64

	// Mordor Spiral (9,957,000) should match Classic Spiral (19,250,000)
	mordorPrecompiles := PrecompiledContractsForConfig(mordor, big.NewInt(9_957_000), &zero)
	classicPrecompiles := PrecompiledContractsForConfig(classic, big.NewInt(19_250_000), &zero)

	// Both should have exactly the same set of addresses
	for a := range classicPrecompiles {
		if mordorPrecompiles[a] == nil {
			t.Errorf("Mordor Spiral missing precompile at %s (present on Classic Spiral)", a.Hex())
		}
	}
	for a := range mordorPrecompiles {
		if classicPrecompiles[a] == nil {
			t.Errorf("Mordor Spiral has extra precompile at %s (not on Classic Spiral)", a.Hex())
		}
	}

	// Verify count matches
	if len(mordorPrecompiles) != len(classicPrecompiles) {
		t.Errorf("precompile count mismatch: Mordor=%d, Classic=%d",
			len(mordorPrecompiles), len(classicPrecompiles))
	}
}
