---
description: "Upgrade an Ethereum Classic node on macOS from any v1.12.x release to the Core-Geth v1.13 line: the steps, the node-key rotation, and how to roll back."
---

# Mac users guide

This guide upgrades a node on macOS from any v1.12.x release to the latest v1.13 release.
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
- **Under launchd:** stop the job you created.

Wait for `INFO [..] Blockchain stopped` before you continue. It confirms the database is consistent.

## 2. Install the latest release

Keep the old binary until the upgrade is verified; step 6 uses it. These commands assume `geth` is in
`/usr/local/bin`; `command -v geth` shows where yours is.

### Method A: release archive

Download the latest release and check it. `curl` does not set the quarantine attribute that makes macOS block a
binary; if you downloaded with a browser instead, [clear it](../../getting-started/installation.md#macos):

```bash
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
curl -LO $BASE/core-geth-osx-arm64-$VERSION.zip
curl -LO $BASE/core-geth-osx-arm64-$VERSION.zip.sha256
shasum -a 256 -c core-geth-osx-arm64-$VERSION.zip.sha256
```

On an Intel Mac, put `osx` where the file name says `osx-arm64`, here and in the next two blocks.

`shasum -a 256 -c` must print `OK`. If it prints `FAILED`, delete the file and download it again rather than
running it. On a node that holds value, also check
[where it came from](../../getting-started/installation.md#verify-where-the-archive-came-from), which needs the
[GitHub CLI](https://cli.github.com/), signed in with `gh auth login`:

```bash
gh attestation verify core-geth-osx-arm64-$VERSION.zip --repo ethereumclassic/core-geth
```

Once signed in, it must print `Verification succeeded`. If it does not, do not run the file.

Then, in the same terminal, unpack it, keep the old binary and install the new one:

```bash
unzip -o core-geth-osx-arm64-$VERSION.zip &&
  { test -e /usr/local/bin/geth-v1.12 || sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12; } &&
  sudo install -m 0755 geth /usr/local/bin/geth &&
  geth version
```

Each line runs only if the one before it worked. Go on only if `geth version` reports a `v1.13` release.
`geth-v1.12` is the old binary, kept so you can roll back. A second run leaves that copy alone, so you can run this
block again after a mistake.

### Method B: build from source

Building needs Go 1.26 and the Xcode Command Line Tools (`xcode-select --install`), which provide the C
compiler and `make`. Build the latest release's tag, not `main`, which keeps moving:

```bash
VERSION=$(curl -sLo /dev/null -w '%{url_effective}' https://github.com/ethereumclassic/core-geth/releases/latest | sed 's|.*/||')
git clone --branch $VERSION --single-branch https://github.com/ethereumclassic/core-geth.git core-geth-$VERSION &&
  cd core-geth-$VERSION &&
  make geth &&
  ./build/bin/geth version
```

**Do not clone with `--recursive`.** The submodules hold test fixtures of several gigabytes that
`make geth` does not read.

Then, in the folder the build left you in, keep the old binary and install the new one. The first line stops the
block anywhere else:

```bash
ls build/bin/geth &&
  { test -e /usr/local/bin/geth-v1.12 || sudo mv /usr/local/bin/geth /usr/local/bin/geth-v1.12; } &&
  sudo install -m 0755 build/bin/geth /usr/local/bin/geth &&
  geth version
```

Go on only if `geth version` reports a `v1.13` release.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why, and
which peers' configuration to update. With the node stopped:

```bash
mv -n ~/Library/Ethereum/classic/geth/nodekey ~/Library/Ethereum/classic/geth/nodekey.old-rotated-$(date +%F)
```

This rename is the required rotation, and the next start writes a fresh key. If it prints an error, look in the
node's data directory before you go on. If `geth/nodekey.old-rotated-` with today's date is there, the rename has
already happened. If only `geth/nodekey` is there, rename that file, because a node started now would start on the
old key.

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags you
start the node with.

## 5. Start and verify

Start the node the way you always do, with the flags you just checked. Then, in this order:

```bash
geth attach --exec 'admin.nodeInfo.name'  ~/Library/Ethereum/classic/geth.ipc
geth attach --exec 'admin.nodeInfo.enode' ~/Library/Ethereum/classic/geth.ipc
geth attach --exec 'net.peerCount'        ~/Library/Ethereum/classic/geth.ipc
geth attach --exec 'eth.blockNumber'      ~/Library/Ethereum/classic/geth.ipc
```

If a line cannot connect, give the node a minute to open its socket and run that line again.

- **The first line** must name a `v1.13` release. It reads the running node, where `geth version` reads the binary
  on your `PATH`, which need not be the one a launchd job starts.
- **The enode ID** must differ from the one you recorded. If it is the same, the node is still on the old key: stop
  it cleanly and rename `geth/nodekey` in the data directory the node really uses. If it differs, it is the new ID
  to send to your peers.
- **The peer count** climbs within minutes. **A count stuck at zero after ten minutes usually means a peer still
  lists your old enode ID.**
- **The block number** resumes where the node stopped. Last, compare it with a block explorer or another node you
  operate. If it settles somewhere else, stop and investigate rather than resyncing.

## 6. If you need to roll back

The data directory works in both directions. Stop the node. If you added a `--mess` flag, take it out of the start
command first: a `v1.12` binary stops with `flag provided but not defined`. Then restore the old binary and start it
again:

```bash
sudo mv /usr/local/bin/geth-v1.12 /usr/local/bin/geth
```

The rotated node key stays rotated: rolling back does not un-expose the old key.
