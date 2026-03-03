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
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

// TestETCForkComplianceClassic verifies that the Classic mainnet chain config
// correctly activates features at each historical fork boundary.
func TestETCForkComplianceClassic(t *testing.T) {
	config := params.ClassicChainConfig

	// Pre-Homestead: DELEGATECALL should not be enabled
	assertForkDisabled(t, config, "Homestead (EIP-7)", 1149999, config.GetEIP7Transition)
	assertForkEnabled(t, config, "Homestead (EIP-7)", 1150000, config.GetEIP7Transition)

	// Tangerine Whistle: gas repricing
	assertForkDisabled(t, config, "Tangerine (EIP-150)", 2499999, config.GetEIP150Transition)
	assertForkEnabled(t, config, "Tangerine (EIP-150)", 2500000, config.GetEIP150Transition)

	// Spurious Dragon: EXP gas repricing
	assertForkDisabled(t, config, "Spurious (EIP-160)", 2999999, config.GetEIP160Transition)
	assertForkEnabled(t, config, "Spurious (EIP-160)", 3000000, config.GetEIP160Transition)

	// Atlantis: Byzantium-equivalent (EIP-161, 170, REVERT, STATICCALL, etc.)
	assertForkDisabled(t, config, "Atlantis (EIP-140)", 8771999, config.GetEIP140Transition)
	assertForkEnabled(t, config, "Atlantis (EIP-140)", 8772000, config.GetEIP140Transition)

	// Agharta: Constantinople-equivalent (SHL/SHR/SAR, CREATE2, EXTCODEHASH)
	assertForkDisabled(t, config, "Agharta (EIP-145)", 9572999, config.GetEIP145Transition)
	assertForkEnabled(t, config, "Agharta (EIP-145)", 9573000, config.GetEIP145Transition)

	// Phoenix: Istanbul-equivalent (CHAINID, SELFBALANCE, Blake2, etc.)
	assertForkDisabled(t, config, "Phoenix (EIP-1344)", 10500838, config.GetEIP1344Transition)
	assertForkEnabled(t, config, "Phoenix (EIP-1344)", 10500839, config.GetEIP1344Transition)

	// Magneto: Berlin-equivalent (EIP-2929 access lists)
	assertForkDisabled(t, config, "Magneto (EIP-2929)", 13189132, config.GetEIP2929Transition)
	assertForkEnabled(t, config, "Magneto (EIP-2929)", 13189133, config.GetEIP2929Transition)

	// Mystique: London-partial (EIP-3529 reduced refunds, EIP-3541)
	assertForkDisabled(t, config, "Mystique (EIP-3529)", 14524999, config.GetEIP3529Transition)
	assertForkEnabled(t, config, "Mystique (EIP-3529)", 14525000, config.GetEIP3529Transition)

	// Spiral: Shanghai-partial (PUSH0, WARM COINBASE, init code limit)
	assertForkDisabled(t, config, "Spiral (EIP-3855/PUSH0)", 19249999, config.GetEIP3855Transition)
	assertForkEnabled(t, config, "Spiral (EIP-3855/PUSH0)", 19250000, config.GetEIP3855Transition)

	// DAO fork must NEVER be active
	assertForkDisabled(t, config, "DAO (EIP-779)", 0, config.GetEthashEIP779Transition)
	assertForkDisabled(t, config, "DAO (EIP-779)", 1920000, config.GetEthashEIP779Transition)
	assertForkDisabled(t, config, "DAO (EIP-779)", 100000000, config.GetEthashEIP779Transition)
}

