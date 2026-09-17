---
title: Public RPC endpoint
description: "Serving Ethereum Classic JSON-RPC to other people from your own Core-Geth node: what the client does itself and what has to sit in front of it."
---

**Looking for an endpoint to use rather than one to run?**
[Public JSON-RPC](../etc-cooperative-transition.md#public-json-rpc) lists what to point an application at today.

A node that answers JSON-RPC for other people lets wallets, explorers and applications use Ethereum
Classic without running a node of their own, and each endpoint run by an independent operator is one
more way to reach the network. This page covers what core-geth does for such an endpoint by itself,
and what has to sit in front of it.

It starts from a synced node run as a service, as
[Run an Ethereum Classic node](../getting-started/run-classic-node.md) sets one up, with the
firewall from [Security and network exposure](../operate/security.md#a-host-firewall) in place.

## Before you start

- **Keep account keys off this node.** Any caller that reaches the `eth` namespace can use an
  unlocked account ([Account keys](../operate/security.md#account-keys)), and serving JSON-RPC needs
  no keystore.
- **Only the peer-to-peer ports face the internet.** The JSON-RPC listener stays on `127.0.0.1` or a
  private network, and other machines reach it through a proxy
  ([Ports and listeners](../operate/security.md#ports-and-listeners)).

## Choose the namespaces

Serve `eth`, `net` and `web3`, and name them on the command line:

```shell
$ geth --classic --datadir <datadir> --http --http.api eth,net,web3
```

`rpc_modules` lists the namespaces an interface serves:

```shell
$ curl -s -X POST -H 'Content-Type: application/json' \
    --data '{"jsonrpc":"2.0","id":1,"method":"rpc_modules","params":[]}' \
    http://127.0.0.1:8545
{"jsonrpc":"2.0","id":1,"result":{"eth":"1.0","net":"1.0","rpc":"1.0","web3":"1.0"}}
```

A method from any other namespace is refused:

```
{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"the method admin_nodeInfo does not exist/is not available"}}
```

[Using JSON-RPC APIs](../JSON-RPC-API/index.md) lists what each interface serves when `--http.api`
or `--ws.api` is left out, and why an empty list is refused. `--ws.api` takes the same names as
`--http.api`.

| Namespace | What it holds | Serve it to others |
| --- | --- | --- |
| `eth` | Blocks, transactions and state; calls and gas estimates; filters and subscriptions; transaction submission | Yes |
| `net` | The network ID, the peer count, and whether the node is listening | Yes |
| `web3` | The client version, and a hash function | Yes |
| `txpool` | The transactions waiting in this node's pool | Optional. [`txpool_content`](../JSON-RPC-API/modules/txpool.md#txpool_content) returns the whole pool in one response |
| `trace` | Traces of blocks and transactions ([the `trace` module](../JSON-RPC-API/trace-module-overview.md)) | Only through a gateway that limits it, since every trace executes the transactions again |
| `debug` | Tracing, and methods that change the node: [`debug_setHead`](../JSON-RPC-API/modules/debug.md#debug_sethead) rewinds its chain | No |
| `admin` | The node's peers and its RPC listeners: [`admin_startHTTP`](../JSON-RPC-API/modules/admin.md#admin_starthttp) opens a listener with any namespaces | No |
| `miner` | Starting and stopping mining, and the mining account | No |
| `ethash` | Mining work, for miners | No |
| `personal` | Deprecated account methods | No |
| `engine` | The Engine API | Cannot be served: named in `--http.api`, it is left out |

## Listening address, hostnames and origins

`--http.addr` and `--ws.addr` default to `localhost`. Keep that default when the proxy runs on the
same machine; when it runs on another, use the address of a private interface the proxy can reach.
In a container, [Docker](../getting-started/installation.md#docker) shows the address to use, and
how publishing the container's port can expose it on every interface of the host.

### Hostnames

The node compares the hostname in each HTTP request's `Host` header with `--http.vhosts`, which
holds `localhost` by default, and refuses any other hostname with status 403:

```
invalid host specified
```

A proxy that passes the public hostname on to the node needs that hostname in the list. The list
then replaces `localhost`:

```shell
$ geth --classic --datadir <datadir> --http --http.api eth,net,web3 --http.vhosts rpc.example.org
```

```shell
$ curl -s -X POST -H 'Content-Type: application/json' -H 'Host: rpc.example.org' \
    --data '{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}' \
    http://127.0.0.1:8545
{"jsonrpc":"2.0","id":1,"result":"0x3d"}
```

The check looks at hostnames only. A request whose `Host` header is an IP address is answered
whatever the list holds, and WebSocket connections are not checked against it. Deciding who may
connect is the job of the firewall and the proxy.

### Browser origins

`--http.corsdomain` names the web origins whose pages a browser lets read the node's answers. The
node answers a request whatever its `Origin` header says. For an origin in the list, it adds an
`Access-Control-Allow-Origin` header naming that origin, and for any other origin it sends the same
answer without the header, which a browser then withholds from the page.

The WebSocket server checks the `Origin` header too. It refuses an origin that `--ws.origins` does
not list, with status 403, and accepts a connection that sends no `Origin` header. Without the flag
it accepts only `http://localhost` and `http://` with the machine's hostname, on any port, and the
flag replaces both.

### Paths

`--http.rpcprefix /rpc` serves JSON-RPC at `/rpc` and answers other paths with status 404, for a
proxy that routes by path. `--ws.rpcprefix /ws` accepts WebSocket connections only at `/ws`.

## Limits the node enforces

| Flag | What it limits | What a caller gets past the limit |
| --- | --- | --- |
| `--rpc.gascap` | Gas for `eth_call` and `eth_estimateGas` | A call that asks for more gas runs with the cap, and the node logs `Caller gas above allowance, capping`. An estimate that needs more fails with `gas required exceeds allowance` |
| `--rpc.evmtimeout` | How long an `eth_call` may run | `execution aborted (timeout = <duration>)` |
| `--rpc.txfeecap` | The fee of a transaction sent through the node | `tx fee (<fee> ether) exceeds the configured cap (<cap> ether)` |
| `--rpc.batch-request-limit` | The number of calls in one batch | `batch too large` |
| `--rpc.batch-response-max-size` | The bytes of results one batch returns | The call that crosses the limit is answered, and every call after it gets `response too large` |
| `--rpc.allow-unprotected-txs` | Transactions without replay protection | Unset, the node refuses them with `only replay-protected (EIP-155) transactions allowed over RPC`. Leave it unset |

[Command-line Options](../getting-started/run-cli.md#command-line-options) gives each default.

The node sets no limit on the block range of an `eth_getLogs` query.

What a node can answer about the past depends on what it keeps. A lookup by transaction hash finds
only transactions in the [transaction index](../operate/sync-modes.md#transaction-index), and the
state and traces of older blocks need an [archive node](archive-node.md).

## What to put in front of the node

The node has no TLS, no authentication and no rate limiting for JSON-RPC over HTTP or WebSocket.
A reverse proxy or an RPC gateway in front of it supplies them:

- TLS for connections from other machines;
- authentication, such as API keys, for an endpoint that is not open to everyone;
- limits for each client on request rate and request size;
- a list of allowed methods, where a namespace mixes methods you offer with methods you do not;
- a limit on the block range of `eth_getLogs`.

GraphQL, turned on by `--graphql`, answers at `/graphql` on the HTTP listener and goes behind the
same proxy. v1.13.0 refuses a GraphQL query that nests fields more than 20 deep, the fix for the
[GraphQL query depth](../audits/2026-03-security-audit.md#graphql-query-depth-dos) finding in the
March 2026 audit.

## What a gateway checks

A gateway or load balancer checks a node with ordinary JSON-RPC calls. These are the calls, and what
a node of this client answers:

| What is checked | Call | Ethereum Classic | Mordor |
| --- | --- | --- | --- |
| Chain ID | `eth_chainId` | `"0x3d"` | `"0x3f"` |
| Network ID | `net_version` | `"1"` | `"7"` |
| Sync | `eth_syncing` | `false` once synced, and an object while syncing | The same |
| Peers | `net_peerCount` | The number of peers, in hexadecimal | The same |
| Finality | `eth_getBlockByNumber` with `"finalized"` or `"safe"` | An error | The same |

A gateway serving this node has to judge finality by block depth, since both finality tags fail
([Why does a request for the `finalized` or `safe` block fail?](../operate/troubleshooting.md#why-does-a-request-for-the-finalized-or-safe-block-fail)).
How it does that, and whether its chain registry includes the network you serve, are in that
gateway's documentation.

A plain HTTP `GET` with no body gets status 200 from any node started with `--http`, synced or not,
with peers or without. A load balancer check built on it shows only that the process is listening.
[Health checks](../operate/monitoring.md#health-checks) defines a ready node; check the same three
conditions with `eth_syncing`, `net_peerCount` and `eth_getBlockByNumber`.

## Several nodes behind one entry point

Some requests depend on state that one node, or one connection, holds:

- **A filter exists only on the node that created it.**
  [`eth_newFilter`](../JSON-RPC-API/modules/eth.md#eth_newfilter),
  [`eth_newBlockFilter`](../JSON-RPC-API/modules/eth.md#eth_newblockfilter) and
  `eth_newPendingTransactionFilter` return an ID that only that node knows. Asked for the same ID,
  another node answers:

    ```
    {"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"filter not found"}}
    ```

    Send each client's filter calls to the node that created its filters.

- **A subscription belongs to its connection.**
  [`eth_subscribe`](../JSON-RPC-API/modules/eth.md#eth_subscribe) works over WebSocket and not over
  HTTP, where it answers `notifications not supported`. A second connection, even to the same node,
  cannot cancel the subscription and gets `subscription not found`. Keep each client's WebSocket
  connection on one node.
- **Each node counts pending transactions in its own pool.**
  [`eth_getTransactionCount`](../JSON-RPC-API/modules/eth.md#eth_gettransactioncount) with
  `"pending"` counts only the transactions that node has received. A node that has not yet received
  a sender's newest transaction returns a lower count, and a wallet that uses that count signs its
  next transaction with a nonce already taken. Send each sender's requests, from the nonce lookup to
  the transaction, to one node.
