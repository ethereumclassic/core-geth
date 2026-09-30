---
title: Production operations
---

What an exchange, custodian, pool or infrastructure provider should have in place before depending on
core-geth. Each item points to the page that covers it.

## Before you depend on a node

| Area | What to have in place | Where |
|---|---|---|
| Sizing | Disk, memory and sync time for your network and sync mode | [Hardware requirements](../getting-started/hardware-requirements.md) |
| Service | The node runs under a service manager and restarts on failure | [Run it as a service](../getting-started/run-classic-node.md#run-it-as-a-service) |
| Data retention | The sync mode and history your queries need | [Sync modes](../operate/sync-modes.md), [Archive node](archive-node.md) |
| Exposure | Only the p2p port is public, and RPC stays private or behind a gateway | [Security and network exposure](../operate/security.md), [Public RPC endpoint](public-rpc-endpoint.md) |
| Clock | NTP keeps the clock synchronized | [Clock](../operate/security.md#clock) |
| Chain selection | Every node you run has the same MESS setting | [MESS](../operate/mess.md) |
| Monitoring | Alerts on peer count, head age and disk space | [Monitoring](../operate/monitoring.md) |
| Backups | Account keys and the node key backed up, and a restore tested | [What to back up](../operate/maintenance.md#what-to-back-up) |
| Upgrades | A tested way to upgrade and to roll back | [Upgrading within v1.13.x](../operate/maintenance.md#upgrading-within-v113x) |

## Redundancy

Run at least two nodes on separate hosts, and put your application behind both:
[Several nodes behind one entry point](public-rpc-endpoint.md#several-nodes-behind-one-entry-point).
Upgrade them one at a time, and keep them on the same release and MESS setting.

## Crediting deposits

- Read confirmations only from a node that is in sync: `eth_syncing` returns `false` and the head is recent.
- Choose the number of confirmations with the [MESS confirmation calculator](mess-calculator.md), and check
  that MESS is in force with `admin.ecbp1100Status()` over IPC.
- A node with few peers can be shown a chain the rest of the network does not follow. Alert when the peer
  count drops, and hold credits while it is low.

## When something goes wrong

- [Troubleshooting](../operate/troubleshooting.md) maps log lines to causes and fixes.
- [After an unclean stop](../operate/maintenance.md#after-an-unclean-stop) covers a node that was killed.
- `Reorg disallowed` in the log means MESS refused a reorganization. Check that your nodes agree on the head
  before you resume crediting.

## Security updates

Follow the [`ethereumclassic/core-geth` releases](https://github.com/ethereumclassic/core-geth/releases),
and report vulnerabilities privately as
[SECURITY.md](https://github.com/ethereumclassic/core-geth/blob/main/SECURITY.md) describes. Exchanges,
pools and service providers can reach the core developers at `security@ethereumclassic.net`.
