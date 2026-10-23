---
title: Archive node
description: "Run a Core-Geth archive node on Ethereum Classic: the full historical state it keeps, what it can answer at any block, and the disk and sync it requires."
---

An archive node keeps the state of every block it executes, so it can answer a question about the
chain's state at any block in its history, such as an account's balance at that block or a trace of
a transaction mined long ago. A node on the default settings keeps the state of its recent blocks
only.

## Who needs one

Most nodes do not. A node on the default settings answers these for every block:

- blocks and their transactions, with [`eth_getBlockByNumber`](../JSON-RPC-API/modules/eth.md#eth_getblockbynumber);
- event logs, with [`eth_getLogs`](../JSON-RPC-API/modules/eth.md#eth_getlogs);
- receipts, for the transactions in its [transaction index](../operate/sync-modes.md#transaction-index);
- the current state: balances, storage, and calls at `latest`.

An archive node is needed to answer these about older blocks:

- **state at a given block:** [`eth_getBalance`](../JSON-RPC-API/modules/eth.md#eth_getbalance),
  [`eth_getStorageAt`](../JSON-RPC-API/modules/eth.md#eth_getstorageat),
  [`eth_getProof`](../JSON-RPC-API/modules/eth.md#eth_getproof) and
  [`eth_call`](../JSON-RPC-API/modules/eth.md#eth_call) with an older block number;
- **traces:** [`debug_traceTransaction`](../JSON-RPC-API/modules/debug.md#debug_tracetransaction),
  [`debug_traceBlockByNumber`](../JSON-RPC-API/modules/debug.md#debug_traceblockbynumber),
  [`trace_block`](../JSON-RPC-API/modules/trace.md#trace_block) and
  [`trace_transaction`](../JSON-RPC-API/modules/trace.md#trace_transaction), described in
  [the `trace` module](../JSON-RPC-API/trace-module-overview.md).

## What each node answers

The default node below runs with no sync or retention flags. The archive node runs as
[Run an archive node](#run-an-archive-node) shows.

| Calls | Default node | Archive node |
| --- | --- | --- |
| `eth_getBlockByNumber`, `eth_getLogs` | Every block | Every block |
| `eth_getTransactionReceipt` | Transactions in the transaction index | Every transaction |
| `eth_getBalance`, `eth_getStorageAt`, `eth_getProof`, `eth_call` at a block | At most the 128 most recent blocks | Every block |
| `debug_traceTransaction`, `debug_traceBlockByNumber`, `trace_block`, `trace_transaction` | The most recent blocks | Every block |

**State.** A default node keeps the state of its 128 most recent blocks while it runs. Just after a
snap sync it has fewer, since its state starts at the sync's pivot
([With snap sync](#with-snap-sync)). Asked about an older block, it answers:

```
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"missing trie node <hash> (path ) state <hash> is not available, not found"}}
```

A restart loses most of those states. The node writes the state of a few recent blocks to disk as
it stops, and has no state for the others when it starts again.

**Traces.** A trace needs the state of the block before the one it traces. A node that lacks that
state goes back up to 128 blocks to the nearest state it has and executes the blocks from there, and
fails when no state is that close:

```
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"required historical state unavailable (reexec=128)"}}
```

**Receipts.** [Transaction index](../operate/sync-modes.md#transaction-index) says which transactions
a default node finds by hash. An archive node indexes every transaction.

## Run an archive node

Start a new data directory with full sync and `--gcmode archive`:

```shell
$ geth --classic --datadir <datadir> --syncmode full --gcmode archive
```

Full sync executes every block from genesis, and in archive mode the node writes the state of each
block to disk. Disk, memory and sync time are on
[Hardware requirements](../getting-started/hardware-requirements.md).

Archive mode also turns on the full transaction index and the recording of key preimages, as
[State of old blocks](../operate/sync-modes.md#state-of-old-blocks) shows. It gives trie pruning no
memory: the `Allocated trie memory caches` line the node logs as it starts reports no dirty cache,
and the share of `--cache` that `--cache.gc` would take goes to the trie and snapshot caches.

### With snap sync

Snap sync, the default, downloads the state of one recent block, the pivot, and executes blocks only
from there. A node started with `--gcmode archive` and no `--syncmode` therefore has no state for any
block before its pivot, and keeps the state of every block from the pivot on. Each start after the
sync logs the pivot:

```
INFO [..] Loaded last snap-sync pivot marker       number=<block>
```

### With the path scheme

Under the path scheme, `--gcmode archive` keeps no older state. A node on `--state.scheme path`
answers state queries for its most recent blocks only, and refuses to trace older ones, even when
it keeps [state history](../operate/sync-modes.md#state-history) for every block:

```
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"historical state not available in path scheme yet"}}
```

A new data directory gets the hash scheme unless `--state.scheme path` names the other
([State scheme](../operate/sync-modes.md#state-scheme)).

## Change an existing node

- **To archive:** restarted with `--gcmode archive`, a node keeps the state of every block it
  executes from then on. The states it did not already have on disk do not come back, so for blocks
  before the restart it answers only where the state was on disk already.
- **Back from archive:** restarted without `--gcmode archive`, a node keeps the states it wrote as
  an archive. For blocks it executes after the restart, it keeps the 128 most recent only.
- **For the whole history:** sync a new data directory with the command above, or remove the chain
  database first ([Resync from scratch](../operate/maintenance.md#resync-from-scratch)).
- **Pruning:** `snapshot prune-state` on an archive node deletes the history it keeps
  ([Offline state pruning](../operate/maintenance.md#offline-state-pruning)).

## Serve it to others

[Public RPC endpoint](public-rpc-endpoint.md) covers serving JSON-RPC to other people. The calls that
need an archive node are in the `eth`, `trace` and `debug` namespaces. `debug` also holds methods
that change the node, such as `debug_setHead`, so offer its tracing methods only through a gateway
that passes them and nothing else from `debug`.
