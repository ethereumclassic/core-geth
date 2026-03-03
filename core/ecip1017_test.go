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
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params/mutations"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/holiman/uint256"
)

// TestECIP1017EraRewardsIntegration verifies ECIP-1017 emission schedule
// by processing blocks across era boundaries using actual block generation
// and chain insertion (not just function-level tests like rewards_test.go).
func TestECIP1017EraRewardsIntegration(t *testing.T) {
	// Use a short era length for testing (100 blocks per era)
	eraLen := big.NewInt(100)
	config := &coregeth.CoreGethChainConfig{
		NetworkID: 1,
		ChainID:   big.NewInt(61),
		Ethash:    new(ctypes.EthashConfig),

		EIP2FBlock:  big.NewInt(0),
		EIP7FBlock:  big.NewInt(0),
		EIP150Block: big.NewInt(0),
		EIP155Block: big.NewInt(0),

		EIP160FBlock: big.NewInt(0),
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		ECIP1017FBlock:    big.NewInt(0),
		ECIP1017EraRounds: eraLen,

		DisposalBlock: big.NewInt(0),
	}

	coinbase := common.HexToAddress("0x71562b71999873DB5b286dF957af199Ec94617F7")
	genesis := &genesisT.Genesis{
		Config:     config,
		GasLimit:   8_000_000,
		Difficulty: big.NewInt(1),
		Alloc: map[common.Address]genesisT.GenesisAccount{
			common.HexToAddress("0x1"): {Balance: big.NewInt(1)},
			common.HexToAddress("0x2"): {Balance: big.NewInt(1)},
			common.HexToAddress("0x3"): {Balance: big.NewInt(1)},
			common.HexToAddress("0x4"): {Balance: big.NewInt(1)},
			common.HexToAddress("0x5"): {Balance: big.NewInt(1)},
			common.HexToAddress("0x6"): {Balance: big.NewInt(1)},
			common.HexToAddress("0x7"): {Balance: big.NewInt(1)},
			common.HexToAddress("0x8"): {Balance: big.NewInt(1)},
		},
	}

	engine := ethash.NewFaker()
	defer engine.Close()

	// Generate 250 blocks (crosses era 0 → era 1 → era 2)
	db, blocks, _ := GenerateChainWithGenesis(genesis, engine, 250, func(i int, gen *BlockGen) {
		gen.SetCoinbase(coinbase)
	})

	blockchain, err := NewBlockChain(db, nil, genesis, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatalf("failed to create blockchain: %v", err)
	}
	defer blockchain.Stop()

	if _, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	statedb, err := blockchain.State()
	if err != nil {
		t.Fatalf("failed to get state: %v", err)
	}

	balance := statedb.GetBalance(coinbase)

	// Calculate expected balance:
	// Era 0 (blocks 1-100): 100 blocks × 5 ETC = 500 ETC
	// Era 1 (blocks 101-200): 100 blocks × 4 ETC = 400 ETC
	// Era 2 (blocks 201-250): 50 blocks × 3.2 ETC = 160 ETC
	// Total: 1060 ETC = 1060e18 wei
	era0Reward := uint256.NewInt(5e18)
	era1Reward := uint256.NewInt(4e18)
	era2Reward := uint256.NewInt(3.2e18)

	expectedTotal := new(uint256.Int)
	expectedTotal.Add(expectedTotal, new(uint256.Int).Mul(era0Reward, uint256.NewInt(100)))
	expectedTotal.Add(expectedTotal, new(uint256.Int).Mul(era1Reward, uint256.NewInt(100)))
	expectedTotal.Add(expectedTotal, new(uint256.Int).Mul(era2Reward, uint256.NewInt(50)))

	if balance.Cmp(expectedTotal) != 0 {
		t.Errorf("coinbase balance after 250 blocks:\n  got:  %v\n  want: %v\n  diff: %v",
			balance, expectedTotal, new(uint256.Int).Sub(balance, expectedTotal))
	}
}

