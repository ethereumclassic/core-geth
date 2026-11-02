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

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure) against the flags you
start the node with. If you pass `--config` a file that
`dumpconfig` wrote on `v1.12`, the checklist also lists the lines to delete.

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

### If other nodes list this node, give them the new ID

The new ID exists only now that this node has started. A node that still names the old ID no longer connects to this
node, and neither node logs the failure at the default verbosity. On each node that lists this one, on `v1.12` or
`v1.13`, run these lines with your own values: the old ID you recorded, the new ID from the block above, and the
address and port the entry already has.

```bash
geth attach --exec "admin.removePeer('enode://<old ID>@<IP>:30303')" ~/Library/Ethereum/classic/geth.ipc
geth attach --exec "admin.addPeer('enode://<new ID>@<IP>:30303')" ~/Library/Ethereum/classic/geth.ipc
geth attach --exec "admin.peers.map(function(p){return p.enode})" ~/Library/Ethereum/classic/geth.ipc
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

```bash
sudo mv /usr/local/bin/geth-v1.12 /usr/local/bin/geth
```

The rotated node key stays rotated: rolling back does not un-expose the old key.
