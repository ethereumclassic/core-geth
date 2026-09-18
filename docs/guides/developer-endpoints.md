---
title: Endpoints for your own project
description: "Run your own Ethereum Classic and Mordor JSON-RPC endpoints for an application you are building: which node answers which call, what to expose, and when you need an archive node."
---

This is for a developer who needs endpoints for their own application, rather than for an operator serving
the public. The difference decides most of the configuration: your consumers are your own services, you
know their addresses, and nobody else needs to reach the node.

If you are serving strangers, [Public RPC endpoint](public-rpc-endpoint.md) is the page for that, and it
covers hostnames, browser origins, gateways and the limits the node enforces.

## Why run your own

A public endpoint is fine for a first contract deployment and becomes the wrong dependency as soon as the
project matters:

- **The methods you need for development are the ones public endpoints disable.** `debug_traceTransaction`,
  `trace_*` and state at an old block are how you find out why a transaction did what it did. No public
  endpoint should expose them, and ours does not.
- **Rate limits apply per source, and your test suite is a single source.** A CI run that replays a few
  thousand calls looks exactly like abuse.
- **Every call tells the operator what you are building**, which addresses you watch and when you deploy.
- **It is one node.** Running it yourself is a package install and a systemd unit:
  [Installation](../getting-started/installation.md), [Running a node](../getting-started/run-a-node.md).

## Start on Mordor

Build against the test network first. Same client, same RPC surface, no real funds, and test coins take
twenty minutes to mine rather than a faucet request that may never be answered:
[Run a Mordor node](../getting-started/run-mordor-node.md), then
[Test coins without a faucet](mordor-mining.md#test-coins-without-a-faucet).

Everything below applies to both networks. Use `--mordor` where it says `--classic`.

## Which node answers which call

| You need | Node |
|---|---|
| Send transactions, read current state and receipts, watch logs from now on | A full node, snap sync |
| `eth_call` or `eth_getBalance` against a block from last month | An archive node |
| `debug_traceTransaction` or `trace_block` on old blocks | An archive node |
| Reindex your application's history from genesis | An archive node |

Most projects need one full node and reach for an archive node only when they start asking about the past.
[Sync modes and data retention](../operate/sync-modes.md) explains what each keeps, and
[Archive node](archive-node.md) covers running one, including the disk it takes.

## A node for your application

```sh
geth --classic \
  --http --http.addr 127.0.0.1 --http.port 8545 --http.api eth,net,web3 \
  --ws --ws.addr 127.0.0.1 --ws.port 8546 --ws.api eth,net,web3
```

- **Keep it on loopback.** If your application runs on another host, put the node on a private network or
  behind a proxy you control rather than binding the RPC to a public address.
  [Security and network exposure](../operate/security.md#rpc-exposure) is the short version of why.
- **`eth,net,web3` is what an application needs.** Add namespaces deliberately, not preemptively.
- **WebSocket is what you want for subscriptions**, `eth_subscribe` for new heads and logs, rather than
  polling `eth_getLogs` on a timer.
- **`personal` and account management do not belong on a node your application talks to.** Sign in your
  application or in a signer, and keep keys out of the node.

## Adding tracing for development

```sh
geth --classic \
  --gcmode archive --syncmode full \
  --http --http.addr 127.0.0.1 --http.api eth,net,web3,debug,trace
```

`debug` and `trace` answer for the blocks whose state the node still has, which is why tracing an old
transaction needs the archive node above rather than a pruned one. Expose them on loopback only, and never
through the proxy that serves your application's users.

[JSON-RPC API](../JSON-RPC-API/index.md) lists every method by module, and
[trace module overview](../JSON-RPC-API/trace-module-overview.md) covers the tracing calls specifically.

## What to hold yourself to

- **Pin the client version your CI runs**, so a test that starts failing tells you about your code rather
  than about an upgrade. Releases are at
  [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth/releases).
- **Back up nothing, and be able to rebuild everything.** A development node holds no state you cannot
  resync. Keep the genesis and config in your repository, not in a server you are afraid to lose.
- **Watch the node the way you watch your application.** `eth_syncing`, peer count and head age catch the
  failure where everything looks healthy and the chain stopped advancing:
  [Monitoring](../operate/monitoring.md).
- **Upgrade with the network.** A node that falls behind a fork stops being a test of anything:
  [Maintenance, backup and upgrades](../operate/maintenance.md).

## When your project outgrows this

Serving your application's users directly from your node makes you an RPC provider, with a different set of
problems: [Public RPC endpoint](public-rpc-endpoint.md) for the exposure question and
[Production operations](production-operations.md) for running more than one node.
