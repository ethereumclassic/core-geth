## CoreGeth: Ethereum Classic Execution Client

> A [go-ethereum](https://github.com/ethereum/go-ethereum) fork providing the production Ethereum Classic (ETC) execution client.

---

## ⚠️ Node operators: upgrade to v1.13.0, and rotate your node key

**Every release in the v1.12.x line carries at least one unpatched CVE, and the earliest
carry all six.** They were built on a Go toolchain that reached end of life in August 2024.
Upgrade, then perform the one cleanup step below.

**1. Upgrade, and change where you track releases.** Releases are cut from
[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth). A node tracking
the previous repository will not see v1.13.0.

**2. Rotate the P2P node key. This is required, not precautionary.** CVE-2026-26315 is an
oracle: an invalid-curve ephemeral key in the RLPx handshake previously reached ECDH and
failed only at MAC verification, which leaks bits of the node key across repeated
handshakes. Any key used by an unpatched node should be treated as exposed.

```bash
# with the node stopped
mv <datadir>/geth/nodekey <datadir>/geth/nodekey.old-rotated-$(date +%F)
```

The client generates a fresh key on the next start. **Your enode ID changes**, so update
every static-peer, trusted-peer or bootnode list that names this node — on your own
machines and with anyone peering with you. Capture the old enode ID before rotating if you
need it to find those references; it cannot be re-derived afterwards.

**3. Re-check your RPC exposure while the node is down.** `--http.addr` should be loopback
unless a trusted proxy sits in front of it, and `admin`, `debug` and `personal` do not
belong in `--http.api` on any reachable interface. `--http.corsdomain="*"` is not a safe
default.

**A resync is not required.** No finding in either audit corrupts chain data, and none is
known to have been exploited against this network. Verify your head matches another source
before concluding otherwise:

```bash
geth attach --exec 'eth.blockNumber' <datadir>/geth.ipc
```

Compare against a block explorer or a second client. Resync only if it diverges — and if it
does, that is a finding worth reporting.

**Nothing here requires touching your keystore.** These are network-layer and handshake
issues; account keys are not implicated by any of them. Rotate account keys only if your
`admin` or `personal` RPC was reachable from an untrusted network, which is its own
exposure rather than one of these CVEs.

Full detail: the [March 2026 audit](docs/audits/2026-03-security-audit.md), the
[August 2026 follow-up](docs/audits/2026-08-security-followup.md), and the
[v1.13.0 migration guide](docs/tutorials/v1.13.0-migration.md).

---

CoreGeth is a production execution client for the Ethereum Classic network. It implements every ETC hard fork from Frontier through Spiral.

**Note:** Upstream go-ethereum has removed support for Ethereum Classic, so ETC consensus rules are maintained here rather than inherited.

## Supported Networks

| Network | Chain ID | Consensus | Flag |
|---------|----------|-----------|------|
| Ethereum Classic (ETC) | 61 | Proof of Work (ETChash) | `--classic` |
| Mordor Testnet | 63 | Proof of Work (ETChash) | `--mordor` |
| MintMe.com Coin | 24734 | Proof of Work | `--mintme` |
| Private chains | configurable | PoW / PoA | genesis config |

Sepolia and Holesky are also registered, inherited from upstream go-ethereum.

### ETC consensus history

Activation blocks, governing specifications and included EIPs, aligned with
[ECIP-1066](https://ecips.ethereumclassic.org/ECIPs/ecip-1066), which is the network description this table follows. The Mordor
column is client detail and is not part of that specification.

| Upgrade | ETC Block | Mordor Block | Date | Specs | Included EIPs |
|---------|----------:|-------------:|------|-------|---------------|
| Spiral <br><sub>Shanghai</sub> | 19,250,000 | 9,957,000 | 2024-02-04 | [ECIP-1109](https://ecips.ethereumclassic.org/ECIPs/ecip-1109) | Shanghai: [EIP-3651](https://eips.ethereum.org/EIPS/eip-3651), [EIP-3855](https://eips.ethereum.org/EIPS/eip-3855), [EIP-3860](https://eips.ethereum.org/EIPS/eip-3860), [EIP-6049](https://eips.ethereum.org/EIPS/eip-6049) |
| MESS Default: Off <br><sub>ECBP-1110</sub> | 19,250,000 | 10,400,000 | 2024-02-04 | [ECIP-1110](https://ecips.ethereumclassic.org/ECIPs/ecip-1110) | None |
| Mystique <br><sub>London</sub> | 14,525,000 | 5,520,000 | 2022-02-12 | [ECIP-1104](https://ecips.ethereumclassic.org/ECIPs/ecip-1104) | London: [EIP-3529](https://eips.ethereum.org/EIPS/eip-3529), [EIP-3541](https://eips.ethereum.org/EIPS/eip-3541) |
| Magneto <br><sub>Berlin</sub> | 13,189,133 | 3,985,893 | 2021-07-23 | [ECIP-1103](https://ecips.ethereumclassic.org/ECIPs/ecip-1103) | Berlin: [EIP-2565](https://eips.ethereum.org/EIPS/eip-2565), [EIP-2718](https://eips.ethereum.org/EIPS/eip-2718), [EIP-2929](https://eips.ethereum.org/EIPS/eip-2929), [EIP-2930](https://eips.ethereum.org/EIPS/eip-2930) |
| Thanos <br><sub>ECIP-1099</sub> | 11,700,000 | 2,520,000 | 2020-11-28 | [ECIP-1099](https://ecips.ethereumclassic.org/ECIPs/ecip-1099) | None |
| MESS Default: On <br><sub>ECBP-1100</sub> | 11,380,000 | 2,380,000 | 2020-10-09 | [ECIP-1100](https://ecips.ethereumclassic.org/ECIPs/ecip-1100) | None |
| Phoenix <br><sub>Istanbul</sub> | 10,500,839 | 999,983 | 2020-06-01 | [ECIP-1088](https://ecips.ethereumclassic.org/ECIPs/ecip-1088) | Istanbul: [EIP-152](https://eips.ethereum.org/EIPS/eip-152), [EIP-1108](https://eips.ethereum.org/EIPS/eip-1108), [EIP-1344](https://eips.ethereum.org/EIPS/eip-1344), [EIP-1884](https://eips.ethereum.org/EIPS/eip-1884), [EIP-2028](https://eips.ethereum.org/EIPS/eip-2028), [EIP-2200](https://eips.ethereum.org/EIPS/eip-2200) |
| Agharta <br><sub>Constantinople+Petersburg</sub> | 9,573,000 | 301,243 | 2020-01-11 | [ECIP-1056](https://ecips.ethereumclassic.org/ECIPs/ecip-1056) | Constantinople+Petersburg: [EIP-145](https://eips.ethereum.org/EIPS/eip-145), [EIP-1014](https://eips.ethereum.org/EIPS/eip-1014), [EIP-1052](https://eips.ethereum.org/EIPS/eip-1052) |
| Atlantis <br><sub>Byzantium</sub> | 8,772,000 | 0 | 2019-09-12 | [ECIP-1054](https://ecips.ethereumclassic.org/ECIPs/ecip-1054) | Spurious Dragon: [EIP-161](https://eips.ethereum.org/EIPS/eip-161), [EIP-170](https://eips.ethereum.org/EIPS/eip-170) <br> Byzantium: [EIP-100](https://eips.ethereum.org/EIPS/eip-100), [EIP-140](https://eips.ethereum.org/EIPS/eip-140), [EIP-196](https://eips.ethereum.org/EIPS/eip-196), [EIP-197](https://eips.ethereum.org/EIPS/eip-197), [EIP-198](https://eips.ethereum.org/EIPS/eip-198), [EIP-211](https://eips.ethereum.org/EIPS/eip-211), [EIP-214](https://eips.ethereum.org/EIPS/eip-214), [EIP-658](https://eips.ethereum.org/EIPS/eip-658) |
| Defuse Difficulty Bomb <br><sub>ECIP-1041</sub> | 5,900,000 | 0 | 2018-05-29 | [ECIP-1041](https://ecips.ethereumclassic.org/ECIPs/ecip-1041) | None |
| Gotham <br><sub>ECIP-1017</sub> | 5,000,000 | 0 | 2017-12-11 | [ECIP-1017](https://ecips.ethereumclassic.org/ECIPs/ecip-1017) <br> [ECIP-1039](https://ecips.ethereumclassic.org/ECIPs/ecip-1039) | None |
| Die Hard <br><sub>Spurious Dragon</sub> | 3,000,000 | 0 | 2017-01-13 | [ECIP-1010](https://ecips.ethereumclassic.org/ECIPs/ecip-1010) | Spurious Dragon: [EIP-155](https://eips.ethereum.org/EIPS/eip-155), [EIP-160](https://eips.ethereum.org/EIPS/eip-160) |
| Gas Reprice <br><sub>Tangerine Whistle</sub> | 2,500,000 | 0 | 2016-10-24 | [ECIP-1015](https://ecips.ethereumclassic.org/ECIPs/ecip-1015) | Tangerine Whistle: [EIP-150](https://eips.ethereum.org/EIPS/eip-150) |
| ~~DAO Fork~~ | 1,920,000 | — | 2016-07-20 | [HFM-779](https://eips.ethereum.org/EIPS/eip-779) | **Rejected on ETC** |
| Homestead | 1,150,000 | 0 | 2016-03-14 | [HFM-606](https://eips.ethereum.org/EIPS/eip-606) | Homestead: [EIP-2](https://eips.ethereum.org/EIPS/eip-2), [EIP-7](https://eips.ethereum.org/EIPS/eip-7), [EIP-8](https://eips.ethereum.org/EIPS/eip-8) |
| Frontier Thawing | 200,000 | 0 | 2015-09-07 | Genesis | None |
| Frontier | 1 | 0 | 2015-07-30 | Genesis | None |

**MESS** (Modified Exponential Subjective Scoring, ECBP-1100) is a chain-selection
defense against deep reorganizations. It is a client-side policy rather than a state
transition, which is why it appears as two rows: the block at which it defaults on, and
the block at which it defaults off.

Ethereum Classic does not implement EIP-1559. Transactions are legacy and EIP-2930
access-list types.

### Wire protocol

**EIP-7642 (`eth/69`) is deliberately excluded.** It removes Total Difficulty from the
protocol handshake, and Ethereum Classic requires that field for proof-of-work chain
selection. Do not restore it.

## Build

```bash
make geth
```

The binary is written to `./build/bin/geth`.

## Test

```bash
# ETC-specific unit tests
go test ./params/... -run TestETC -v
go test ./core/... -run "TestGasLimit|TestECIP1017|TestETCForkCompliance" -v
go test ./consensus/ethash/... -run "TestDifficultyETC|TestDifficultyECIP" -v

# Live-network tests, build-tagged; requires a running Mordor or ETC node
go test -tags live ./tests/live_etc/ -v
```

`tests/evm-benchmarks` is a submodule. Run `git submodule update --init --recursive`
before any suite that reads the shared consensus fixtures.

## Run a node

```bash
./build/bin/geth --classic --datadir <path>     # Ethereum Classic mainnet
./build/bin/geth --mordor  --datadir <path>     # Mordor testnet
```

`--http.addr` defaults to loopback. Do not widen it, enable `--http.corsdomain`, or add
`admin` or `debug` to `--http.api` on a node reachable from an untrusted network.

## Mining

CoreGeth supports Ethash/ETChash proof-of-work mining:

```bash
./build/bin/geth --classic --mine --miner.etherbase <address>
```

For testing with fake PoW, which skips DAG generation:

```bash
./build/bin/geth --classic --mine --miner.etherbase <address> --fakepow
```

## Documentation

- [CoreGeth docs](https://etclabscore.github.io/core-geth) — general documentation
- [go-ethereum docs](https://geth.ethereum.org/docs/) — upstream reference

## Security

This release backports six CVE fixes from upstream go-ethereum and adds a GraphQL query
depth limit. Each is a separate commit; run `git log --grep CVE-` for the full detail,
including the upstream commit every fix derives from.

| Identifier | Component | Issue |
|-----|-----------|-------|
| CVE-2026-26313 | `eth/protocols` | Unbounded RLP message decoding |
| CVE-2026-26314 | `crypto/secp256k1` | `IsOnCurve` accepted coordinates outside the field |
| CVE-2026-26315 | `crypto/ecies` | Invalid-curve attack on the RLPx handshake |
| CVE-2026-22862 | `crypto/ecies` | Minimum ciphertext length check |
| CVE-2026-22868 | `eth` | Invalid KZG blob proofs |
| CVE-2025-24883 | `crypto` | `UnmarshalPubkey` accepted off-curve points |
| — | `graphql` | Unbounded query depth |

**Rotate P2P node keys after upgrading.** CVE-2026-26315 is an oracle against the node key
itself, so a key used by an unpatched node should be treated as exposed. Stop the node,
remove `<datadir>/geth/nodekey`, and restart; the client generates a new one. The node's
enode ID changes, so update any static-peer or trusted-peer list that names it.

`SECURITY.md` carries the disclosure policy and PGP key.

## License

The core-geth library (outside `cmd/`) is licensed under [LGPL-3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html).
The core-geth binaries (`cmd/`) are licensed under [GPL-3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).

`COPYING` and `COPYING.LESSER` carry the full texts, and `NOTICE` carries attribution.
There is deliberately no `LICENSE` file: the split is inherited from go-ethereum, and a
per-file header states which of the two applies.
