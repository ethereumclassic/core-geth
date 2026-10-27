---
title: Mining pool node
description: "Set up the Core-Geth node behind an Ethereum Classic mining pool: serving work over JSON-RPC, node configuration, and where the pool software fits."
---

A mining pool runs pool software between its miners and one or more Core-Geth nodes. The node builds
blocks and hands out work over JSON-RPC. The pool software splits that work among its miners, checks their
shares, and submits solutions back to the node. This page covers the node side.

## Start the node

```sh
geth --classic --mine --miner.etherbase 0xPOOL_ADDRESS \
  --http --http.addr 127.0.0.1 --http.port 8545 --http.api eth,net,web3
```

- `--miner.etherbase` is the pool's reward address.
- `--mine` with the default `--miner.threads 0` builds work without mining on the CPU.
- The pool software calls `eth_getWork` for work and `eth_submitWork` with each solution.

Anyone who can reach the HTTP endpoint can request work and submit results, so keep it reachable only by
the pool software: [RPC exposure](../operate/security.md#rpc-exposure).

## Get new work to the pool quickly

- `--miner.notify http://127.0.0.1:<port>/` makes the node post each new work package to the pool software
  as soon as it has one, instead of waiting to be polled. Separate several URLs with commas.
- `--miner.notify.full` sends the pending block's header instead of a work package, for pool software that
  expects it. Set it on the command line: the config file's `NotifyFull` key is not read.
- `--miner.recommit`, 2s by default, is how often the node rebuilds the block template to take in new
  transactions.

## Stay connected

A pool node that falls behind builds work on an old head, and its miners' solutions go stale.

- Keep the p2p port open on TCP and UDP: [Ports and listeners](../operate/security.md#ports-and-listeners).
- Peer the pool's nodes with each other and with nodes you trust:
  [Static and trusted peers](bootnodes-and-discovery.md#static-and-trusted-peers).
- Keep the clock synchronized: [Clock](../operate/security.md#clock).
- Alert when the peer count drops or the head stops advancing: [Monitoring](../operate/monitoring.md).

## Run more than one node

Run at least two nodes on separate hosts and point the pool software at both, so a restart or an upgrade
does not stop the pool. Keep them on the same release and the same [MESS](../operate/mess.md) setting, so
they agree on the chain during a deep reorganization.

## Before the pool pays out

Payouts are only as final as the blocks behind them. Wait for enough confirmations before crediting
miners. The [MESS confirmation calculator](mess-calculator.md) shows how much hashrate it takes to
reorganize a block of a given age on nodes that run MESS.
