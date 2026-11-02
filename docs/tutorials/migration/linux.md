---
description: "Upgrade an Ethereum Classic node on Linux from any v1.12.x release to the Core-Geth v1.13 line: the steps, the node-key rotation, and how to roll back."
---

# Linux users guide

This guide upgrades a node on Linux from any v1.12.x release to the latest v1.13 release.
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

Record what you are running, so you can tell afterwards whether anything changed that you did not intend:

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

Wait for this line before you continue. It confirms the database is consistent:

```
INFO [..] Blockchain stopped
```

[Stop it safely](../../getting-started/run-classic-node.md#stop-it-safely) says what each kind of stop does.

## 2. Install the latest release

Keep the old binary until the upgrade is verified; step 6 uses it. These commands assume `geth` is in
`/usr/local/bin`; `command -v geth` shows where yours is.

### Method A: release archive

Download the latest release and check it:

```bash
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
curl -LO $BASE/core-geth-linux-$VERSION.zip
curl -LO $BASE/core-geth-linux-$VERSION.zip.sha256
sha256sum -c core-geth-linux-$VERSION.zip.sha256
```

On an ARM board, put the name for your processor where the file name says `linux`, here and in the next two blocks:
`arm64` if `uname -m` prints `aarch64` or `arm64`, `arm7` for `armv7l`, `arm6` for `armv6l`, and `arm5` for `armv5tel`.

`sha256sum -c` must print `OK`. If it prints `FAILED`, delete the file and download it again rather than
running it. On a node that holds value, also check
[where it came from](../../getting-started/installation.md#verify-where-the-archive-came-from), which needs the
[GitHub CLI](https://cli.github.com/), signed in with `gh auth login`:

```bash
gh attestation verify core-geth-linux-$VERSION.zip --repo ethereumclassic/core-geth
```

Once signed in, it must print `Verification succeeded`. If it does not, do not run the file.

Then, in the same terminal, unpack it, keep the old binary and install the new one:

```bash
unzip -o core-geth-linux-$VERSION.zip &&
  { test -e /usr/local/bin/geth-v1.12 || sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12; } &&
  sudo install -m 0755 geth /usr/local/bin/geth &&
  geth version
```

Each line runs only if the one before it worked. Go on only if `geth version` reports a `v1.13` release.
`geth-v1.12` is the old binary, kept so you can roll back. A second run leaves that copy alone, so you can run this
block again after a mistake.

### Method B: build from source

Building needs Go 1.26 and a C compiler. Build the latest release's tag, not `main`, which keeps moving:

```bash
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
git clone --branch $VERSION --single-branch https://github.com/ethereumclassic/core-geth.git core-geth-$VERSION &&
  cd core-geth-$VERSION &&
  make geth &&
  ./build/bin/geth version
```

**Do not clone with `--recursive`.** The submodules hold consensus test fixtures of several gigabytes, and
`make geth` does not read them. [Build from source](../../developers/build-from-source.md) covers the
dependencies.

Then, in the folder the build left you in, keep the old binary and install the new one. The first line stops the
block anywhere else:

```bash
ls build/bin/geth &&
  { test -e /usr/local/bin/geth-v1.12 || sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12; } &&
  sudo install -m 0755 build/bin/geth /usr/local/bin/geth &&
  geth version
```

Go on only if `geth version` reports a `v1.13` release.

The executable is still named `geth`, so your units and scripts keep working unchanged.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why, and
which peers' configuration to update. With the node stopped:

```bash
mv -n ~/.ethereum/classic/geth/nodekey ~/.ethereum/classic/geth/nodekey.old-rotated-$(date +%F)
```

On the unit in the Linux run guide, the key belongs to `geth`, so run this instead:

```bash
sudo -u geth mv -n /var/lib/core-geth/geth/nodekey /var/lib/core-geth/geth/nodekey.old-rotated-$(date +%F)
```

This rename is the required rotation, and the next start writes a fresh key. If it prints an error, look in the
node's data directory before you go on. If `geth/nodekey.old-rotated-` with today's date is there, the rename has
already happened. If only `geth/nodekey` is there, rename that file, because a node started now would start on the
old key.

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags in
your unit file or start script.

## 5. Start and verify

Start the node the way you always do, with the flags you just checked. For a systemd unit, the reload picks up any
flag you changed in it:

```bash
sudo systemctl daemon-reload && sudo systemctl start core-geth
```

Then, in this order:

```bash
geth attach --exec 'admin.nodeInfo.name'  ~/.ethereum/classic/geth.ipc
geth attach --exec 'admin.nodeInfo.enode' ~/.ethereum/classic/geth.ipc
geth attach --exec 'net.peerCount'        ~/.ethereum/classic/geth.ipc
geth attach --exec 'eth.blockNumber'      ~/.ethereum/classic/geth.ipc
```

If a line cannot connect, give the node a minute to open its socket and run that line again.

- **The first line** must name a `v1.13` release. It reads the running node, where `geth version` reads the binary
  on your `PATH`, which need not be the one a unit starts.
- **The enode ID** must differ from the one you recorded. If it is the same, the node is still on the old key: stop
  it cleanly and rename `geth/nodekey` in the data directory the node really uses. If it differs, it is the new ID
  to send to your peers.
- **The peer count** climbs within minutes. **A count stuck at zero after ten minutes usually means a peer still
  lists your old enode ID.**
- **The block number** resumes where the node stopped. Last, compare it with a block explorer or another node you
  operate. If it settles somewhere else, stop and investigate rather than resyncing.

## 6. If you need to roll back

The data directory works in both directions, so rolling back swaps the binary. Stop the node. If you added a
`--mess` flag, take it out of the start command first: a `v1.12` binary stops with
`flag provided but not defined`. Under systemd, run `sudo systemctl daemon-reload` after you edit the unit. Then:

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
