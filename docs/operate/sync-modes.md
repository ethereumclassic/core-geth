---
title: Sync modes and data retention
description: "Core-Geth sync modes and data retention for an Ethereum Classic node: snap and full sync, the hash and path state schemes, and gcmode, with their disk, memory and sync-time trade-offs."
---

A node's sync mode decides how a new data directory gets the chain's state, and its retention
settings decide how much state and history the node keeps. Started with none of these flags, a node
uses snap sync, the hash state scheme on a new data directory, and `--gcmode full`. Figures for
disk, memory and sync time are on [Hardware requirements](../getting-started/hardware-requirements.md).

## The settings

| Setting | Flag | Default | What it decides |
| --- | --- | --- | --- |
| Sync mode | `--syncmode` | `snap` | How a new data directory gets the chain's state |
| State scheme | `--state.scheme` | `hash` on a new data directory | How the node stores state |
| Garbage collection | `--gcmode` | `full` | Under the hash scheme, whether the state of every block is kept |
| State history | `--history.state` | `90000` blocks | Under the path scheme, how far back state can be restored |
| Transaction index | `--history.transactions` | `2350000` blocks | Which transactions can be looked up by hash |
| Ancient store | `--datadir.ancient` | inside `chaindata` | Where older headers, bodies and receipts are kept |

[Command-line Options](../getting-started/run-cli.md#command-line-options) lists every flag.

## Sync modes

`--syncmode` takes `snap` or `full`:

| Mode | Downloaded from peers | Executed locally |
| --- | --- | --- |
| `snap`, the default | Headers, bodies and receipts, and the state at a recent block, the pivot | The blocks after the pivot |
| `full` | Headers and bodies | Every block from genesis, which rebuilds the state |

To use full sync, give the flag to a new data directory:

```shell
$ geth --classic --datadir <datadir> --syncmode full
```

A snap sync logs its chain and state downloads, then state healing, as
[a first sync of Ethereum Classic](../getting-started/run-classic-node.md#what-a-healthy-first-sync-looks-like)
shows. When it finishes, the node logs `Snap sync complete, auto disabling` and executes every new
block from then on.

**Snap sync needs the state snapshot.** `--snapshot=false` turns the snapshot off only together with
`--syncmode full`. With the default `--syncmode snap`, the node keeps the snapshot and warns:

```
WARN [09-13|12:01:08.011] Snap sync requested, enabling --snapshot
```

**The sync mode applies to a data directory that has not synced.** A data directory that already
holds blocks and the state of its latest block executes every new block, whatever `--syncmode`
says. Started with the default `snap` on such a directory, the node logs:

```
WARN [09-13|12:06:24.924] Switch sync mode from snap sync to full sync reason="snap sync complete"
```

If a snap sync stopped after it had stored blocks, the data directory goes on with snap sync, even
under `--syncmode full`, and logs `Switch sync mode from full sync to snap sync`. To sync in the
other mode, start from a new data directory.

**There is no light sync.** `--syncmode light` stops the node as it starts:

```
Fatal: Failed to register the Ethereum service: can't run eth.Ethereum in light sync mode, light mode has been deprecated
```

The `--light.*` flags are still accepted, and have no effect.

## State scheme

`--state.scheme` is `hash` or `path`. A new data directory gets `hash` unless the flag names `path`,
and logs the choice:

```
INFO [09-13|12:01:02.881] State schema set to default              scheme=hash
```

A data directory keeps the scheme it started with. Started again without the flag, the node uses
the stored scheme:

```
INFO [09-13|12:01:29.446] State scheme set to already existing     scheme=hash
```

On Ethereum Classic, started with the other scheme, it stops:

```
Fatal: Failed to register the Ethereum service: incompatible state scheme, stored: hash, provided: path
```

To use the path scheme, start from a new data directory:

```shell
$ geth --classic --datadir <datadir> --state.scheme path
```

The two schemes keep old state differently: `--gcmode` changes how state is kept only under `hash`,
and `--history.state` applies only under `path`.

## Retention

### State of old blocks

Under the hash scheme, `--gcmode` decides how much state reaches the disk. With `full`, the default,
the node keeps the state of recent blocks in memory and writes state to disk only from time to time,
so the state of the blocks in between is not kept. With `archive`, it writes the state of every
block it executes. Under the path scheme, `--gcmode` does not change how state is kept.

```shell
$ geth --classic --datadir <datadir> --gcmode archive
```

Archive mode also indexes every transaction and records key preimages, the addresses and storage
keys behind the hashed keys in the state. The node logs both as it starts:

```
INFO [09-13|12:01:33.378] Enabling recording of key preimages since archive mode is used
WARN [09-13|12:01:33.378] Disabled transaction unindexing for archive node
```

On an archive node, `--history.transactions` therefore has no effect.
[Archive node](../guides/archive-node.md) covers running one from genesis, and what a node
answers without one.

### State history

Under the path scheme, the node keeps a history of state changes for the most recent
`--history.state` blocks, and uses it to restore the state of an earlier block when the chain is
rewound or reorganized. The default is 90,000 blocks, and `0` keeps the history for the entire
chain. The hash scheme does not use this setting.

### Transaction index

The node indexes the transactions in the most recent `--history.transactions` blocks. The default
is 2,350,000 blocks, and `0` indexes every block. A lookup by transaction hash, with
[`eth_getTransactionByHash`](../JSON-RPC-API/modules/eth.md#eth_gettransactionbyhash) or
[`eth_getTransactionReceipt`](../JSON-RPC-API/modules/eth.md#eth_gettransactionreceipt), finds only
indexed transactions, and returns `null` for an older one.

`--txlookuplimit` is the old name for this flag. It still sets the same value, and logs:

```
WARN [09-13|12:01:02.535] The flag --txlookuplimit is deprecated and will be removed, please use --history.transactions
```

### The ancient store

The `ancient` store keeps older headers, bodies and receipts as flat files, inside `chaindata`
unless `--datadir.ancient` names another directory.
[The `ancient` store on another disk](../getting-started/hardware-requirements.md#the-ancient-store-on-another-disk)
covers using that flag, and [History on another disk](../operate/maintenance.md#history-on-another-disk)
moves an existing store.
