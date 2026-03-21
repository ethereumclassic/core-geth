## CoreGeth: Ethereum Classic Execution Client

> A [go-ethereum](https://github.com/ethereum/go-ethereum) fork providing the primary execution client for the Ethereum Classic (ETC) network.

**Status: Maintenance mode.** CoreGeth is the only production ETC execution client. It is maintained by community contributors on a volunteer basis. There is no dedicated funding or full-time development team. The last official release (v1.12.20) was June 2024. Contributions, sponsorship, and review are welcome.

CoreGeth supports all ETC hard forks from Frontier through Spiral (ECIP-1109), and implements ETC-specific consensus rules including ETChash, ECIP-1017 emission schedule, and ECBP-1100 (MESS) artificial finality.

Upstream go-ethereum development is merged when feasible. Compatible RPC, CLI APIs, data storage schemas, and P2P protocols are maintained.

## Supported Networks

| Network | Chain ID | Flag | Consensus |
|---------|----------|------|-----------|
| Ethereum Classic | 61 | `--classic` | Proof of Work (ETChash) |
| Mordor Testnet | 63 | `--mordor` | Proof of Work (ETChash) |

### ETC Hard Fork History ([ECIP-1066](https://ecips.ethereumclassic.org/ECIPs/ecip-1066))

| Fork | Block | Mordor | Included EIPs / ECIPs | Spec |
|------|-------|--------|----------------------|------|
| Frontier | 1 | 0 | Genesis | — |
| Frontier Thawing | 200,000 | 0 | Ice Age introduction | — |
| Homestead | 1,150,000 | 0 | [EIP-2](https://eips.ethereum.org/EIPS/eip-2), [EIP-7](https://eips.ethereum.org/EIPS/eip-7), [EIP-8](https://eips.ethereum.org/EIPS/eip-8) | [HFM-606](https://eips.ethereum.org/EIPS/eip-606) |
| Gas Reprice | 2,500,000 | 0 | [EIP-150](https://eips.ethereum.org/EIPS/eip-150) | [ECIP-1015](https://ecips.ethereumclassic.org/ECIPs/ecip-1015) |
| Die Hard | 3,000,000 | 0 | [ECIP-1010](https://ecips.ethereumclassic.org/ECIPs/ecip-1010), [EIP-155](https://eips.ethereum.org/EIPS/eip-155), [EIP-160](https://eips.ethereum.org/EIPS/eip-160) | — |
| Gotham | 5,000,000 | 0 | [ECIP-1017](https://ecips.ethereumclassic.org/ECIPs/ecip-1017), [ECIP-1039](https://ecips.ethereumclassic.org/ECIPs/ecip-1039) | — |
| Defuse Difficulty Bomb | 5,900,000 | 0 | [ECIP-1041](https://ecips.ethereumclassic.org/ECIPs/ecip-1041) | — |
| Atlantis | 8,772,000 | 0 | [EIP-100](https://eips.ethereum.org/EIPS/eip-100), [EIP-140](https://eips.ethereum.org/EIPS/eip-140), [EIP-196](https://eips.ethereum.org/EIPS/eip-196), [EIP-197](https://eips.ethereum.org/EIPS/eip-197), [EIP-198](https://eips.ethereum.org/EIPS/eip-198), [EIP-211](https://eips.ethereum.org/EIPS/eip-211), [EIP-214](https://eips.ethereum.org/EIPS/eip-214), [EIP-658](https://eips.ethereum.org/EIPS/eip-658) | [ECIP-1054](https://ecips.ethereumclassic.org/ECIPs/ecip-1054) |
| Agharta | 9,573,000 | 301,243 | [EIP-145](https://eips.ethereum.org/EIPS/eip-145), [EIP-1014](https://eips.ethereum.org/EIPS/eip-1014), [EIP-1052](https://eips.ethereum.org/EIPS/eip-1052) | [ECIP-1056](https://ecips.ethereumclassic.org/ECIPs/ecip-1056) |
| Phoenix | 10,500,839 | 999,983 | [EIP-152](https://eips.ethereum.org/EIPS/eip-152), [EIP-1108](https://eips.ethereum.org/EIPS/eip-1108), [EIP-1344](https://eips.ethereum.org/EIPS/eip-1344), [EIP-1884](https://eips.ethereum.org/EIPS/eip-1884), [EIP-2028](https://eips.ethereum.org/EIPS/eip-2028), [EIP-2200](https://eips.ethereum.org/EIPS/eip-2200) | [ECIP-1088](https://ecips.ethereumclassic.org/ECIPs/ecip-1088) |
| Thanos | 11,700,000 | 2,520,000 | [ECIP-1099](https://ecips.ethereumclassic.org/ECIPs/ecip-1099) (ETChash, 60K epoch length) | — |
| Magneto | 13,189,133 | 3,985,893 | [EIP-2565](https://eips.ethereum.org/EIPS/eip-2565), [EIP-2718](https://eips.ethereum.org/EIPS/eip-2718), [EIP-2929](https://eips.ethereum.org/EIPS/eip-2929), [EIP-2930](https://eips.ethereum.org/EIPS/eip-2930) | [ECIP-1103](https://ecips.ethereumclassic.org/ECIPs/ecip-1103) |
| Mystique | 14,525,000 | 5,520,000 | [EIP-3529](https://eips.ethereum.org/EIPS/eip-3529), [EIP-3541](https://eips.ethereum.org/EIPS/eip-3541) (partial London, no EIP-1559) | [ECIP-1104](https://ecips.ethereumclassic.org/ECIPs/ecip-1104) |
| Spiral | 19,250,000 | 9,957,000 | [EIP-3651](https://eips.ethereum.org/EIPS/eip-3651), [EIP-3855](https://eips.ethereum.org/EIPS/eip-3855), [EIP-3860](https://eips.ethereum.org/EIPS/eip-3860), [EIP-6049](https://eips.ethereum.org/EIPS/eip-6049) (Shanghai equivalent) | [ECIP-1109](https://ecips.ethereumclassic.org/ECIPs/ecip-1109) |

## Build

Requires Go 1.21+.

```bash
make geth
```

The binary is built to `./build/bin/geth`.

## Run

```bash
# ETC Mainnet
./build/bin/geth --classic

# Mordor Testnet
./build/bin/geth --mordor

# With RPC enabled
./build/bin/geth --classic --http --http.api eth,net,web3
```

## Test

```bash
# Core tests
go test ./core/... -count=1 -timeout 10m

# Consensus tests
go test ./consensus/... -count=1 -timeout 5m

# Chain config tests
go test ./params/... -count=1 -timeout 2m

# ETC-specific tests
go test ./params/... -run TestETC -v
go test ./core/... -run "TestGasLimit|TestForkCompliance|TestECIP1017" -v
go test ./consensus/ethash/... -run "TestDifficultyETC|TestDifficultyECIP" -v
go test ./core/vm/... -run TestETC -v

# Live RPC tests (requires running node)
go test -tags live ./tests/live_etc/ -v
```

## Key ETC Files

| File | Purpose |
|------|---------|
| `params/config_classic.go` | ETC mainnet chain config (fork blocks, chain ID) |
| `params/config_mordor.go` | Mordor testnet chain config |
| `consensus/ethash/consensus.go` | ETChash PoW consensus + ECIP-1017 emission |
| `core/blockchain_af.go` | ECBP-1100 (MESS) artificial finality |
| `core/vm/contracts.go` | Precompile registry |

## Docker

```bash
docker build -t coregeth:latest -f Dockerfile .
docker run --rm coregeth:latest version
```

## Documentation

- [CoreGeth docs](https://etclabscore.github.io/core-geth) — Installation, CLI, JSON-RPC API
- [go-ethereum docs](https://geth.ethereum.org/docs/) — Upstream reference
- [ECIP repository](https://github.com/ethereumclassic/ECIPs) — ETC protocol specifications

## Contributing

Contributions are welcome. Please fork, fix, commit, and send a pull request.

- Code must use [gofmt](https://golang.org/cmd/gofmt/) formatting
- Pull requests should target the `master` branch
- Commit messages should be prefixed with the package(s) they modify (e.g., `eth, rpc: make trace configs optional`)
- Run `go test ./...` before submitting

## License

The core-geth library (outside `cmd/`) is licensed under [LGPL-3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html).
The core-geth binaries (`cmd/`) are licensed under [GPL-3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).
