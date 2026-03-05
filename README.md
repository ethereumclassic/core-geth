## CoreGeth: Ethereum Classic Execution Client

> A [go-ethereum](https://github.com/ethereum/go-ethereum) fork providing the production Ethereum Classic (ETC) execution client.

CoreGeth is the current production client for the Ethereum Classic network. It supports all ETC hard forks from Frontier through Spiral, and implements the upcoming **Olympia** hard fork (ECIP-1111, ECIP-1112, ECIP-1121).

**Note:** Upstream go-ethereum has deprecated support for Ethereum Classic. Long-term, ETC will migrate to [Fukuii](https://github.com/chris-mercer/fukuii) as the native ETC client. CoreGeth remains the recommended production client during this transition.

## Supported Networks

| Network | Chain ID | Consensus | Status |
|---------|----------|-----------|--------|
| Ethereum Classic (ETC) | 61 | Proof of Work (ETChash) | Production |
| Mordor Testnet | 63 | Proof of Work (ETChash) | Active |
| Private chains | configurable | PoW / PoA | Supported |

### ETC Hard Fork History

| Fork | Block | ECIPs |
|------|-------|-------|
| Atlantis | 8,772,000 | ECIP-1054 (Byzantium equivalent) |
| Agharta | 9,573,000 | ECIP-1056 (Constantinople equivalent) |
| Phoenix | 10,500,839 | ECIP-1088 (Istanbul equivalent) |
| Thanos | 11,700,000 | ECIP-1099 (ETChash 60K epochs) |
| Magneto | 13,189,133 | ECIP-1103 (Berlin equivalent) |
| Mystique | 14,525,000 | ECIP-1104 (partial London, no EIP-1559) |
| Spiral | 19,250,000 | ECIP-1109 (Shanghai equivalent) |
| **Olympia** | **TBD** | ECIP-1111 + ECIP-1121 (EIP-1559 + treasury + EVM modernization) |

**Gas limit:** 8M (current) converging to 60M post-Olympia (EIP-7935).

## Olympia Hard Fork

The `olympia` branch implements the Olympia upgrade:

- **ECIP-1111:** EIP-1559 dynamic base fee + EIP-3198 BASEFEE opcode. Base fee is redirected to the Olympia Treasury (not burned).
- **ECIP-1112:** Deterministic, immutable Treasury contract receiving all base fee revenue.
- **ECIP-1121:** 13 execution-layer EIPs for EVM modernization (EIP-5656/MCOPY, EIP-1153/transient storage, EIP-6780/SELFDESTRUCT, EIP-2537/BLS12-381, EIP-7951/P256VERIFY, EIP-7702/EOA code delegation, EIP-2935/block hashes in state, and more).

**EIP-7642 (eth/69) is intentionally excluded** — it removes Total Difficulty from the protocol handshake, which ETC requires for Proof-of-Work chain selection.

### Activation

- **Mordor testnet:** Block 15,800,850 (~March 28, 2026)
- **ETC mainnet:** ~Block 24,751,337 (~mid-June 2026, before Era 6)

### Cross-Client Alignment

| Client | Pre-Olympia | Post-Olympia | Role |
|--------|-------------|--------------|------|
| [core-geth](https://github.com/chris-mercer/core-geth) | `etc` branch | `olympia` branch | Production client |
| [Fukuii](https://github.com/chris-mercer/fukuii) | `alpha` branch | `olympia` branch | Native ETC client (migration target) |
| [Besu](https://github.com/chris-mercer/besu) | `etc` branch | `olympia` branch | Reference/testing client |

## Build

```bash
make geth
```

## Test

```bash
# ETC-specific unit tests
go test ./params/... -run TestETC -v
go test ./core/... -run "TestGasLimit|TestForkCompliance|TestECIP1017|TestTreasury|TestOlympia" -v
go test ./consensus/ethash/... -run "TestDifficultyETC|TestDifficultyECIP" -v

# Live tests (requires running Mordor/ETC node)
go test -tags live ./tests/live_etc/ -v
```

## Mining

CoreGeth supports full Ethash/ETChash Proof-of-Work mining:

```bash
./build/bin/geth --classic --mine --miner.etherbase <address>
```

For testing with fake PoW (no DAG generation):

```bash
./build/bin/geth --classic --mine --miner.etherbase <address> --fakepow
```

## Documentation

- [ETC-HANDOFF.md](./ETC-HANDOFF.md) — Pre-Olympia branch documentation
- [OLYMPIA-HANDOFF.md](./OLYMPIA-HANDOFF.md) — Olympia hard fork documentation
- [CoreGeth docs](https://etclabscore.github.io/core-geth) — General documentation
- [go-ethereum docs](https://geth.ethereum.org/docs/) — Upstream reference

## Security Patches

The `etc` branch includes 5 security patches applied before any feature work:

| CVE | Severity | Component |
|-----|----------|-----------|
| CVE-2026-26314 | High | secp256k1 coordinate validation |
| CVE-2026-26315 | Moderate | ECIES invalid-curve attack |
| CVE-2026-22862 | High | ECIES decrypt length |
| CVE-2025-24883 | Moderate | UnmarshalPubkey curve check |
| x/crypto, x/net | — | Dependency updates |

Nodes should rotate P2P keys after upgrading: `rm <datadir>/geth/nodekey`

## License

The core-geth library (outside `cmd/`) is licensed under [LGPL-3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html).
The core-geth binaries (`cmd/`) are licensed under [GPL-3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).
