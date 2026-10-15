---
description: "Build a private network with Core-Geth on one machine: a custom genesis block, proof-of-work mining with a CPU miner, and connecting your own nodes."
---

# Tutorial: Operating a Private Network

This tutorial builds a private network on one machine: a chain with its own genesis block that
only your nodes run. The numbered steps use proof of work, with a CPU miner.
[Proof of authority with Clique](#proof-of-authority-with-clique) replaces the miner with signers.
Every command on this page was
[run before it was published](#this-tutorial-was-run-before-it-was-published), and the output
shown is from that run. To build a network into the client instead, with a flag of its own, see
[Add a network](../developers/add-network.md).

**The steps:** [1. Genesis](#1-define-the-genesis-state) ·
[2. Initialize](#2-initialize-every-data-directory) ·
[3. Bootnode](#3-start-a-bootnode) ·
[4. Member node](#4-start-a-member-node) ·
[5. Miner](#5-start-a-miner) ·
[6. Confirm](#6-confirm-the-network-works) ·
[7. Stop](#7-stop-the-network)

**Signers instead of a miner:** [Proof of authority with Clique](#proof-of-authority-with-clique)

## Before you start

You need `geth` and `bootnode` from the same release. `bootnode` ships in the `alltools` archive
([What else is in a release](../getting-started/installation.md#what-else-is-in-a-release)), and a
[build from source](../developers/build-from-source.md#source) produces it too.

**A private network runs with no network flag.** With no network flag, `geth` runs Ethereum
Classic. A data directory already initialized with another genesis keeps running that chain
instead ([Supported networks](../index.md#supported-networks)). Three things follow, and every
command below is written for them:

- **Every command for a node names its `--datadir`.** `geth init` without `--datadir` writes your
  genesis into the default Ethereum Classic data directory. A node started on a directory nobody
  initialized writes Ethereum Classic's genesis there, and becomes an Ethereum Classic node, chain
  ID 61.
- **Ethereum Classic's network defaults still apply**, so each node overrides them.
  `--bootnodes` replaces its bootnodes with yours, `--discovery.dns ""` switches off its DNS
  discovery lists, and `--nat none` stops `geth` looking up its public address and asking your
  router to forward its ports. A `geth` node that is itself the bootnode takes `--bootnodes ""`,
  so it contacts no Ethereum Classic bootnode.
- **`--networkid` is the chain ID from `genesis.json`.** Without it, `geth` takes the network ID
  from the chain ID anyway, but its first log line reads `Starting Core-Geth on Ethereum Classic...`.

**Every node on one machine needs its own ports.** A node started on a port another node holds
exits with `bind: address already in use`. The ports on
this page are arbitrary, and stay clear of the ports a node takes by default
([Ports and listeners](../operate/security.md#ports-and-listeners)).

## 1. Define the genesis state

Every node on the network starts from the same genesis block, defined in a JSON file. Save this
as `genesis.json`:

```json
{
  "config": {
    "chainId": 12345,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "ethash": {}
  },
  "alloc": {},
  "coinbase": "0x0000000000000000000000000000000000000000",
  "difficulty": "0x20000",
  "extraData": "",
  "gasLimit": "0x2fefd8",
  "nonce": "0x0000000000000042",
  "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp": "0x00"
}
```

- **`chainId`** is your chain's ID, and every node passes it as `--networkid`. Choose a value no
  public network uses.
- **`config`** lists the protocol upgrades the chain runs from its first block. This example stops
  at Petersburg; the JSON genesis on [Add a network](../developers/add-network.md#write-the-genesis-as-json)
  activates, from block 0, every EIP that `ClassicChainConfig` sets. **`"ethash": {}`** selects
  proof of work, which step 2 confirms.
- **`nonce`** makes the genesis block yours. Nodes whose genesis block differs from yours cannot
  peer with your nodes (at `--verbosity 4`, both sides log `genesis mismatch`), and nodes with an
  identical one can, so change it to a value of your own.
- **`alloc`** funds accounts at genesis, with balances in wei:

```json
"alloc": {
  "<address>": { "balance": "1000000000000000000" }
}
```

## 2. Initialize every data directory

Write the genesis into each node's data directory before that node first starts. This tutorial
runs two nodes, `node1` and `node2`:

```bash
geth init --datadir node1 genesis.json
geth init --datadir node2 genesis.json
```

Each `init` reports the genesis it wrote, then repeats the lines for `lightchaindata`:

```text
INFO [09-13|02:27:24.587] Writing custom genesis block
INFO [09-13|02:27:24.668] Wrote genesis block OK                   config="Chain ID:  12345 (unknown)\nConsensus: Ethash (proof-of-work)\n\n_ Block-based Forks: []_ Time-based Forks: []_ TTD: <nil>"
INFO [09-13|02:27:24.669] Successfully wrote genesis state         database=chaindata hash=5e1fc7..d790e0
```

**Check for `Chain ID:  12345` and `Consensus: Ethash (proof-of-work)`.** Without the `ethash`
object, the same file reports `Consensus: unknown`.

## 3. Start a bootnode

[Bootnodes and peer discovery](../guides/bootnodes-and-discovery.md) covers bootnodes on the public networks.

The bootnode is how the nodes find each other. Generate its key once, then start it, and leave it
running:

```bash
bootnode --genkey=boot.key
bootnode --nodekey=boot.key --addr=:30501
```

It prints its [`enode` URL](https://ethereum.org/developers/docs/networking-layer/network-addresses/#enode):

```text
enode://<node ID>@127.0.0.1:0?discport=30501
Note: you're using cmd/bootnode, a developer tool.
We recommend using a regular node as bootstrap node for production deployments.
```

- **The URL names `127.0.0.1`**, which reaches nodes on the same machine. For nodes on other
  machines, replace it with an address they can reach this one on.
- **The bootnode listens on UDP only.** Allow UDP on its port through any firewall between the
  machines. A TCP check such as `telnet` cannot connect to it, whether it is running or not.
- **The same `boot.key` gives the same URL** every time the bootnode starts.

## 4. Start a member node

In a new terminal, start `node1`, pointed at the bootnode:

```bash
geth --datadir node1 --networkid 12345 --port 30502 --nat none --bootnodes "<enode>" --discovery.dns ""
```

**Quote the `enode` URL.** It contains `?`, and zsh, by default, refuses to run a command with an
unquoted `?` that matches no file.

The log shows the private genesis loading:

```text
INFO [09-13|02:27:26.717] Initialising Ethereum protocol           network=12345 dbversion=<nil>
INFO [09-13|02:27:26.724] Found stored genesis block               config="Chain ID:  12345 (unknown)\nConsensus: Ethash (proof-of-work)\n\n_ Block-based Forks: []_ Time-based Forks: []_ TTD: <nil>"
```

## 5. Start a miner

[Mining](../guides/mining.md) covers the mining flags in more detail.

The miner is a node like `node1`, with mining switched on. Its rewards go to an account whose key
you hold, so create one in `node2`'s data directory:

```bash
geth account new --datadir node2
```

```text
Your new key was generated

Public address of the key:   <address>
Path of the secret key file: node2/keystore/UTC--2026-09-13T08-27-31.511586361Z--<address, lowercase, without 0x>
```

Then, in another terminal, start `node2`:

```bash
geth --datadir node2 --networkid 12345 --port 30503 --nat none --bootnodes "<enode>" --discovery.dns "" \
  --mine --miner.threads=1 --miner.etherbase=<address>
```

- **`--miner.etherbase` is required.** Set to the zero address, `geth` exits with
  `Fatal: Failed to start mining: etherbase missing: etherbase must be explicitly specified`.
- **`--miner.threads=1` is what mines.** `--mine` alone logs `Updated mining threads threads=0`,
  runs no CPU search and seals no block: a node started that way was still at block 0 after 90
  seconds. One thread is enough for a private network.
- **`--miner.gaslimit`** is the gas ceiling blocks move toward;
  [Configuration](../getting-started/run-cli.md#configuration) gives the ceiling `geth` uses when you
  do not set it. **`--miner.gasprice`** is the lowest gas price the miner accepts; locally submitted
  transactions are exempt.

**Before its first block, the miner generates the Ethash dataset (the DAG) for the current epoch**,
and it generates the next epoch's while it mines. Each is about 1 GB, is written to `~/.etchash`
on Linux rather than to the data directory (`--ethash.dagdir` moves it), and occupies every CPU
core while it is generated:

```text
INFO [09-13|02:27:35.247] Generating DAG in progress               epoch=0 epochLength=30000 percentage=0 elapsed=281.437ms
...
INFO [09-13|02:28:08.241] Generated ethash verification dataset    epoch=0 epochLength=30000 elapsed=33.275s
INFO [09-13|02:28:09.325] Generating DAG in progress               epoch=1 epochLength=30000 percentage=0 elapsed=319.257ms
...
INFO [09-13|02:28:13.318] Successfully sealed new block            number=1 sealhash=59dee0..29f76b hash=cde5d9..fcdb3b elapsed=38.890s
```

**This is Ethash, not Etchash.** Ethereum Classic and Mordor mine Etchash, the variant ECIP-1099
introduced, whose epochs are 60,000 blocks (`epochLengthECIP1099` in
`consensus/ethash/algorithm.go`). A chain from this genesis mines the original Ethash, with the
30,000-block epochs its log shows (`epochLengthDefault` in the same file).

## 6. Confirm the network works

Attach to `node1`, which does not mine, and check three things: it has a peer, its chain grows,
and the miner's account is being paid:

```bash
geth attach --exec 'net.peerCount' node1/geth.ipc
geth attach --exec 'eth.blockNumber' node1/geth.ipc
geth attach --exec 'eth.getBalance("<address>")' node1/geth.ipc
```

```text
1
6
14000000000000000000
```

The peer is `node2`, the blocks are the ones it mined, and the balance is in wei. Your numbers
will differ, and grow as the miner runs.

## 7. Stop the network

Press Ctrl-C in each node's terminal, and in the bootnode's. A node logs
`Got interrupt, shutting down...` and exits after `Blockchain stopped`:

```text
INFO [09-13|02:28:37.872] Got interrupt, shutting down...
...
INFO [09-13|02:28:37.910] Blockchain stopped
```

**To restart, start the bootnode and both nodes with the same commands**, skipping
`bootnode --genkey` and `geth account new`. The data directories keep the chain: `node1` restarted
at block 8, where it stopped.

## Proof of authority with Clique

**Clique replaces the miner with signers.** The accounts listed in the genesis block take turns
sealing a block every `period` seconds, and a majority of them can add or remove a signer
([EIP-225](https://eips.ethereum.org/EIPS/eip-225)). This section builds a two-signer network with
the same bootnode and network flags as the steps above. Work in a new directory: it is a separate
chain.

### Create the signer accounts

```bash
geth account new --datadir signer1
geth account new --datadir signer2
```

Each prints its address, as in step 5. Save each password in a file readable only by you, in that
data directory: `signer1/password.txt` and `signer2/password.txt`. A signer reads it when it
starts, to unlock its account.

### Define a Clique genesis

Save this as `genesis.json`, with the two signer addresses in `extraData`:

```json
{
  "config": {
    "chainId": 12345,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "clique": {
      "period": 5,
      "epoch": 30000
    }
  },
  "alloc": {},
  "difficulty": "0x1",
  "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000<first signer address><second signer address>0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "gasLimit": "0x2fefd8"
}
```

- **`extraData` lists the initial signers:** 32 zero bytes (64 hex zeros), then each signer's
  address without `0x`, in ascending order, then 65 zero bytes (130 hex zeros).
- **`clique`, and no `ethash`.** With both objects present, this file runs Clique, while `init` and
  the node log report `Consensus: Ethash (proof-of-work)`, so the check in step 2 would mislead you.
- **`period`** is the number of seconds between blocks. **`epoch`** is the number of blocks between
  checkpoints, which reset pending votes.

Initialize both data directories:

```bash
geth init --datadir signer1 genesis.json
geth init --datadir signer2 genesis.json
```

**Check for `Consensus: Clique (proof-of-authority)`:**

```text
INFO [09-13|02:38:44.008] Wrote genesis block OK                   config="Chain ID:  12345 (unknown)\nConsensus: Clique (proof-of-authority)\n\n_ Block-based Forks: []_ Time-based Forks: []_ TTD: <nil>"
```

### Start the signers

Start a bootnode as in [step 3](#3-start-a-bootnode). Then, in a new terminal, start `signer1`:

```bash
geth --datadir signer1 --networkid 12345 --port 30502 --nat none --bootnodes "<enode>" --discovery.dns "" \
  --mine --miner.etherbase=<signer1 address> --unlock <signer1 address> --password signer1/password.txt
```

In another terminal, start `signer2`:

```bash
geth --datadir signer2 --networkid 12345 --port 30503 --nat none --bootnodes "<enode>" --discovery.dns "" \
  --mine --miner.etherbase=<signer2 address> --unlock <signer2 address> --password signer2/password.txt
```

- **`--unlock` and `--password` are required.** A signer started without them logs
  `Block sealing failed err="authentication needed: password or unlock"` and seals no block.
- **Keep a signer's RPC private.** `geth` refuses to unlock an account while HTTP or WebSocket RPC
  is enabled, unless `--allow-insecure-unlock` is set, and exits with
  `Fatal: Account unlock with HTTP access is forbidden!`

Each signer logs its account unlocking, then the blocks it seals:

```text
INFO [09-13|02:39:14.742] Unlocked account                         address=<signer1 address>
...
INFO [09-13|02:39:14.761] Successfully sealed new block            number=1 sealhash=884c0d..455948 hash=c02b9f..43216e elapsed=18.463ms
```

### Confirm both signers seal

```bash
geth attach --exec 'clique.status()' signer1/geth.ipc
```

```text
{
  inturnPercent: 100,
  numBlocks: 7,
  sealerActivity: {
    <signer1 address>: 3,
    <signer2 address>: 4
  }
}
```

**Both signers are sealing when both counts in `sealerActivity` are above zero.** The counts cover
the most recent blocks, up to 64. `clique.getSigners()` lists the signers the chain currently
accepts, in address order.

### Change the signer set

**Each signer votes from its own node** with `clique.propose`: `true` to add an address, `false` to
remove one. It casts the vote in the blocks it seals. With two signers, one vote is not a majority,
so a few blocks after `signer1` votes, the signer set is unchanged:

```bash
geth attach --exec 'clique.propose("<new signer address>", true)' signer1/geth.ipc
geth attach --exec 'clique.getSigners()' signer1/geth.ipc
```

```text
null
["<signer1 address>", "<signer2 address>"]
```

A few blocks after `signer2` votes too, the new signer is in:

```bash
geth attach --exec 'clique.propose("<new signer address>", true)' signer2/geth.ipc
geth attach --exec 'clique.getSigners()' signer1/geth.ipc
```

```text
null
["<signer1 address>", "<new signer address>", "<signer2 address>"]
```

`clique.proposals` shows the proposals a node votes for. Once a change has taken effect, clear its
proposal on each signer with `clique.discard("<new signer address>")`. Adding an address starts no
node for it: the two running signers kept the chain going, and `clique.status()` counted blocks
for both of them and none for the new signer.

Stop the signers and the bootnode as in [step 7](#7-stop-the-network).

## This tutorial was run before it was published

Every command above was run on September 13, 2026, on one Linux x86_64 machine with an Intel Core
i7-10710U (12 threads) and 31 GiB of RAM, with `geth` and `bootnode` built from source. The output
is quoted from that run, with account addresses and enode IDs replaced by placeholders.
