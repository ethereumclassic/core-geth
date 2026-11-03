---
description: "Upgrade an Ethereum Classic node on Windows from any v1.12.x release to the Core-Geth v1.13 line: the steps, the node-key rotation, and how to roll back."
---

# Windows users guide

This guide upgrades a node on Windows from any v1.12.x release to the latest v1.13 release.
[Migrating to v1.13](../v1.13.0-migration.md) explains why, and what changes for every node.

The commands are for PowerShell. They use the default Ethereum Classic data directory,
`%LOCALAPPDATA%\Ethereum\classic`. The node uses `%USERPROFILE%\AppData\Roaming\Ethereum\classic` instead if
an older `AppData\Roaming\Ethereum` folder exists and is not empty. If you pass `--datadir`, use yours; for
Mordor, use `mordor` in place of `classic`.

`geth attach` reaches the node on a named pipe, `\\.\pipe\geth.ipc`.

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

- **In a console window:** press Ctrl+C once and wait. A window opened from the Startup folder then asks
  `Terminate batch job (Y/N)?`; type `Y`.
- **Under a service wrapper:** stop it with that tool.

Wait until the log prints `Blockchain stopped` and the process exits.

## 2. Install the latest release

Keep the old binary until the upgrade is verified; step 6 uses it.

### Method A: release archive

In the folder that holds your current `geth.exe`, download the latest release and check it:

```powershell
$VERSION = (Invoke-RestMethod "https://api.github.com/repos/ethereumclassic/core-geth/releases/latest").tag_name
$BASE = "https://github.com/ethereumclassic/core-geth/releases/download/$VERSION"
Invoke-WebRequest "$BASE/core-geth-win64-$VERSION.zip" -OutFile "core-geth-win64-$VERSION.zip"
Invoke-WebRequest "$BASE/core-geth-win64-$VERSION.zip.sha256" -OutFile "core-geth-win64-$VERSION.zip.sha256"
(Get-FileHash "core-geth-win64-$VERSION.zip" -Algorithm SHA256).Hash
Get-Content "core-geth-win64-$VERSION.zip.sha256"
```

Compare the last two lines: the letters must match, and the case does not matter. If they differ, delete the file and
download it again. On a node that holds value, also check
[where it came from](../../getting-started/installation.md#verify-where-the-archive-came-from), which needs the
[GitHub CLI](https://cli.github.com/), signed in with `gh auth login`:

```powershell
gh attestation verify "core-geth-win64-$VERSION.zip" --repo ethereumclassic/core-geth
```

It must print `Verification succeeded`. If it does not, do not run the file.

Then, in the same window, keep the old binary and unpack the new one in its place:

```powershell
Rename-Item .\geth.exe geth-v1.12.exe
Expand-Archive "core-geth-win64-$VERSION.zip" -DestinationPath . -Force
.\geth.exe version
```

Go on only if `.\geth.exe version` reports a `v1.13` release.

### Method B: build from source

Building needs Go 1.26 and a C compiler on `PATH`, such as MinGW-w64. Build the latest release's tag, not
`main`, which keeps moving:

```powershell
$VERSION = (Invoke-RestMethod "https://api.github.com/repos/ethereumclassic/core-geth/releases/latest").tag_name
git clone --branch $VERSION --single-branch https://github.com/ethereumclassic/core-geth.git "core-geth-$VERSION"
cd "core-geth-$VERSION"
go run build/ci.go install ./cmd/geth
.\build\bin\geth.exe version
```

Rename your current `geth.exe` to `geth-v1.12.exe`, then copy `build\bin\geth.exe` in its place. Do not clone with
`--recursive`: the submodules are several gigabytes, and the build does not read them.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why. With the node
stopped:

```powershell
$key = "$env:LOCALAPPDATA\Ethereum\classic\geth\nodekey"
Rename-Item $key "nodekey.old-rotated-$(Get-Date -Format yyyy-MM-dd)"
```

The next start writes a fresh key. If you start the node with `--nodekey` or `--nodekeyhex`, or set `GETH_NODEKEY` or
`GETH_NODEKEYHEX`, take that out too: it overrides the file.

If the rename prints an error, check the data directory: `geth\nodekey` must not be there when the node starts.

### If this node lists static peers, update the list

They are under `[Node.P2P]` in the file you pass with `--config`. Edit it while the node is down:

```toml
[Node.P2P]
StaticNodes = ["enode://<node ID>@<IP>:30303"]
```

A peer that has already moved has a new ID, so use that. Separate entries with commas. `TrustedNodes` takes the same
entries. With no file yet, save these two lines and add `--config` and the file's path to the start command. The client
ignores `static-nodes.json`, so move its entries here.

Save the file as UTF-8 with no byte order mark. A byte order mark, or UTF-16, stops the start with
`invalid TOML syntax`.

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags you start the
node with. If you use a config file that `dumpconfig` wrote on `v1.12`,
[delete the lines listed here](../v1.13.0-migration.md#if-your-config-file-came-from-dumpconfig).

## 5. Start and verify

Start the node the way you always do. Then, in this order:

```powershell
.\geth.exe attach --exec 'admin.nodeInfo.name'  \\.\pipe\geth.ipc
.\geth.exe attach --exec 'admin.nodeInfo.enode' \\.\pipe\geth.ipc
.\geth.exe attach --exec 'net.peerCount'        \\.\pipe\geth.ipc
.\geth.exe attach --exec 'eth.blockNumber'      \\.\pipe\geth.ipc
```

If a line cannot connect, wait a minute and run it again.

- **The first line** must name a `v1.13` release.
- **The enode ID** must differ from the one you recorded. If it is the same, the node is still on the old key: stop it
  and redo step 3.
- **The peer count** climbs within minutes. If it stays at zero, see
  [Why does the node have no peers?](../../operate/troubleshooting.md#why-does-the-node-have-no-peers)
- **The block number** resumes where the node stopped, then catches up. Compare it with a block explorer or another
  node you run. If it settles somewhere else, stop and investigate rather than resync.

### If other nodes list this node, give them the new ID

They can no longer connect to this node, and neither side logs it at the default verbosity. On each one, on `v1.12` or
`v1.13`, run these lines with the old ID you recorded, the new ID from the block above, and the address and port the
entry already has:

```powershell
.\geth.exe attach --exec "admin.removePeer('enode://<old ID>@<IP>:30303')" \\.\pipe\geth.ipc
.\geth.exe attach --exec "admin.addPeer('enode://<new ID>@<IP>:30303')" \\.\pipe\geth.ipc
.\geth.exe attach --exec "admin.peers.map(function(p){return p.enode})" \\.\pipe\geth.ipc
```

The first two lines each print `true` and take effect with no restart. The third lists the new ID once the nodes
connect. The change lasts only until that node restarts, so put the new ID in its config file too. Send the new ID to
anyone else who lists this node.

For an entry in `TrustedNodes`, use `admin.removeTrustedPeer` and `admin.addTrustedPeer`. They do not open a connection.
[Static and trusted peers](../../guides/bootnodes-and-discovery.md#static-and-trusted-peers) explains the two lists.

## 6. If you need to roll back

The data directory works in both directions. Stop the node, and take any `--mess` flag out of the start command, because
a `v1.12` binary refuses it. Then restore the old binary and start it again:

```powershell
Remove-Item .\geth.exe
Rename-Item .\geth-v1.12.exe geth.exe
```

Keep the new key: rolling back does not make the old key safe.
