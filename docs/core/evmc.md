---
hide:
  - toc        # Hide table of contents
description: "Run Core-Geth with an external EVM via EVMC: connecting a third-party interpreter such as evmone or Hera to an Ethereum Classic node."
---

# Running Geth with an External VM

Geth supports the __[EVMC](https://github.com/ipsilon/evmc) VM connector API Version 7__ as an experimental feature. This interface provides support for external EVM and EWASM interpreters.

!!! danger "`--vm.evm` and `--vm.ewasm` are not supported on Ethereum Classic or Mordor"
    `--vm.evm` replaces the built-in interpreter for ordinary EVM bytecode. `getRevision`
    (shown below) returns no EVMC revision past Istanbul, and panics once EIP-2565 is
    active: on Ethereum Classic at block 13,189,133
    (`ClassicChainConfig.EIP2565FBlock`, `params/config_classic.go`) and on Mordor at
    block 3,985,893 (`MordorChainConfig.EIP2565FBlock`, `params/config_mordor.go`). A
    node already past that block panics on its first contract execution; a node syncing
    from genesis panics on reaching it.

    `--vm.ewasm` is also unsupported on both networks: EWASM is not part of either
    network's consensus rules, and ordinary EVM bytecode does not execute through it.

External interpreters can be configured on the command line via
a `--vm.`-prefixed flag for normal instantiation, and `--evmc.` for testing.

```
  --vm.evm value                      External EVM configuration (default = built-in interpreter)
  --vm.ewasm value                    External ewasm configuration (default = built-in interpreter)
```

Provided to these flags should be EWASM and EVM shared object libraries, as follows:

- `--vm.ewasm=<path/to/interpreter.so`
- `--vm.evm=<path/to/interpreter.so`

## Testing EVMC Support

This implementation may be tested by following the command defined in the Makefile as `test-evmc`, which
tests the implementation against the EWASM and EVM libraries mentioned above, against the `/tests/` StateTest suite.

These tests run exclusively via Github Actions, configured at `.github/workflows/evmc.yml`.

## Discussion: Customizing EVMC Configuration

While Core-Geth supports highly granular EIP/ECIP/xIP chain feature configuration (ie fork feature configs),
EVMC does not. EVMC only supports the Fork configurations supported by ethereum/go-ethereum (eg. Byzantium, Constantinople, &c).
Thus, the implementation at Core-Geth of EVMC requires a somewhat arbitrary mapping of granular features as keys toggling
entire Ethereum fork configurations.

The following code snippet, taken from [`./core/vm/evmc.go`](https://github.com/ethereumclassic/core-geth/blob/main/core/vm/evmc.go#L339-L367), handles this translation.

```go
// getRevision translates ChainConfig's HF block information into EVMC revision.
func getRevision(env *EVM) evmc.Revision {
	n := env.Context.BlockNumber
	conf := env.ChainConfig()
	switch {
	// This is an example of choosing to use an "abstracted" idea
	// about chain config, where I'm choosing to prioritize "indicative" features
	// as identifiers for Fork-Feature-Groups. Note that this is very different
	// than using Feature-complete sets to assert "did Forkage."
	case conf.IsEnabled(conf.GetEIP2565Transition, n):
		panic("berlin is unsupported by EVMCv7")
	case conf.IsEnabled(conf.GetEIP1884Transition, n):
		return evmc.Istanbul
	case conf.IsEnabled(conf.GetEIP1283DisableTransition, n):
		return evmc.Petersburg
	case conf.IsEnabled(conf.GetEIP145Transition, n):
		return evmc.Constantinople
	case conf.IsEnabled(conf.GetEIP198Transition, n):
		return evmc.Byzantium
	case conf.IsEnabled(conf.GetEIP155Transition, n):
		return evmc.SpuriousDragon
	case conf.IsEnabled(conf.GetEIP150Transition, n):
		return evmc.TangerineWhistle
	case conf.IsEnabled(conf.GetEIP7Transition, n):
		return evmc.Homestead
	default:
		return evmc.Frontier
	}
}
```

As you can see, individual features, like EIP1884, are translated as proxy signifiers for entire fork configurations
(in this case, an Istanbul-featured VM revision).
This approach, rather than requiring a complete set of the compositional features for any of these given Ethereum forks,
trades a descriptive 1:1 mapping for application flexibility. Pursuing a necessarily complete feature-set -> fork
map would presume chain features that are not necessarily relevant to the virtual machine, like block reward configurations,
or difficulty configurations, for example. This approach allows applications to use advanced opcodes with the fewest
number of incidental restrictions.

This approach is not without risk or nuance however; without a solid understanding of customizations here,
experiments in customization can result in foot shooting.
