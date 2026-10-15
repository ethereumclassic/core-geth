---
title: About
description: "About Core-Geth, an Ethereum Classic execution client: a go-ethereum downstream that keeps chain configuration data-driven, and the networks it supports."
---

# Core-Geth

Core-Geth is an Ethereum Classic execution client. It is a downstream of
[ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) that keeps chain configuration as
data rather than code, so one `geth` binary runs the networks in the
[Supported networks](../index.md#supported-networks) table, private networks among them.
Core-Geth is developed and released from [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth),
the canonical home for the client in the Ethereum Classic community organization.
[Project history](../about/project-history.md) tells where the project came from.

This page lists what Core-Geth adds to go-ethereum, how its chain configuration is built, and what
go-ethereum has that Core-Geth does not.

## Additional features

### JSON-RPC API

- **OpenRPC service discovery.** The `rpc.discover` method returns an
  [OpenRPC](https://open-rpc.org/) document describing the client's methods, including some the
  interface it answers on does not serve, with their parameters and results. See
  [OpenRPC discovery](../JSON-RPC-API/openrpc.md).
- **The `trace` module.** `trace_block`, `trace_transaction`, `trace_call` and `trace_callMany`
  return trace entries in the format of OpenEthereum's trace module, and `trace_filter` works
  only as a subscription. See the
  [trace module overview](../JSON-RPC-API/trace-module-overview.md) and the
  [`trace` method reference](../JSON-RPC-API/modules/trace.md).

The [JSON-RPC API overview](../JSON-RPC-API/index.md) covers the rest of the API.

### External EVMs through EVMC

Core-Geth implements version 7 of the [EVMC](https://github.com/ipsilon/evmc) connector API as an
experimental feature. `--vm.evm` and `--vm.ewasm` load an external EVM or ewasm interpreter from a
shared library, and neither is supported on Ethereum Classic or Mordor: see
[Running Geth with an External VM](evmc.md).

### Development and testing

- **`--dev.pow`** starts an ephemeral proof-of-work network with a pre-funded developer account. It
  switches mining on with no CPU threads, so it seals no block unless `--miner.threads` is above 0;
  `--fakepow.poisson --miner.threads 2` seals blocks without generating an Ethash DAG.
- **`--fakepow.poisson`** disables proof-of-work verification and adds a Poisson mining delay based on
  `--miner.threads`.
- **`make test-coregeth`** runs the tests specific to Core-Geth, including imports of simulated
  canonical chains. See [Testing](../developers/testing.md).

### Ethereum Classic proposals

Ethereum Classic's own proposals, ECIPs and ECBPs, are chain configuration fields of their own,
alongside the EIPs. Among them:

- ECIP-1010, the difficulty bomb pause: `ECIP1010PauseBlock` and `ECIP1010Length`
- ECIP-1017, the monetary policy: `ECIP1017FBlock` and `ECIP1017EraRounds`
- ECIP-1041, the difficulty bomb disposal: `DisposalBlock`
- ECIP-1099, Etchash: `ECIP1099FBlock`
- ECBP-1100, MESS (Modified Exponential Subjective Scoring): `ECBP1100FBlock`

`CoreGethChainConfig` in `params/types/coregeth/chain_config.go` defines the fields, and
`ClassicChainConfig` in `params/config_classic.go` sets Ethereum Classic's activation blocks.

## Divergent design

### Chain configuration by feature

At the code level, Core-Geth and go-ethereum differ in how code asks the chain configuration whether
a feature is active.

In go-ethereum, it asks about a named upgrade:

```go
blockNumber := big.NewInt(0)
config := params.MainnetChainConfig
if config.IsByzantium(blockNumber) {
	// do a special thing for post-Byzantium chains
}
```

Byzantium is a group of EIPs, and the check does not say which of them the code depends on.
Core-Geth asks about the EIP itself:

```go
blockNumber := big.NewInt(0)
config := params.MainnetChainConfig
if config.IsEnabled(config.GetEIP658Transition, blockNumber) {
	// do a special thing for post-EIP658 chains
}
```

The `ChainConfigurator` interface in
[`params/types/ctypes/configurator_iface.go`](https://github.com/ethereumclassic/core-geth/blob/main/params/types/ctypes/configurator_iface.go)
carries a transition getter and setter for each supported feature, in `ProtocolSpecifier`, and
`IsEnabled`, in `Forker`.

Because each EIP is a setting of its own, a chain can adopt some of an upgrade's EIPs and not others.
Ethereum Classic runs EIP-658, which embeds a transaction status code in receipts, without EIP-649,
Byzantium's difficulty bomb delay and block reward reduction. `TestClassicIs649` in
`params/config_test.go` asserts that `ClassicChainConfig` has no EIP-649 transition.

!!! note "Genesis file formats"
    `geth init` reads a genesis file in Core-Geth's format, as
    [Adding a network](../developers/add-network.md) shows, or in go-ethereum's. An OpenEthereum
    chain specification is rejected with `invalid configurator schema`.

## Limitations

- **Release builds keep source file paths.** go-ethereum builds with `-trimpath`; Core-Geth's
  `build/ci.go` does not, and its comment gives the reason: the OpenRPC discovery document relies on
  reflection and source parsing that break when paths are trimmed. A `geth` binary therefore contains
  the source paths of the machine that built it.
- **Some go-ethereum features are absent**, among them the `--vmtrace` flag, which records internal
  VM operations with a named tracer, and the `blsync` beacon light syncer.
- **Ethereum's networks are not supported.** See
  [Networks this client does not support](../index.md#networks-this-client-does-not-support).
- **`geth version-check` checks go-ethereum's vulnerability feed, not this client's.** See
  [`SECURITY.md`](https://github.com/ethereumclassic/core-geth/blob/main/SECURITY.md).
