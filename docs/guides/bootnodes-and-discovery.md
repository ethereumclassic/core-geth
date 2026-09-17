---
title: Bootnodes and peer discovery
description: "How Core-Geth finds peers on Ethereum Classic: bootnodes, the DNS discovery trees, how the lists are built from a crawl, and how to publish your own."
---

A node needs peers to receive blocks. This page covers where core-geth finds them, how to make your node a
useful peer for others, and how the network's discovery lists are built and published.

## Where a node finds peers

| Source | Set by | Default on Ethereum Classic and Mordor |
|---|---|---|
| Bootnodes | `--bootnodes` | The list compiled into the client, in `params/bootnodes_classic.go` and `params/bootnodes_mordor.go` |
| DNS discovery trees | `--discovery.dns` | The trees in the same files: `ClassicDNSNetwork1` to `3` or `MordorDNSNetwork1` to `3` |
| Discovery v4 | `--discovery.v4` | On |
| Discovery v5 | `--discovery.v5` | Off |
| Static and trusted peers | `[Node.P2P]` in the config file, or `admin_addPeer` | None |

Bootnodes and DNS trees give the node its first peers. Discovery then finds more.

- `--bootnodes enode://...,enode://...` replaces the compiled list.
- `--discovery.dns enrtree://<public key>@<domain>` replaces the trees, and `--discovery.dns ""` turns DNS
  discovery off.
- `--nodiscover` turns discovery off. The node then connects only to static and trusted peers, and to peers
  that connect to it.
- `--maxpeers`, 50 by default, caps the number of peers. Trusted peers can connect past it.

## Make your node a useful peer

- **Open the p2p port.** Port 30303 on TCP for connections and on UDP for discovery, set by `--port` and
  `--discovery.port`: [Ports and listeners](../operate/security.md#ports-and-listeners).
- **Advertise the right address.** Behind NAT, pass `--nat extip:<public IP>` so the node advertises the
  address other nodes can reach.
- **Keep the node key.** Your enode ID comes from `<datadir>/geth/nodekey`, and a new key is a new enode
  ID. Back it up with the rest of the node: [What to back up](../operate/maintenance.md#what-to-back-up).
- **Name your node.** `--identity <your name>` puts a name in the client string other operators see. Use
  one that identifies you, not one that suggests someone else runs the node.

## How the discovery lists are built

The DNS trees this client uses by default are published from
[`ethereumclassic/discv4-dns-lists`](https://github.com/ethereumclassic/discv4-dns-lists), whose README
describes how often they are rebuilt and what gets a node listed. The lists come from crawls of the
discovery network, so a node can be listed when it:

- answers discovery on a reachable address, and
- advertises an `eth` entry for the network in its node record, which
  `devp2p nodeset filter -eth-network classic` (or `mordor`) checks.

## Run your own bootnode

The `bootnode` tool, in the `alltools` archive, runs discovery only, without syncing a chain.

```sh
bootnode -genkey boot.key
bootnode -nodekey boot.key -addr :30301 -nat extip:<public IP>
```

| Flag | Default | Meaning |
|---|---|---|
| `-addr` | `:30301` | Listen address |
| `-genkey` | none | Write a new node key to this file |
| `-nodekey` | none | The node key file to use |
| `-writeaddress` | false | Print the node's public key and exit |
| `-nat` | `none` | Port mapping: `any`, `none`, `upnp`, `pmp`, `pmp:<IP>` or `extip:<IP>` |
| `-netrestrict` | none | Restrict communication to these CIDR ranges |
| `-v5` | false | Run a discovery v5 bootnode |
| `-verbosity` | 3 | Log verbosity, from 0 to 5 |

[Private network, step 3](../tutorials/private-network.md#3-start-a-bootnode) walks through one. Adding a
bootnode to the list compiled into the client is a change to `params/bootnodes_classic.go` or
`params/bootnodes_mordor.go`, proposed by pull request.

## Publish your own discovery tree

The `devp2p` tool, also in `alltools`, does each step:

1. **Crawl** the network with `devp2p discv4 crawl`.
2. **Filter** the result with `devp2p nodeset filter`, using `-eth-network classic` or `mordor`.
3. **Sign** a tree with a key from `devp2p key generate`, using `devp2p dns sign`.
4. **Deploy** it to DNS with `devp2p dns to-cloudflare` or `devp2p dns to-route53`.

[The devp2p command](https://github.com/ethereumclassic/core-geth/blob/main/cmd/devp2p/README.md) gives the
arguments for each. Nodes then use the tree with `--discovery.dns enrtree://<public key>@<domain>`. Keep the
signing key offline: anyone holding it can change the peers your tree gives out.

## Static and trusted peers

Static peers are nodes the client keeps connecting to. Trusted peers can also connect when the node already
has `--maxpeers` peers.

```toml
[Node.P2P]
StaticNodes = ["enode://<node ID>@<IP>:30303"]
TrustedNodes = ["enode://<node ID>@<IP>:30303"]
```

`admin_addPeer` and `admin_addTrustedPeer` add them to a running node, until it restarts. A
`static-nodes.json` or `trusted-nodes.json` file in the data directory is not read:
[Troubleshooting](../operate/troubleshooting.md#why-does-the-node-ignore-static-nodesjson).

`--netrestrict` limits peers to a list of CIDR ranges: [Peer restrictions](../operate/security.md#peer-restrictions).
A node with no peers is covered in [Troubleshooting](../operate/troubleshooting.md#why-does-the-node-have-no-peers).