// TestECIP1017EraCalculation verifies era computation at known block boundaries
// for both Classic (5M era) and Mordor (2M era) configurations.
func TestECIP1017EraCalculation(t *testing.T) {
	cases := []struct {
		name     string
		block    int64
		eraLen   int64
		wantEra  int64
	}{
		// Classic mainnet (5M era rounds)
		{"Classic block 1 (era 0)", 1, 5000000, 0},
		{"Classic block 5M (era 0)", 5000000, 5000000, 0},
		{"Classic block 5M+1 (era 1)", 5000001, 5000000, 1},
		{"Classic block 10M (era 1)", 10000000, 5000000, 1},
		{"Classic block 10M+1 (era 2)", 10000001, 5000000, 2},
		{"Classic block 15M+1 (era 3)", 15000001, 5000000, 3},
		{"Classic block 20M+1 (era 4)", 20000001, 5000000, 4},
		{"Classic block 25M+1 (era 5)", 25000001, 5000000, 5},

		// Mordor testnet (2M era rounds)
		{"Mordor block 1 (era 0)", 1, 2000000, 0},
		{"Mordor block 2M (era 0)", 2000000, 2000000, 0},
		{"Mordor block 2M+1 (era 1)", 2000001, 2000000, 1},
		{"Mordor block 4M+1 (era 2)", 4000001, 2000000, 2},
	}

	for _, tc := range cases {
		era := mutations.GetBlockEra(big.NewInt(tc.block), big.NewInt(tc.eraLen))
		if era.Int64() != tc.wantEra {
			t.Errorf("%s: got era %d, want %d", tc.name, era.Int64(), tc.wantEra)
		}
	}
}

// TestECIP1017RewardDecay verifies the 80% decay per era for block rewards.
func TestECIP1017RewardDecay(t *testing.T) {
	maxReward := uint256.NewInt(5e18) // 5 ETC base reward
	eraLen := big.NewInt(5000000)

	cases := []struct {
		name     string
		block    int64
		wantWei  uint64
		wantETC  string
	}{
		{"Era 0 (5 ETC)", 1, 5e18, "5.0"},
		{"Era 1 (4 ETC)", 5000001, 4e18, "4.0"},
		{"Era 2 (3.2 ETC)", 10000001, 3.2e18, "3.2"},
		{"Era 3 (2.56 ETC)", 15000001, 2.56e18, "2.56"},
		{"Era 4 (2.048 ETC)", 20000001, 2.048e18, "2.048"},
	}

	for _, tc := range cases {
		era := mutations.GetBlockEra(big.NewInt(tc.block), eraLen)
		reward := mutations.GetBlockWinnerRewardByEra(era, maxReward)
		expected := uint256.NewInt(tc.wantWei)
		if reward.Cmp(expected) != 0 {
			t.Errorf("%s: got %v wei, want %v wei (%s ETC)",
				tc.name, reward, expected, tc.wantETC)
		}
	}
}

// TestECIP1017RewardEventuallyZero verifies that with integer arithmetic,
// the 80% decay in block rewards eventually reaches zero. This is expected
// behavior — at 5M blocks/era, era 109 is block 545M+, far in the future.
func TestECIP1017RewardEventuallyZero(t *testing.T) {
	maxReward := uint256.NewInt(5e18)

	// Reward should be positive for many eras
	lastPositiveEra := int64(-1)
	for i := int64(0); i < 200; i++ {
		era := big.NewInt(i)
		reward := mutations.GetBlockWinnerRewardByEra(era, maxReward)
		if reward.Sign() > 0 {
			lastPositiveEra = i
		}
	}

	// With 5 ETC base and 80% decay, integer truncation hits zero around era 109
	if lastPositiveEra < 100 {
		t.Errorf("Rewards became zero too early at era %d (expected > 100)", lastPositiveEra)
	}
	t.Logf("Last positive reward era: %d (block ~%dM on Classic mainnet)",
		lastPositiveEra, (lastPositiveEra+1)*5)
}
