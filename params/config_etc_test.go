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
	"math/big"
	"testing"
)

// TestETCClassicChainID verifies ETC mainnet uses chain ID 61.
func TestETCClassicChainID(t *testing.T) {
	id := ClassicChainConfig.GetChainID()
	if id.Cmp(big.NewInt(61)) != 0 {
		t.Fatalf("Classic chain ID: got %v, want 61", id)
	}
}

// TestETCMordorChainID verifies Mordor testnet uses chain ID 63.
func TestETCMordorChainID(t *testing.T) {
	id := MordorChainConfig.GetChainID()
	if id.Cmp(big.NewInt(63)) != 0 {
		t.Fatalf("Mordor chain ID: got %v, want 63", id)
	}
}

// TestETCClassicForkOrdering verifies that all Classic mainnet forks are
// activated in strictly ascending block order.
func TestETCClassicForkOrdering(t *testing.T) {
	type fork struct {
		name  string
		block *big.Int
	}
	forks := []fork{
		{"Homestead (EIP-2/7)", ClassicChainConfig.EIP2FBlock},
		{"Tangerine Whistle (EIP-150)", big.NewInt(2500000)},
		{"Spurious Dragon (EIP-155/160)", ClassicChainConfig.EIP155Block},
		{"ECIP-1017 (emission)", ClassicChainConfig.ECIP1017FBlock},
		{"Atlantis (EIP-161/170)", ClassicChainConfig.EIP161FBlock},
		{"Agharta (EIP-145)", ClassicChainConfig.EIP145FBlock},
		{"Phoenix (EIP-152)", ClassicChainConfig.EIP152FBlock},
		{"ECIP-1099 (Etchash)", ClassicChainConfig.ECIP1099FBlock},
		{"Magneto (EIP-2565)", ClassicChainConfig.EIP2565FBlock},
		{"Mystique (EIP-3529)", ClassicChainConfig.EIP3529FBlock},
		{"Spiral (EIP-3651)", ClassicChainConfig.EIP3651FBlock},
	}

	for i := 1; i < len(forks); i++ {
		if forks[i].block == nil {
			t.Fatalf("Classic fork %q has nil activation block", forks[i].name)
		}
		if forks[i-1].block == nil {
			t.Fatalf("Classic fork %q has nil activation block", forks[i-1].name)
		}
		if forks[i].block.Cmp(forks[i-1].block) < 0 {
			t.Errorf("Classic fork ordering violation: %q (block %v) < %q (block %v)",
				forks[i].name, forks[i].block,
				forks[i-1].name, forks[i-1].block)
		}
	}
}

// TestETCMordorForkOrdering verifies that all Mordor testnet forks are
// activated in ascending block order.
func TestETCMordorForkOrdering(t *testing.T) {
	type fork struct {
		name  string
		block *big.Int
	}
	forks := []fork{
		{"Genesis forks (EIP-2/7/150/155/161)", MordorChainConfig.EIP2FBlock},
		{"Agharta (EIP-145)", MordorChainConfig.EIP145FBlock},
		{"Phoenix (EIP-152)", MordorChainConfig.EIP152FBlock},
		{"ECIP-1099 (Etchash)", MordorChainConfig.ECIP1099FBlock},
		{"Magneto (EIP-2565)", MordorChainConfig.EIP2565FBlock},
		{"Mystique (EIP-3529)", MordorChainConfig.EIP3529FBlock},
		{"Spiral (EIP-3651)", MordorChainConfig.EIP3651FBlock},
	}

	for i := 1; i < len(forks); i++ {
		if forks[i].block == nil {
			t.Fatalf("Mordor fork %q has nil activation block", forks[i].name)
		}
		if forks[i-1].block == nil {
			t.Fatalf("Mordor fork %q has nil activation block", forks[i-1].name)
		}
		if forks[i].block.Cmp(forks[i-1].block) < 0 {
			t.Errorf("Mordor fork ordering violation: %q (block %v) < %q (block %v)",
				forks[i].name, forks[i].block,
				forks[i-1].name, forks[i-1].block)
		}
	}
}

// TestETCClassicKnownForkBlocks verifies canonical fork activation blocks on Classic mainnet.
func TestETCClassicKnownForkBlocks(t *testing.T) {
	cases := map[string]struct {
		got  *big.Int
		want int64
	}{
		"Homestead": {ClassicChainConfig.EIP2FBlock, 1150000},
		"EIP-155":   {ClassicChainConfig.EIP155Block, 3000000},
		"ECIP-1017": {ClassicChainConfig.ECIP1017FBlock, 5000000},
		"Atlantis":  {ClassicChainConfig.EIP161FBlock, 8772000},
		"Agharta":   {ClassicChainConfig.EIP145FBlock, 9573000},
		"Phoenix":   {ClassicChainConfig.EIP152FBlock, 10500839},
		"ECIP-1099": {ClassicChainConfig.ECIP1099FBlock, 11700000},
		"Magneto":   {ClassicChainConfig.EIP2565FBlock, 13189133},
		"Mystique":  {ClassicChainConfig.EIP3529FBlock, 14525000},
		"Spiral":    {ClassicChainConfig.EIP3651FBlock, 19250000},
	}
	for name, tc := range cases {
		if tc.got == nil {
			t.Errorf("Classic %s fork block is nil", name)
			continue
		}
		if tc.got.Int64() != tc.want {
			t.Errorf("Classic %s fork block: got %d, want %d", name, tc.got.Int64(), tc.want)
		}
	}
}

