---
title: Maintenance, backup and upgrades
description: "Maintain a Core-Geth Ethereum Classic node: what to back up, moving and upgrading a node within the v1.13.x line, the offline database commands, state pruning, and recovery after an unclean stop."
---

This page covers keeping a node in service: what to back up, moving a node and upgrading it within
the v1.13.x releases, and what the node does after an unclean stop. It also covers the offline
commands that inspect, prune, move, export and remove the chain database.

**Stop the node before you run any command on this page that opens its database.** Run against a
data directory a node is using, `geth` stops at once:

```
Fatal: Failed to create the protocol stack: datadir already used by another process
```

[Stop it safely](../getting-started/run-classic-node.md#stop-it-safely) says how to stop a node and
how to tell the stop has finished. For a node run as a service, stop the service, run the command as
the user and with the `--datadir` the unit uses, and start the service again afterwards.

## What to back up

| What | Where it is | Why |
| --- | --- | --- |
| Account keyfiles | `keystore/` in the data directory, or the directory `--keystore` names | The node holds no other copy of these keys. Each keyfile is encrypted with its passphrase, so keep the passphrases as well, somewhere else |
| The configuration file and the service unit | where you created them | They hold the network flag, the data directory and every other setting the node runs with |
| The node key | `geth/nodekey` in the data directory | Optional. A node without one writes a new key when it starts, which gives it a new identity on the network ([The node key](security.md#the-node-key)) |
| The chain data | `geth/chaindata/`, and the ancient store if `--datadir.ancient` moved it | Optional. A node without it syncs the chain again |

If you copy the chain data, stop the node first and copy the whole data directory. A copy made that
way opens at the same head block as the original.
[Where your data lives](../getting-started/run-classic-node.md#where-your-data-lives) lists what else
the data directory holds.

## Moving a node to another host

1. Stop the node.
2. Copy its data directory to the new host, with the ancient store if it lives elsewhere, and the
   configuration file and service unit.
3. Start the node on the new host and check it with
   [How to tell it is done](../getting-started/run-classic-node.md#how-to-tell-it-is-done).

**Never run nodes on two hosts from copies of one data directory.** Both would load the same
`geth/nodekey` and use the same enode ID ([The node key](security.md#the-node-key)). If the old host
has to keep running a node, delete its `geth/nodekey` before that node starts again, so it writes a
new key.

## Upgrading within v1.13.x

A node on a v1.12.x release follows [Upgrading](../upgrading.md)
instead.

1. **Read what changed.** Releases and their notes are on the
   [releases page](https://github.com/ethereumclassic/core-geth/releases), and published
   [security advisories](https://github.com/ethereumclassic/core-geth/security/advisories) describe
   the vulnerabilities a release fixes.
   [`SECURITY.md`](https://github.com/ethereumclassic/core-geth/blob/main/SECURITY.md) says how they
   are disclosed.
2. **Download and verify the new release**, following your route in
   [Choose your route](../getting-started/installation.md#choose-your-route) and then
   [Verify where the archive came from](../getting-started/installation.md#verify-where-the-archive-came-from).
3. **Stop the node** and wait for the process to exit
   ([Stop it safely](../getting-started/run-classic-node.md#stop-it-safely)).
4. **Replace the binary**, or the image tag in a container, and start the node.
5. **Check the result.** `geth version` reports the new release, and
   [How to tell it is done](../getting-started/run-classic-node.md#how-to-tell-it-is-done) confirms
   the node is back in sync.

## After an unclean stop

A node whose process was killed or crashed, and so skipped its shutdown sequence, needs nothing done
by hand. Start it as usual. As it starts, it logs:

```
WARN [..] Unclean shutdown detected                booted=<time> age=<age>
```

The node records its ten most recent unclean stops, and every later start logs this line once for
each of them, including starts that follow a clean stop. One line on a start is therefore not proof
that the run before it crashed: compare `booted` with the lines the previous start logged.

If the node was importing blocks when it stopped, the state it held in memory for its newest blocks
is gone. It logs `Head state missing, repairing` and moves its head back to a lower block whose
state it still has. Check it with
[How to tell it is done](../getting-started/run-classic-node.md#how-to-tell-it-is-done).

A node that shuts itself down because its disk is nearly full stops cleanly
([Running out of disk](../getting-started/hardware-requirements.md#running-out-of-disk)).

## The database commands

### See what the database holds

```shell
$ geth --classic --datadir <datadir> db inspect
```

`db inspect` reads the whole database and prints a table with one row for each kind of data it holds.
These rows are an excerpt of its output, with the size and item count in each row replaced by
placeholders:

```
│       DATABASE        │        CATEGORY         │    SIZE    │  ITEMS   │
│ Key-Value store       │ Hash trie nodes         │ <size>     │ <count>  │
│ Key-Value store       │ Path trie account nodes │ <size>     │ <count>  │
│ Key-Value store       │ Account snapshot        │ <size>     │ <count>  │
│ Ancient store (Chain) │ Headers                 │ <size>     │ <count>  │
│ Ancient store (Chain) │ Bodies                  │ <size>     │ <count>  │
```

- `Hash trie nodes` holds the state of a database on the hash scheme, whose `Path trie` rows stay
  empty. On the path scheme the `Path trie` rows fill ([State scheme](sync-modes.md#state-scheme)).
- `Account snapshot` and `Storage snapshot` are the state snapshot, which offline pruning reads.
- The `Ancient store (Chain)` rows are the older chain data kept as flat files
  ([The ancient store](sync-modes.md#the-ancient-store)).

[Hardware requirements](../getting-started/hardware-requirements.md) gives the sizes measured on
synced nodes.

### Read the storage engine's statistics

```shell
$ geth --classic --datadir <datadir> db stats
```

`db stats` prints the storage engine's own statistics: a table of its storage levels, then lines
about its log, flushes, compactions, caches and files. It prints the same statistics twice.

### Compact the database

```shell
$ geth --classic --datadir <datadir> db compact
```

`db compact` rewrites the database's files, which drops data the node has deleted or overwritten. It
logs three lines, the first and last each followed by the statistics `db stats` prints:

```
INFO [..] Stats before compaction
INFO [..] Triggering compaction
INFO [..] Stats after compaction
```

`geth db compact --help` warns that interrupting the command can corrupt the database, so let it
finish.

## Offline state pruning

Under the hash scheme ([State scheme](sync-modes.md#state-scheme)), the database keeps state the node
no longer needs. `snapshot prune-state` deletes it and keeps two states: the genesis state and the
state of the block 127 below the head.

**Check these first:**

- **The database uses the hash scheme.** On the path scheme the command stops:

    ```
    CRIT [..] Offline pruning is not required for path scheme
    ```

- **The state snapshot has 128 blocks on top of it.** The node needs to have imported at least 128
  blocks since its state snapshot finished generating. Otherwise the command stops with the number of
  blocks still missing:

    ```
    ERROR[..] Failed to prune state                    err="snapshot not old enough yet: need 86 more blocks"
    ```

- **The machine has the memory for the bloom filter** the command builds, whose size
  `--bloomfilter.size` sets in megabytes
  ([Command-line Options](../getting-started/run-cli.md#command-line-options) gives the default).
- **The node is not an archive node.** On a node run with `--gcmode archive`, the command is not
  refused and deletes the history the archive exists to keep. A query for the state of an older block
  afterwards fails with `missing trie node`.
  [Archive node](../guides/archive-node.md#what-each-node-answers) covers what it keeps.

**Prune**, with the node stopped:

```shell
$ geth --classic --datadir <datadir> snapshot prune-state
```

The command reads the state it keeps from the snapshot, writes a bloom filter of that state to disk,
then deletes every other state entry. On a large database it logs its progress as
`Pruning state data` lines while it deletes. These lines mark its stages, in order, with placeholders
for the values:

```
INFO [..] Selecting bottom-most difflayer as the pruning target root=<root> height=<block>
INFO [..] State bloom filter committed             name=<datadir>/geth/statebloom.<root>.bf.gz
INFO [..] Pruned state data                        nodes=<count> size=<size> elapsed=<time>
INFO [..] State pruning successful                 pruned=<size> elapsed=<time>
```

`height` is the block whose state it keeps. When the command has logged `State pruning successful`
and exited, start the node. On that start the node logs `Head state missing, repairing`, and its
head is the block pruning kept, 127 below where it stopped.

## History on another disk

`--datadir.ancient` keeps the ancient store, the older headers, bodies and receipts, in a directory
you choose, such as one on another disk. To move the store of a synced node:

1. Stop the node.
2. Move the store:

    ```shell
    $ mv <datadir>/geth/chaindata/ancient <ancient dir>
    ```

3. Start the node with `--datadir.ancient`, and add the flag to its service unit:

    ```shell
    $ geth --classic --datadir <datadir> --datadir.ancient <ancient dir>
    ```

    As it starts, it logs where it opened the store:

    ```
    INFO [..] Opened ancient database                  database=<ancient dir>/chain readonly=false
    ```

4. Pass the same flag to every command on this page that opens the database, `removedb` among
   them.

Started without the flag, the node refuses, as
[The `ancient` store on another disk](../getting-started/hardware-requirements.md#the-ancient-store-on-another-disk)
shows. That start leaves a new, empty store at the old path and does not touch the moved one. Delete
the empty store before you move the real one back.

## Export and import blocks

`export` writes blocks to a file:

```shell
$ geth --classic --datadir <datadir> export <file> <first block> <last block>
```

```
INFO [..] Exporting blockchain                     file=<file>
INFO [..] Exporting batch of blocks                count=1000
INFO [..] Exported blockchain to                   file=<file>
```

Given a range, `export` appends to a file that already exists, so a chain can be exported in pieces.
For a file name ending in `.gz` the output is compressed, and `import` reads it the same way.

`import` reads such a file into a data directory of the same network and executes each block, as a
full sync does:

```shell
$ geth --classic --datadir <other datadir> import <file>
```

While it works it logs `Imported new chain segment` lines. When the file is done it prints its memory
use and compacts the database.

- **A new data directory needs the chain from its start.** A file that does not begin at block 0 or
  block 1 stops the import, and so do blocks from another network:

    ```
    ERROR[..] Import error                             err="invalid block 5000: unknown ancestor"
    ```

- **Blocks the data directory already has are skipped:**

    ```
    INFO [..] Skipping batch as all blocks present     batch=0 first=<hash> last=<hash>
    ```

## Resync from scratch

`removedb` deletes the chain database, so the node syncs the chain again from its peers when it next
starts:

```shell
$ geth --classic --datadir <datadir> removedb
```

It asks about two parts in turn, and deletes a part only when you answer `y`. For a resync, answer `y`
to both. After a few startup lines, the prompts and their results look like this:

```
Location(s) of 'state data': 
	- <datadir>/geth/chaindata
	- <datadir>/geth/chaindata/ancient/state

Remove 'state data'? [y/n] y
INFO [..] Folder is not existent                   path=<datadir>/geth/chaindata/ancient/state
INFO [..] Database successfully deleted            kind="state data" paths=[<datadir>/geth/chaindata] elapsed=<time>
Location(s) of 'ancient chain': 
	- <datadir>/geth/chaindata/ancient/chain

Remove 'ancient chain'? [y/n] y
INFO [..] Database successfully deleted            kind="ancient chain" paths=[<datadir>/geth/chaindata/ancient/chain] elapsed=<time>
```

- **`state data` is the database outside the ancient store**, and `ancient chain` is the ancient
  store.
- **Do not keep the ancient chain alone.** A data directory whose state data was removed and whose
  ancient chain was kept does not start: the node stops logging after `Empty database, resetting
  chain` and never finishes starting.
- **`keystore/` and `geth/nodekey` stay**, so the node keeps its accounts and its identity.
- **`--remove.state` and `--remove.chain` answer the questions without prompting.**
- **With `--datadir.ancient`, pass it to `removedb` too.** Without it, `removedb` looks for the
  ancient store inside `chaindata`.
