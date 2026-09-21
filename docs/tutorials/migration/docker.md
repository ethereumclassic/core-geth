# Docker users guide

This guide upgrades a node running in Docker from any v1.12.x image to v1.13.0.
[Migrating to v1.13.0](../v1.13.0-migration.md) explains why, and what changes for every node.

The image moves from `etclabscore/core-geth` to images built from `ethereumclassic/core-geth`. Your volume,
ports and flags stay the same. The commands assume a container named `core-geth` with the host directory
`$LOCAL_DATADIR` mounted at `/root`, as [Installation, Run the container](../../getting-started/installation.md#run-the-container)
shows.

## Choose how you get the image

| Method | Image name it gives you |
|---|---|
| [Pull from the registry](../../getting-started/installation.md#docker) | `ghcr.io/ethereumclassic/core-geth:v1.13.0` |
| [Load the release tarball](../../getting-started/installation.md#load-the-image-from-the-release) | `ghcr.io/ethereumclassic/core-geth:v1.13.0` |
| [Build an image from the release archive](../../getting-started/installation.md#build-an-image-from-the-release-archive) | the tag you give it, for example `core-geth:v1.13.0` |
| [Build an image from source](../../getting-started/installation.md#build-an-image-from-source) | the tag you give it, for example `core-geth:local` |

The pull needs no account and no token. The other three methods are there for a machine with no
route to the registry, or if you would rather build what you run.

## Before you start

Record what you are running:

```bash
docker exec core-geth geth version
docker exec core-geth geth attach --exec 'eth.blockNumber'      /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec 'admin.nodeInfo.enode' /root/.ethereum/classic/geth.ipc
```

**Keep that last line.** It is your current enode ID, and step 3 changes it.

## 1. Stop the node cleanly

```bash
docker stop -t 300 core-geth
docker logs core-geth 2>&1 | grep 'Blockchain stopped'
```

**`-t 300` is not optional.** `docker stop` sends SIGTERM, waits 10 seconds, then kills the process. A clean
shutdown writes cached state to disk first and takes longer than that on a real node, so a plain
`docker stop` kills it mid-write, which is the one thing in this upgrade that costs a resync. Continue only
once `Blockchain stopped` is in the log.

## 2. Get the v1.13.0 image

Use one of the methods above, and set `IMAGE` to the name it gives you:

```bash
IMAGE=ghcr.io/ethereumclassic/core-geth:v1.13.0
```

Keep your previous image; step 6 uses it.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why, and
which peer lists to update. The key lives in the mounted volume, so rotate it on the host. The container
wrote it as root, so this may need `sudo`:

```bash
mv "$LOCAL_DATADIR/.ethereum/classic/geth/nodekey" \
   "$LOCAL_DATADIR/.ethereum/classic/geth/nodekey.old-rotated-$(date +%F)"
```

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure). In Docker, publish the
RPC port on `127.0.0.1` only, as the warning under
[Run the container](../../getting-started/installation.md#run-the-container) explains.

## 5. Start and verify

```bash
docker rm core-geth
docker run -d --name core-geth \
  -v $LOCAL_DATADIR:/root \
  -p 30303:30303 -p 30303:30303/udp -p 127.0.0.1:8545:8545 \
  $IMAGE \
  --classic --http --http.addr 0.0.0.0 --http.port 8545
```

Then, in this order:

```bash
docker exec core-geth geth version                                                     # Core-Geth 1.13.0
docker exec core-geth geth attach --exec 'admin.nodeInfo.enode' /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec 'net.peerCount'        /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec 'eth.blockNumber'      /root/.ethereum/classic/geth.ipc
```

Send the new enode ID to your peers. **A peer count stuck at zero after ten minutes usually means a peer
still lists your old enode ID.** Last, compare `eth.blockNumber` with a block explorer or another node you
operate.

## 6. If you need to roll back

The data directory works in both directions, so rolling back swaps the image:

```bash
docker stop -t 300 core-geth
docker rm core-geth
# run your previous command again, with your etclabscore/core-geth v1.12.x image
```

The rotated node key stays rotated: rolling back does not un-expose the old key.
