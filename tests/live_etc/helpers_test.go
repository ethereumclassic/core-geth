//go:build live

// Package live_etc provides pre-fork integration tests that verify current
// ETC and Mordor chain state via JSON-RPC. These tests require a running
// core-geth node and are excluded from normal test runs.
//
// Run with: go test -tags live ./tests/live_etc/ -v
//
// Environment variables:
//
//	MORDOR_RPC  - Mordor RPC endpoint (default: http://localhost:8545)
//	ETC_RPC     - ETC mainnet RPC endpoint (default: https://etc.rivet.link)
package live_etc

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// Pre-olympia chain constants (NO olympia-specific values here)
const (
	MordorChainID     = 63
	ETCMainnetChainID = 61

	// ECIP-1017 era parameters
	ClassicEraLength = 5_000_000
	MordorEraLength  = 2_000_000

	// ECIP-1099 etchash fork blocks
	ClassicECIP1099Block = 11_700_000
	MordorECIP1099Block  = 2_520_000

	// ECBP-1100 (MESS) activation windows
	MordorECBP1100Activate   = 2_380_000
	MordorECBP1100Deactivate = 10_400_000

	ClassicECBP1100Activate   = 11_380_000
	ClassicECBP1100Deactivate = 19_250_000

	// Spiral fork blocks
	MordorSpiralBlock  = 9_957_000
	ClassicSpiralBlock = 19_250_000

	// Gas limits (pre-olympia)
	ETCGasLimit = 8_000_000

	// Etchash epoch lengths
	EpochLengthDefault  = 30_000
	EpochLengthECIP1099 = 60_000
)

// Known genesis hashes
var (
	MordorGenesisHash = common.HexToHash("0xa68ebde7932f0bf2579b075499416f0a693de84c26b05cd01de86e60aad05ec0")
	ETCGenesisHash    = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
)

// rpcBlock is a minimal block structure for JSON-RPC responses.
type rpcBlock struct {
	Number     *hexutil.Big   `json:"number"`
	Hash       common.Hash    `json:"hash"`
	ParentHash common.Hash    `json:"parentHash"`
	GasUsed    *hexutil.Big   `json:"gasUsed"`
	GasLimit   *hexutil.Big   `json:"gasLimit"`
	StateRoot  common.Hash    `json:"stateRoot"`
	Miner      common.Address `json:"miner"`
	Difficulty *hexutil.Big   `json:"difficulty"`
	Nonce      string         `json:"nonce"`
	MixHash    common.Hash    `json:"mixHash"`
	Timestamp  *hexutil.Big   `json:"timestamp"`
}

// getMordorRPC returns the Mordor RPC endpoint.
func getMordorRPC() string {
	if v := os.Getenv("MORDOR_RPC"); v != "" {
		return v
	}
	return "http://localhost:8545"
}

// getETCRPC returns the ETC mainnet RPC endpoint.
func getETCRPC() string {
	if v := os.Getenv("ETC_RPC"); v != "" {
		return v
	}
	return "https://etc.rivet.link"
}

// dialRPC connects to an RPC endpoint with a 10-second timeout.
func dialRPC(t *testing.T, endpoint string) *rpc.Client {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		t.Skipf("cannot connect to %s: %v", endpoint, err)
	}
	return client
}

// getChainID returns the chain ID from eth_chainId.
func getChainID(t *testing.T, client *rpc.Client) uint64 {
	t.Helper()
	var result hexutil.Big
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.CallContext(ctx, &result, "eth_chainId"); err != nil {
		t.Fatalf("eth_chainId failed: %v", err)
	}
	return result.ToInt().Uint64()
}

// getBlockByNumber fetches a block by number.
func getBlockByNumber(t *testing.T, client *rpc.Client, num *big.Int) *rpcBlock {
	t.Helper()
	var block rpcBlock
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var numArg string
	if num == nil {
		numArg = "latest"
	} else {
		numArg = hexutil.EncodeBig(num)
	}
	if err := client.CallContext(ctx, &block, "eth_getBlockByNumber", numArg, false); err != nil {
		t.Fatalf("eth_getBlockByNumber(%s) failed: %v", numArg, err)
	}
	return &block
}

// getNetVersion returns the network version from net_version.
func getNetVersion(t *testing.T, client *rpc.Client) string {
	t.Helper()
	var result string
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.CallContext(ctx, &result, "net_version"); err != nil {
		t.Fatalf("net_version failed: %v", err)
	}
	return result
}

// getSyncing returns whether the node is syncing.
func getSyncing(t *testing.T, client *rpc.Client) json.RawMessage {
	t.Helper()
	var result json.RawMessage
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.CallContext(ctx, &result, "eth_syncing"); err != nil {
		t.Fatalf("eth_syncing failed: %v", err)
	}
	return result
}
