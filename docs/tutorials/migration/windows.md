---
description: "Upgrade an Ethereum Classic node on Windows from any v1.12.x release to Core-Geth v1.13.0: the steps, the node-key rotation, and how to roll back."
---

# Windows users guide

This guide upgrades a node on Windows from any v1.12.x release to v1.13.0.
[Migrating to v1.13.0](../v1.13.0-migration.md) explains why, and what changes for every node.

The commands are for PowerShell. They use the default Ethereum Classic data directory,
`%LOCALAPPDATA%\Ethereum\classic`. The node uses `%USERPROFILE%\AppData\Roaming\Ethereum\classic` instead if
an older `AppData\Roaming\Ethereum` folder exists and is not empty. If you pass `--datadir`, use yours; for
Mordor, use `mordor` in place of `classic`.

On Windows the node listens for `geth attach` on a named pipe, `\\.\pipe\geth.ipc`, not on a file in the data
directory.

## Choose how you update

| Method | Use it if |
|---|---|
| [Release archive](#method-a-release-archive) | You run the published binary, `core-geth-win64` |
| [Build from source](#method-b-build-from-source) | You compile the client yourself |

## Before you start

Record what you are running, from the folder that holds `geth.exe`:

```powershell
.\geth.exe version
.\geth.exe attach --exec 'eth.blockNumber'      \\.\pipe\geth.ipc
.\geth.exe attach --exec 'net.peerCount'        \\.\pipe\geth.ipc
.\geth.exe attach --exec 'admin.nodeInfo.enode' \\.\pipe\geth.ipc
```

**Keep that last line.** It is your current enode ID, and step 3 changes it.

## 1. Stop the node cleanly

- **In a console window:** press Ctrl+C once and wait.
- **As a service:** stop it with the tool you installed it with, and give it time to shut down cleanly.

Wait for `INFO [..] Blockchain stopped` before you continue. It confirms the database is consistent.

## 2. Install v1.13.0

Keep the old binary until the upgrade is verified; step 6 uses it.

### Method A: release archive

Download and check `core-geth-win64-v1.13.0.zip` as [Installation, Windows](../../getting-started/installation.md#windows)
shows, then [verify where it came from](../../getting-started/installation.md#verify-where-the-archive-came-from).
In the folder that holds your current `geth.exe`:

```powershell
Rename-Item .\geth.exe geth-v1.12.exe
Expand-Archive <download folder>\core-geth-win64-v1.13.0.zip -DestinationPath .
.\geth.exe version
```

### Method B: build from source

Building needs Go 1.26 and a C compiler on `PATH`, such as MinGW-w64. Build from the release branch, not
from `main`:

```powershell
git clone --branch archive-release-v1.13.0 --single-branch https://github.com/ethereumclassic/core-geth.git
cd core-geth
go run build/ci.go install ./cmd/geth
.\build\bin\geth.exe version
```

That is the command `make geth` runs on other systems. Rename your current `geth.exe` to
`geth-v1.12.exe`, then copy `build\bin\geth.exe` in its place. **Do not clone with `--recursive`:** the
submodules hold test fixtures of several gigabytes that the build does not read.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why, and
which peer lists to update. With the node stopped:

```powershell
$key = "$env:LOCALAPPDATA\Ethereum\classic\geth\nodekey"
Rename-Item $key "nodekey.old-rotated-$(Get-Date -Format yyyy-MM-dd)"
```

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags you
start the node with.

## 5. Start and verify

Start the node the way you always do. Then, in this order:

```powershell
.\geth.exe version                                                # Core-Geth 1.13.0
.\geth.exe attach --exec 'admin.nodeInfo.enode' \\.\pipe\geth.ipc   # the new ID: send it to your peers
.\geth.exe attach --exec 'net.peerCount'        \\.\pipe\geth.ipc   # climbs within minutes
.\geth.exe attach --exec 'eth.blockNumber'      \\.\pipe\geth.ipc   # resumes where it stopped
```

**A peer count stuck at zero after ten minutes usually means a peer still lists your old enode ID.** Last,
compare `eth.blockNumber` with a block explorer or another node you operate.

## 6. If you need to roll back

The data directory works in both directions. Stop the node, then restore the old binary and start it again:

```powershell
Remove-Item .\geth.exe
Rename-Item .\geth-v1.12.exe geth.exe
```

The rotated node key stays rotated: rolling back does not un-expose the old key.
