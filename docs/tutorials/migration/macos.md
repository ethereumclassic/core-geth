---
description: "Upgrade an Ethereum Classic node on macOS from any v1.12.x release to Core-Geth v1.13.0: the steps, the node-key rotation, and how to roll back."
---

# Mac users guide

This guide upgrades a node on macOS from any v1.12.x release to v1.13.0.
[Migrating to v1.13.0](../v1.13.0-migration.md) explains why, and what changes for every node.

The commands use the default Ethereum Classic data directory on macOS, `~/Library/Ethereum/classic`. If you
pass `--datadir`, use yours; for Mordor, use `~/Library/Ethereum/mordor`.

## Choose how you update

| Method | Use it if |
|---|---|
| [Release archive](#method-a-release-archive) | You run a published binary: `core-geth-osx-arm64` on Apple Silicon, `core-geth-osx` on Intel |
| [Build from source](#method-b-build-from-source) | You compile the client yourself |

`uname -m` prints `arm64` on Apple Silicon and `x86_64` on Intel.

## Before you start

Record what you are running:

```bash
geth version
geth attach --exec 'eth.blockNumber'      ~/Library/Ethereum/classic/geth.ipc
geth attach --exec 'net.peerCount'        ~/Library/Ethereum/classic/geth.ipc
geth attach --exec 'admin.nodeInfo.enode' ~/Library/Ethereum/classic/geth.ipc
```

**Keep that last line.** It is your current enode ID, and step 3 changes it.

## 1. Stop the node cleanly

- **In a terminal:** press Ctrl-C once and wait. A second interrupt does not kill it.
- **Under launchd:** stop the job you created.

Wait for `INFO [..] Blockchain stopped` before you continue. It confirms the database is consistent.

## 2. Install v1.13.0

Keep the old binary until the upgrade is verified; step 6 uses it.

### Method A: release archive

Download the archive for your processor and check it, as [Installation, macOS](../../getting-started/installation.md#macos)
shows, including clearing the quarantine attribute if a browser downloaded it, then
[verify where it came from](../../getting-started/installation.md#verify-where-the-archive-came-from). Replace
the binary:

```bash
sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12
sudo install -m 0755 geth /usr/local/bin/geth
geth version
```

### Method B: build from source

Building needs Go 1.26 and the Xcode Command Line Tools (`xcode-select --install`), which provide the C
compiler and `make`. Build from the release branch, not from `main`:

```bash
git clone --branch archive-release-v1.13.0 --single-branch https://github.com/ethereumclassic/core-geth.git
cd core-geth
make geth
./build/bin/geth version
sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12
sudo install -m 0755 build/bin/geth /usr/local/bin/geth
```

**Do not clone with `--recursive`.** The submodules hold test fixtures of several gigabytes that
`make geth` does not read.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why, and
which peer lists to update. With the node stopped:

```bash
mv ~/Library/Ethereum/classic/geth/nodekey ~/Library/Ethereum/classic/geth/nodekey.old-rotated-$(date +%F)
```

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags you
start the node with.

## 5. Start and verify

Start the node the way you always do. Then, in this order:

```bash
geth version                                                                 # Core-Geth 1.13.0
geth attach --exec 'admin.nodeInfo.enode' ~/Library/Ethereum/classic/geth.ipc   # the new ID: send it to your peers
geth attach --exec 'net.peerCount'        ~/Library/Ethereum/classic/geth.ipc   # climbs within minutes
geth attach --exec 'eth.blockNumber'      ~/Library/Ethereum/classic/geth.ipc   # resumes where it stopped
```

**A peer count stuck at zero after ten minutes usually means a peer still lists your old enode ID.** Last,
compare `eth.blockNumber` with a block explorer or another node you operate.

## 6. If you need to roll back

The data directory works in both directions. Stop the node, restore the old binary, and start it again:

```bash
sudo mv /usr/local/bin/geth-v1.12 /usr/local/bin/geth
```

The rotated node key stays rotated: rolling back does not un-expose the old key.
