---
title: Migrating to v1.13 on Linux
description: "Upgrade an Ethereum Classic node on Linux from any v1.12.x release to the Core-Geth v1.13 line: the steps, the node-key rotation, and how to roll back."
---

# Linux users guide

This guide upgrades a node on Linux from any v1.12.x release to the v1.13 line.
[Migrating to v1.13](../v1.13.0-migration.md) explains why, and what changes for every node.

The commands use the default Ethereum Classic data directory, `~/.ethereum/classic`. If you pass
`--datadir`, use yours; for Mordor, use `~/.ethereum/mordor`.

## Choose how you update

| Method | Use it if |
|---|---|
| [Release archive](#method-a-release-archive) | You run a published binary: `core-geth-linux` on x86-64, `core-geth-arm64` on 64-bit ARM, or `core-geth-arm5`, `arm6` or `arm7` on 32-bit ARM |
| [Build from source](#method-b-build-from-source) | You compile the client yourself |

Both methods share every other step.

## Before you start

Record what you are running:

```bash
geth version
geth attach --exec 'eth.blockNumber'      ~/.ethereum/classic/geth.ipc
geth attach --exec 'net.peerCount'        ~/.ethereum/classic/geth.ipc
geth attach --exec 'admin.nodeInfo.enode' ~/.ethereum/classic/geth.ipc
```

**Keep that last line.** It is your current enode ID, and step 3 changes it.

On the unit in the [Linux run guide](../../getting-started/run/linux.md#6-keep-it-running-a-systemd-service), the
`geth` user owns `/var/lib/core-geth`, so attach as `geth`, here and in every later attach line:

```bash
sudo -u geth -H geth --classic attach --exec 'admin.nodeInfo.enode' /var/lib/core-geth/geth.ipc
```

## 1. Stop the node cleanly

- **In a terminal:** press Ctrl-C once and wait. A second interrupt does not kill it.
- **Under systemd:** stop the unit you wrote, for example `sudo systemctl stop core-geth`.

Wait until the log prints `Blockchain stopped` and the process exits. If the process exits without that line, the
stop was not clean: the node repairs its database at its next start, and its block number can start lower.
[Stop it safely](../../getting-started/run-classic-node.md#stop-it-safely) has the details.

## 2. Install v1.13

Keep the old binary: step 6 uses it to roll back. These commands assume `geth` is in `/usr/local/bin`;
`command -v geth` shows where yours is, and if it is elsewhere, put that path in place of `/usr/local/bin/geth` below.

### Method A: release archive

Install `curl` and `unzip`, which a minimal server can lack. On Debian or Ubuntu:

```bash
sudo apt-get update && sudo apt-get install -y curl unzip
```

On Fedora or RHEL, run `sudo dnf install -y --refresh curl unzip && sudo dnf upgrade -y curl unzip` instead, and on
Arch, `sudo pacman -Syu --needed curl unzip`. Each line also updates curl and unzip if they are already installed.

Then download the latest release and check it:

```bash
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
curl -LO $BASE/core-geth-linux-$VERSION.zip
curl -LO $BASE/core-geth-linux-$VERSION.zip.sha256
sha256sum -c core-geth-linux-$VERSION.zip.sha256
```

On an ARM machine, put the name for your processor where the file name says `linux`, here and in the next two blocks:
`arm64` if `uname -m` prints `aarch64` or `arm64`, `arm7` for `armv7l`, `arm6` for `armv6l`, and `arm5` for `armv5tel`.

`sha256sum -c` must print `OK`. If it prints `FAILED`, delete the file and download it again. On a node that holds
account keys or that others rely on, also check [where it came
from](../../getting-started/installation.md#verify-where-the-archive-came-from), which needs the GitHub CLI 2.49.0 or
later, signed in, as [Check the GitHub CLI](../../getting-started/installation.md#check-the-github-cli) shows:

```bash
gh attestation verify core-geth-linux-$VERSION.zip --repo ethereumclassic/core-geth
```

It must print `Verification succeeded`. If it does not, do not run the file.

Then, in the same terminal, unpack it, keep the old binary and install the new one:

```bash
unzip -o core-geth-linux-$VERSION.zip &&
  { test -e /usr/local/bin/geth-v1.12 || sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12; } &&
  sudo install -m 0755 geth /usr/local/bin/geth &&
  geth version
```

Each line runs only if the one before it worked. Go on only if `geth version` reports a `v1.13` release.

### Method B: build from source

Building needs Git, Go, a C compiler and `make`: [Dependencies](../../developers/build-from-source.md#dependencies)
installs and checks each one, with Go at the release `main` pins. This builds `main`, which takes fixes and dependency
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
`Git Commit` line. Note that commit: it identifies the code you run. The executable is still named `geth`, so your units
and scripts keep working.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why. With the node
stopped:

```bash
mv -n ~/.ethereum/classic/geth/nodekey ~/.ethereum/classic/geth/nodekey.old-rotated-$(date +%F)
```

On the unit in the Linux run guide, the key belongs to `geth`, so run this instead:

```bash
sudo -u geth mv -n /var/lib/core-geth/geth/nodekey /var/lib/core-geth/geth/nodekey.old-rotated-$(date +%F)
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

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags in your unit
file or start script, and the config file if you use one. Check that file too for [the lines listed
here](../v1.13.0-migration.md#if-your-config-file-came-from-dumpconfig), whoever wrote it. If you added
`--discovery.dns` or `--bootnodes` on `v1.12`, [take them
out](../v1.13.0-migration.md#if-you-added-discoverydns-or-bootnodes-on-v112).

## 5. Start and verify

Start the node with the flags it runs with, your MESS setting among them: `--mess` or `--mess=false`, as [Decide your
MESS setting](../v1.13.0-migration.md#decide-your-mess-setting) explains. A node on `--datadir` started without that
flag starts an empty chain in the default directory. Under the systemd unit in the [Linux run
guide](../../getting-started/run/linux.md#6-keep-it-running-a-systemd-service), start the unit. `core-geth` stands for
the unit you wrote, and the reload picks up any flag you changed in it:

```bash
sudo systemctl daemon-reload && sudo systemctl start core-geth
```

Otherwise, in a terminal, add those flags to this line:

```bash
geth --classic
```

The node runs in this terminal until you stop it, so run the next steps in a second terminal. For Mordor, use
`--mordor` in place of `--classic`.

Then, in this order:

```bash
geth attach --exec 'admin.nodeInfo.name'  ~/.ethereum/classic/geth.ipc
geth attach --exec 'admin.nodeInfo.enode' ~/.ethereum/classic/geth.ipc
geth attach --exec 'net.peerCount'        ~/.ethereum/classic/geth.ipc
geth attach --exec 'eth.blockNumber'      ~/.ethereum/classic/geth.ipc
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
geth attach --exec "admin.removePeer('enode://<old ID>@<IP>:30303')" ~/.ethereum/classic/geth.ipc
geth attach --exec "admin.addPeer('enode://<new ID>@<IP>:30303')" ~/.ethereum/classic/geth.ipc
geth attach --exec "admin.peers.map(function(p){return p.enode})" ~/.ethereum/classic/geth.ipc
```

The first two lines each print `true` and take effect with no restart. The third lists the new ID once the nodes
connect. The change lasts only until that node restarts, so put the new ID in that node's config file too. Send the new
ID to anyone else who lists this node.

For an entry in `TrustedNodes`, run `admin.removeTrustedPeer` and `admin.addTrustedPeer` the same way, as well as the
first two lines if that node also lists this one as a static peer. They do not open a connection.
[Static and trusted peers](../../guides/bootnodes-and-discovery.md#static-and-trusted-peers) explains the two lists.

## 6. If you need to roll back

The data directory works in both directions, so rolling back swaps the binary. Stop the node, and take any `--mess` flag
out of the start command, because a `v1.12` binary refuses it. Under systemd, run `sudo systemctl daemon-reload`. Then:

```bash
sudo mv /usr/local/bin/geth-v1.12 /usr/local/bin/geth
```

Start it again. Keep the new key: rolling back does not make the old key safe. A release before `v1.12.21` leaks bits of
the key it runs with, so rolling back to one exposes the new key too.
