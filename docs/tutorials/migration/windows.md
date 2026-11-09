---
title: Migrating to v1.13 on Windows
description: "Upgrade an Ethereum Classic node on Windows from any v1.12.x release to the Core-Geth v1.13 line: the steps, the node-key rotation, and how to roll back."
---

# Windows users guide

This guide upgrades a node on Windows from any v1.12.x release to the v1.13 line.
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

Wait until the log prints `Blockchain stopped` and the process exits. If the process exits without that line, the
stop was not clean: the node repairs its database at its next start, and its block number can start lower.

## 2. Install v1.13

Keep the old binary: step 6 uses it to roll back.

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
download it again. On a node that holds account keys or that others rely on, also check
[where it came from](../../getting-started/installation.md#verify-where-the-archive-came-from), which needs the GitHub
CLI 2.49.0 or later, signed in, as [Check the GitHub CLI](../../getting-started/installation.md#check-the-github-cli)
shows:

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

Building needs Git, Go and a C compiler: [Dependencies](../../developers/build-from-source.md#dependencies) installs
and checks each one, with Go at the release `main` pins. This builds `main`, which takes fixes and dependency updates
before a release does. It is not a release, and it keeps moving: two builds a day apart can be different code.

In the PowerShell window you checked the tools in, one line at a time:

```powershell
cd $HOME
git clone --branch main --single-branch https://github.com/ethereumclassic/core-geth.git
cd core-geth
go run build/ci.go install ./cmd/geth
.\build\bin\geth.exe version
```

The last line must print a version ending in `-unstable`, such as `Core-Geth 1.13.1-unstable`, and a `Git Commit` line.
Note that commit: it identifies the code you run. Do not clone with `--recursive`: the submodules are several
gigabytes, and the build does not read them. If a `core-geth` folder is already there from an earlier build, run
`git -C core-geth pull` in place of the clone. [Build from source](../../developers/build-from-source.md#source) says
how changes reach `main`.

Then go back to the folder that holds the node's `geth.exe`, keep the old binary, and copy the build in its place:

```powershell
Rename-Item .\geth.exe geth-v1.12.exe
Copy-Item "$HOME\core-geth\build\bin\geth.exe" .\geth.exe
.\geth.exe version
```

Go on only if `.\geth.exe version` prints the same first line and `Git Commit` line the build printed. The second line
copies the build from `core-geth` in your home folder, where the build cloned it.

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
form. With no file yet, save these two lines and add `--config` and the file's path to the start command. The client
ignores `static-nodes.json`, so move its entries here.

Save the file, after this edit or any other, as UTF-8 with no byte order mark. A byte order mark, or UTF-16, stops the
start with `invalid TOML syntax`.

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags you start the
node with, and the config file if you use one. Check that file too for [the lines listed
here](../v1.13.0-migration.md#if-your-config-file-came-from-dumpconfig), whoever wrote it. If you added
`--discovery.dns` or `--bootnodes` on `v1.12`, [take them
out](../v1.13.0-migration.md#if-you-added-discoverydns-or-bootnodes-on-v112).

## 5. Start and verify

Start the node with the flags it runs with, your MESS setting among them: `--mess` or `--mess=false`, as [Decide your
MESS setting](../v1.13.0-migration.md#decide-your-mess-setting) explains. A node on `--datadir` started without that
flag starts an empty chain in the default directory. If the node starts from the Startup folder, as the run guide's
[Keep it running](../../getting-started/run/windows.md#6-keep-it-running-start-it-when-you-sign-in) sets up, run that
command file instead. Otherwise, in PowerShell, in the folder that holds the node's `geth.exe`, add those flags to this
line:

```powershell
.\geth.exe --classic
```

The node runs in this window until you stop it, so run the next steps in a second PowerShell window, in the same
folder. For Mordor, use `--mordor` in place of `--classic`.

Then, in this order:

```powershell
.\geth.exe attach --exec 'admin.nodeInfo.name'  \\.\pipe\geth.ipc
.\geth.exe attach --exec 'admin.nodeInfo.enode' \\.\pipe\geth.ipc
.\geth.exe attach --exec 'net.peerCount'        \\.\pipe\geth.ipc
.\geth.exe attach --exec 'eth.blockNumber'      \\.\pipe\geth.ipc
```

If a line cannot connect, wait a minute and run it again.

- **The first line** must name the version you installed: a `v1.13` release, or the `-unstable` version a build
  of `main` printed.
- **The enode ID** must differ from the one you recorded. If it is the same, the node is still on the old key: stop it
  and redo step 3.
- **The peer count** climbs within minutes. If it stays at zero, see
  [Why does the node have no peers?](../../operate/troubleshooting.md#why-does-the-node-have-no-peers)
- **The block number** starts near where the node stopped, lower after an unclean stop, then climbs. Compare it with
  a block explorer or another node you run. If it stops climbing short of the chain head, stop and investigate
  rather than resync.

### If other nodes list this node, give them the new ID

They can no longer connect to this node, and neither side logs it at the default verbosity. On each one, on `v1.12` or
`v1.13`, run these lines with the old ID you recorded, the new ID that step 5 printed, and the address and port the
entry already has:

```powershell
.\geth.exe attach --exec "admin.removePeer('enode://<old ID>@<IP>:30303')" \\.\pipe\geth.ipc
.\geth.exe attach --exec "admin.addPeer('enode://<new ID>@<IP>:30303')" \\.\pipe\geth.ipc
.\geth.exe attach --exec "admin.peers.map(function(p){return p.enode})" \\.\pipe\geth.ipc
```

The first two lines each print `true` and take effect with no restart. The third lists the new ID once the nodes
connect. The change lasts only until that node restarts, so put the new ID in that node's config file too. Send the new
ID to anyone else who lists this node.

For an entry in `TrustedNodes`, run `admin.removeTrustedPeer` and `admin.addTrustedPeer` the same way, as well as the
first two lines if that node also lists this one as a static peer. They do not open a connection.
[Static and trusted peers](../../guides/bootnodes-and-discovery.md#static-and-trusted-peers) explains the two lists.

## 6. If you need to roll back

Roll back only for as long as it takes to find what went wrong: every `v1.12` release carries the flaws the audits
found.

The data directory works in both directions. Stop the node, and take any `--mess` flag out of the start command, because
a `v1.12` binary refuses it. Then restore the old binary and start it again:

```powershell
Remove-Item .\geth.exe
Rename-Item .\geth-v1.12.exe geth.exe
```

Keep the new key: rolling back does not make the old key safe. A release before `v1.12.21` leaks bits of the key it runs
with, so rolling back to one exposes the new key too.

Once the cause is fixed, move again from step 1, renaming the key again in step 3. A second rename on the same day finds
its target name taken, so add `-2` after `$(Get-Date -Format yyyy-MM-dd)` in it. If the node has to stay on `v1.12`
longer, [If a v1.12 node cannot move yet](../../etc-cooperative-transition.md#if-a-v112-node-cannot-move-yet) keeps it
finding peers.

## 7. Once the node is stable, delete the v1.12 binary

Once the node has caught up, holds peers, and has run on the `v1.13` line long enough that you will not roll back, in
PowerShell, in the folder that holds the node's `geth.exe`:

```powershell
Remove-Item .\geth-v1.12.exe
```

Every `v1.12` release carries the flaws the audits found, so a copy left on the machine is one a later mistake can
start.
