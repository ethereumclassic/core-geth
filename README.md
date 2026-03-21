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

### ETC Hard Fork History

| Fork | ETC Mainnet Block | Mordor Block | ECIPs |
|------|-------------------|--------------|-------|
| Atlantis | 8,772,000 | 252,500 | ECIP-1054 (Byzantium equivalent) |
| Agharta | 9,573,000 | 301,243 | ECIP-1056 (Constantinople equivalent) |
| Phoenix | 10,500,839 | 999,983 | ECIP-1088 (Istanbul equivalent) |
| Thanos | 11,700,000 | 2,520,000 | ECIP-1099 (ETChash, 60K epoch length) |
| Magneto | 13,189,133 | 3,985,893 | ECIP-1103 (Berlin equivalent) |
| Mystique | 14,525,000 | 5,520,000 | ECIP-1104 (partial London, no EIP-1559) |
| Spiral | 19,250,000 | 9,957,000 | ECIP-1109 (Shanghai equivalent, PUSH0) |

## Build

Requires Go 1.21+ (1.24+ recommended).

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
