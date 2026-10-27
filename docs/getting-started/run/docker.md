---
description: "Run a Core-Geth Ethereum Classic or Mordor node in Docker: the image, a chain volume, ports, and a restart policy that brings the node back after a failure or reboot."
---

# Docker users guide

This guide runs a node in a Docker container, on Ethereum Classic or on the Mordor test network, with
its chain in a volume and a restart policy that brings it back after a failure or a reboot.
[Running a node](../run-a-node.md) covers what every platform shares.

Commands are shown for both networks. Use the lines for the network you chose, and keep its flag on
every command.

## 1. Get the image

[Installation, Docker](../installation.md#docker) gives the ways to get it: pull it, load it from the
release, or build it. This guide uses the release image's name,
`ghcr.io/ethereumclassic/core-geth:latest`; for an image you built, use your own tag. Check it:

```bash
docker run --rm ghcr.io/ethereumclassic/core-geth:latest version    # prints Core-Geth and its version
```

## 2. Choose the network

| Network | Flag | Container and volume in this guide | Data directory, inside the container |
|---|---|---|---|
| Ethereum Classic | `--classic` | `core-geth-classic` | `/root/.ethereum/classic` |
| Mordor | `--mordor` | `core-geth-mordor` | `/root/.ethereum/mordor` |

The names are this guide's choice. Each node gets its own container, and a volume of the same name
mounted at `/root`, which keeps the chain when the container is replaced.

## 3. Start the node

```bash
# Ethereum Classic
docker run -d --name core-geth-classic \
  --restart unless-stopped --stop-timeout 300 \
  -v core-geth-classic:/root \
  -p 30303:30303 -p 30303:30303/udp \
  ghcr.io/ethereumclassic/core-geth:latest --classic

# Mordor
docker run -d --name core-geth-mordor \
  --restart unless-stopped --stop-timeout 300 \
  -v core-geth-mordor:/root \
  -p 30303:30303 -p 30303:30303/udp \
  ghcr.io/ethereumclassic/core-geth:latest --mordor
```

- **`-p 30303:30303 -p 30303:30303/udp`** lets other nodes reach yours, over TCP and UDP.
- **`--restart unless-stopped`** starts the container again after the node fails, and when the Docker
  daemon starts, such as after a reboot. A container you stopped yourself stays stopped
  ([restart policies](https://docs.docker.com/engine/containers/start-containers-automatically/)).
- **`--stop-timeout 300`** gives the node five minutes to stop before Docker kills it. Docker's
  default, 10 seconds, is shorter than a stop during the first sync can take
  ([stop times](../hardware-requirements.md#stopping)).
- **`-v core-geth-mordor:/root`** keeps the chain in a named volume. To keep it in a directory of
  your choice, such as on another disk, mount that directory instead: `-v /data/core-geth-mordor:/root`.

**For a [configuration](../run-a-node.md#choose-your-configuration)**, add its flags after the
network flag. A configuration that serves HTTP also needs `--http.addr 0.0.0.0` and
`-p 127.0.0.1:8545:8545`, which step 5 explains.

To run both networks on one host, give the Mordor node its own peer port: `--port 30304` after
`--mordor`, published as `-p 30304:30304 -p 30304:30304/udp`.

## 4. Check that it syncs

Follow the log; Ctrl-C stops following it, and the node keeps running. The first line names the
network: `Starting Core-Geth on Ethereum Classic...` or `Starting Core-Geth on Mordor testnet...`.

```bash
docker logs -f core-geth-classic      # Ethereum Classic
docker logs -f core-geth-mordor       # Mordor
```

[What a healthy first sync looks like](../run-classic-node.md#what-a-healthy-first-sync-looks-like)
explains the lines that follow. To ask the node itself, run `attach` inside its container:

```bash
# Ethereum Classic
docker exec core-geth-classic geth --classic attach --exec 'eth.syncing'       # an object while it syncs, false once synced
docker exec core-geth-classic geth --classic attach --exec 'net.peerCount'     # above 0 within minutes
docker exec core-geth-classic geth --classic attach --exec 'eth.blockNumber'   # once synced, compare with a block explorer

# Mordor
docker exec core-geth-mordor geth --mordor attach --exec 'eth.syncing'
docker exec core-geth-mordor geth --mordor attach --exec 'net.peerCount'
docker exec core-geth-mordor geth --mordor attach --exec 'eth.blockNumber'
```

`attach` looks for the node's socket in the data directory of the network you name, so a Mordor node
needs `--mordor` here too. [How to tell it is done](../run-classic-node.md#how-to-tell-it-is-done)
lists every sign of a finished sync.

## 5. Use it

`attach` without `--exec` opens an interactive console on the node; `-it` gives it your terminal.
Type `exit` to leave it; the node keeps running.

```bash
docker exec -it core-geth-mordor geth --mordor attach
```

For wallets and programs on the host, the node serves JSON-RPC when it is started with
`--http --http.addr 0.0.0.0`, and `-p 127.0.0.1:8545:8545` publishes it on the host's loopback
interface only. [Run the container](../installation.md#run-the-container) explains why both are
needed. A container keeps the flags it was created with, so replace it; the volume keeps the chain:

```bash
docker stop core-geth-mordor && docker rm core-geth-mordor
docker run -d --name core-geth-mordor \
  --restart unless-stopped --stop-timeout 300 \
  -v core-geth-mordor:/root \
  -p 30303:30303 -p 30303:30303/udp \
  -p 127.0.0.1:8545:8545 \
  ghcr.io/ethereumclassic/core-geth:latest --mordor --http --http.addr 0.0.0.0
```

Once the log shows the node has started, ask it for the chain ID from the host:

```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}' http://127.0.0.1:8545
```

The `result` is `0x3d` on Ethereum Classic and `0x3f` on Mordor. Read
[RPC exposure](../../operate/security.md#rpc-exposure) before you publish it anywhere else.

## 6. Keep it running

The restart policy from step 3 already does: the node comes back after a failure and after a reboot.
Check a container's policy and stop timeout with:

```bash
docker inspect -f '{{.HostConfig.RestartPolicy.Name}} {{.Config.StopTimeout}}' core-geth-mordor    # unless-stopped 300
```

## 7. Stop it cleanly

```bash
docker stop core-geth-classic      # Ethereum Classic
docker stop core-geth-mordor       # Mordor
```

`docker stop` sends the node SIGTERM and waits for it to exit, for up to the stop timeout. Check that
its log ends with `Blockchain stopped`:

```bash
docker logs --tail 5 core-geth-mordor
```

The container stays stopped, across reboots too, until `docker start core-geth-mordor`.
[Stop it safely](../run-classic-node.md#stop-it-safely) shows what a clean stop logs.

## 8. Update it

Get the new release's image as in step 1, then replace the container. The volume keeps the chain, so
the node resumes where it stopped:

```bash
docker stop core-geth-mordor && docker rm core-geth-mordor
```

Pull `:latest` again and run the step 3 command, keeping any flags you added; to pin a release instead, use its tag from the releases page.
[Upgrading within v1.13.x](../../operate/maintenance.md#upgrading-within-v113x) lists what to read
and check. A node on a v1.12.x release follows the
[Docker migration guide](../../tutorials/migration/docker.md) instead.
