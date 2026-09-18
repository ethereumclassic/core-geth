---
title: Mining
---

Ethereum Classic is secured by proof of work. Its mining algorithm is Etchash (ECIP-1099), a modified
Ethash that doubles the epoch length to 60,000 blocks so the DAG grows more slowly. Etchash applies from
block 11,700,000 on Ethereum Classic and from block 2,520,000 on Mordor. Block rewards follow ECIP-1017:
the reward drops by 20% every 5,000,000 blocks.

core-geth mines in two ways:

- **It hands work to external mining software** over JSON-RPC. GPU and ASIC miners and pool software
  work this way. [Mining pool node](mining-pool-node.md) covers a node that serves a pool.
- **It mines on its own CPU.** CPU hashrate is far below what the network's GPU and ASIC miners produce,
  so this is for [Mordor](mordor-mining.md) and [private networks](../tutorials/private-network.md).

## Before you start

- A node that has finished syncing: [How to tell it is done](../getting-started/run-classic-node.md#how-to-tell-it-is-done).
- An address for the rewards. The node needs no key for it: the address only goes into the blocks it mines.
- For CPU mining, disk and memory for the Etchash DAG of the current epoch, which is several gigabytes.
  External mining software builds its own DAG.

## Mine with external mining software

```sh
geth --classic --mine --miner.etherbase 0xYOUR_ADDRESS \
  --http --http.addr 127.0.0.1 --http.api eth,net,web3
```

`--mine` with the default `--miner.threads 0` builds blocks for external miners without mining on the CPU.
Mining software fetches work with `eth_getWork` and submits solutions with `eth_submitWork`. Keep the HTTP
endpoint on the loopback address or a private network: [RPC exposure](../operate/security.md#rpc-exposure).

To push new work instead of having miners poll for it, list their URLs, separated by commas, with
`--miner.notify`. The node posts each new work package to them. `--miner.notify.full` sends the pending
block's header instead. Set it on the command line: the config file's `NotifyFull` key is not read.

## Mine on the CPU

**This is for Mordor, not for Ethereum Classic mainnet.** Mainnet is mined by ASICs, so CPU threads
find nothing there and the DAG generation and memory they cost buy nothing. On the test network the
difficulty is low enough that a CPU does find blocks, which is what makes it useful for trying a
setup end to end.

```sh
geth --mordor --mine --miner.threads 2 --miner.etherbase 0xYOUR_ADDRESS
```

`--miner.threads` is the number of CPU threads that mine. The DAG is written to `--ethash.dagdir`,
`~/.ethash` by default. The node keeps up to two DAGs on disk (`--ethash.dagsondisk`) and one in memory
(`--ethash.dagsinmem`).

On mainnet, leave `--miner.threads` at its default of 0 and let mining hardware fetch work over the
RPC endpoint, as the section above describes.

## Check that it is mining

```sh
geth --classic attach --exec 'eth.mining' <datadir>/geth.ipc
geth --classic attach --exec 'eth.hashrate' <datadir>/geth.ipc
```

`eth.mining` is `true` once mining is on. `eth.hashrate` is the hashrate the node itself measures;
external mining software reports its own.

## Block template settings

| Flag | Default | What it sets |
|---|---|---|
| `--miner.gaslimit` | 8000000 on Ethereum Classic and Mordor, 30000000 on other networks | The gas ceiling the node moves each mined block toward |
| `--miner.gasprice` | 1000000000 (1 gwei) | The lowest gas price a transaction needs to be included in a mined block |
| `--miner.recommit` | 2s | How often the node rebuilds the block it is mining |
| `--miner.extradata` | the client version | Extra data written into mined blocks, at most 32 bytes |

## Related

- [Mine on Mordor](mordor-mining.md)
- [Mining pool node](mining-pool-node.md)
- [Security and network exposure](../operate/security.md)
- [Monitoring](../operate/monitoring.md)
