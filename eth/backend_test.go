package eth

import (
	"bytes"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/triedb"
)

func TestMakeExtraDataDefault(t *testing.T) {
	if !bytes.Contains(makeExtraData(nil), []byte("CoreGeth")) {
		t.Error("missing extra data default client identifier")
	}
}

// TestNewResolvesDefaultGenesis checks that a backend given no genesis, on an
// empty database, settles on Ethereum Classic's before it builds the consensus
// engine. The engine takes its ECIP-1099 epoch schedule from that genesis; left
// nil, it would verify every block after the transition with the wrong epoch.
func TestNewResolvesDefaultGenesis(t *testing.T) {
	stack, err := node.New(&node.Config{P2P: p2p.Config{ListenAddr: "127.0.0.1:0", NoDiscovery: true}})
	if err != nil {
		t.Fatal(err)
	}
	defer stack.Close()

	config := ethconfig.Defaults
	config.Ethash.PowMode = ethash.ModeFake
	backend, err := New(stack, &config)
	if err != nil {
		t.Fatal(err)
	}
	if config.Genesis == nil {
		t.Fatal("genesis left nil, want Ethereum Classic's")
	}
	want := params.ClassicChainConfig.GetEthashECIP1099Transition()
	if have := config.Genesis.GetEthashECIP1099Transition(); have == nil || *have != *want {
		t.Errorf("ECIP-1099 transition %v, want %d", have, *want)
	}
	// Ethereum Classic's network ID is 1, not its chain ID. A node that took
	// the chain ID would fail the status handshake with every peer.
	if want := *params.ClassicChainConfig.GetNetworkID(); backend.networkID != want {
		t.Errorf("network ID %d, want %d", backend.networkID, want)
	}
}

// TestCheckUnmaintainedChain checks which stored chains are refused when no
// network flag chooses a genesis: those of Ethereum mainnet, Sepolia and Holesky,
// and no other. Ethereum mainnet and Ethereum Classic share a genesis, so the
// stored chain ID decides between them, and a genesis stored without a config
// is taken for Ethereum Classic's, as core.DefaultGenesisFor takes it.
func TestCheckUnmaintainedChain(t *testing.T) {
	private := &genesisT.Genesis{
		Config: &goethereum.ChainConfig{ChainID: big.NewInt(4269), HomesteadBlock: big.NewInt(0)},
		Alloc:  genesisT.GenesisAlloc{{1}: {Balance: big.NewInt(1)}},
	}
	for _, test := range []struct {
		name    string
		prepare func(ethdb.Database)
		refused string // the network the refusal names, "" for none
	}{
		{"empty database", func(ethdb.Database) {}, ""},
		{"Ethereum Classic", commit(params.DefaultClassicGenesisBlock()), ""},
		{"Mordor", commit(params.DefaultMordorGenesisBlock()), ""},
		{"private network", commit(private), ""},
		{"shared genesis without a stored config", func(db ethdb.Database) {
			block := core.GenesisToBlock(params.DefaultClassicGenesisBlock(), nil)
			rawdb.WriteBlock(db, block)
			rawdb.WriteCanonicalHash(db, block.Hash(), 0)
		}, ""},
		{"Ethereum mainnet", commit(params.DefaultGenesisBlock()), "Ethereum mainnet"},
		{"Sepolia", commit(params.DefaultSepoliaGenesisBlock()), "Sepolia"},
		{"Holesky", commit(params.DefaultHoleskyGenesisBlock()), "Holesky"},
	} {
		db := rawdb.NewMemoryDatabase()
		test.prepare(db)
		err := CheckUnmaintainedChain(db)
		switch {
		case test.refused == "" && err != nil:
			t.Errorf("%s: refused with %v", test.name, err)
		case test.refused != "" && (err == nil || !strings.Contains(err.Error(), "holds "+test.refused+"'s chain")):
			t.Errorf("%s: error %v, want %s's chain refused", test.name, err, test.refused)
		}
	}
}

// commit returns a function committing g to a database, for a test table.
func commit(g *genesisT.Genesis) func(ethdb.Database) {
	return func(db ethdb.Database) { core.MustCommitGenesis(db, triedb.NewDatabase(db, nil), g) }
}

// TestNewRefusesEthereumChain checks that a backend given no genesis refuses a
// database holding Ethereum mainnet's chain, as geth refuses --ethereum. The two
// networks share a genesis, so the backend once ran Ethereum's rules there,
// which this client follows only through the Cancun upgrade.
func TestNewRefusesEthereumChain(t *testing.T) {
	stack, err := node.New(&node.Config{DataDir: t.TempDir(), P2P: p2p.Config{ListenAddr: "127.0.0.1:0", NoDiscovery: true}})
	if err != nil {
		t.Fatal(err)
	}
	defer stack.Close()

	db, err := stack.OpenDatabaseWithFreezer("chaindata", 16, 16, "", "", false)
	if err != nil {
		t.Fatal(err)
	}
	core.MustCommitGenesis(db, triedb.NewDatabase(db, nil), params.DefaultGenesisBlock())
	db.Close()

	config := ethconfig.Defaults
	config.Ethash.PowMode = ethash.ModeFake
	if _, err := New(stack, &config); err == nil || !strings.Contains(err.Error(), "holds Ethereum mainnet's chain") {
		t.Fatalf("error %v, want the database refused", err)
	}
}
