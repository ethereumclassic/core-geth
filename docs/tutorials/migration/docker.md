---
title: Migrating to v1.13 in Docker
description: "Upgrade an Ethereum Classic node in Docker from any v1.12.x image to the Core-Geth v1.13 line: the steps, the node-key rotation, and how to roll back."
---

# Docker users guide

This guide upgrades a node running in Docker from any v1.12.x image to the v1.13 line.
[Migrating to v1.13](../v1.13.0-migration.md) explains why, and what changes for every node.

The image moves from `etclabscore/core-geth` to images built from `ethereumclassic/core-geth`. Your volume, ports and
flags carry over, apart from what steps 3 and 4 change. The commands assume a container named `core-geth` with its data
mounted at `/root`, from a host directory as [Installation, Run the
container](../../getting-started/installation.md#run-the-container) shows or from a named volume as the [Docker run
guide](../../getting-started/run/docker.md) shows. Put your own container name in place of `core-geth`. For Mordor, use
`.ethereum/mordor` in the paths.

## Choose how you get the image

| Method | Image name it gives you |
|---|---|
| [Pull from the registry](../../getting-started/installation.md#docker) | `ghcr.io/ethereumclassic/core-geth:latest` |
| [Load the release tarball](../../getting-started/installation.md#load-the-image-from-the-release) | `ghcr.io/ethereumclassic/core-geth:<tag>` |
| [Build an image from the release archive](../../getting-started/installation.md#build-an-image-from-the-release-archive) | the tag you give it, for example `core-geth:local` |
| [Build an image from source](../../getting-started/installation.md#build-an-image-from-source) | the tag you give it, for example `core-geth:local` |

The pull needs no account and no token. The other three methods are for a machine with no route to the registry, or if
you would rather build what you run.

## Before you start

Record what you are running:

```bash
docker exec core-geth geth version
docker exec core-geth geth attach --exec 'eth.blockNumber'      /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec 'admin.nodeInfo.enode' /root/.ethereum/classic/geth.ipc
docker inspect core-geth --format '{{.HostConfig.RestartPolicy.Name}}'
```

**Keep the last two lines' output.** The first is your current enode ID, which step 3 changes, and the second is the
container's restart policy, which step 6 restores if you roll back.

Then check where the container keeps its data:

```bash
docker inspect core-geth --format '{{range .Mounts}}{{.Destination}} {{end}}'
```

It prints `/root`, or `/root/.ethereum`, when the data is in a host directory or a named volume, which the steps below
reach through `--volumes-from`. If it prints neither, as when the only mount is a config file, the data is inside the
container itself: follow [If the data is inside the container](#if-the-data-is-inside-the-container) in step 3, or the
new container starts with no chain and resyncs.

## 1. Stop the node cleanly

```bash
docker stop -t 300 core-geth
docker logs --tail 50 core-geth 2>&1 | grep 'Blockchain stopped'
```

Continue once the second command prints `Blockchain stopped`. If it prints nothing, the node was killed before it
stopped cleanly. You can still continue: the node repairs its database at its next start, and its block number can
start lower.

**Use `-t 300`.** A plain `docker stop` kills the node if it has not stopped within 10 seconds. A killed node that was
importing blocks moves its head back on its next start and imports them again
([After an unclean stop](../../operate/maintenance.md#after-an-unclean-stop)).

## 2. Get the latest image

This pulls the latest release from the registry, and sets `IMAGE` for the steps below:

```bash
IMAGE=ghcr.io/ethereumclassic/core-geth:latest
docker pull $IMAGE
```

With another of the methods above, set `IMAGE` to the name it gives you. On a node that holds account keys or that
others rely on, also check where a pulled image came from, which needs the [GitHub CLI](https://cli.github.com/), signed
in with `gh auth login`:

```bash
gh attestation verify oci://$IMAGE --repo ethereumclassic/core-geth
```

It must print `Verification succeeded`. If it does not, do not run the image.

Keep your previous image; step 6 uses it.

## 3. Rotate the P2P node key

**Required.** [Rotate the P2P node key](../v1.13.0-migration.md#rotate-the-p2p-node-key) explains why. This runs in a
throwaway container from the new image, which reaches the data through the old container's mounts:

```bash
docker run --rm --volumes-from core-geth --entrypoint sh $IMAGE \
  -c "mv -n /root/.ethereum/classic/geth/nodekey /root/.ethereum/classic/geth/nodekey.old-rotated-$(date +%F)"
```

The next start writes a fresh key. If you start the node with `--nodekey` or `--nodekeyhex`, or set `GETH_NODEKEY` or
`GETH_NODEKEYHEX`, leave that out of the `docker run` in step 5: it overrides the file.

If the rename prints an error, check the data directory: `geth/nodekey` must not be there when the node starts.
`No such file or directory` means the data is not where the command looks: check the path, and see
[If the data is inside the container](#if-the-data-is-inside-the-container).

### If the data is inside the container

When the check in [Before you start](#before-you-start) printed neither, copy the data out of the stopped container,
rename the key in the copy, and run the new container on the copy. The copy needs as much free space as the chain takes:

```bash
mkdir -p $HOME/core-geth-data
docker cp core-geth:/root/. $HOME/core-geth-data/
docker run --rm -v $HOME/core-geth-data:/root --entrypoint sh $IMAGE \
  -c "mv -n /root/.ethereum/classic/geth/nodekey /root/.ethereum/classic/geth/nodekey.old-rotated-$(date +%F)"
```

In step 5, give the new container `-v $HOME/core-geth-data:/root` in place of `--volumes-from core-geth-v1.12`. The
chain carries over with the copy, so the node does not resync.

### If this node lists static peers, update the list

They are under `[Node.P2P]` in the file you pass with `--config`. Edit it while the node is down:

```toml
[Node.P2P]
StaticNodes = ["enode://<node ID>@<IP>:30303"]
```

A peer that has already moved has a new ID, so use that. Separate entries with commas. `TrustedNodes` takes the same
form. The client ignores `static-nodes.json`, so move its entries here. With no file yet, save these two lines.

To give the new container the file, add `-v <path>/config.toml:/config.toml:ro` to the `docker run` in step 5, and
`--config /config.toml` after the image name.

## 4. Re-check RPC exposure

Go through the [RPC exposure checklist](../v1.13.0-migration.md#re-check-rpc-exposure). In Docker, publish the RPC port
on `127.0.0.1` only, as the warning under [Run the container](../../getting-started/installation.md#run-the-container)
explains. If you use a config file, check it for [the lines listed
here](../v1.13.0-migration.md#if-your-config-file-came-from-dumpconfig), whoever wrote it. If the old container ran with
`--discovery.dns` or `--bootnodes`, [leave them
out](../v1.13.0-migration.md#if-you-added-discoverydns-or-bootnodes-on-v112) of the new one.

## 5. Start and verify

Keep the old container, stopped, and start a new one from the new image. First read the old container's flags, port
mappings and environment, for the last lines of the `docker run` below:

```bash
docker inspect core-geth --format '{{json .Args}}'
docker inspect core-geth --format '{{json .HostConfig.PortBindings}}'
docker inspect core-geth --format '{{json .Config.Env}}'
```

Put its flags after the image name, keep its peer-to-peer port mappings, and pass each `GETH_` variable with `-e`. Keep
the RPC port published on `127.0.0.1`, as step 4 says, even if the old container published it on every interface. For
Mordor, use `--mordor` in place of `--classic`. Then:

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

Then, in this order:

```bash
docker exec core-geth geth version
docker exec core-geth geth attach --exec 'admin.nodeInfo.enode' /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec 'net.peerCount'        /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec 'eth.blockNumber'      /root/.ethereum/classic/geth.ipc
```

If a line cannot connect, wait a minute and run it again.

- **The first line** must name the version you installed: a `v1.13` release, or, for an image built from source,
  the `-unstable` version of `main` it was built from.
- **The enode ID** must differ from the one you recorded. If it is the same, the node is still on the old key: stop it
  and redo step 3.
- **The peer count** climbs within minutes. If it stays at zero, see
  [Why does the node have no peers?](../../operate/troubleshooting.md#why-does-the-node-have-no-peers)
- **The block number** starts near where the node stopped, lower after an unclean stop, then climbs. Compare it with
  a block explorer or another node you run. If it stops climbing short of the chain head, stop and investigate
  rather than resync.

### If other nodes list this node, give them the new ID

They can no longer connect to this node, and neither side logs it at the default verbosity. On each one, on `v1.12` or
`v1.13`, run these lines in that node's own container, with its name in place of `core-geth`, using the old ID you
recorded, the new ID that step 5 printed, and the address and port the entry already has. On a node that is not in a
container, run the `geth attach` part of each line on its host:

```bash
docker exec core-geth geth attach --exec "admin.removePeer('enode://<old ID>@<IP>:30303')" /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec "admin.addPeer('enode://<new ID>@<IP>:30303')" /root/.ethereum/classic/geth.ipc
docker exec core-geth geth attach --exec "admin.peers.map(function(p){return p.enode})" /root/.ethereum/classic/geth.ipc
```

The first two lines each print `true` and take effect with no restart. The third lists the new ID once the nodes
connect. The change lasts only until that node restarts, so put the new ID in that node's config file too. Send the new
ID to anyone else who lists this node.

For an entry in `TrustedNodes`, run `admin.removeTrustedPeer` and `admin.addTrustedPeer` the same way, as well as the
first two lines if that node also lists this one as a static peer. They do not open a connection.
[Static and trusted peers](../../guides/bootnodes-and-discovery.md#static-and-trusted-peers) explains the two lists.

## 6. If you need to roll back

The data directory works in both directions, so rolling back swaps the container for the old one:

```bash
docker stop -t 300 core-geth &&
  docker rm core-geth &&
  docker rename core-geth-v1.12 core-geth &&
  docker update --restart unless-stopped core-geth &&
  docker start core-geth
```

Put the restart policy you recorded in place of `unless-stopped`. Keep the new key: rolling back does not make the old
key safe. A release before `v1.12.21` leaks bits of the key it runs with, so rolling back to one exposes the new key
too.

If you copied the data out, run only the first two lines, which stop and remove the new container, and do not start the
old one: its own copy of the data still holds the old key. Run the image it used,
`docker inspect core-geth-v1.12 --format '{{.Config.Image}}'`, on the copy instead, with
`-v $HOME/core-geth-data:/root` and the old container's ports and flags.
