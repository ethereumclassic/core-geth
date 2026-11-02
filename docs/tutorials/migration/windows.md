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

PowerShell has no equivalent of `sha256sum -c`, so compare the last two lines yourself: `Get-FileHash` prints
uppercase hex and the `.sha256` file holds lowercase followed by the file name. The letters must match; the case
does not. If they differ, delete the file and download it again rather than running it. On a node that holds
value, also check [where it came from](../../getting-started/installation.md#verify-where-the-archive-came-from),
which needs the [GitHub CLI](https://cli.github.com/), signed in with `gh auth login`:

```powershell
gh attestation verify "core-geth-win64-$VERSION.zip" --repo ethereumclassic/core-geth
```

Once signed in, it must print `Verification succeeded`. If it does not, do not run the file.

Then, in the same window, keep the old binary and unpack the new one in its place:

```powershell
Rename-Item .\geth.exe geth-v1.12.exe
Expand-Archive "core-geth-win64-$VERSION.zip" -DestinationPath . -Force
.\geth.exe version
```

Go on only if `.\geth.exe version` reports a `v1.13` release. `geth-v1.12.exe` is the old binary, kept so you can
roll back.

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

That is the command `make geth` runs on other systems. Rename your current `geth.exe` to
`geth-v1.12.exe`, then copy `build\bin\geth.exe` in its place. **Do not clone with `--recursive`:** the
submodules hold test fixtures of several gigabytes that the build does not read.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why, and
which peers' configuration to update. With the node stopped:

```powershell
$key = "$env:LOCALAPPDATA\Ethereum\classic\geth\nodekey"
Rename-Item $key "nodekey.old-rotated-$(Get-Date -Format yyyy-MM-dd)"
```

This rename is the required rotation, and the next start writes a fresh key. If it prints an error, look in the
node's data directory before you go on. If `geth/nodekey.old-rotated-` with today's date is there, the rename has
already happened. If only `geth/nodekey` is there, rename that file, because a node started now would start on the
old key.

### If this node lists static peers, update the list

Static peers are nodes a node keeps connecting to, such as the rest of a mining pool's fleet. They sit under
`[Node.P2P]` in the file you pass with `--config`, so edit that file while the node is down:

```toml
[Node.P2P]
StaticNodes = ["enode://<node ID>@<IP>:30303"]
```

A node that has already moved has a new ID, so put the new ID in its entry and keep the address and port. Separate
entries with commas. `TrustedNodes`, on a line of its own in the same section, takes the same entries. With no file
yet, save these two lines with your own values, then add `--config` and the file's path to the start command. The
client ignores `static-nodes.json` and logs an error at the start, so move any entries from that file into this one.
To change the list on a running node instead, use the lines
[at the end of step 5](#if-other-nodes-list-this-node-give-them-the-new-id).

Save the file as UTF-8 with no byte order mark. A byte order mark, or UTF-16, stops the start with
`Fatal: <file>, line 1: invalid TOML syntax`.

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags you
start the node with. If you pass `--config` a file that
`dumpconfig` wrote on `v1.12`, the checklist also lists the lines to delete.

## 5. Start and verify

Start the node the way you always do, with the flags you just checked. Then, in this order:

```powershell
.\geth.exe attach --exec 'admin.nodeInfo.name'  \\.\pipe\geth.ipc
.\geth.exe attach --exec 'admin.nodeInfo.enode' \\.\pipe\geth.ipc
.\geth.exe attach --exec 'net.peerCount'        \\.\pipe\geth.ipc
.\geth.exe attach --exec 'eth.blockNumber'      \\.\pipe\geth.ipc
```

If a line cannot connect, give the node a minute to open its socket and run that line again.

- **The first line** must name a `v1.13` release. It reads the running node, where `.\geth.exe version` reads the file in
  this folder, which need not be the one a service starts.
- **The enode ID** must differ from the one you recorded. If it is the same, the node is still on the old key: stop
  it cleanly and rename `geth/nodekey` in the data directory the node really uses. If it differs, it is the new ID
  to send to your peers.
- **The peer count** climbs within minutes. **A count stuck at zero after ten minutes usually means a peer still
  lists your old enode ID.**
- **The block number** resumes where the node stopped. Last, compare it with a block explorer or another node you
  operate. If it settles somewhere else, stop and investigate rather than resyncing.

### If other nodes list this node, give them the new ID

The new ID exists only now that this node has started. A node that still names the old ID no longer connects to this
node, and neither node logs the failure at the default verbosity. On each node that lists this one, on `v1.12` or
`v1.13`, run these lines with your own values: the old ID you recorded, the new ID from the block above, and the
address and port the entry already has.

```powershell
.\geth.exe attach --exec "admin.removePeer('enode://<old ID>@<IP>:30303')" \\.\pipe\geth.ipc
.\geth.exe attach --exec "admin.addPeer('enode://<new ID>@<IP>:30303')" \\.\pipe\geth.ipc
.\geth.exe attach --exec "admin.peers.map(function(p){return p.enode})" \\.\pipe\geth.ipc
```

The first two lines each print `true` and take effect with no restart, and the third lists the new ID once the two
nodes connect. A running node keeps this change only until it restarts, so put the new ID in that node's configuration
file too. Send the new ID to anyone else who lists this node.

For an entry in `TrustedNodes`, use `admin.removeTrustedPeer` and `admin.addTrustedPeer`, which also print `true`.
They do not open a connection, so a node given only those stays unconnected until the other node dials it.
[Static and trusted peers](../../guides/bootnodes-and-discovery.md#static-and-trusted-peers) explains the two lists.

## 6. If you need to roll back

The data directory works in both directions. Stop the node. If you added a `--mess` flag, take it out of the start
command first: a `v1.12` binary stops with `flag provided but not defined`. Then restore the old binary and start it
again:

```powershell
Remove-Item .\geth.exe
Rename-Item .\geth-v1.12.exe geth.exe
```

The rotated node key stays rotated: rolling back does not un-expose the old key.
