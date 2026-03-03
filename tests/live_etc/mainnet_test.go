//go:build live

package live_etc

import (
	"math/big"
	"testing"
)

// TestETCMainnetChainID verifies the ETC mainnet reports chain ID 61.
func TestETCMainnetChainID(t *testing.T) {
	client := dialRPC(t, getETCRPC())
	defer client.Close()

	chainID := getChainID(t, client)
	if chainID != ETCMainnetChainID {
		t.Errorf("ETC mainnet chain ID = %d, want %d", chainID, ETCMainnetChainID)
	}
}

// TestETCMainnetGenesisHash verifies the ETC mainnet genesis block hash
// (same as original Ethereum genesis — the chain split happened at block 1,920,000).
func TestETCMainnetGenesisHash(t *testing.T) {
	client := dialRPC(t, getETCRPC())
	defer client.Close()

	genesis := getBlockByNumber(t, client, big.NewInt(0))
	if genesis.Hash != ETCGenesisHash {
		t.Errorf("ETC genesis hash = %s, want %s", genesis.Hash.Hex(), ETCGenesisHash.Hex())
	}
}

// TestETCMainnetPoW verifies that ETC mainnet blocks have valid PoW fields.
func TestETCMainnetPoW(t *testing.T) {
	client := dialRPC(t, getETCRPC())
	defer client.Close()

	block := getBlockByNumber(t, client, nil) // latest
	if block.Difficulty == nil || block.Difficulty.ToInt().Sign() <= 0 {
		t.Error("latest block has zero or nil difficulty — not PoW")
	}
}

// TestETCMainnetECIP1017Era verifies the current ECIP-1017 era based on
// the latest block number. As of 2026, ETC is in era 4 (blocks 20M-25M).
func TestETCMainnetECIP1017Era(t *testing.T) {
	client := dialRPC(t, getETCRPC())
	defer client.Close()

	block := getBlockByNumber(t, client, nil)
	blockNum := block.Number.ToInt().Int64()

	if blockNum < ClassicEraLength {
		t.Skipf("block %d is before era 1 start", blockNum)
	}

	era := (blockNum - 1) / ClassicEraLength
	t.Logf("ETC mainnet block %d is in era %d", blockNum, era)

	// As of early 2026, should be in era 4 (20M+) heading toward era 5 (25M+)
	if era < 3 {
		t.Errorf("expected era >= 3 for current ETC mainnet, got era %d (block %d)", era, blockNum)
	}
}

// TestETCMainnetGasLimit verifies the ETC mainnet gas limit is around 8M.
func TestETCMainnetGasLimit(t *testing.T) {
	client := dialRPC(t, getETCRPC())
	defer client.Close()

	block := getBlockByNumber(t, client, nil)
	gasLimit := block.GasLimit.ToInt().Uint64()

	if gasLimit < 7_000_000 || gasLimit > 9_000_000 {
		t.Errorf("ETC mainnet gas limit = %d, expected ~8M", gasLimit)
	}
}

// TestETCMainnetNetVersion verifies net_version returns "1" (ETC network ID).
func TestETCMainnetNetVersion(t *testing.T) {
	client := dialRPC(t, getETCRPC())
	defer client.Close()

	version := getNetVersion(t, client)
	if version != "1" {
		t.Errorf("ETC mainnet net_version = %q, want %q", version, "1")
	}
}

// TestETCMainnetDAOForkBlock verifies that ETC did NOT execute the DAO fork.
// Block 1,920,000 should have the "classic" state root (no irregular state change).
func TestETCMainnetDAOForkBlock(t *testing.T) {
	client := dialRPC(t, getETCRPC())
	defer client.Close()

	// The DAO fork block
	block := getBlockByNumber(t, client, big.NewInt(1920000))

	if block.Number.ToInt().Int64() != 1920000 {
		t.Fatalf("expected block 1920000, got %d", block.Number.ToInt().Int64())
	}

	// ETC's block 1920000 hash — confirms this chain rejected the DAO fork
	expectedHash := "0x94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f"
	if block.Hash.Hex() != expectedHash {
		t.Errorf("DAO fork block hash = %s, want %s (ETC classic chain)", block.Hash.Hex(), expectedHash)
	}
}
