---
title: Monitoring
description: "Monitor a Core-Geth Ethereum Classic node: the metrics it serves over HTTP, reading its log, and building a health check from JSON-RPC calls on the node's own machine."
---

A node reports on itself through metrics served over HTTP and through its log. It has no health
endpoint: a health check is built from ordinary JSON-RPC calls made on the node's own machine.

## Metrics

Metrics are off by default. `--metrics` turns collection on, and `--metrics.addr` starts the HTTP
server that serves the values:

```shell
$ geth --classic --datadir <datadir> --metrics --metrics.addr 127.0.0.1 --metrics.port <port>
```

```
INFO [..] Starting metrics server                  addr=http://127.0.0.1:<port>/debug/metrics
```

`--metrics.port` can be left out. [Ports and listeners](security.md#ports-and-listeners) gives its
default, and what to change when pprof runs on the same node.

The server answers at two paths:

| Path | Format |
| --- | --- |
| `/debug/metrics` | JSON, one value per key |
| `/debug/metrics/prometheus` | The Prometheus text format |

A metric keeps its name in both formats, with each `/` of the JSON key written as `_` in the
Prometheus one. The peer count and the head block, from a node that has just started:

```shell
$ curl -s http://127.0.0.1:<port>/debug/metrics/prometheus | grep -E '^(# TYPE )?(p2p_peers|chain_head_block) '
# TYPE chain_head_block gauge
chain_head_block 0
# TYPE p2p_peers gauge
p2p_peers 0
```

In `/debug/metrics` the same two values are under `p2p/peers` and `chain/head/block`. With
`--metrics.expensive`, the node also collects the metrics its code marks as costly to gather.

**The configuration file does not reach metrics.** A `[Metrics]` section, including the one
`geth dumpconfig` writes, neither turns collection on nor starts the server. Set the flags on the
command line, as in a service unit's `ExecStart`, or through the environment variables
`geth --help` names beside each flag, such as `GETH_METRICS` and `GETH_METRICS_ADDR`.

**Keep the server on `127.0.0.1` or a private network.** Besides the metrics, `/debug/metrics`
returns Go's memory statistics and the node's full command line, including any password given there
as a flag.

### Pushing to InfluxDB

With `--metrics` set, the node can also push its metrics to InfluxDB:

| Flag | InfluxDB version | What it sets |
| --- | --- | --- |
| `--metrics.influxdb` | 1 | Turns on the push |
| `--metrics.influxdbv2` | 2 | Turns on the push |
| `--metrics.influxdb.endpoint` | both | The InfluxDB API endpoint |
| `--metrics.influxdb.tags` | both | Tags added to every measurement |
| `--metrics.influxdb.database` | 1 | The database |
| `--metrics.influxdb.username`, `--metrics.influxdb.password` | 1 | The credentials |
| `--metrics.influxdb.bucket`, `--metrics.influxdb.organization`, `--metrics.influxdb.token` | 2 | The bucket, organization and token |

Give the password or token in the environment, as `GETH_METRICS_INFLUXDB_PASSWORD` or
`GETH_METRICS_INFLUXDB_TOKEN`. A value given as a flag is part of the command line, which
`/debug/metrics` returns. When the endpoint cannot be reached, the node logs `Unable to send to
InfluxDB` warnings and keeps running.

## What to watch

| Signal | Where | What it means |
| --- | --- | --- |
| `p2p_peers` | Metrics | The number of connected peers. At zero, the node receives no new blocks ([no peers](troubleshooting.md#why-does-the-node-have-no-peers)) |
| `chain_head_block` | Metrics | The head block number. On a synced node it rises as blocks arrive; when it stops rising, the node has stopped importing |
| `Unclean shutdown detected` | Log, as the node starts | The node has recorded an unclean stop ([After an unclean stop](maintenance.md#after-an-unclean-stop)) |
| `Disk space is running low` | Log | Free space is nearing the level at which the node stops itself ([Running out of disk](../getting-started/hardware-requirements.md#running-out-of-disk)) |
| `Disabled artificial finality features` | Log | MESS has switched itself off on this node ([MESS on this node](../getting-started/run-classic-node.md#mess-on-this-node)) |
| `Reorg disallowed` | Log | MESS refused a reorganization to another chain ([Troubleshooting](troubleshooting.md#what-do-disabled-artificial-finality-features-and-reorg-disallowed-mean)) |
| `NRestarts` | systemd | How many times systemd has restarted the node's service |

For a node run as the `core-geth` service
([Linux users guide](../getting-started/run/linux.md#6-keep-it-running-a-systemd-service)):

```shell
$ systemctl show core-geth -p NRestarts
NRestarts=1
```

## Health checks

A check runs on the node's machine and talks to the node over its IPC socket.

**Live** means the node answers:

```shell
$ geth --classic attach --exec 'eth.blockNumber' <datadir>/geth.ipc
```

It prints the head block number and exits with status 0. With no node behind the socket, it exits
with status 1:

```
Fatal: Unable to attach to remote geth: dial unix <socket>: connect: no such file or directory
```

**Ready** means all three of these hold:

1. `eth.syncing` is `false`. It stays an object until the node has finished syncing
   ([How to tell it is done](../getting-started/run-classic-node.md#how-to-tell-it-is-done)).
2. `net.peerCount` is at least the number of peers you require.
3. The latest block is recent: its timestamp is no older than the age you allow.

The last two catch what the first misses. A synced node that has lost its peers, and has stopped
receiving blocks, still reports `eth.syncing` as `false`.

This script checks all three and exits with status 0 when the node is ready. `MIN_PEERS` and
`MAX_HEAD_AGE` are examples to set for your own node.

```sh
#!/bin/sh
# Readiness check for a core-geth node, over its IPC socket. Exits 0 when the node is ready.
IPC="$1"            # the node's socket, <datadir>/geth.ipc
MIN_PEERS=3         # the fewest peers you accept
MAX_HEAD_AGE=600    # the oldest head you accept, in seconds

state=$(geth --classic attach --exec '[eth.syncing === false, net.peerCount, Math.floor(Date.now() / 1000) - eth.getBlock("latest").timestamp].join(" ")' "$IPC") \
  || { echo "not live: no answer on $IPC"; exit 2; }
set -- $(echo "$state" | tr -d '"')
[ "$1" = true ] || { echo "not ready: still syncing"; exit 1; }
[ "$2" -ge "$MIN_PEERS" ] || { echo "not ready: too few peers"; exit 1; }
[ "$3" -le "$MAX_HEAD_AGE" ] || { echo "not ready: head too old"; exit 1; }
echo "ready"
```

Saved as `readiness.sh` and run against a node that is still syncing:

```shell
$ sh readiness.sh <datadir>/geth.ipc
not ready: still syncing
```

Over local HTTP the same checks are the
[`eth_syncing`](../JSON-RPC-API/modules/eth.md#eth_syncing),
[`net_peerCount`](../JSON-RPC-API/modules/net.md#net_peercount) and
[`eth_getBlockByNumber`](../JSON-RPC-API/modules/eth.md#eth_getblockbynumber) calls.

## Logs

The node logs to standard error. Under systemd that goes to the journal, which the
[Linux users guide](../getting-started/run/linux.md#6-keep-it-running-a-systemd-service) shows how to follow.

| Flag | What it does |
| --- | --- |
| `--verbosity` | How much to log, from `0`, nothing, to `5`, the most detail |
| `--log.vmodule` | A different verbosity per module, such as `eth/*=5,p2p=4` |
| `--log.format` | The format: `terminal`, `logfmt` or `json` |
| `--log.file` | Also writes the log to a file |
| `--log.rotate` | Rotates the log file |
| `--log.maxsize` | The size, in megabytes, at which a log file is rotated |
| `--log.maxbackups` | How many rotated files to keep |
| `--log.maxage` | How many days to keep a rotated file |
| `--log.compress` | Compresses rotated files |

[Command-line Options](../getting-started/run-cli.md#command-line-options) gives each default.

- **`--log.file` appends to the file**, and the node still logs to standard error.
- **`--log.maxsize`, `--log.maxbackups`, `--log.maxage` and `--log.compress` apply only with
  `--log.rotate`.** Give `--log.rotate` a `--log.file` as well. Without one, it writes to
  `geth-lumberjack.log` in the system's temporary directory.
- **Use `json` or `logfmt` for a file.** In the default `terminal` format, a node started from an
  interactive terminal writes color codes into the file too.

With `--log.format json`, each line is one JSON object:

```shell
$ geth --classic --datadir <datadir> --log.format json
```

```
{"t":"<time>","lvl":"info","msg":"Starting Core-Geth on Ethereum Classic..."}
```

## ethstats

`--ethstats` reports the node to an ethstats server, a dashboard of the nodes that report to it. The
value is the name this node reports under, followed by the server's secret and address:

```shell
$ geth --classic --datadir <datadir> --ethstats <node name>:<secret>@<host>:<port>
```

When the server cannot be reached, the node logs a warning and keeps running:

```
WARN [..] Stats server unreachable                 err="dial tcp <host>:<port>: connect: connection refused"
```

## Profiling

`--pprof` starts Go's profiling server, which serves data for diagnosing a problem in the node. It
listens on `127.0.0.1` unless `--pprof.addr` names another address. Keep it there
([Ports and listeners](security.md#ports-and-listeners)).

```
INFO [..] Starting pprof server                    addr=http://127.0.0.1:<port>/debug/pprof
```
