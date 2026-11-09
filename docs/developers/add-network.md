---
description: "Add a network to Core-Geth: define its chain configuration and genesis block with a worked example, and start a node on it."
---

# Adding a network to Core-Geth

A network is its chain configuration and its genesis block. This page defines an example network,
AlphaBeta Coin (ABC), and starts a node on it:

- **ABC is proof of work**, with Ethash.
- **Its genesis block allocates 42 wei to one address.**
- **Every EIP that `ClassicChainConfig` in `params/config_classic.go` sets is active from block 0.**
  Ethereum Classic's own choices are left out: Etchash, the ECIP-1017 monetary policy and MESS stay
  off, and the difficulty bomb is disposed of at block 0.

Core-Geth reads a network definition in two forms:

1. **Go values in `params/`**, beside the Ethereum Classic and Mordor presets. A network flag, such as
   `--mordor`, selects a preset. The references to `MordorFlag` in `cmd/utils/flags.go` and
   `cmd/geth/main.go` show where a network flag is wired; adding one is outside this page.
2. **A JSON genesis file**, which `geth init` writes into a new database. A node runs a network from
   this file alone, with no change to the code.

## Define the network in `params/`

The files follow the Mordor preset, `config_mordor.go`, `genesis_mordor.go` and
`bootnodes_mordor.go`, and the MintMe preset's genesis hash test, `config_mintme_test.go`.

### Chain configuration

`params/config_abc.go`:

```go
package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

var (
	// ABCChainConfig is the chain parameters to run a node on the AlphaBeta Coin network (PoW).
	ABCChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID:                 4269,
		ChainID:                   big.NewInt(4269),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,
		Ethash:                    new(ctypes.EthashConfig),

		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		EIP150Block: big.NewInt(0),

		EIP155Block: big.NewInt(0),

		// EIP158 eq
		EIP160FBlock: big.NewInt(0),
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(0),
		EIP1014FBlock: big.NewInt(0),
		EIP1052FBlock: big.NewInt(0),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0), // RePetersburg (== re-1283)

		// Berlin eq, aka Magneto
		EIP2565FBlock: big.NewInt(0),
		EIP2718FBlock: big.NewInt(0),
		EIP2929FBlock: big.NewInt(0),
		EIP2930FBlock: big.NewInt(0),

		// London (partially), aka Mystique
		EIP3529FBlock: big.NewInt(0),
		EIP3541FBlock: big.NewInt(0),

		// Spiral, aka Shanghai (partially)
		EIP3651FBlock: big.NewInt(0), // Warm COINBASE (gas reprice)
		EIP3855FBlock: big.NewInt(0), // PUSH0 instruction
		EIP3860FBlock: big.NewInt(0), // Limit and meter initcode
		EIP6049FBlock: big.NewInt(0), // Deprecate SELFDESTRUCT (noop)

		ECIP1099FBlock: nil, // Etchash

		DisposalBlock:      big.NewInt(0), // Dispose difficulty bomb
		ECIP1017FBlock:     nil,           // Ethereum Classic's disinflationary monetary policy
		ECIP1017EraRounds:  nil,
		ECIP1010PauseBlock: nil, // No need to delay difficulty bomb, is defused by default
		ECIP1010Length:     nil,
		ECBP1100FBlock:     nil, // ECBP1100 (MESS artificial finality)
		RequireBlockHashes: map[uint64]common.Hash{},
	}
)
```

A field left `nil` is never active. `CoreGethChainConfig` in `params/types/coregeth/chain_config.go`
defines every field and its JSON name.

### Genesis block

`params/genesis_abc.go`:

```go
package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

var ABCGenesisHash = common.HexToHash("0x5f32ce1ed875a04d74361164bcfcc24af721df8e616642486338300a691fe582")

func DefaultABCGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     ABCChainConfig,
		Nonce:      hexutil.MustDecodeUint64("0x0"),
		ExtraData:  hexutil.MustDecode("0x42"),
		GasLimit:   hexutil.MustDecodeUint64("0x2fefd8"),
		Difficulty: hexutil.MustDecodeBig("0x20000"),
		Timestamp:  hexutil.MustDecodeUint64("0x6048d57c"),
		Alloc: genesisT.GenesisAlloc{
			common.HexToAddress("366ae7da62294427c764870bd2a460d7ded29d30"): genesisT.GenesisAccount{
				Balance: big.NewInt(42),
			},
		},
	}
}
```

### Bootnodes

`params/bootnodes_abc.go`:

```go
package params

// ABCBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the AlphaBeta Coin network.
var ABCBootnodes = []string{
	"enode://82ce42238ce90803f2506ae98c433ddb1fa6d35c48786e0f830ce331b170624e2d130fac0dffab6af25e4a6940cbf689766cd0e09e83d295f27c2404f8ec7ac9@203.0.113.10:30303",
}
```

`203.0.113.10` is an address reserved for documentation; a real network lists the enode URLs of its
own bootnodes. `ABCBootnodes` is read only once a network flag selects it, in `setBootstrapNodes`
(`cmd/utils/flags.go`). A node started from the JSON genesis takes its bootnodes from `--bootnodes`,
and without that flag uses Ethereum Classic's, as [Establish a network](#establish-a-network) shows.

