---
title: Run a Mordor node
---

Mordor is Ethereum Classic's proof-of-work test network. Its chain is a fraction of mainnet's
size, which makes it the place to try the client, test contracts, and rehearse a mainnet
deployment before it matters.

Everything on [Run an Ethereum Classic node](run-classic-node.md) applies here: stopping safely,
reading the log, and telling when the sync is done. This page covers what differs.
[Running a node](run-a-node.md) gives the commands for each platform, including Windows and
Docker, with Mordor's beside Ethereum Classic's.

## What Mordor is

`--mordor` selects it. Its chain ID is listed under
[supported networks](../index.md#supported-networks), its network ID is `MordorChainConfig.NetworkID`
in `params/config_mordor.go`, and its genesis is in `params/genesis_mordor.go`. A running node
reports both IDs:

```shell
$ geth --mordor attach --exec 'eth.chainId()'
"0x3f"
$ geth --mordor attach --exec 'net.version'
"7"
```

## Pass `--mordor` on every command

Give `--mordor` to every `geth` command for a Mordor node, including `attach`, `dumpconfig`,
`export` and `import`, and pass it beside `--datadir` and `--config` too. Without the flag, `geth`
does not treat a Mordor data directory as Mordor.

Started without it, a Mordor data directory logs `Starting Core-Geth on Ethereum Classic...` and
reports network ID `63` instead of Mordor's `7`. Mordor nodes refuse to peer with a node whose
network ID differs from theirs. A configuration file does not carry the network either: see the
warning under [Configuration](run-cli.md#configuration).

## Start it, and where the data lives

```shell
$ geth --mordor --datadir <datadir>
```

Without `--datadir`, the node uses a `mordor` directory beside Ethereum Classic's `classic`
directory, `~/.ethereum/mordor` on Linux. Inside it, the layout is the one described under
[where your data lives](run-classic-node.md#where-your-data-lives).

An account is not tied to a network: the same key controls the same address on Ethereum
Classic and on Mordor alike. What differs by default is the directory `geth` reads it from,
since the keystore sits inside each network's own data directory, so a keyfile created under
Ethereum Classic's default directory does not appear here unless you copy it into this
directory's `keystore/` or point `--keystore` at it directly.

A new Mordor node's first lines name the network. Check them before anything else:

```
INFO [09-13|01:26:34.012] Starting Core-Geth on Mordor testnet...
INFO [09-13|01:26:34.287] Initialising Ethereum protocol           network=7 dbversion=<nil>
INFO [09-13|01:26:34.288] Writing custom genesis block
INFO [09-13|01:26:34.304] Loaded most recent local block           number=0 hash=a68ebd..59a3f1 td=131,072 age=7y2w2d
```

`network=7` is Mordor's network ID, and `a68ebd..59a3f1` is its genesis hash, shortened. The sync
ends the same way it does on Ethereum Classic:

```
INFO [09-13|01:45:24.101] Snap sync complete, auto disabling
INFO [09-13|01:45:24.101] Enabled artificial finality features     reason=synced peers=13
```

As on Ethereum Classic, the second line appears only when MESS is on
([MESS on this node](run-classic-node.md#mess-on-this-node)).

The warnings explained for
[a first sync of Ethereum Classic](run-classic-node.md#warnings-a-first-sync-logs) can appear on
Mordor too.

## Disk and time

A Mordor node's disk, memory and sync time are measured under
[Mordor](hardware-requirements.md#mordor) on the hardware page.

## Test coins

[Mine on Mordor](../guides/mordor-mining.md) is the other way to get coins. Two public faucets are set up to give out Mordor coins. This is what each showed on 2026-09-13:

| Faucet | How you ask | State on 2026-09-13 |
| --- | --- | --- |
| [`mordortestnet/mordor-public-faucet`](https://github.com/mordortestnet/mordor-public-faucet) | Open a pull request | No request received and no distribution recorded; no funding address published |
| [faucet.etcmc-monitor.org](https://faucet.etcmc-monitor.org/) | Enter an address and mine on the site, then claim what you collected | Funded; 109 claims waiting, the oldest from 2026-08-08, and none sent |

## Mordor beside Ethereum Classic on one host

Two nodes on one host need separate data directories and separate ports. Each flag below moves
one listener:

| Listener | Flag | [Default](../operate/security.md#ports-and-listeners) | Mordor node in this example |
| --- | --- | --- | --- |
| Peer connections (TCP) | `--port` | 30303 | 30304 |
| Discovery (UDP) | follows `--port`; `--discovery.port` sets it separately | 30303 | 30304 |
| HTTP JSON-RPC, with `--http` | `--http.port` | 8545 | 8547 |
| WebSocket JSON-RPC, with `--ws` | `--ws.port` | 8546 | 8548 |
| Metrics, with `--metrics` and `--metrics.addr` | `--metrics.port` | 6060 | 6070 |
| Profiling, with `--pprof` | `--pprof.port` | 6060 | 6071 |

A second node that reuses the first one's peer, HTTP or WebSocket port exits at
startup with `bind: address already in use`
([Troubleshooting](../operate/troubleshooting.md#why-does-the-node-stop-with-bind-address-already-in-use)).
A metrics port that is already taken does not stop
the node: its metrics server fails to start, and it logs `Failure in running metrics server`.
[Ports and listeners](../operate/security.md#ports-and-listeners) covers the shared metrics and
profiling default.

With the Ethereum Classic node on the defaults, the Mordor node needs only the ports it would
otherwise share:

```shell
$ geth --classic --datadir <classic-datadir> --http
$ geth --mordor --datadir <mordor-datadir> --port 30304 --http --http.port 8547
```

On Linux and macOS, each node's IPC socket is inside its own data directory. On Windows, every
node's IPC endpoint is the named pipe `\\.\pipe\geth.ipc` unless `--ipcpath` names another, so
give the second node its own, such as `--ipcpath geth-mordor.ipc`.

On Linux, this lists what each node listens on:

```shell
$ ss -lntup
```

## Attach

`geth --mordor attach` finds a node that uses the default data directory:

```shell
$ geth --mordor attach --exec 'eth.blockNumber'
```

For a node with its own `--datadir`, give the socket's path:

```shell
$ geth --mordor attach --exec 'eth.blockNumber' <datadir>/geth.ipc
```

On Windows, give the pipe's name instead of a path. `geth --mordor attach` with no name connects
to `\\.\pipe\geth.ipc`, which with two nodes running is whichever node holds that name.

Without `--mordor`, `geth attach` looks for Ethereum Classic's socket and does not find a Mordor
node:

```
Fatal: Unable to attach to remote geth: dial unix <home>/.ethereum/classic/geth.ipc: connect: no such file or directory
```

For a node that runs as a systemd service, attach as the service's user, as the
[Linux users guide](run/linux.md#6-keep-it-running-a-systemd-service) shows.
