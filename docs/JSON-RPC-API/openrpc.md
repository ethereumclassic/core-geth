---
title: OpenRPC Discovery
description: "OpenRPC service discovery in Core-Geth: the rpc.discover method that returns a machine-readable description of the JSON-RPC API over HTTP, WebSocket, IPC and console."
hide:
  - toc        # Hide table of contents
---

# OpenRPC

!!! tldr "TLDR"
    The `rpc.discover` method returns an API service description structured per the [OpenRPC specification](https://spec.open-rpc.org/).

## Discovery

Core-Geth supports [OpenRPC's Service Discovery method](https://spec.open-rpc.org/#service-discovery-method), enabling efficient and well-spec'd JSON RPC interfacing and tooling. This method follows the established JSON RPC patterns, and is accessible via HTTP, WebSocket, IPC, and console servers. It answers as `rpc.discover` or `rpc_discover`. To use this method over HTTP, on a node started with `--http`:

!!! Example

    ```shell
    $ curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"rpc_discover","params":[],"id":1}' http://localhost:8545
    ```

    An excerpt of the response from a v1.13 node, formatted for reading, with `[...]` marking
    what is left out:

    ```
    {
      "jsonrpc": "2.0",
      "id": 1,
      "result": {
        "openrpc": "1.2.6",
        "info": {
          "title": "Core-Geth RPC API",
    [...]
        "methods": [
          {
            "name": "eth_accounts",
    [...]
            "summary": "Accounts returns the collection of accounts this node manages.\n",
            "paramStructure": "by-position",
            "params": [],
    [...]
    ```

!!! Tip "The namespace pages in this section"

    The namespace pages in this section are generated from the client's discovery document.