### Tests

`params/config_abc_test.go` checks the genesis hash, with the `genesisToBlock` helper in
`params/config_test.go`, and checks that every bootnode entry parses:

```go
package params

import (
	"testing"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

// TestGenesisHashABC tests that ABCGenesisHash is the correct value for the genesis configuration.
func TestGenesisHashABC(t *testing.T) {
	genesis := DefaultABCGenesisBlock()
	block := genesisToBlock(genesis, nil)
	if block.Hash() != ABCGenesisHash {
		t.Errorf("want: %s, got: %s", ABCGenesisHash.Hex(), block.Hash().Hex())
	}
}

// TestABCBootnodes tests that every ABCBootnodes entry parses as an enode URL.
func TestABCBootnodes(t *testing.T) {
	for _, url := range ABCBootnodes {
		if _, err := enode.Parse(enode.ValidSchemes, url); err != nil {
			t.Errorf("%s: %v", url, err)
		}
	}
}
```

`params/example_abc_test.go` checks that `DefaultABCGenesisBlock()` encodes to exactly the JSON in
the next section. `go test` vets example names, so the example is named for the function it
documents: an example whose name matches no identifier fails the build.

```go
package params

import (
	"encoding/json"
	"fmt"
	"log"
)

func ExampleDefaultABCGenesisBlock() {
	genesis := DefaultABCGenesisBlock()
	jsonBytes, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "config": {
	//     "networkId": 4269,
	//     "chainId": 4269,
	//     "supportedProtocolVersions": [
	//       68
	//     ],
	//     "eip2FBlock": 0,
	//     "eip7FBlock": 0,
	//     "eip150Block": 0,
	//     "eip155Block": 0,
	//     "eip160Block": 0,
	//     "eip161FBlock": 0,
	//     "eip170FBlock": 0,
	//     "eip100FBlock": 0,
	//     "eip140FBlock": 0,
	//     "eip198FBlock": 0,
	//     "eip211FBlock": 0,
	//     "eip212FBlock": 0,
	//     "eip213FBlock": 0,
	//     "eip214FBlock": 0,
	//     "eip658FBlock": 0,
	//     "eip145FBlock": 0,
	//     "eip1014FBlock": 0,
	//     "eip1052FBlock": 0,
	//     "eip152FBlock": 0,
	//     "eip1108FBlock": 0,
	//     "eip1344FBlock": 0,
	//     "eip1884FBlock": 0,
	//     "eip2028FBlock": 0,
	//     "eip2200FBlock": 0,
	//     "eip2565FBlock": 0,
	//     "eip2718FBlock": 0,
	//     "eip2929FBlock": 0,
	//     "eip2930FBlock": 0,
	//     "eip3541FBlock": 0,
	//     "eip3529FBlock": 0,
	//     "eip3651FBlock": 0,
	//     "eip3855FBlock": 0,
	//     "eip3860FBlock": 0,
	//     "eip6049FBlock": 0,
	//     "disposalBlock": 0,
	//     "ethash": {},
	//     "requireBlockHashes": {}
	//   },
	//   "nonce": "0x0",
	//   "timestamp": "0x6048d57c",
	//   "extraData": "0x42",
	//   "gasLimit": "0x2fefd8",
	//   "difficulty": "0x20000",
	//   "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	//   "coinbase": "0x0000000000000000000000000000000000000000",
	//   "alloc": {
	//     "366ae7da62294427c764870bd2a460d7ded29d30": {
	//       "balance": "0x2a"
	//     }
	//   },
	//   "number": "0x0",
	//   "gasUsed": "0x0",
	//   "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
	// }
}
```

From the repository root:

```sh
go test ./params/ -run 'TestGenesisHashABC|TestABCBootnodes|ExampleDefaultABCGenesisBlock' -count=1 -v
```

```
=== RUN   TestGenesisHashABC
--- PASS: TestGenesisHashABC (0.00s)
=== RUN   TestABCBootnodes
--- PASS: TestABCBootnodes (0.00s)
=== RUN   ExampleDefaultABCGenesisBlock
--- PASS: ExampleDefaultABCGenesisBlock (0.00s)
PASS
ok  	github.com/ethereum/go-ethereum/params	0.029s
```

## Write the genesis as JSON

Save the example test's output as `abc_genesis.json`:

