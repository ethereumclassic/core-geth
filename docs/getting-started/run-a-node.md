---
description: "Run an Ethereum Classic or Mordor node with Core-Geth, with the flags, ports and service setup each kind of node needs."
---

# Running a node

This section takes a `geth` you have [installed](installation.md) to a node that syncs, keeps
running, and stops cleanly. The steps are the same on every platform while the commands differ, so
each platform has its own guide.

## Choose your guide

| You run the node on | Guide | How the guide keeps it running |
|---|---|---|
| Linux | **[Linux users guide](run/linux.md)** | A systemd service |
| macOS | **[Mac users guide](run/macos.md)** | A launchd agent |
| Windows | **[Windows users guide](run/windows.md)** | A window that opens when you sign in |
| Docker | **[Docker users guide](run/docker.md)** | A container restart policy |

Every guide takes the same steps: install `geth`, choose the network, start the node, check that it
syncs, use it, keep it running, stop it cleanly, and update it.

## Choose your network

**The network flag decides which chain the node follows.** Pass it every time you start the node.
On Linux, macOS and Docker, pass it to `geth attach` as well: `attach` looks for the node inside that
network's data directory.

| Network | Use it for | Flag | Chain ID |
|---|---|---|---|
| Ethereum Classic | The main network | `--classic`, or no flag | 61 |
| Mordor | Testing: trying the client, testing contracts, rehearsing a deployment | `--mordor` | 63 |

MintMe.com Coin runs with `--mintme`, and a private network runs from its own genesis
([Private network](../tutorials/private-network.md)).

Each network keeps its own data directory:

| Platform | Ethereum Classic | Mordor |
|---|---|---|
| Linux | `~/.ethereum/classic` | `~/.ethereum/mordor` |
| macOS | `~/Library/Ethereum/classic` | `~/Library/Ethereum/mordor` |
| Windows | `%LOCALAPPDATA%\Ethereum\classic` | `%LOCALAPPDATA%\Ethereum\mordor` |
| Docker | `/root/.ethereum/classic`, inside the container | `/root/.ethereum/mordor`, inside the container |

**To move a node to the other network, stop it and start it with the other flag.** It syncs that
network into that network's directory and leaves the first network's data where it was. `--datadir`
puts the data somewhere else; give each network its own. [Run a Mordor node](run-mordor-node.md)
covers running both networks on one machine.

## Choose your configuration

**A configuration is the set of flags that suits what the node is for.** Add them to the start
command, after the network flag. Your platform guide shows where they go, both when you start the
node yourself and in the service that keeps it running. The flags below name Ethereum Classic; on
Mordor, use `--mordor` in place of `--classic`.

