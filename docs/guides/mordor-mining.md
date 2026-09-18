---
title: Mine on Mordor
---

Mordor is Ethereum Classic's proof-of-work test network. Mining there tests a mining or pool setup without
real funds, and it is how a developer gets test coins without waiting on a faucet. Mordor uses Etchash from
block 2,520,000.

## Test coins without a faucet

**You do not need a faucet to fund development on Mordor.** Point a node at the test network with a CPU
thread or two and leave it running: twenty minutes is usually more coin than a project needs for testing.
The [faucets](../getting-started/run-mordor-node.md#test-coins) are worth trying, and they are
frequently dry or unattended, which is why this page exists.

It works because the test network's difficulty is low enough that an ordinary machine is a real share of
it. Measured on 18 September 2026 at block 17,006,119: the previous twenty blocks averaged 10.7 seconds
apart at a difficulty near 7.2 million, which puts the whole network around 0.67 MH/s. The block reward in
the current era is 0.8388608 METC, since ECIP-1017 reduces the reward by a fifth every 2,000,000 blocks on
Mordor and era 9 begins at block 16,000,001.

So a machine contributing 0.2 MH/s is roughly a quarter of the network, which at that spacing is on the
order of 80 blocks an hour, or about 70 METC. Your own rate decides the rest, and as hashrate joins the
network the difficulty retargets upward and your share settles lower.

**These figures move, so measure rather than trust them.** `eth.hashrate` in the console reports what your
node is contributing, and the current network numbers come from any Mordor endpoint:

```sh
geth --mordor attach --exec 'eth.hashrate' <datadir>/geth.ipc
geth --mordor attach --exec 'eth.getBlock("latest").difficulty' <datadir>/geth.ipc
```

Mordor coins have no value and are not exchangeable for ETC. Nothing here applies to mainnet, where the
difficulty is set by mining hardware and a CPU finds nothing.

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
DAG takes a few minutes and several gigabytes before the first hash is tried, so the first blocks arrive
later than the arithmetic above suggests. For more hashrate than a CPU gives, point GPU mining software at
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