```json
{
  "config": {
    "networkId": 4269,
    "chainId": 4269,
    "supportedProtocolVersions": [
      68
    ],
    "eip2FBlock": 0,
    "eip7FBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip160Block": 0,
    "eip161FBlock": 0,
    "eip170FBlock": 0,
    "eip100FBlock": 0,
    "eip140FBlock": 0,
    "eip198FBlock": 0,
    "eip211FBlock": 0,
    "eip212FBlock": 0,
    "eip213FBlock": 0,
    "eip214FBlock": 0,
    "eip658FBlock": 0,
    "eip145FBlock": 0,
    "eip1014FBlock": 0,
    "eip1052FBlock": 0,
    "eip152FBlock": 0,
    "eip1108FBlock": 0,
    "eip1344FBlock": 0,
    "eip1884FBlock": 0,
    "eip2028FBlock": 0,
    "eip2200FBlock": 0,
    "eip2565FBlock": 0,
    "eip2718FBlock": 0,
    "eip2929FBlock": 0,
    "eip2930FBlock": 0,
    "eip3541FBlock": 0,
    "eip3529FBlock": 0,
    "eip3651FBlock": 0,
    "eip3855FBlock": 0,
    "eip3860FBlock": 0,
    "eip6049FBlock": 0,
    "disposalBlock": 0,
    "ethash": {},
    "requireBlockHashes": {}
  },
  "nonce": "0x0",
  "timestamp": "0x6048d57c",
  "extraData": "0x42",
  "gasLimit": "0x2fefd8",
  "difficulty": "0x20000",
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "coinbase": "0x0000000000000000000000000000000000000000",
  "alloc": {
    "366ae7da62294427c764870bd2a460d7ded29d30": {
      "balance": "0x2a"
    }
  },
  "number": "0x0",
  "gasUsed": "0x0",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
}
```

## Initialize a database

Build `geth` as [Build from source](build-from-source.md) describes. Then, from the repository root:

```sh
./build/bin/geth --datadir ./abc-datadir init abc_genesis.json
```

Among its output:

```
INFO [09-13|02:11:42.995] Writing custom genesis block
INFO [09-13|02:11:43.086] Wrote genesis block OK                   config="NetworkID: 4269, ChainID: 4269 Engine: ethash EIP1014: 0 EIP1052: 0 EIP1108: 0 EIP1344: 0 EIP140: 0 EIP145: 0 EIP150: 0 EIP152: 0 EIP155: 0 EIP160: 0 EIP161abc: 0 EIP161d: 0 EIP170: 0 EIP1884: 0 EIP198: 0 EIP2028: 0 EIP211: 0 EIP212: 0 EIP213: 0 EIP214: 0 EIP2200: 0 EIP2565: 0 EIP2718: 0 EIP2929: 0 EIP2930: 0 EIP2: 0 EIP3529: 0 EIP3541: 0 EIP3651: 0 EIP3855: 0 EIP3860: 0 EIP6049: 0 EIP658: 0 EIP7: 0 EthashECIP1041: 0 EthashEIP100B: 0 EthashHomestead: 0 "
INFO [09-13|02:11:43.087] Successfully wrote genesis state         database=chaindata hash=5f32ce..1fe582
```

The hash is `ABCGenesisHash`, shortened. To read back what was stored, without starting a node:

```sh
./build/bin/geth --datadir ./abc-datadir dumpgenesis
```

It prints the stored genesis, with the same values as `abc_genesis.json`.

## Start a node

```sh
./build/bin/geth --datadir ./abc-datadir --networkid 4269 --nodiscover --port 30471
```

- **`--networkid 4269`** sets the network ID. Without it, and with no network flag, the node takes its
  network ID from the genesis `chainId` rather than its `networkId`, and its first log line reads
  `Starting Core-Geth on Ethereum Classic...`.
- **`--nodiscover`** keeps this node from looking for peers.
- **`--port`** moves the peer-to-peer listener off its
  [default](../operate/security.md#ports-and-listeners), so the node can run beside another node on
  the same host.

Among its output:

```
INFO [09-13|02:11:43.587] Initialising Ethereum protocol           network=4269 dbversion=<nil>
INFO [09-13|02:11:43.592] Found stored genesis block               config="NetworkID: 4269, ChainID: 4269 Engine: ethash EIP1014: 0 EIP1052: 0 EIP1108: 0 EIP1344: 0 EIP140: 0 EIP145: 0 EIP150: 0 EIP152: 0 EIP155: 0 EIP160: 0 EIP161abc: 0 EIP161d: 0 EIP170: 0 EIP1884: 0 EIP198: 0 EIP2028: 0 EIP211: 0 EIP212: 0 EIP213: 0 EIP214: 0 EIP2200: 0 EIP2565: 0 EIP2718: 0 EIP2929: 0 EIP2930: 0 EIP2: 0 EIP3529: 0 EIP3541: 0 EIP3651: 0 EIP3855: 0 EIP3860: 0 EIP6049: 0 EIP658: 0 EIP7: 0 EthashECIP1041: 0 EthashEIP100B: 0 EthashHomestead: 0 "
INFO [09-13|02:11:43.593] Loaded most recent local block           number=0 hash=5f32ce..1fe582 td=131,072 age=5y7mo2d
```

Stop the node with Ctrl-C.

## Establish a network

With no network flag, a node applies Ethereum Classic's peer discovery defaults.
[Private network](../tutorials/private-network.md#before-you-start) gives the flags that override
them; on this network every node also takes `--networkid 4269`.

`geth dumpconfig` with those flags shows the peer settings a node will use, under
`BootstrapNodes` and `EthDiscoveryURLs`. [Private network](../tutorials/private-network.md) walks
through running a bootnode and member nodes.
