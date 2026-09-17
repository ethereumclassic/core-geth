---
title: Security and network exposure
description: "Securing an Ethereum Classic node: which ports and RPC modules to expose, which to keep closed, and how to set up a host firewall."
---

A node needs only its peer-to-peer ports reachable from the internet. Everything else it listens on
is for the machine itself or for a network you trust, and the table below says which listener is
which. Before a binary runs on a node that holds value,
[verify where the archive came from](../getting-started/installation.md#verify-where-the-archive-came-from).

## Ports and listeners

| Listener | Flag and default | Default address | Starts by default | May face the internet |
| --- | --- | --- | --- | --- |
| Peer connections, TCP | `--port`, 30303 | every interface | yes | yes |
| Peer discovery, UDP | `--discovery.port`, 30303; without it, follows `--port` | every interface | yes, unless `--nodiscover` | yes |
| HTTP JSON-RPC, TCP | `--http.port`, 8545 | `localhost` (`--http.addr`) | no: `--http` | no |
| GraphQL | none: it answers at `/graphql` on the HTTP listener | the HTTP listener's | no: `--graphql`, which needs `--http` | no |
| WebSocket JSON-RPC, TCP | `--ws.port`, 8546 | `localhost` (`--ws.addr`) | no: `--ws` | no |
| Metrics, TCP | `--metrics.port`, 6060 | none: `--metrics.addr` sets it | no: `--metrics` and `--metrics.addr` | no |
| Profiling (pprof), TCP | `--pprof.port`, 6060 | `127.0.0.1` (`--pprof.addr`) | no: `--pprof` | no |
| IPC | `--ipcpath`, `geth.ipc` | a socket in the data directory; a named pipe on Windows | yes, unless `--ipcdisable` | no |

The Engine API, the authenticated interface Ethereum consensus clients use, starts only on a chain
configured for the merge, which no network built into the client is. There it listens on `localhost`,
port 8551, set by `--authrpc.addr` and `--authrpc.port`.

`localhost` means the node listens on `127.0.0.1`. Metrics and pprof share a default port, so a node
that runs both needs one of them moved. On Linux and macOS the IPC socket can be opened only by the
user the node runs as, and by root.

On Linux, `ss` lists what each process listens on:

```shell
$ ss -lntup
```

A node with the default settings has two lines, for peer connections and discovery on every
interface:

```
udp UNCONN 0      0                                          *:30303       *:* users:(("geth",pid=<pid>,fd=70))
tcp LISTEN 0      4096                                       *:30303       *:* users:(("geth",pid=<pid>,fd=67))
```

HTTP, WebSocket and pprof each add a line on `127.0.0.1`, and metrics adds one on the address
`--metrics.addr` names. These four lines come from a node that ran all of them, with metrics on
`127.0.0.1` and every port moved off its default: metrics on 19514, pprof on 19515, HTTP on 19512 and
WebSocket on 19513.

```
tcp LISTEN 0      4096                               127.0.0.1:19514 0.0.0.0:* users:(("geth",pid=<pid>,fd=6))
tcp LISTEN 0      4096                               127.0.0.1:19515 0.0.0.0:* users:(("geth",pid=<pid>,fd=3))
tcp LISTEN 0      4096                               127.0.0.1:19512 0.0.0.0:* users:(("geth",pid=<pid>,fd=70))
tcp LISTEN 0      4096                               127.0.0.1:19513 0.0.0.0:* users:(("geth",pid=<pid>,fd=71))
```

With the default `--nat any`, a node behind a router that supports UPnP or NAT-PMP asks the router to
forward its peer-to-peer ports, and no other port. `--nat none` turns that off.

## A host firewall

This nftables ruleset drops every new inbound connection except peer-to-peer traffic and SSH.
Replies to connections the host opened still arrive, and so does the host's own traffic on the
loopback interface, so RPC on `127.0.0.1` keeps answering. Save it as `core-geth.nft`:

```
# Inbound firewall for a core-geth node: peer-to-peer and SSH in, nothing else.

# Replace this table if it is already loaded, and leave other tables alone.
table inet core-geth
delete table inet core-geth

table inet core-geth {
    chain input {
        type filter hook input priority 0; policy drop;

        # Replies to connections this host opened
        ct state vmap { established : accept, related : accept, invalid : drop }

        # The host's own traffic, including local RPC
        iifname "lo" accept

        # IPv6 neighbor discovery, without which IPv6 stops working
        icmpv6 type { nd-neighbor-solicit, nd-router-advert, nd-neighbor-advert } accept

        # SSH: change 22 to your SSH port, or delete the line if you do not use SSH
        tcp dport 22 accept

        # Peer connections (TCP) and discovery (UDP): change 30303 if you set --port
        tcp dport 30303 accept
        udp dport 30303 accept
    }
}
```

Load it as root:

```shell
$ sudo nft -f core-geth.nft
```

- **Other machines then reach only TCP and UDP 30303 and SSH.** That holds for an RPC listener bound
  to every interface too.
- **Loading the file again replaces the `core-geth` table** and leaves any other table in place.
- **An existing firewall still applies.** A packet this table accepts is dropped if another table on
  the same host drops it, so open the peer-to-peer ports there as well.
- **Match the node's ports.** If you set `--port`, use that port on both peer-to-peer lines; if you
  set `--discovery.port`, use it on the `udp` line.
- **Load it at boot** with your distribution's nftables service.

## RPC exposure

The JSON-RPC and GraphQL listeners have no TLS, authentication or rate limiting of their own. Keep
them on `127.0.0.1`, and put a proxy that provides those in front of any that other machines must
reach. Before you widen one:

- [Docker](../getting-started/installation.md#docker) shows how publishing a container's RPC port
  can expose it on every interface of the host.
- [4. Re-check RPC exposure](../tutorials/v1.13.0-migration.md#re-check-rpc-exposure) lists what
  to check in `--http.addr`, `--http.api`, `--http.corsdomain`, `--ws.addr` and `--ws.api`.
- [Using JSON-RPC APIs](../JSON-RPC-API/index.md) says which namespaces each interface serves,
  including why an empty `--http.api` is refused.
- [Public RPC endpoint](../guides/public-rpc-endpoint.md) covers which namespaces to serve to other
  machines, and what a proxy in front must supply.

## Keys

### The node key

`<datadir>/geth/nodekey` holds the node's private key for the peer-to-peer network, and the node's
enode ID comes from it. The node reads the file at every start, and writes a new key there if it
cannot load one. The key is stored unencrypted, in a file only its owner can read or write.

A copy of a data directory carries its `nodekey`, so two hosts started from copies of one data
directory have the same enode ID. `--nodekey` loads the key from another file instead.

A node upgraded from a v1.12.x release needs a new key:
[3. Rotate the P2P node key](../tutorials/v1.13.0-migration.md#rotate-the-p2p-node-key) explains why
and how.

### Account keys

Account keys are keyfiles in `<datadir>/keystore`, or in the directory `--keystore` names, each
encrypted with its passphrase.
[Key derivation on a small machine](../getting-started/run-cli.md#key-derivation-on-a-small-machine)
covers how hard that encryption is to break.

[`clef`](../core/alltools.md#executables) is a signer that runs as its own process and holds account
keys outside the node. Started with `--signer` and the path of clef's IPC socket, the node sends
signing requests to `clef` and opens no keystore of its own. Start `clef` first: a node that cannot
reach its signer stops as it starts.

```
Fatal: Failed to set account manager backends: error connecting to external signer: dial unix <signer path>: connect: no such file or directory
```

`--unlock` with `--password` unlocks keystore accounts as the node starts, and they stay unlocked
until it stops. While an account is unlocked, any caller that reaches the `eth` namespace can make
it sign, with `eth_sign` or `eth_sendTransaction`, and no passphrase. On a node whose RPC other
machines can reach, leave accounts locked, never set `--allow-insecure-unlock`, and sign with
`clef`. [Start the signers](../tutorials/private-network.md#start-the-signers) shows the refusal
`geth` gives when `--unlock` meets HTTP or WebSocket RPC.

## Peer restrictions

`--netrestrict` takes a comma-separated list of networks in CIDR notation. The node then dials only
peers in those networks, closes incoming connections from any other address as they arrive, and
drops nodes outside those networks from discovery replies and from the bootnode list. Its discovery
port still answers packets from any address; `--nodiscover` closes that port.

```shell
$ geth --classic --datadir <datadir> --netrestrict 10.0.0.0/8,192.168.0.0/16
```

`--maxpeers` caps the number of peers, 50 by default. `--maxpeers 0` does not close the
peer-to-peer port, which still listens.

## Clock

Peer discovery depends on the clock. Every discovery packet carries an expiry time, set a fixed
interval ahead of the sender's clock, and a node drops any packet whose expiry time has already
passed by its own clock. When two nodes' clocks differ by more than that interval, the node whose
clock is behind sends packets the other drops, so the request and reply that discovery needs
between them never complete.

Keep the clock synchronized with NTP. On a host that runs systemd, this reports whether it is:

```shell
$ timedatectl show -p NTPSynchronized
NTPSynchronized=yes
```

## Reporting a vulnerability

To report a vulnerability, follow
[`SECURITY.md`](https://github.com/ethereumclassic/core-geth/blob/main/SECURITY.md).
`geth version-check` checks go-ethereum's advisory feed, which does not track this client, so its
result says nothing about this client's vulnerabilities.
