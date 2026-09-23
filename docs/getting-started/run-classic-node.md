---
title: Run an Ethereum Classic node
---

This page follows an Ethereum Classic node in depth, from its first start to a clean stop. Its
commands are written for Linux and macOS. [Running a node](run-a-node.md) has a guide for each
platform, including Windows and Docker, and the flags for each kind of node.
[Run a Mordor node](run-mordor-node.md) covers what differs on the test network.

## What this runs, and what it does not need

`geth` is the whole node. Ethereum Classic is proof of work, so there is no consensus client or
beacon node to run beside it, and no JWT secret to configure.

The Engine API, the interface Ethereum's consensus clients use, starts only on a chain configured for
the merge, and neither Ethereum Classic nor Mordor is. The node opens no port 8551 and writes no JWT
secret. Given `--authrpc.addr`, `--authrpc.port`, `--authrpc.vhosts` or `--authrpc.jwtsecret`, it logs
that it ignored them.

## Before you start

- **`geth version` reports a 1.13 release.** See [Installation](installation.md).
- **The machine has the disk and memory the network needs.** See
  [Hardware requirements](hardware-requirements.md).
- **Other nodes can reach this one.** Allow the peer-to-peer port through your firewall, over TCP
  and UDP. [Ports and listeners](../operate/security.md#ports-and-listeners) gives the port, and
  [A host firewall](../operate/security.md#a-host-firewall) a ruleset.

## Start it

```shell
$ geth --classic --datadir <datadir>
```

`<datadir>` is a directory on the disk that will hold the chain, and the node creates it on first
start. The node runs in the foreground and logs to the terminal until you
[stop it](#stop-it-safely).

**Every command on this page passes `--classic`.** On a new or Ethereum Classic data directory,
`geth` with no network flag also runs Ethereum Classic ([what that changes for a node upgraded
from v1.12.x](../tutorials/v1.13.0-migration.md#if-you-start-the-node-without-a-network-flag)). The
flag is still worth passing. It records the network in your service file, and it makes the node
refuse a data directory that holds another network instead of running it. Pointed at a Mordor
data directory, `geth --classic` stops at once:

```
Fatal: Failed to register the Ethereum service: database contains incompatible genesis (have a68ebde7932eccb177d38d55dcc6461a019dd795a681e59b5a3e4f3a7259a3f1, new d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3)
```

## Where your data lives

| Path, inside `<datadir>` | What it holds |
| --- | --- |
| `geth/chaindata/` | The chain database |
| `geth/chaindata/ancient/` | Older headers, bodies and receipts, kept as flat files |
| `geth/nodekey` | The node's private key, which sets its identity on the network; written on first start |
| `geth/nodes/` | The database of other nodes this node has found |
| `geth/etchash/` | Etchash verification caches, written during sync |
| `geth/blobpool/`, `geth/transactions.rlp` | Transaction pool state |
| `keystore/` | Account keyfiles; empty until you create an account |
| `tmp/` | Temporary files |
| `geth.ipc` | The socket `geth attach` connects to; it exists while the node runs |

Delete nothing inside it by hand. Removing `geth/nodekey`, for one, gives the node a new identity
on the network.

Without `--datadir`, the node uses a `classic` directory under a default location:

| Platform | Default data directory |
| --- | --- |
| Linux | `~/.ethereum/classic` |
| macOS | `~/Library/Ethereum/classic` |
| Windows | `%LOCALAPPDATA%\Ethereum\classic`, unless `AppData\Roaming\Ethereum` in your home directory already exists and is not empty, in which case its `classic` directory |

## What a healthy first sync looks like

The lines below come from a first sync of Ethereum Classic with the default settings. They are in
order, with the data directory shown as `<datadir>`; the timestamps show where the log skips
ahead.

**Starting:**

```
INFO [09-13|08:13:49.579] Starting Core-Geth on Ethereum Classic...
INFO [09-13|08:13:49.854] State schema set to default              scheme=hash
INFO [09-13|08:13:50.497] Initialising Ethereum protocol           network=1 dbversion=<nil>
INFO [09-13|08:13:50.502] Writing custom genesis block
WARN [09-13|08:13:50.986] Failed to load snapshot                  err="missing or corrupted snapshot"
INFO [09-13|08:13:50.988] Rebuilding state snapshot
INFO [09-13|08:13:51.048] IPC endpoint opened                      url=<datadir>/geth.ipc
```

- `Starting Core-Geth on Ethereum Classic...` names the network.
- `scheme=hash` is the state scheme a new database gets
  ([state scheme](../operate/sync-modes.md#state-scheme)).
- `network=1` is Ethereum Classic's network ID, and `dbversion=<nil>` means the database is new.
- `Writing custom genesis block` appears once, on a new data directory.
- `Failed to load snapshot` is expected on a new data directory: it has no state snapshot yet, so
  the node starts building one.
- `IPC endpoint opened` gives the socket that `geth attach` connects to.

**Syncing:**

```
INFO [09-13|08:14:01.050] Block synchronisation started
INFO [09-13|08:14:01.515] Looking for peers                        peercount=2 tried=7 static=0
INFO [09-13|08:14:03.224] Syncing: chain download in progress      synced=0.00% chain=884.00B headers=384@812.00B bodies=2@39.00B receipts=2@33.00B eta=165h21m12.866s
INFO [09-13|08:14:10.350] Syncing: state download in progress      synced=0.71% state=211.50MiB accounts=610,064@152.50MiB slots=302,393@57.55MiB codes=592@1.44MiB eta=18m35.254s
WARN [09-13|08:27:34.493] Pivot seemingly stale, moving            old=25,336,690 new=25,336,754
INFO [09-13|09:03:28.930] Syncing: state healing in progress       accounts=9@759.00B            slots=0@0.00B              codes=0@0.00B         nodes=573@281.61KiB pending=824
```

- Once it has peers, the node downloads two things side by side and logs each one: the chain
  (headers, bodies and receipts) and the state (accounts, storage slots and contract code).
  `synced=` is each one's progress. `eta` is the node's own estimate, and it swings widely at
  first.
- `state healing` lines follow the state download, while the node fetches state that changed as
  it downloaded. `pending` counts what is left, and healing lines can reappear briefly later in
  the sync.
- The `ERROR` and `WARN` lines are expected. The [table below](#warnings-a-first-sync-logs) says
  why each appears.

**Synced:**

```
INFO [09-13|09:23:07.990] Snap sync complete, auto disabling
INFO [09-13|09:23:07.990] Snapshot generation in progress, snap/1 serving unavailable until complete diskroot=7b0087..46160e headroot=cc7503..8ca53c
INFO [09-13|09:23:07.990] Enabled artificial finality features     reason=synced peers=16
INFO [09-13|09:23:20.834] Imported new chain segment               number=25,337,079 hash=ca8921..521cd4 blocks=1         txs=0         mgas=0.000 elapsed=166.553ms mgasps=0.000 snapdiffs=394.00B triedirty=485.77KiB af=true
INFO [09-13|09:23:35.590] Indexed transactions                     blocks=2,350,000 txs=5,065,894 tail=22,987,076 elapsed=27.785s
```

- `Snap sync complete, auto disabling` marks the end of the first sync.
- `Enabled artificial finality features reason=synced` is MESS switching on
  ([MESS on this node](#mess-on-this-node)).
- From here the node imports new blocks as they arrive, usually one block per
  `Imported new chain segment` line. `af=true` on those lines means MESS is in force.
- The node goes on building its state snapshot in the background, and serves no snapshot data to
  other nodes until that finishes.
- `Indexed transactions` means the transaction index is built.

### Warnings a first sync logs

A first sync logs many of these. None of them stops it.

| Line | Level | Why it appears |
| --- | --- | --- |
| `Failed to load snapshot` | WARN | A new data directory has no state snapshot, so the node builds one |
| `Pivot seemingly stale, moving` | WARN | The chain moved on while the state downloaded, so the node moved its sync target forward |
| `Synchronisation failed, dropping peer` | WARN | Syncing from one peer failed, for example with `err=timeout`, and the node dropped that peer |

## How to tell it is done

A new node is synced when all four of these hold.

1. **The log shows the sync finishing:** `Snap sync complete, auto disabling`, then
   `Enabled artificial finality features reason=synced`.
2. **`eth.syncing` returns `false`:**

    ```shell
    $ geth --classic attach --exec 'eth.syncing' <datadir>/geth.ipc
    ```

    Until then it returns an object, even on a node that has not started syncing, and it stays an
    object until `Indexed transactions` is logged. In that object `currentBlock` climbs toward
    `highestBlock`, but `highestBlock` is the head the sync set out for, and during a long sync
    the network moves past it. It is not the network's head.

3. **The node has peers.** This returns a number above zero:

    ```shell
    $ geth --classic attach --exec 'net.peerCount' <datadir>/geth.ipc
    ```

4. **Its head matches the network's.** Compare this with a block explorer or another node you
   operate; the two should match, give or take the blocks mined between the readings:

    ```shell
    $ geth --classic attach --exec 'eth.blockNumber' <datadir>/geth.ipc
    ```

After an upgrade from v1.12.x, use the checks in
[5. Start, and verify](../tutorials/migration/linux.md#5-start-and-verify) instead.

## Run it as a service

Each platform guide keeps the node running, and gives it time to stop cleanly:

| Platform | How | Where |
| --- | --- | --- |
| Linux | A systemd service | [Linux users guide, step 6](run/linux.md#6-keep-it-running-a-systemd-service) |
| macOS | A launchd agent | [Mac users guide, step 6](run/macos.md#6-keep-it-running-a-launchd-agent) |
| Windows | A window that opens when you sign in | [Windows users guide, step 6](run/windows.md#6-keep-it-running-start-it-when-you-sign-in) |
| Docker | A restart policy and a stop timeout | [Docker users guide, step 3](run/docker.md#3-start-the-node) |

## Stop it safely

Send the node one interrupt, then wait. In the foreground, press Ctrl-C once.
`sudo systemctl stop core-geth` and `docker stop` send SIGTERM, which the node handles the same
way. A synced node logs:

```
INFO [09-13|09:25:19.159] Got interrupt, shutting down...
INFO [09-13|09:25:19.159] IPC endpoint closed                      url=<datadir>/geth.ipc
INFO [09-13|09:25:19.176] Writing cached state to disk             block=25,337,086 hash=ef6268..d9aa3f root=5dea04..dcdb69
INFO [09-13|09:25:19.189] Writing cached state to disk             block=25,337,085 hash=f0d143..8cb13c root=7d7c0e..00c685
INFO [09-13|09:25:19.194] Writing cached state to disk             block=25,336,959 hash=ef6426..3f1c55 root=019e6a..05edeb
INFO [09-13|09:25:19.194] Writing snapshot state to disk           root=cc7503..8ca53c
INFO [09-13|09:25:19.201] Blockchain stopped
```

`Blockchain stopped` means the node has written its cached state to disk. The process exits after
that line; wait for it to exit before you start anything else on the data directory.

**Pressing Ctrl-C again does not speed it up.** Each further interrupt only logs a warning:

```
WARN [09-13|09:54:06.042] Already shutting down, interrupt more to panic. times=9
```

The tenth further interrupt makes the node panic instead of finishing, which is an unclean stop.

**A stop during the first sync takes longer**, because the node first rolls back the part of the
chain it was downloading ([stop times](hardware-requirements.md#stopping)). It also logs an error
that you can ignore:

```
ERROR[09-12|11:08:36.601] Failed to journal state snapshot         err="snapshot [0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421] missing"
```

Started again, the node resumes the sync.

**A killed node** logs `Unclean shutdown detected` when it next starts
([After an unclean stop](../operate/maintenance.md#after-an-unclean-stop)).

A node started with the `console` subcommand stops differently; see
[Command line](run-cli.md#a-node-on-an-ethereum-classic-network).

## Use your node

`geth attach` without `--exec` opens an interactive JavaScript console on the running node:

```shell
$ geth --classic attach <datadir>/geth.ipc
```

For programs and wallets, turn on the HTTP endpoint:

```shell
$ geth --classic --datadir <datadir> --http
```

It listens on `127.0.0.1:8545` only, and serves the `eth`, `net` and `web3` namespaces; `admin`
calls are refused. Ask it for the chain ID:

```shell
$ curl -s -X POST -H 'Content-Type: application/json' \
    --data '{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}' \
    http://127.0.0.1:8545
{"jsonrpc":"2.0","id":1,"result":"0x3d"}
```

`0x3d` is Ethereum Classic's chain ID ([supported networks](../index.md#supported-networks)). A
wallet uses the same URL and chain ID.

Keep the endpoint on `127.0.0.1`. Before you open it wider, read
[RPC exposure](../operate/security.md#rpc-exposure).

## MESS on this node

This client implements MESS, Modified Exponential Subjective Scoring
([ECIP-1100](https://ecips.ethereumclassic.org/ECIPs/ecip-1100)), a chain-selection defense
against deep reorganizations. **Whether your node applies it depends on the release and on the
flags you pass**, and `admin_ecbp1100Status` answers that without changing anything.
[MESS](../operate/mess.md) has the timeline per release and the trade either way:

```shell
$ geth --classic attach --exec 'admin.ecbp1100Status()' <datadir>/geth.ipc
```

`nodeSwitch` is the switch that log line turns on. `enabled` is `true` when that switch is on and
the head has reached `activatedAtBlock`, the block MESS applies from.
[MESS](../operate/mess.md) covers turning it off and back on, and when the node switches it off by itself.

## Next steps

- [Running a node](run-a-node.md), for your platform's commands and the flags for each kind of node.
- [Run a Mordor node](run-mordor-node.md), to rehearse on the test network.
- [Hardware requirements](hardware-requirements.md), for disk, memory and sync times.
- [Sync modes and data retention](../operate/sync-modes.md), for how a node syncs and what it
  keeps.
- [Maintenance, backup and upgrades](../operate/maintenance.md), to back up, upgrade and prune
  the node.
- [Monitoring](../operate/monitoring.md), for metrics, logs and health checks.
- [Troubleshooting](../operate/troubleshooting.md), for what a problem looks like in the log and
  how to fix it.
- [Choose your role](../guides/index.md), for the guides that match what you run the node for.
- [Security and network exposure](../operate/security.md), for ports, the firewall and RPC
  exposure.
- [Configuration](run-cli.md#configuration), to keep settings in a file, and
  [Command-line Options](run-cli.md#command-line-options), for every flag.
- [JSON-RPC API](../JSON-RPC-API/index.md), and the [`eth` namespace](../JSON-RPC-API/modules/eth.md)
  for `eth_syncing` and the rest of it.
- [Migrating to v1.13.0](../tutorials/v1.13.0-migration.md), if this machine ran a v1.12.x node.
