---
title: Hardware requirements
description: "Hardware requirements for a Core-Geth Ethereum Classic node: measured disk, memory and sync-time figures for snap and full sync on ETC and Mordor, and what happens when disk runs low."
---

These figures were measured on the machine below, with default settings unless a row says otherwise. The chain grows
with every block, so a figure describes what its date needed. Full sync and archive mode were not measured.

**The reference machine** used for every figure below:

- **CPU:** Intel Core i7-10710U, 6 cores, 12 threads
- **RAM:** 31.1 GiB
- **Disk:** NVMe SSD, formatted btrfs
- **[Default settings](../operate/sync-modes.md):** snap sync, `--gcmode full`, the hash state
  scheme for a new database, and `--cache 1024`

## Ethereum Classic

| Mode | Time to sync | `chaindata`, of which `ancient` | Whole data directory | Peak memory | Measured |
| --- | --- | --- | --- | --- | --- |
| Snap sync (default) | 69 min | 50.98 GiB, of which 27.24 GiB | 51.18 GiB | 2.33 GiB | 2026-09-13; v1.13.0; synced at block 25,337,078; reference machine, `--cache 1024`; state snapshot still generating when disk was read |
| Snap sync (default) | 76 min | not measured | not measured | not measured | 2026-09-12; v1.13.0; synced at block 25,331,662; reference machine, `--cache 1024`, `--http` on |

## Mordor

| Mode | Time to sync | `chaindata`, of which `ancient` | Whole data directory | Peak memory | Measured |
| --- | --- | --- | --- | --- | --- |
| Snap sync (default) | 19 min | 9.61 GiB, of which 7.09 GiB | 9.76 GiB | 1.13 GiB | 2026-09-13; v1.13.0; synced at block 16,972,715; reference machine, `--cache 1024`; state snapshot generated when disk was read |

## Stopping

How long a clean stop took, from the interrupt to `Blockchain stopped` in the log, and to the
process exiting.

| Network | Node when stopped | To `Blockchain stopped` | To process exit | Measured |
| --- | --- | --- | --- | --- |
| Ethereum Classic | synced | 42 ms | 6 s | 2026-09-13; v1.13.0; block 25,337,086; reference machine |
| Ethereum Classic | synced | 54 ms | not measured | 2026-09-12; v1.13.0; block 25,331,716; reference machine, `--http` on |
| Mordor | synced | 82 ms | 4 s | 2026-09-13; v1.13.0; block 16,972,720; reference machine |
| Mordor | synced | 35 ms | not measured | 2026-09-12; v1.13.0; block 16,968,782; reference machine, `--http` on |
| Mordor | first sync, downloading the chain | 1 min 21 s | not measured | 2026-09-12; v1.13.0; headers at block 2,556,480; reference machine, `--http` on |
| Mordor | first sync, downloading the chain | 11 s | not measured | 2026-09-12; v1.13.0; headers at block 7,250,672; reference machine, `--http` on |

## How each figure was measured

- **Time to sync** runs from the node's first log line to `Snap sync complete, auto disabling`.
- **Disk** is `du -sb` on each directory, read within two minutes of that line. The state snapshot, still generating on
  Ethereum Classic, adds to `chaindata` as it grows.
- **Peak memory** is `VmHWM` in `/proc/<pid>/status`, read after the sync finished.
- **Stop times** run from `Got interrupt, shutting down...` to `Blockchain stopped`. Process exit is to the nearest
  second.

## Memory and `--cache`

`--cache` is the memory, in MiB, that the node gives its internal caches. The peak memory
figures above use the default, 1024, which the node divides among `--cache.database`,
`--cache.trie`, `--cache.gc` and `--cache.snapshot`. A default start logs the first three:

```
INFO [09-13|08:13:49.777] Allocated trie memory caches             clean=154.00MiB dirty=256.00MiB
INFO [09-13|08:13:49.777] Allocated cache and file handles         database=<datadir>/geth/chaindata cache=512.00MiB handles=524,288
```

The node lowers a `--cache` larger than one third of the machine's memory, and logs the change.
Asked for 20000 on the reference machine, it used 10601:

```
WARN [09-13|09:55:06.321] Sanitizing cache to Go's GC limits       provided=20000 updated=10601
```

On a 32-bit platform it counts no more than 2 GiB of memory when it sets that limit.

## The `ancient` store on another disk

`--datadir.ancient` puts the `ancient` store, the older headers, bodies and receipts that
otherwise live in `chaindata/ancient`, in a directory you choose. The tables above show how much
of `chaindata` it holds. [History on another disk](../operate/maintenance.md#history-on-another-disk)
moves the store of a synced node.

**Pass the same `--datadir.ancient` on every start.** Started without it, the node looks for the
store inside `chaindata` instead. A node whose first blocks have already moved into its `ancient`
store then refuses to start, with
`ancient chain segments already extracted, please set --datadir.ancient to the correct path`. A
node that has not started syncing creates a new, empty store there, and logs no error.

## Running out of disk

The node checks the free space in `<datadir>/geth` every 30 seconds. Below twice its threshold, it
warns:

```
WARN [09-13|09:55:10.945] Disk space is running low. Geth will shutdown if disk space runs below critical level. available=55.12GiB critical_level=35.80GiB path=<datadir>/geth
```

Below the threshold, it shuts itself down cleanly, as it would on an interrupt:

```
ERROR[09-13|09:55:10.091] Low disk space. Gracefully shutting down Geth to prevent database corruption. available=55.12GiB path=<datadir>/geth
INFO [09-13|09:55:10.091] Got interrupt, shutting down...
```

The threshold is twice the share of `--cache` given to `--cache.gc`, which is 512 MiB at the
defaults. `--datadir.minfreedisk` sets it in MiB, and `0` turns the check off.

## Small machines

Nothing on this page was measured on a smaller machine. Two limits matter on one, and both have
their own section:
[key derivation on a small machine](run-cli.md#key-derivation-on-a-small-machine), and
[which ARM archive to install](installation.md#linux-arm).
