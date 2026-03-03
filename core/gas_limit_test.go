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
	"testing"

	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/params/vars"
)

// TestGasLimitBoundDivisor verifies the 1/1024 adjustment rule constant.
func TestGasLimitBoundDivisor(t *testing.T) {
	if vars.GasLimitBoundDivisor != 1024 {
		t.Fatalf("GasLimitBoundDivisor: got %d, want 1024", vars.GasLimitBoundDivisor)
	}
}

// TestGasLimitMinimum verifies the minimum gas limit is 5000.
func TestGasLimitMinimum(t *testing.T) {
	if vars.MinGasLimit != 5000 {
		t.Fatalf("MinGasLimit: got %d, want 5000", vars.MinGasLimit)
	}
}

// TestGasLimitVerifyValid verifies that gas limit changes within 1/1024 are accepted.
func TestGasLimitVerifyValid(t *testing.T) {
	cases := []struct {
		name   string
		parent uint64
		header uint64
	}{
		{"unchanged at 8M", 8_000_000, 8_000_000},
		{"increase by 1", 8_000_000, 8_000_001},
		{"decrease by 1", 8_000_000, 7_999_999},
		// max delta = parent/1024 - 1 = 8000000/1024 - 1 = 7811
		{"max increase", 8_000_000, 8_007_811},
		{"max decrease", 8_000_000, 7_992_189},
		{"at minimum", 5000, 5000},
	}
	for _, tc := range cases {
		if err := misc.VerifyGaslimit(tc.parent, tc.header); err != nil {
			t.Errorf("VerifyGaslimit(%s): unexpected error: %v", tc.name, err)
		}
	}
}

// TestGasLimitVerifyInvalid verifies that gas limit changes exceeding 1/1024 are rejected.
func TestGasLimitVerifyInvalid(t *testing.T) {
	cases := []struct {
		name   string
		parent uint64
		header uint64
	}{
		// max allowed delta = parent/1024 - 1 = 7811, so 7812 is invalid
		{"too large increase", 8_000_000, 8_007_813},
		{"too large decrease", 8_000_000, 7_992_187},
		{"below minimum", 5000, 4999},
		{"zero gas limit", 8_000_000, 0},
	}
	for _, tc := range cases {
		if err := misc.VerifyGaslimit(tc.parent, tc.header); err == nil {
			t.Errorf("VerifyGaslimit(%s): expected error for parent=%d header=%d",
				tc.name, tc.parent, tc.header)
		}
	}
}

// TestCalcGasLimitConvergence verifies CalcGasLimit moves toward the desired target.
func TestCalcGasLimitConvergence(t *testing.T) {
	// Starting at 8M, targeting 8M — should stay at 8M
	result := CalcGasLimit(8_000_000, 8_000_000)
	if result != 8_000_000 {
		t.Errorf("CalcGasLimit(8M, 8M): got %d, want 8000000", result)
	}

	// Starting at 8M, targeting 10M — should increase
	result = CalcGasLimit(8_000_000, 10_000_000)
	if result <= 8_000_000 {
		t.Errorf("CalcGasLimit(8M, 10M): got %d, should be > 8000000", result)
	}

	// Starting at 8M, targeting 6M — should decrease
	result = CalcGasLimit(8_000_000, 6_000_000)
	if result >= 8_000_000 {
		t.Errorf("CalcGasLimit(8M, 6M): got %d, should be < 8000000", result)
	}
}

// TestCalcGasLimitMinimumFloor verifies CalcGasLimit never goes below MinGasLimit.
func TestCalcGasLimitMinimumFloor(t *testing.T) {
	// Even if desired is below minimum, should not go below MinGasLimit
	result := CalcGasLimit(vars.MinGasLimit, 1)
	if result < vars.MinGasLimit {
		t.Errorf("CalcGasLimit at minimum: got %d, should be >= %d", result, vars.MinGasLimit)
	}
}

// TestCalcGasLimitDelta verifies the per-block adjustment is bounded by 1/1024.
func TestCalcGasLimitDelta(t *testing.T) {
	parent := uint64(8_000_000)
	// delta = parent/GasLimitBoundDivisor - 1 = 8000000/1024 - 1 = 7811
	maxDelta := parent/vars.GasLimitBoundDivisor - 1

	// Targeting much higher — increase should be bounded
	result := CalcGasLimit(parent, 100_000_000)
	actualDelta := result - parent
	if actualDelta > maxDelta {
		t.Errorf("CalcGasLimit upward delta %d exceeds max %d", actualDelta, maxDelta)
	}
	if actualDelta != maxDelta {
		t.Errorf("CalcGasLimit upward delta %d should equal max %d when target is far away", actualDelta, maxDelta)
	}

	// Targeting much lower — decrease should be bounded
	result = CalcGasLimit(parent, vars.MinGasLimit)
	actualDelta = parent - result
	if actualDelta > maxDelta {
		t.Errorf("CalcGasLimit downward delta %d exceeds max %d", actualDelta, maxDelta)
	}
}
