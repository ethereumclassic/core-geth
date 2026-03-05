# ETC Branch Handoff — core-geth

**Branch:** `etc` (on `chris-mercer/core-geth`)
**Base:** `master` (v1.12.20)
**Date:** 2026-03-03
**Status:** Pre-Olympia stabilization COMPLETE

## Purpose

The `etc` branch brings core-geth up to date with present-day ETC/Mordor chain rules, fixes 4 CVEs, and adds comprehensive pre-Olympia testing. It serves as the stable base for the `olympia` branch (hard fork features).

## Branch Topology

```
master (v1.12.20)
  └── etc (14 commits: 5 security + 9 test)
       └── olympia (25 hard fork commits)
```

## Commits (14 total)

### Security Patches (5)
| Commit | Description |
|--------|-------------|
| `c62b31054` | CVE-2025-24883: IsOnCurve check in UnmarshalPubkey |
| `9dc4fadf3` | CVE-2026-22862: ECIES decrypt length check |
| `c9fc287b8` | CVE-2026-26315: ECIES GenerateShared validation |
| `21f515da2` | CVE-2026-26314: secp256k1 coordinate validation |
| `ac56a1cc9` | golang.org/x/crypto + x/net dependency updates |

### Test Suite (9)
| Commit | File(s) | Tests | Coverage |
|--------|---------|-------|----------|
| `bfdf5845f` | `params/config_etc_test.go` | 13 | Chain config validation (fork ordering, chain IDs) |
| `ba3def965` | `core/gas_limit_test.go` | 8 | Gas limit rules (1/1024 adjustment, 8M target) |
| `52765042f` | `core/etc_fork_compliance_test.go` | 10 | Fork compliance + EVM opcodes per fork |
| `35be1a1b7` | `core/ecip1017_test.go` | 4 | ECIP-1017 emission (era transitions, miner rewards) |
| `c0a7447ab` | `consensus/ethash/difficulty_etc_test.go` | 6+subtests | ETChash difficulty (ECIP-1010/1041/1099, min floor) |
| `27a6e068b` | `tests/live_etc/` (3 files) | 13 | Live RPC for Mordor + ETC mainnet |
| `31175757c` | `core/vm/contracts_etc_test.go` | 4+subtests | Precompile activation per fork, gas repricing |
| `bf6a802c3` | `tests/live_etc/` (extend) | 9 | ECBP-1100 deactivation, ECIP-1099 epoch, Spiral fork |
| `a0b4f5959` | `tests/live_etc/testdata/` | — | Shared ETC consensus vectors (JSON) |

## Test Coverage Summary

| Area | Tests | Status |
|------|-------|--------|
| Chain config (Classic + Mordor) | 13 | PASS |
| Gas limit rules | 8 | PASS |
| Fork compliance + EVM opcodes | 10 | PASS |
| ECIP-1017 emission schedule | 4 | PASS |
| ETChash difficulty (ECIP-1010/1041/1099) | 22 (6 + subtests) | PASS |
| ECBP-1100 (MESS) artificial finality | 18+ | PASS (existing in master) |
| Precompile activation per fork | 10 (4 + subtests) | PASS |
| Live RPC (Mordor) | 11 | Requires node |
| Live RPC (ETC mainnet) | 11 | Requires node or public RPC |
| Shared test vectors | 30+ vectors | JSON reference |
| **Total new tests** | **~67** | — |

## Cross-Client Alignment

| Test Area | core-geth `etc` | besu `etc` | fukuii `alpha` |
|-----------|-----------------|------------|----------------|
| Chain config validation | 13 | 20+ | 13 |
| Gas limit rules | 8 | implicit | 11 |
| Fork compliance | 10 | 15+ | 14 |
| ECIP-1017 emission | 4 | 27 | 4 |
| Difficulty (ECIP-1010/1041/1099) | 22 | 24+ | 29+ |
| ECBP-1100 (MESS) polynomial | 18+ | 14+ | 25 |
| Precompile per fork | 10 | — | — |
| Live RPC | 22 | 28+ | 16 |

### ECIP-1100 MESS Alignment (All 3 Clients)

All three ETC clients implement the ECIP-1100 cubic polynomial antigravity curve identically:

```
polynomialV(x) = DENOMINATOR + (3x² − 2x³/xcap) × HEIGHT / xcap²
where DENOMINATOR=128, xcap=25132, HEIGHT=3840
```

- **core-geth:** `core/blockchain_af.go` — reference implementation
- **Besu:** `ArtificialFinality.java` — matches core-geth
- **Fukuii:** `ArtificialFinality.scala` — rewritten (March 2026) from exponential decay to match spec

MESS is deactivated at Spiral (ECIP-1110): ETC mainnet block 19,250,000, Mordor block 10,400,000.

### EIP-7642 Exclusion (Intentional)

EIP-7642 (eth/69 — removes TD from the `Status` handshake message) is deliberately excluded from all three ETC clients. ETC uses Proof of Work and requires Total Difficulty for chain selection. Removing TD from `Status` would break PoW peer negotiation. This will be removed from the ECIP-1121 draft before finalization.

### Treasury Address Note

The Mordor treasury address (`0xCfE1e0ECbff745e6c800fF980178a8dDEf94bEe2`) is a placeholder in Chris Mercer's development branches. Once merged into the official `ethereumclassic` organization, the address will be updated through proper coordination with the core development team.

### Olympia Activation Timeline

- **Mordor testnet:** block 15,800,850 (~March 28, 2026)
- **ETC mainnet:** block ~24,751,337 (~mid-June 2026)

## How to Run

```bash
# Unit tests
go test ./params/... -run TestETC -v
go test ./core/... -run "TestGasLimit|TestForkCompliance|TestECIP1017" -v
go test ./consensus/ethash/... -run "TestDifficultyETC|TestDifficultyECIP" -v
go test ./core/vm/... -run TestETC -v

# Security patches regression
go test ./crypto/... -v

# Live tests (requires running Mordor node at localhost:8545)
go test -tags live ./tests/live_etc/ -v

# Live tests with custom RPC
MORDOR_RPC=http://localhost:8545 ETC_RPC=https://etc.rivet.link go test -tags live ./tests/live_etc/ -v

# Build
make geth
```

## Key Files

| File | Purpose |
|------|---------|
| `params/config_etc_test.go` | Chain config validation |
| `core/gas_limit_test.go` | Gas limit rules |
| `core/etc_fork_compliance_test.go` | Fork compliance + EVM opcodes |
| `core/ecip1017_test.go` | ECIP-1017 emission |
| `consensus/ethash/difficulty_etc_test.go` | ETChash difficulty |
| `core/vm/contracts_etc_test.go` | Precompile activation |
| `tests/live_etc/helpers_test.go` | Live RPC helpers + constants |
| `tests/live_etc/mordor_test.go` | Mordor live tests |
| `tests/live_etc/mainnet_test.go` | ETC mainnet live tests |
| `tests/live_etc/testdata/etc_consensus_vectors.json` | Shared test vectors |

## Next Steps

1. Push `etc` branch to origin
2. Merge etc additions forward into `olympia` branch
3. Deploy Mordor node with etc branch binary for live testing
4. Set mainnet activation after Mordor soak period
