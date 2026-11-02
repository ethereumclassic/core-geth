---
description: "Upgrade an Ethereum Classic node in Docker from any v1.12.x image to the Core-Geth v1.13 line: the steps, the node-key rotation, and how to roll back."
---

# Docker users guide

This guide upgrades a node running in Docker from any v1.12.x image to the latest v1.13 release.
[Migrating to v1.13](../v1.13.0-migration.md) explains why, and what changes for every node.

The image moves from `etclabscore/core-geth` to images built from `ethereumclassic/core-geth`. Your volume,
ports and flags stay the same. The commands assume a container named `core-geth` with its data mounted at `/root`,
from a host directory as [Installation, Run the container](../../getting-started/installation.md#run-the-container)
shows or from a named volume as the [Docker run guide](../../getting-started/run/docker.md) shows. Put your own
container name in place of `core-geth`. For Mordor, use `.ethereum/mordor` in the paths.

## Choose how you get the image

| Method | Image name it gives you |
|---|---|
| [Pull from the registry](../../getting-started/installation.md#docker) | `ghcr.io/ethereumclassic/core-geth:latest` |
| [Load the release tarball](../../getting-started/installation.md#load-the-image-from-the-release) | `ghcr.io/ethereumclassic/core-geth:<tag>` |
| [Build an image from the release archive](../../getting-started/installation.md#build-an-image-from-the-release-archive) | the tag you give it, for example `core-geth:local` |
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
docker logs --tail 50 core-geth 2>&1 | grep 'Blockchain stopped'
```

**`-t 300` is not optional.** `docker stop` sends SIGTERM, waits 10 seconds, then kills the process. A clean
shutdown writes cached state to disk first and takes longer than that on a real node, so a plain
`docker stop` kills it mid-write, which is the one thing in this upgrade that costs a resync. Continue only
once `Blockchain stopped` is in the log.

## 2. Get the latest image

This pulls the latest release from the registry, and sets `IMAGE` for the steps below:

```bash
IMAGE=ghcr.io/ethereumclassic/core-geth:latest
docker pull $IMAGE
```

With another of the methods above, set `IMAGE` to the name it gives you. On a node that holds value, also check
where a pulled image came from, which needs the [GitHub CLI](https://cli.github.com/), signed in with
`gh auth login`:

```bash
gh attestation verify oci://$IMAGE --repo ethereumclassic/core-geth
```

Once signed in, it must print `Verification succeeded`. If it does not, do not run the image.

Keep your previous image; step 6 uses it.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why, and
which peers' configuration to update. This runs in a throwaway container from the new image, so it reaches the key
however the old container mounts the data directory, with no `sudo`:

```bash
docker run --rm --volumes-from core-geth --entrypoint sh $IMAGE \
  -c "mv -n /root/.ethereum/classic/geth/nodekey /root/.ethereum/classic/geth/nodekey.old-rotated-$(date +%F)"
```

This rename is the required rotation, and the next start writes a fresh key. If it prints an error, look in the
node's data directory before you go on. If `geth/nodekey.old-rotated-` with today's date is there, the rename has
already happened. If only `geth/nodekey` is there, rename that file, because a node started now would start on the
old key.

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure). In Docker, publish the
RPC port on `127.0.0.1` only, as the warning under
[Run the container](../../getting-started/installation.md#run-the-container) explains.

## 5. Start and verify

Keep the old container, stopped, and start a new one from the new image:

```bash
docker rename core-geth core-geth-v1.12
docker update --restart no core-geth-v1.12
docker run -d --name core-geth \
  --restart unless-stopped --stop-timeout 300 \
  --volumes-from core-geth-v1.12 \
  -p 30303:30303 -p 30303:30303/udp -p 127.0.0.1:8545:8545 \
  $IMAGE \
  --classic --http --http.addr 0.0.0.0 --http.port 8545
```

Put your own ports and flags in the last lines of `docker run`: `docker inspect core-geth-v1.12` shows what the old
container used. For Mordor, use `--mordor` in place of `--classic`. The new container takes the old one's volumes, so
the chain carries over. The old container is what you roll back to.

Then, in this order:

```bash
docker exec core-geth geth version
docker exec core-geth geth attach --exec 'admin.nodeInfo.enode' /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec 'net.peerCount'        /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec 'eth.blockNumber'      /root/.ethereum/classic/geth.ipc
```

If a line cannot connect, give the node a minute to open its socket and run that line again.

- **The first line** must name a `v1.13` release.
- **The enode ID** must differ from the one you recorded. If it is the same, the node is still on the old key: stop
  it cleanly and rename `geth/nodekey` in the data directory the node really uses. If it differs, it is the new ID
  to send to your peers.
- **The peer count** climbs within minutes. **A count stuck at zero after ten minutes usually means a peer still
  lists your old enode ID.**
- **The block number** resumes where the node stopped. Last, compare it with a block explorer or another node you
  operate. If it settles somewhere else, stop and investigate rather than resyncing.

## 6. If you need to roll back

The data directory works in both directions, so rolling back swaps the container for the old one:

```bash
docker stop -t 300 core-geth &&
  docker rm core-geth &&
  docker rename core-geth-v1.12 core-geth &&
  docker update --restart unless-stopped core-geth &&
  docker start core-geth
```

The old container keeps its own flags, so it has no `--mess` flag to take out, and it takes the same volumes. Put its
restart policy in place of `unless-stopped`. The rotated node key stays rotated: rolling back does not un-expose the
old key.