// TestETCMordorKnownForkBlocks verifies canonical fork activation blocks on Mordor testnet.
func TestETCMordorKnownForkBlocks(t *testing.T) {
	cases := map[string]struct {
		got  *big.Int
		want int64
	}{
		"Agharta":   {MordorChainConfig.EIP145FBlock, 301243},
		"Phoenix":   {MordorChainConfig.EIP152FBlock, 999983},
		"ECIP-1099": {MordorChainConfig.ECIP1099FBlock, 2520000},
		"Magneto":   {MordorChainConfig.EIP2565FBlock, 3985893},
		"Mystique":  {MordorChainConfig.EIP3529FBlock, 5520000},
		"Spiral":    {MordorChainConfig.EIP3651FBlock, 9957000},
	}
	for name, tc := range cases {
		if tc.got == nil {
			t.Errorf("Mordor %s fork block is nil", name)
			continue
		}
		if tc.got.Int64() != tc.want {
			t.Errorf("Mordor %s fork block: got %d, want %d", name, tc.got.Int64(), tc.want)
		}
	}
}

// TestETCClassicDAORejection verifies ETC rejects the DAO fork.
func TestETCClassicDAORejection(t *testing.T) {
	// EIP-779 (DAO fork) must never be enabled on Classic
	blockNumbers := []*big.Int{
		big.NewInt(0),
		big.NewInt(1920000),
		big.NewInt(10000000),
		big.NewInt(20000000),
	}
	for _, bn := range blockNumbers {
		if ClassicChainConfig.IsEnabled(ClassicChainConfig.GetEthashEIP779Transition, bn) {
			t.Errorf("Classic has DAO fork enabled at block %v", bn)
		}
	}
}

// TestETCClassicECIP1017Config verifies ECIP-1017 emission schedule configuration.
func TestETCClassicECIP1017Config(t *testing.T) {
	if ClassicChainConfig.ECIP1017FBlock.Int64() != 5000000 {
		t.Errorf("ECIP-1017 fork block: got %d, want 5000000", ClassicChainConfig.ECIP1017FBlock.Int64())
	}
	if ClassicChainConfig.ECIP1017EraRounds.Int64() != 5000000 {
		t.Errorf("ECIP-1017 era rounds: got %d, want 5000000", ClassicChainConfig.ECIP1017EraRounds.Int64())
	}
}

// TestETCMordorECIP1017Config verifies Mordor's ECIP-1017 configuration (shorter eras for testing).
func TestETCMordorECIP1017Config(t *testing.T) {
	if MordorChainConfig.ECIP1017FBlock.Int64() != 0 {
		t.Errorf("Mordor ECIP-1017 fork block: got %d, want 0", MordorChainConfig.ECIP1017FBlock.Int64())
	}
	if MordorChainConfig.ECIP1017EraRounds.Int64() != 2000000 {
		t.Errorf("Mordor ECIP-1017 era rounds: got %d, want 2000000", MordorChainConfig.ECIP1017EraRounds.Int64())
	}
}

// TestETCClassicECIP1010Config verifies difficulty bomb pause/continue config.
func TestETCClassicECIP1010Config(t *testing.T) {
	if ClassicChainConfig.ECIP1010PauseBlock == nil {
		t.Fatal("ECIP-1010 PauseBlock is nil")
	}
	if ClassicChainConfig.ECIP1010PauseBlock.Int64() != 3000000 {
		t.Errorf("ECIP-1010 PauseBlock: got %d, want 3000000", ClassicChainConfig.ECIP1010PauseBlock.Int64())
	}
	if ClassicChainConfig.ECIP1010Length == nil {
		t.Fatal("ECIP-1010 Length is nil")
	}
	if ClassicChainConfig.ECIP1010Length.Int64() != 2000000 {
		t.Errorf("ECIP-1010 Length: got %d, want 2000000", ClassicChainConfig.ECIP1010Length.Int64())
	}
}

// TestETCClassicEthashConfig verifies Classic uses Ethash (PoW).
func TestETCClassicEthashConfig(t *testing.T) {
	if ClassicChainConfig.Ethash == nil {
		t.Fatal("Classic chain is not configured for Ethash")
	}
}

// TestETCMordorNoBombPause verifies Mordor does not use ECIP-1010 bomb pause.
func TestETCMordorNoBombPause(t *testing.T) {
	if MordorChainConfig.ECIP1010PauseBlock != nil {
		t.Errorf("Mordor should not have ECIP-1010 bomb pause, got block %v", MordorChainConfig.ECIP1010PauseBlock)
	}
}

// TestETCClassicECBP1100Config verifies MESS (ECBP-1100) activation and deactivation.
func TestETCClassicECBP1100Config(t *testing.T) {
	if ClassicChainConfig.ECBP1100FBlock == nil {
		t.Fatal("ECBP-1100 activation block is nil")
	}
	if ClassicChainConfig.ECBP1100FBlock.Int64() != 11380000 {
		t.Errorf("ECBP-1100 activation: got %d, want 11380000", ClassicChainConfig.ECBP1100FBlock.Int64())
	}
	if ClassicChainConfig.ECBP1100DeactivateFBlock == nil {
		t.Fatal("ECBP-1100 deactivation block is nil")
	}
	if ClassicChainConfig.ECBP1100DeactivateFBlock.Int64() != 19250000 {
		t.Errorf("ECBP-1100 deactivation: got %d, want 19250000", ClassicChainConfig.ECBP1100DeactivateFBlock.Int64())
	}
}
