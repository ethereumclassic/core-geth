# Core-Geth — Ethereum Classic Execution Client

**Status:** Production — syncing Mordor testnet, ETC mainnet capable
**Language:** Go 1.24
**Build:** Makefile + build/ci.go orchestration
**License:** LGPL-3.0
**Origin:** Fork of etclabscore/core-geth (itself a go-ethereum fork)
**Branch:** `etc` (pre-Olympia stabilization), `olympia` (hard fork implementation)

---

## Quick Commands

```bash
make geth                    # Build geth binary (→ ./build/bin/geth or ./geth)
make all                     # Build all packages/executables
go test ./core/... -count=1 -timeout 10m    # Core tests
go test ./consensus/... -count=1 -timeout 5m # Consensus tests
go test ./params/... -count=1 -timeout 2m    # Params/config tests
go run build/ci.go install -static ./cmd/geth # Static build (Docker)
```

### Run on Mordor Testnet

```bash
bash run-mordor.sh
```

### Run on ETC Mainnet

```bash
bash run-classic.sh
```

---

## Data Directories

| Network | Path |
|---------|------|
| Mordor | `/media/dev/2tb/data/blockchain/core-geth/mordor/` |
| ETC Mainnet | `/media/dev/2tb/data/blockchain/core-geth/classic/` |

---

## Network Ports

| Port | Protocol |
|------|----------|
| 8545 | HTTP JSON-RPC |
| 8546 | WS JSON-RPC |
| 30303 | P2P (TCP+UDP) |

---

## ETC Support

Core-geth has built-in ETC network support via `--mordor` and `--classic` flags.

Key chain config: `params/config_classic.go` (fork blocks, chain ID, ECBP1100 MESS)

---

## Project Structure

```
cmd/geth/               # CLI entry point
core/                   # Block processing, state, genesis
  vm/                   # EVM (runtime, interpreter, opcodes)
consensus/              # Consensus engines (ethash, clique)
  ethash/               # ETChash PoW (ETC-specific difficulty)
params/                 # Chain configs, protocol versions
  config_classic.go     # ETC mainnet config
  config_mordor.go      # Mordor testnet config
internal/build/         # Build toolchain (gotool.go, ci.go)
p2p/                    # P2P networking, discovery
eth/                    # Wire protocol, sync
accounts/               # Account management
crypto/                 # Cryptographic primitives
```

---

## Key Files for ETC Work

| File | Purpose |
|------|---------|
| `params/config_classic.go` | ETC mainnet fork blocks, chain ID, MESS |
| `params/config_mordor.go` | Mordor testnet fork blocks |
| `consensus/ethash/ethash.go` | ETChash PoW consensus |
| `core/vm/contracts.go` | Precompile registry |
| `core/vm/opcodes.go` | EVM opcode definitions |
| `internal/build/gotool.go` | CGO_CFLAGS, build flags (blst -std=gnu11 fix here) |
| `Dockerfile` | Multi-stage Docker build (golang:1.24-alpine) |

---

## Docker

```bash
# Build image
docker build -t coregeth-etc:local -f Dockerfile .

# Smoke test
docker run --rm coregeth-etc:local version
```

Multi-stage: `golang:1.24-alpine` builder + `alpine:latest` runtime. Static binary via `build/ci.go install -static`.

**Note:** `gotool.go:58` hardcodes `CGO_CFLAGS=-O2 -g -D__BLST_PORTABLE__ -std=gnu11`. This overrides any Docker ENV CGO_CFLAGS.

---

## Boundaries

### Always Do
- Run `go test ./core/... ./consensus/... ./params/...` before commits
- Use `make geth` for local builds
- Test against Mordor before mainnet
- Respect consensus-critical code boundaries

### Ask First
- Changes to EVM opcodes or precompiles
- Modifying consensus engine behavior
- Docker image or CI/CD changes
- Dependency upgrades (especially blst, c-kzg)

### Never Do
- Break backwards compatibility with ETC network protocol
- Commit private keys or mnemonics
- Remove or bypass tests
- Use `latest` tags in Docker images
- Modify `gotool.go` CGO_CFLAGS without understanding the blst/C23 implications