| The node is for | Configuration | Full guide |
|---|---|---|
| Yourself, or a wallet on the same machine | [A personal node](#a-personal-node) | This section |
| Crediting deposits and sending withdrawals | [An exchange or custodian](#an-exchange-or-custodian) | [Production operations](../guides/production-operations.md) |
| Serving JSON-RPC to other people | [A public RPC endpoint](#a-public-rpc-endpoint) | [Public RPC endpoint](../guides/public-rpc-endpoint.md) |
| Explorers, indexers and analytics | [An archive node](#an-archive-node) | [Archive node](../guides/archive-node.md) |
| Mining with your own hardware | [A solo mining node](#a-solo-mining-node) | [Mining](../guides/mining.md) |
| Supplying work to a mining pool | [A mining pool node](#a-mining-pool-node) | [Mining pool node](../guides/mining-pool-node.md) |
| Helping other nodes find peers | [A bootnode](#a-bootnode) | [Bootnodes and peer discovery](../guides/bootnodes-and-discovery.md) |

### A personal node

```text
--classic
```

The platform guide's commands, unchanged. For a wallet on the same machine that connects over HTTP,
add `--http`. The node then serves JSON-RPC on `127.0.0.1:8545`, and to nothing beyond the machine.

### An exchange or custodian

```text
--classic --http --http.api eth,net,web3
```

- Keep the HTTP listener on `127.0.0.1`, or on a private network address set with `--http.addr`, where
  only your own systems reach it ([RPC exposure](../operate/security.md#rpc-exposure)).
- Run at least two nodes on separate hosts, on the same release and the same MESS setting
  ([Redundancy](../guides/production-operations.md#redundancy)).
- Credit a deposit only from a node that is in sync, after the number of confirmations the
  [MESS confirmation calculator](../guides/mess-calculator.md) gives
  ([Crediting deposits](../guides/production-operations.md#crediting-deposits)).

### A public RPC endpoint

```text
--classic --http --http.api eth,net,web3 --http.vhosts <hostname>
```

- The node listens on `127.0.0.1`, and a reverse proxy or RPC gateway faces the internet in front of
  it ([What to put in front of the node](../guides/public-rpc-endpoint.md#what-to-put-in-front-of-the-node)).
- `--http.vhosts` names the public hostname the proxy passes on
  ([Hostnames](../guides/public-rpc-endpoint.md#hostnames)).
- Serve only `eth`, `net` and `web3`, and keep account keys off the node. The flags that cap gas,
  call time and batch size are under
  [Limits the node enforces](../guides/public-rpc-endpoint.md#limits-the-node-enforces).

### An archive node

```text
--classic --syncmode full --gcmode archive
```

- Start it on a new data directory. Full sync executes every block from genesis, and archive mode
  writes the state of every block to disk, so it needs more disk and time than the default
  ([Hardware requirements](hardware-requirements.md)).
- Keep the hash scheme, which a new data directory gets by default. Under `--state.scheme path`, an
  archive node keeps no older state ([With the path scheme](../guides/archive-node.md#with-the-path-scheme)).
- Serve it to your applications with `--http --http.api eth,net,web3`. Tracing calls are in the `debug`
  namespace, which also holds methods that change the node
  ([Serve it to others](../guides/archive-node.md#serve-it-to-others)).

### A solo mining node

```text
--classic --mine --miner.etherbase <your address> --http --http.api eth,net,web3
```

- Your mining software fetches work from the node with `eth_getWork`, on `127.0.0.1:8545` or a private
  network address. `--mine` does not mine on the node's own CPU unless `--miner.threads` asks it to.
- The node needs no key for the reward address: the address only goes into the blocks it builds.
- Mine with a node that has finished syncing. To practice first, [mine on Mordor](../guides/mordor-mining.md).

### A mining pool node

```text
--classic --mine --miner.etherbase <pool address> --http --http.api eth,net,web3 --miner.notify http://127.0.0.1:<port>/
```

- Only the pool software should reach the HTTP endpoint: anyone who reaches it can request work and
  submit results.
- `--miner.notify` posts each new work package to the pool software as soon as the node has it.
- Run at least two pool nodes, peered with each other, on the same release and MESS setting
  ([Run more than one node](../guides/mining-pool-node.md#run-more-than-one-node)).

### A bootnode

```text
--classic --nat extip:<public IP> --identity <your name>
```

- Open port 30303 on TCP and UDP. `--nat extip` advertises the address other nodes can reach.
- Keep `geth/nodekey`. It sets the enode ID other operators list, and a new key is a new ID.
- The `bootnode` tool runs discovery without a chain
  ([Run your own bootnode](../guides/bootnodes-and-discovery.md#run-your-own-bootnode)).

## The same on every platform

- [What a healthy first sync looks like](run-classic-node.md#what-a-healthy-first-sync-looks-like),
  and the [warnings a first sync logs](run-classic-node.md#warnings-a-first-sync-logs)
- [How to tell it is done](run-classic-node.md#how-to-tell-it-is-done)
- [Where your data lives](run-classic-node.md#where-your-data-lives), inside the data directory
- [Stop it safely](run-classic-node.md#stop-it-safely)
- [MESS on this node](run-classic-node.md#mess-on-this-node)
- [Hardware requirements](hardware-requirements.md): disk, memory and sync time for each network
- [Security and network exposure](../operate/security.md), before you open a port
