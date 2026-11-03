---
hide:
  - toc        # Hide table of contents
description: "Use the JSON-RPC, WebSocket and IPC interfaces of a Core-Geth Ethereum Classic node: the namespaces each one serves and how to enable them."
---

# Using JSON-RPC APIs

## Programmatically interfacing `geth` nodes

`geth` serves JSON-RPC ([standard APIs](https://ethereum.org/developers/docs/apis/json-rpc/) and `geth`-specific ones;
most namespaces have a page of their own in this section) over HTTP, WebSocket and IPC (a UNIX socket, or a named pipe on
Windows).

IPC is on by default. On Ethereum Classic and Mordor it serves `admin`,
`debug`, `engine`, `eth`, `ethash`, `miner`, `net`, `rpc`, `trace`, `txpool` and `web3`, and
adds the deprecated `personal` namespace only with `--rpc.enabledeprecatedpersonal`. HTTP and WebSocket must be turned on, and serve only the namespaces named in `--http.api` and `--ws.api`, plus `rpc`.
An empty list is not an empty set. `--http.api ""` and `--ws.api ""` stop the node at startup,
because an empty list would serve every namespace except `engine`, `admin`, `debug` and
`personal` among them. An empty `HTTPModules` or `WSModules` in a configuration file still does
serve them, so name the namespaces there too.
**`--rpc.enabledeprecatedpersonal` does not gate HTTP or WS:** naming `personal` in
`--http.api` or `--ws.api` serves it there whether or not the flag is set.

HTTP based JSON-RPC API options:

  * `--http` Enable the HTTP-RPC server
  * `--http.addr` HTTP-RPC server listening interface (default: `localhost`)
  * `--http.port` HTTP-RPC server listening port ([default](../operate/security.md#ports-and-listeners))
  * `--http.api` API's offered over the HTTP-RPC interface (default: `eth,net,web3`)
  * `--http.corsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
  * `--ws` Enable the WS-RPC server
  * `--ws.addr` WS-RPC server listening interface (default: `localhost`)
  * `--ws.port` WS-RPC server listening port ([default](../operate/security.md#ports-and-listeners))
  * `--ws.api` API's offered over the WS-RPC interface (default: `eth,net,web3`)
  * `--ws.origins` Origins from which to accept websockets requests
  * `--graphql` Enable GraphQL on the HTTP-RPC server. Note that GraphQL can only be started if an HTTP server is started as well.
  * `--graphql.corsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
  * `--graphql.vhosts` Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (default: "localhost")
  * `--ipcdisable` Disable the IPC-RPC server
  * `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

All three speak [JSON-RPC](https://www.jsonrpc.org/specification), and one connection can carry many requests.
[OpenRPC discovery](openrpc.md) shows how to ask a running node for its OpenRPC document. The
document can list methods an interface does not serve, such as `personal_*` over IPC without
`--rpc.enabledeprecatedpersonal` and the `trace_filter` subscription, so `rpc_modules` is the check
for which namespaces an interface serves.

!!! Attention

        Read [RPC exposure](../operate/security.md#rpc-exposure) before you turn on HTTP or WebSocket. Attackers scan the
    internet for nodes with exposed APIs, and any browser tab can reach a web server on the same machine.
    [Public RPC endpoint](../guides/public-rpc-endpoint.md) covers serving JSON-RPC to other people.