// TestETCForkComplianceMordor verifies Mordor testnet fork activation.
func TestETCForkComplianceMordor(t *testing.T) {
	config := params.MordorChainConfig

	// Genesis-activated forks (block 0)
	assertForkEnabled(t, config, "Homestead (EIP-7)", 0, config.GetEIP7Transition)
	assertForkEnabled(t, config, "Tangerine (EIP-150)", 0, config.GetEIP150Transition)
	assertForkEnabled(t, config, "Spurious (EIP-160)", 0, config.GetEIP160Transition)
	assertForkEnabled(t, config, "Atlantis (EIP-140)", 0, config.GetEIP140Transition)

	// Agharta
	assertForkDisabled(t, config, "Agharta (EIP-145)", 301242, config.GetEIP145Transition)
	assertForkEnabled(t, config, "Agharta (EIP-145)", 301243, config.GetEIP145Transition)

	// Phoenix
	assertForkDisabled(t, config, "Phoenix (EIP-1344)", 999982, config.GetEIP1344Transition)
	assertForkEnabled(t, config, "Phoenix (EIP-1344)", 999983, config.GetEIP1344Transition)

	// Magneto
	assertForkDisabled(t, config, "Magneto (EIP-2929)", 3985892, config.GetEIP2929Transition)
	assertForkEnabled(t, config, "Magneto (EIP-2929)", 3985893, config.GetEIP2929Transition)

	// Mystique
	assertForkDisabled(t, config, "Mystique (EIP-3529)", 5519999, config.GetEIP3529Transition)
	assertForkEnabled(t, config, "Mystique (EIP-3529)", 5520000, config.GetEIP3529Transition)

	// Spiral
	assertForkDisabled(t, config, "Spiral (EIP-3855/PUSH0)", 9956999, config.GetEIP3855Transition)
	assertForkEnabled(t, config, "Spiral (EIP-3855/PUSH0)", 9957000, config.GetEIP3855Transition)
}

func assertForkEnabled(t *testing.T, config ctypes.ChainConfigurator, name string, block uint64, getter func() *uint64) {
	t.Helper()
	bn := new(big.Int).SetUint64(block)
	if !config.IsEnabled(getter, bn) {
		t.Errorf("%s should be enabled at block %d", name, block)
	}
}

func assertForkDisabled(t *testing.T, config ctypes.ChainConfigurator, name string, block uint64, getter func() *uint64) {
	t.Helper()
	bn := new(big.Int).SetUint64(block)
	if config.IsEnabled(getter, bn) {
		t.Errorf("%s should NOT be enabled at block %d", name, block)
	}
}

