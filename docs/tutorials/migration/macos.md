---
title: Migrating to v1.13 on macOS
description: "Upgrade an Ethereum Classic node on macOS from any v1.12.x release to the Core-Geth v1.13 line: the steps, the node-key rotation, and how to roll back."
---

# Mac users guide

This guide upgrades a node on macOS from any v1.12.x release to the v1.13 line.
[Migrating to v1.13](../v1.13.0-migration.md) explains why, and what changes for every node.

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
- **Under launchd:** unload the job. [Stop it cleanly](../../getting-started/run/macos.md#7-stop-it-cleanly) has the
  command for the run guide's agent.

Wait until the log prints `Blockchain stopped` and the process exits. If the process exits without that line, the
stop was not clean: the node repairs its database at its next start, and its block number can start lower.

## 2. Install v1.13

Keep the old binary: step 6 uses it to roll back. These commands assume `geth` is in `/usr/local/bin`;
`command -v geth` shows where yours is, and if it is elsewhere, put that path in place of `/usr/local/bin/geth` below.

### Method A: release archive

Download the latest release and check it:

```bash
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
curl -LO $BASE/core-geth-osx-arm64-$VERSION.zip
curl -LO $BASE/core-geth-osx-arm64-$VERSION.zip.sha256
shasum -a 256 -c core-geth-osx-arm64-$VERSION.zip.sha256
```

On an Intel Mac, put `osx` where the file name says `osx-arm64`, here and in the next two blocks. If you download with a
browser instead, [clear the quarantine attribute](../../getting-started/installation.md#macos).

`shasum -a 256 -c` must print `OK`. If it prints `FAILED`, delete the file and download it again. On a node that holds
account keys or that others rely on, also check [where it came
from](../../getting-started/installation.md#verify-where-the-archive-came-from), which needs the GitHub CLI 2.49.0 or
later, signed in, as [Check the GitHub CLI](../../getting-started/installation.md#check-the-github-cli) shows:

```bash
gh attestation verify core-geth-osx-arm64-$VERSION.zip --repo ethereumclassic/core-geth
```

It must print `Verification succeeded`. If it does not, do not run the file.

Then, in the same terminal, unpack it, keep the old binary and install the new one:

```bash
unzip -o core-geth-osx-arm64-$VERSION.zip &&
  { test -e /usr/local/bin/geth-v1.12 || sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12; } &&
  sudo install -m 0755 geth /usr/local/bin/geth &&
  geth version
```

Each line runs only if the one before it worked. Go on only if `geth version` reports a `v1.13` release.

### Method B: build from source

Building needs the Xcode Command Line Tools and Go: [Dependencies](../../developers/build-from-source.md#dependencies)
installs and checks both, with Go at the release `main` pins. This builds `main`, which takes fixes and dependency
updates before a release does. It is not a release, and it keeps moving: two builds a day apart can be different code.

```bash
git clone --branch main --single-branch https://github.com/ethereumclassic/core-geth.git &&
  cd core-geth &&
  make geth &&
  ./build/bin/geth version
```

Do not clone with `--recursive`: the submodules are several gigabytes, and `make geth` does not read them. If a
`core-geth` folder is already there from an earlier build, run `git -C core-geth pull` in place of the clone.
[Build from source](../../developers/build-from-source.md#source) says how changes reach `main`.

Then, in that folder, keep the old binary and install the new one:

```bash
ls build/bin/geth &&
  { test -e /usr/local/bin/geth-v1.12 || sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12; } &&
  sudo install -m 0755 build/bin/geth /usr/local/bin/geth &&
  geth version
```

Go on only if `geth version` reports a version ending in `-unstable`, such as `Core-Geth 1.13.1-unstable`, and a
`Git Commit` line. Note that commit: it identifies the code you run.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why. With the node
stopped:

```bash
mv -n ~/Library/Ethereum/classic/geth/nodekey ~/Library/Ethereum/classic/geth/nodekey.old-rotated-$(date +%F)
```

The next start writes a fresh key. If you start the node with `--nodekey` or `--nodekeyhex`, or set `GETH_NODEKEY` or
`GETH_NODEKEYHEX`, take that out too: it overrides the file.

If the rename prints an error, check the data directory: `geth/nodekey` must not be there when the node starts.

### If this node lists static peers, update the list

They are under `[Node.P2P]` in the file you pass with `--config`. Edit it while the node is down:

```toml
[Node.P2P]
StaticNodes = ["enode://<node ID>@<IP>:30303"]
```

A peer that has already moved has a new ID, so use that. Separate entries with commas. `TrustedNodes` takes the same
form. With no file yet, save these two lines and add `--config` and the file's path to the start command. The client
ignores `static-nodes.json`, so move its entries here.

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags you start the
node with, and the config file if you use one. Check that file too for [the lines listed
here](../v1.13.0-migration.md#if-your-config-file-came-from-dumpconfig), whoever wrote it. If you added
`--discovery.dns` or `--bootnodes` on `v1.12`, [take them
out](../v1.13.0-migration.md#if-you-added-discoverydns-or-bootnodes-on-v112).

## 5. Start and verify

Start the node with the flags it runs with, your MESS setting among them: `--mess` or `--mess=false`, as [Decide your
MESS setting](../v1.13.0-migration.md#decide-your-mess-setting) explains. A node on `--datadir` started without that
flag starts an empty chain in the default directory. If the node runs under the run guide's launchd agent, start the
agent instead, as [Keep it running](../../getting-started/run/macos.md#6-keep-it-running-a-launchd-agent) shows.
Otherwise, in a terminal, add those flags to this line:

```bash
geth --classic
```

The node runs in this terminal until you stop it, so run the next steps in a second terminal. For Mordor, use
`--mordor` in place of `--classic`.

Then, in this order:

```bash
geth attach --exec 'admin.nodeInfo.name'  ~/Library/Ethereum/classic/geth.ipc
geth attach --exec 'admin.nodeInfo.enode' ~/Library/Ethereum/classic/geth.ipc
geth attach --exec 'net.peerCount'        ~/Library/Ethereum/classic/geth.ipc
geth attach --exec 'eth.blockNumber'      ~/Library/Ethereum/classic/geth.ipc
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

```bash
geth attach --exec "admin.removePeer('enode://<old ID>@<IP>:30303')" ~/Library/Ethereum/classic/geth.ipc
geth attach --exec "admin.addPeer('enode://<new ID>@<IP>:30303')" ~/Library/Ethereum/classic/geth.ipc
geth attach --exec "admin.peers.map(function(p){return p.enode})" ~/Library/Ethereum/classic/geth.ipc
```

The first two lines each print `true` and take effect with no restart. The third lists the new ID once the nodes
connect. The change lasts only until that node restarts, so put the new ID in that node's config file too. Send the new
ID to anyone else who lists this node.

For an entry in `TrustedNodes`, run `admin.removeTrustedPeer` and `admin.addTrustedPeer` the same way, as well as the
first two lines if that node also lists this one as a static peer. They do not open a connection.
[Static and trusted peers](../../guides/bootnodes-and-discovery.md#static-and-trusted-peers) explains the two lists.

## 6. If you need to roll back

The data directory works in both directions, so rolling back swaps the binary. Stop the node, and take any `--mess` flag
out of the start command, because a `v1.12` binary refuses it. Then:

```bash
sudo mv /usr/local/bin/geth-v1.12 /usr/local/bin/geth
```

Start it again. Keep the new key: rolling back does not make the old key safe. A release before `v1.12.21` leaks bits of
the key it runs with, so rolling back to one exposes the new key too.
