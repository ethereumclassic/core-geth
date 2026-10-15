---
title: Upgrading
description: "How to upgrade a Core-Geth node to the latest release, and the extra steps a few upgrades need, listed by the release line a node is leaving."
---

# Upgrading

**Run the latest release.** Most upgrades are the five steps below and nothing else. A few need more,
and those are listed further down by the release line a node is leaving. [Releases](releases.md) links
the notes for every release.

## Upgrade to the latest release

1. **Read what changed.** Each release's notes are on the
   [releases page](https://github.com/ethereumclassic/core-geth/releases), and published
   [security advisories](https://github.com/ethereumclassic/core-geth/security/advisories) describe
   the vulnerabilities a release fixes.
   [`SECURITY.md`](https://github.com/ethereumclassic/core-geth/blob/main/SECURITY.md) says how they
   are disclosed.
2. **Download and verify the new release**, following your route in
   [Choose your route](getting-started/installation.md#choose-your-route) and then
   [Verify where the archive came from](getting-started/installation.md#verify-where-the-archive-came-from).
3. **Stop the node** and wait for the process to exit
   ([Stop it safely](getting-started/run-classic-node.md#stop-it-safely)).
4. **Replace the binary**, or the image tag in a container, and start the node.
5. **Check the result.** `geth version` reports the new release, and
   [How to tell it is done](getting-started/run-classic-node.md#how-to-tell-it-is-done) confirms
   the node is back in sync.

## Upgrades that need more than a new binary

Find the release line the node is leaving. If it is not listed here, the five steps above are the
whole upgrade.

### Leaving the v1.12.x line

A node on any `v1.12.x` release has more to do than replace the binary:

- **Rotate the P2P node key.** This is required. One of the vulnerabilities in that line leaks bits
  of the key, and a fix cannot recall what already leaked.
  [Rotate the P2P node key](tutorials/v1.13.0-migration.md#rotate-the-p2p-node-key) says why, and which
  peer lists to update.
- **Check how the node is started.** With no network flag the client runs Ethereum Classic, where a
  `v1.12.x` node ran Ethereum mainnet.
  [If you start the node without a network flag](tutorials/v1.13.0-migration.md#if-you-start-the-node-without-a-network-flag)
  lists what that changes.
- **Decide your MESS setting**, and give every node in a fleet the same one. An upgrade can change
  whether a node runs MESS. [MESS](operate/mess.md) sets out what the setting decides.
- **Check a custom genesis or chainspec file**, if the node runs a private network.
  [If your chain uses its own genesis or chainspec file](tutorials/v1.13.0-migration.md#if-your-chain-uses-its-own-genesis-or-chainspec-file)
  lists three kinds of value the client reads differently.
- **Re-check RPC exposure** while the node is down.
  [Re-check RPC exposure](tutorials/v1.13.0-migration.md#re-check-rpc-exposure) is the checklist.

[Migrating to v1.13.0](tutorials/v1.13.0-migration.md) is the full guide, with every step for
[Linux](tutorials/migration/linux.md), [macOS](tutorials/migration/macos.md),
[Windows](tutorials/migration/windows.md) and [Docker](tutorials/migration/docker.md), and how to roll
back. Where the guide names `v1.13.0`, install the latest release instead.
