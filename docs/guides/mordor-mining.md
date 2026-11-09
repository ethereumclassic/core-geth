---
title: Mine on Mordor
description: "Mine on Mordor, Ethereum Classic's proof-of-work test network, with Core-Geth: rehearse a mining or pool setup and get test ETC without a faucet."
---

Mordor is Ethereum Classic's proof-of-work test network. Mining there tests a mining or pool setup without
real funds, and it is how a developer gets test coins without waiting on a faucet. Mordor uses Etchash from
block 2,520,000.

## Test coins without a faucet

**You do not need a faucet to fund development on Mordor.** Point a node at the test network with a CPU
thread or two and leave it running: twenty minutes is usually more coin than a project needs for testing.

The difficulty is low enough that an ordinary machine is a real share of the network. On September 18, 2026, the whole
network ran at around 0.67 MH/s, so a machine contributing 0.2 MH/s earned about 70 METC an hour. These figures move. To
read your node's hashrate and the current difficulty:

```sh
geth --mordor attach --exec 'eth.hashrate' <datadir>/geth.ipc
geth --mordor attach --exec 'eth.getBlock("latest").difficulty' <datadir>/geth.ipc
```

Mordor coins have no value and are not exchangeable for ETC.

## Before you start

- A synced Mordor node: [Run a Mordor node](../getting-started/run-mordor-node.md).
- An address for the rewards. The node needs no key for it.
- For CPU mining, disk and memory for the Etchash DAG, which is several gigabytes. The node generates it
  before it starts mining.

## Mine on the CPU

```sh
geth --mordor --mine --miner.threads 1 --miner.etherbase 0xYOUR_ADDRESS
```

The node writes the DAG to `--ethash.dagdir`, `~/.ethash` by default, then starts mining. Generating the
DAG takes a few minutes and several gigabytes before the first hash is tried. For more hashrate than a CPU gives, point GPU mining software at
the node as [Mining](mining.md#mine-with-external-mining-software) describes, with `--mordor` in place of
`--classic`.

## Check the rewards

```sh
geth --mordor attach --exec 'eth.getBalance("0xYOUR_ADDRESS")' <datadir>/geth.ipc
```

The balance is in wei, and it grows as the node imports blocks mined to the address. For a readable figure:

```sh
geth --mordor attach --exec 'web3.fromWei(eth.getBalance("0xYOUR_ADDRESS"), "ether")' <datadir>/geth.ipc
```

## Stop mining

Stop the node, and start it again without `--mine`.
[Stop it safely](../getting-started/run-classic-node.md#stop-it-safely) applies to Mordor too.