// testETCGenesis returns a minimal genesis for ETC chain tests.
func testETCGenesis(config ctypes.ChainConfigurator) *genesisT.Genesis {
	return &genesisT.Genesis{
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
}

// TestETCClassicBlockGeneration verifies that block generation works
// using the Classic chain config with ethash.NewFaker().
func TestETCClassicBlockGeneration(t *testing.T) {
	engine := ethash.NewFaker()
	defer engine.Close()

	db, blocks, _ := GenerateChainWithGenesis(testETCGenesis(params.ClassicChainConfig), engine, 10, nil)
	if len(blocks) != 10 {
		t.Fatalf("expected 10 blocks, got %d", len(blocks))
	}

	for i, block := range blocks {
		if block.NumberU64() != uint64(i+1) {
			t.Errorf("block %d: wrong number %d", i, block.NumberU64())
		}
		if block.GasLimit() == 0 {
			t.Errorf("block %d: zero gas limit", i)
		}
		if block.Difficulty().Sign() <= 0 {
			t.Errorf("block %d: non-positive difficulty %v", i, block.Difficulty())
		}
	}

	blockchain, err := NewBlockChain(db, nil, testETCGenesis(params.ClassicChainConfig), nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatalf("failed to create blockchain: %v", err)
	}
	defer blockchain.Stop()

	if _, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	head := blockchain.CurrentBlock()
	if head.Number.Uint64() != 10 {
		t.Errorf("head block number: got %d, want 10", head.Number.Uint64())
	}
}

// TestETCMordorBlockGeneration verifies block generation with Mordor config.
func TestETCMordorBlockGeneration(t *testing.T) {
	engine := ethash.NewFaker()
	defer engine.Close()

	db, blocks, _ := GenerateChainWithGenesis(testETCGenesis(params.MordorChainConfig), engine, 10, nil)
	if len(blocks) != 10 {
		t.Fatalf("expected 10 blocks, got %d", len(blocks))
	}

	blockchain, err := NewBlockChain(db, nil, testETCGenesis(params.MordorChainConfig), nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatalf("failed to create blockchain: %v", err)
	}
	defer blockchain.Stop()

	if _, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	head := blockchain.CurrentBlock()
	if head.Number.Uint64() != 10 {
		t.Errorf("head block number: got %d, want 10", head.Number.Uint64())
	}
}

// TestETCCoinbaseRewardAtGenesis verifies that block rewards are credited
// to the coinbase address on the Classic chain.
func TestETCCoinbaseRewardAtGenesis(t *testing.T) {
	coinbase := common.HexToAddress("0x71562b71999873DB5b286dF957af199Ec94617F7")

	engine := ethash.NewFaker()
	defer engine.Close()

	db, blocks, _ := GenerateChainWithGenesis(testETCGenesis(params.ClassicChainConfig), engine, 1, func(i int, gen *BlockGen) {
		gen.SetCoinbase(coinbase)
	})

	blockchain, err := NewBlockChain(db, nil, testETCGenesis(params.ClassicChainConfig), nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatalf("failed to create blockchain: %v", err)
	}
	defer blockchain.Stop()

	if _, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	// Block 1 is in Era 0, reward should be 5 ETC (5e18 wei)
	statedb, err := blockchain.State()
	if err != nil {
		t.Fatalf("failed to get state: %v", err)
	}

	balance := statedb.GetBalance(coinbase)
	expectedReward := new(big.Int).SetUint64(5e18) // 5 ETC in wei
	if balance.ToBig().Cmp(expectedReward) != 0 {
		t.Errorf("coinbase balance: got %v, want %v (5 ETC)", balance, expectedReward)
	}
}

// TestETCRequireBlockHashes verifies that Classic and Mordor have checkpoint hashes.
func TestETCRequireBlockHashes(t *testing.T) {
	classicHashes := params.ClassicChainConfig.RequireBlockHashes
	if len(classicHashes) < 2 {
		t.Errorf("Classic should have at least 2 checkpoint hashes, got %d", len(classicHashes))
	}
	if _, ok := classicHashes[1920000]; !ok {
		t.Error("Classic missing checkpoint at block 1920000 (DAO fork)")
	}
	if hash := classicHashes[1920000]; hash != (common.HexToHash("0x94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f")) {
		t.Errorf("Classic block 1920000 hash mismatch: got %s", hash.Hex())
	}

	mordorHashes := params.MordorChainConfig.RequireBlockHashes
	if len(mordorHashes) < 2 {
		t.Errorf("Mordor should have at least 2 checkpoint hashes, got %d", len(mordorHashes))
	}
}

// TestETCClassicNetworkID verifies Classic uses network ID 1.
func TestETCClassicNetworkID(t *testing.T) {
	nid := params.ClassicChainConfig.GetNetworkID()
	if nid == nil || *nid != 1 {
		t.Errorf("Classic network ID: got %v, want 1", nid)
	}
}

// TestETCMordorNetworkID verifies Mordor uses network ID 7.
func TestETCMordorNetworkID(t *testing.T) {
	nid := params.MordorChainConfig.GetNetworkID()
	if nid == nil || *nid != 7 {
		t.Errorf("Mordor network ID: got %v, want 7", nid)
	}
}

// TestETCPrecompileAvailability verifies precompile addresses are available
// at the correct fork blocks on Classic.
func TestETCPrecompileAvailability(t *testing.T) {
	config := params.ClassicChainConfig

	// ModExp (0x05) — EIP-198, activated at Atlantis (8772000)
	assertForkDisabled(t, config, "ModExp (EIP-198)", 8771999, config.GetEIP198Transition)
	assertForkEnabled(t, config, "ModExp (EIP-198)", 8772000, config.GetEIP198Transition)

	// ECAdd (0x06) — EIP-212, activated at Atlantis
	assertForkDisabled(t, config, "ECAdd (EIP-212)", 8771999, config.GetEIP212Transition)
	assertForkEnabled(t, config, "ECAdd (EIP-212)", 8772000, config.GetEIP212Transition)

	// Blake2F (0x09) — EIP-152, activated at Phoenix (10500839)
	assertForkDisabled(t, config, "Blake2F (EIP-152)", 10500838, config.GetEIP152Transition)
	assertForkEnabled(t, config, "Blake2F (EIP-152)", 10500839, config.GetEIP152Transition)
}

// TestETCTransactionTypes verifies transaction type support at fork boundaries.
func TestETCTransactionTypes(t *testing.T) {
	config := params.ClassicChainConfig

	// Type 1 (EIP-2930 access list) — activated at Magneto (13189133)
	assertForkDisabled(t, config, "EIP-2930 (Type 1)", 13189132, config.GetEIP2930Transition)
	assertForkEnabled(t, config, "EIP-2930 (Type 1)", 13189133, config.GetEIP2930Transition)

	// Type 2 (EIP-1559) — NOT activated pre-Olympia on Classic
	if config.IsEnabled(config.GetEIP1559Transition, big.NewInt(19250000)) {
		t.Error("EIP-1559 should NOT be active at Spiral on Classic (pre-Olympia)")
	}
}
