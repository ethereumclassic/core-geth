---
title: Hardware requirements
description: "Hardware requirements for a Core-Geth Ethereum Classic node: measured disk, memory and sync-time figures for snap and full sync on ETC and Mordor, and what happens when disk runs low."
---

This page holds the disk, memory and time figures for running a node. Each one was measured on the
machine described below, with the client's default settings unless its row says otherwise, and
each row records the date and the block it was measured at. The chain grows with every block, so a
figure describes what that date needed.

**The reference machine** used for every figure below:

- **CPU:** Intel Core i7-10710U, 6 cores, 12 threads
- **RAM:** 31.1 GiB
- **Disk:** NVMe SSD, formatted btrfs
- **[Default settings](../operate/sync-modes.md):** snap sync, `--gcmode full`, the hash state
  scheme for a new database, and `--cache 1024`

## Ethereum Classic

| Mode | Time to sync | `chaindata`, of which `ancient` | Whole data directory | Peak memory | Measured |
| --- | --- | --- | --- | --- | --- |
| Snap sync (default) | 69 min | 50.98 GiB, of which 27.24 GiB | 51.18 GiB | 2.33 GiB | 2026-09-13; v1.13.0; synced at block 25,337,078; reference machine, `--cache 1024`; state snapshot still generating when disk was read; load average peaked at 12.76 |
| Snap sync (default) | 76 min | not measured | not measured | not measured | 2026-09-12; v1.13.0; synced at block 25,331,662; reference machine, `--cache 1024`, `--http` on; load average not recorded |
| Full sync (`--syncmode full`) | not measured | not measured | not measured | not measured | not measured |
| Archive (`--gcmode archive`) | not measured | not measured | not measured | not measured | not measured |

## Mordor

| Mode | Time to sync | `chaindata`, of which `ancient` | Whole data directory | Peak memory | Measured |
| --- | --- | --- | --- | --- | --- |
| Snap sync (default) | 19 min | 9.61 GiB, of which 7.09 GiB | 9.76 GiB | 1.13 GiB | 2026-09-13; v1.13.0; synced at block 16,972,715; reference machine, `--cache 1024`; state snapshot generated when disk was read; load average peaked at 6.08 |
| Full sync (`--syncmode full`) | not measured | not measured | not measured | not measured | not measured |
| Archive (`--gcmode archive`) | not measured | not measured | not measured | not measured | not measured |

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

**Time to sync** runs from the node's first log line to `Snap sync complete, auto disabling`, with
both timestamps read from the log the node wrote to standard error:

```shell
$ grep -m1 -F 'Starting Core-Geth' geth.log
$ grep -m1 -F 'Snap sync complete, auto disabling' geth.log
```

**Disk** is `du -sb` on each directory, read within two minutes of that line. `-b` reports the
apparent size in bytes, which the tables convert to GiB. Give `du` one directory per command: it
counts each directory once per command, so nested directories passed together leave one of them
out of a total.

```shell
$ du -sb <datadir>/geth/chaindata
$ du -sb <datadir>/geth/chaindata/ancient
$ du -sb <datadir>
```

On Ethereum Classic the node was still generating its state snapshot when the disk was read, and
that snapshot is written into `chaindata` as it grows. On Mordor, generation had finished.

**Peak memory** is the process's peak resident set size as the kernel reports it, read after the
sync finished:

```shell
$ grep VmHWM /proc/<pid>/status
```

**Load average** is the first field of `/proc/loadavg`, sampled once a minute. It includes other
work running on the machine at the same time.

**Stop times** run from `Got interrupt, shutting down...` to `Blockchain stopped` in the log.
Process exit was checked once a second, so those times are to the nearest second.

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
