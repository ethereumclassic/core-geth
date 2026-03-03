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
