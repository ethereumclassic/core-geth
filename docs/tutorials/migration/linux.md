---
description: "Upgrade an Ethereum Classic node on Linux from any v1.12.x release to Core-Geth v1.13.0: the steps, the node-key rotation, and how to roll back."
---

# Linux users guide

This guide upgrades a node on Linux from any v1.12.x release to v1.13.0.
[Migrating to v1.13.0](../v1.13.0-migration.md) explains why, and what changes for every node.

The commands use the default Ethereum Classic data directory, `~/.ethereum/classic`. If you pass
`--datadir`, use yours; for Mordor, use `~/.ethereum/mordor`.

## Choose how you update

| Method | Use it if |
|---|---|
| [Release archive](#method-a-release-archive) | You run a published binary: `core-geth-linux` on x86-64, `core-geth-arm64` on 64-bit ARM, or `core-geth-arm5`, `arm6` or `arm7` on 32-bit ARM |
| [Build from source](#method-b-build-from-source) | You compile the client yourself |

Both methods share every other step.

## Before you start

Record what you are running, so you can tell afterwards whether anything changed that you did not intend:

```bash
geth version
geth attach --exec 'eth.blockNumber'      ~/.ethereum/classic/geth.ipc
geth attach --exec 'net.peerCount'        ~/.ethereum/classic/geth.ipc
geth attach --exec 'admin.nodeInfo.enode' ~/.ethereum/classic/geth.ipc
```

**Keep that last line.** It is your current enode ID, and step 3 changes it.

## 1. Stop the node cleanly

- **In a terminal:** press Ctrl-C once and wait. A second interrupt does not kill it.
- **Under systemd:** stop the unit you wrote, for example `sudo systemctl stop core-geth`.

Wait for this line before you continue. It confirms the database is consistent:

```
INFO [..] Blockchain stopped
```

[Stop it safely](../../getting-started/run-classic-node.md#stop-it-safely) says what each kind of stop does.

## 2. Install v1.13.0

Keep the old binary until the upgrade is verified; step 6 uses it.

### Method A: release archive

Download the archive for your processor and check it, as
[Installation, Linux x86_64](../../getting-started/installation.md#linux-x86_64) or
[Linux, ARM](../../getting-started/installation.md#linux-arm) shows, then
[verify where it came from](../../getting-started/installation.md#verify-where-the-archive-came-from). Replace
the binary:

```bash
sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12
sudo install -m 0755 geth /usr/local/bin/geth
geth version
```

### Method B: build from source

Building needs Go 1.26 and a C compiler. Build from the release branch, not from `main`, which keeps moving:

```bash
git clone --branch archive-release-v1.13.0 --single-branch https://github.com/ethereumclassic/core-geth.git
cd core-geth
make geth
./build/bin/geth version
sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12
sudo install -m 0755 build/bin/geth /usr/local/bin/geth
```

**Do not clone with `--recursive`.** The submodules hold consensus test fixtures of several gigabytes, and
`make geth` does not read them. [Build from source](../../developers/build-from-source.md) covers the
dependencies.

The executable is still named `geth`, so your units and scripts keep working unchanged.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why, and
which peer lists to update. With the node stopped:

```bash
mv ~/.ethereum/classic/geth/nodekey ~/.ethereum/classic/geth/nodekey.old-rotated-$(date +%F)
```

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags in
your unit file or start script.

## 5. Start and verify

Start the node the way you always do, for example `sudo systemctl start core-geth`. Then, in this order:

```bash
geth version                                                          # Core-Geth 1.13.0
geth attach --exec 'admin.nodeInfo.enode' ~/.ethereum/classic/geth.ipc   # the new ID: send it to your peers
geth attach --exec 'net.peerCount'        ~/.ethereum/classic/geth.ipc   # climbs within minutes
geth attach --exec 'eth.blockNumber'      ~/.ethereum/classic/geth.ipc   # resumes where it stopped
```

**A peer count stuck at zero after ten minutes usually means a peer still lists your old enode ID.** Last,
compare `eth.blockNumber` with a block explorer or another node you operate. If it settles somewhere else,
stop and investigate rather than resyncing.

## 6. If you need to roll back

The data directory works in both directions, so rolling back swaps the binary. Stop the node, then:

```bash
sudo mv /usr/local/bin/geth-v1.12 /usr/local/bin/geth
```

Start it again. The rotated node key stays rotated: rolling back does not un-expose the old key.

## This procedure was executed before it was published

The claims above are measured, not inferred. The steps in this guide were run end to end
against a real Mordor data directory carrying a full chain, using the published `v1.12.23`
binary and a `v1.13.0` build:

| What was tested | Result |
|---|---|
| In-place upgrade, `v1.12.23` → `v1.13.0`, same datadir | Head preserved exactly; no resync |
| Rollback, `v1.13.0` datadir reopened by `v1.12.23` | Reopened cleanly; head preserved |
| Node-key rotation (step 3) | New key generated; chain data untouched |
| RPC exposure (step 4) | `--http` alone binds `127.0.0.1`; `--http.addr 0.0.0.0` binds all interfaces |

The rollback arm is the one worth calling out. It is the step most likely to be wrong and the
most damaging if it is, because it is the one an operator reaches for when something has
already gone badly, and nothing had previously tested whether the older client reopens a data
directory the newer one has written. It does, with the head intact.

This was run with peer discovery disabled, so any change in head across a restart would have
been the client rewinding rather than the network advancing.
