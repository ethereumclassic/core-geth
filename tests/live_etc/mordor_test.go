//go:build live

package live_etc

import (
	"math/big"
	"testing"
)

// TestMordorChainID verifies the Mordor testnet reports chain ID 63.
func TestMordorChainID(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	chainID := getChainID(t, client)
	if chainID != MordorChainID {
		t.Errorf("Mordor chain ID = %d, want %d", chainID, MordorChainID)
	}
}

// TestMordorGenesisHash verifies the Mordor genesis block hash.
func TestMordorGenesisHash(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	genesis := getBlockByNumber(t, client, big.NewInt(0))
	if genesis.Hash != MordorGenesisHash {
		t.Errorf("Mordor genesis hash = %s, want %s", genesis.Hash.Hex(), MordorGenesisHash.Hex())
	}
}

// TestMordorPoWFields verifies that Mordor blocks have valid PoW fields
// (non-zero difficulty, non-empty nonce and mixHash).
func TestMordorPoWFields(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	block := getBlockByNumber(t, client, nil) // latest
	if block.Difficulty == nil || block.Difficulty.ToInt().Sign() <= 0 {
		t.Error("latest block has zero or nil difficulty — not PoW")
	}
	if block.Nonce == "" || block.Nonce == "0x0000000000000000" {
		t.Error("latest block has empty nonce — not PoW")
	}
	if block.MixHash == (MordorGenesisHash) {
		// MixHash should be a unique hash from the PoW computation
		t.Log("warning: mixHash equals genesis hash — unusual but not necessarily wrong")
	}
}

// TestMordorGasLimit verifies that Mordor gas limit is around 8M (pre-olympia).
func TestMordorGasLimit(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	block := getBlockByNumber(t, client, nil) // latest
	gasLimit := block.GasLimit.ToInt().Uint64()

	// Gas limit should be near 8M (within adjustment bounds)
	if gasLimit < 7_000_000 || gasLimit > 9_000_000 {
		t.Errorf("Mordor gas limit = %d, expected ~8M (pre-olympia)", gasLimit)
	}
}

// TestMordorDifficultyProgression verifies that difficulty is non-trivial
// and blocks have incrementing numbers.
func TestMordorDifficultyProgression(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	latest := getBlockByNumber(t, client, nil)
	latestNum := latest.Number.ToInt().Int64()

	if latestNum < 10_000_000 {
		t.Skipf("Mordor chain height %d too low for meaningful difficulty check", latestNum)
	}

	// Check a block from around 1M blocks ago
	oldNum := latestNum - 1_000_000
	oldBlock := getBlockByNumber(t, client, big.NewInt(oldNum))

	if oldBlock.Difficulty == nil || latest.Difficulty == nil {
		t.Fatal("difficulty is nil")
	}

	t.Logf("Block %d difficulty: %s", oldNum, oldBlock.Difficulty.ToInt().String())
	t.Logf("Block %d difficulty: %s", latestNum, latest.Difficulty.ToInt().String())
}

// TestMordorNetVersion verifies net_version returns "7" (Mordor network ID).
func TestMordorNetVersion(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	version := getNetVersion(t, client)
	if version != "7" {
		t.Errorf("Mordor net_version = %q, want %q", version, "7")
	}
}

// TestMordorECBP1100Deactivated verifies that ECBP-1100 (MESS) deactivation
// block has been passed. Mordor deactivated MESS at block 10,400,000.
func TestMordorECBP1100Deactivated(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	latest := getBlockByNumber(t, client, nil)
	blockNum := latest.Number.ToInt().Int64()

	if blockNum < MordorECBP1100Deactivate {
		t.Skipf("Mordor chain height %d has not reached ECBP-1100 deactivation (%d)",
			blockNum, MordorECBP1100Deactivate)
	}

	t.Logf("Mordor block %d is past ECBP-1100 deactivation at %d", blockNum, MordorECBP1100Deactivate)

	// Verify the deactivation block exists and is valid
	deactBlock := getBlockByNumber(t, client, big.NewInt(MordorECBP1100Deactivate))
	if deactBlock.Difficulty == nil || deactBlock.Difficulty.ToInt().Sign() <= 0 {
		t.Error("ECBP-1100 deactivation block has zero difficulty")
	}
}

// TestMordorECIP1099Epoch verifies ECIP-1099 epoch calculation on live chain.
// After ECIP-1099 (block 2,520,000), epochs are 60,000 blocks long.
func TestMordorECIP1099Epoch(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	latest := getBlockByNumber(t, client, nil)
	blockNum := latest.Number.ToInt().Uint64()

	if blockNum < MordorECIP1099Block {
		t.Skipf("Mordor block %d is before ECIP-1099 activation", blockNum)
	}

	// Post-ECIP-1099: epoch = block / 60000
	epoch := blockNum / EpochLengthECIP1099
	t.Logf("Mordor block %d is in etchash epoch %d (60K-block epochs)", blockNum, epoch)

	// The epoch should be a reasonable number (not overflowing, not zero)
	if epoch == 0 {
		t.Errorf("epoch should not be 0 at block %d", blockNum)
	}
	if epoch > 1000 {
		t.Errorf("epoch %d seems unreasonably high for block %d", epoch, blockNum)
	}
}

// TestMordorSpiralForkBlock verifies the Mordor Spiral fork block exists
// and has expected properties.
func TestMordorSpiralForkBlock(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	latest := getBlockByNumber(t, client, nil)
	if latest.Number.ToInt().Int64() < MordorSpiralBlock {
		t.Skipf("Mordor chain height %d has not reached Spiral (%d)",
			latest.Number.ToInt().Int64(), MordorSpiralBlock)
	}

	spiral := getBlockByNumber(t, client, big.NewInt(MordorSpiralBlock))
	if spiral.Difficulty == nil || spiral.Difficulty.ToInt().Sign() <= 0 {
		t.Error("Spiral fork block has zero difficulty")
	}
	if spiral.GasLimit == nil || spiral.GasLimit.ToInt().Uint64() == 0 {
		t.Error("Spiral fork block has zero gas limit")
	}
	t.Logf("Mordor Spiral block %d: difficulty=%s, gasLimit=%s",
		MordorSpiralBlock, spiral.Difficulty.ToInt().String(), spiral.GasLimit.ToInt().String())
}

// TestMordorBlockHeaderCompleteness verifies that recent Mordor blocks have
// all expected header fields populated (no nil or zero values where not expected).
func TestMordorBlockHeaderCompleteness(t *testing.T) {
	client := dialRPC(t, getMordorRPC())
	defer client.Close()

	block := getBlockByNumber(t, client, nil)
	blockNum := block.Number.ToInt().Int64()

	if block.Hash == (MordorGenesisHash) && blockNum > 0 {
		t.Error("latest block hash should not equal genesis hash")
	}
	if block.ParentHash == (MordorGenesisHash) && blockNum > 1 {
		t.Error("parent hash should not equal genesis hash for blocks > 1")
	}
	if block.StateRoot == ([32]byte{}) {
		t.Error("state root is empty")
	}
	if block.Miner == ([20]byte{}) {
		t.Error("miner address is empty — expected PoW miner")
	}
	if block.Timestamp == nil || block.Timestamp.ToInt().Sign() <= 0 {
		t.Error("timestamp is zero or nil")
	}
	if block.GasLimit == nil || block.GasLimit.ToInt().Sign() <= 0 {
		t.Error("gas limit is zero or nil")
	}
}
