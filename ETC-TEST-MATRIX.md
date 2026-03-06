# Ethereum Classic Cross-Client Test Matrix

**Last Updated:** 2026-03-05
**Companion to:** `ETC-HANDOFF.md`
**Specification Source:** [ECIP-1066 — Ethereum Classic Network Description](https://ecips.ethereumclassic.org/ECIPs/ecip-1066)
**Clients:** core-geth (Go), Besu (Java), Fukuii (Scala)

---

## How to Read This Matrix

Each table maps EIP/ECIP specifications from a fork era to test coverage across all 3 clients.

**Coverage levels:**
- **DIRECT** — Dedicated test file/function targeting this specific EIP/ECIP
- **IMPLICIT** — Tested via fork compliance or integration tests (not a dedicated test)
- **INHERITED** — Covered by upstream go-ethereum/hyperledger test suites
- **NONE** — No identified test coverage (gap)

**File paths** are relative to each repo's root directory.

---

## Pre-Olympia Test Matrix (ECIP-1066)

### Frontier (Block 1) — Base EVM

| Feature | core-geth | Besu | Fukuii |
|---------|-----------|------|--------|
| Base EVM execution | INHERITED: upstream go-ethereum tests | INHERITED: upstream Besu EVM tests | DIRECT: `OpCodeFunSpec`, `OpCodeGasSpec` |
| Genesis block | DIRECT: `params/config_etc_test.go` `TestGenesisHashes` | DIRECT: `GenesisConfigClassicTest` `classicChainId()` | DIRECT: `GenesisBlockResponseSpec` |

### Homestead (Block 1,150,000) — EIP-2, EIP-7, EIP-8

| EIP | Feature | core-geth | Besu | Fukuii |
|-----|---------|-----------|------|--------|
| EIP-2 | Homestead hard fork changes | INHERITED | INHERITED | IMPLICIT: `PreOlympiaForkComplianceSpec` `HomesteadFeeSchedule` |
| EIP-7 | DELEGATECALL | DIRECT: `etc_fork_compliance_test.go` `TestETCForkComplianceClassic` | INHERITED | IMPLICIT: `CallOpcodesSpec` |
| EIP-8 | devp2p forward compatibility | INHERITED | INHERITED | INHERITED |

### DAO Fork (Block 1,920,000) — REJECTED

| Feature | core-geth | Besu | Fukuii |
|---------|-----------|------|--------|
| EIP-779 DAO fork rejection | DIRECT: `params/config_classic_test.go` `TestClassicDAO`, `config_etc_test.go` `TestETCClassicDAORejection` | IMPLICIT: Classic config uses no DAO fork block | IMPLICIT: ETC config has no DAO fork |

### Gas Reprice / Tangerine Whistle (Block 2,500,000) — EIP-150

| EIP | Feature | core-geth | Besu | Fukuii |
|-----|---------|-----------|------|--------|
| EIP-150 | Gas cost repricing | DIRECT: `etc_fork_compliance_test.go` EIP-150 compliance | DIRECT: `ClassicProtocolSpecsTest` `dieHardUsesDieHardGasCalculator()` | IMPLICIT: `PreOlympiaForkComplianceSpec`, `OpCodeGasSpec` |

### Die Hard / Spurious Dragon (Block 3,000,000) — ECIP-1010, EIP-155, EIP-160

| Spec | Feature | core-geth | Besu | Fukuii |
|------|---------|-----------|------|--------|
| ECIP-1010 | Difficulty bomb pause | DIRECT: `config_etc_test.go` `TestETCClassicECIP1010Config`, `difficulty_etc_test.go` `TestDifficultyECIP1010BombPause` | DIRECT: `ClassicDifficultyCalculatorsTest` `bombPaused*()` (4 tests) | IMPLICIT: `EthashDifficultyCalculatorSpec` |
| EIP-155 | Replay protection (chain ID) | DIRECT: `etc_fork_compliance_test.go` | IMPLICIT: `ClassicProtocolSpecsTest` fork identification | DIRECT: `EIP155BigIntChainIdSpec`, `SignedLegacyTransactionSpec` |
| EIP-160 | EXP gas cost increase | DIRECT: `etc_fork_compliance_test.go` | DIRECT: `ClassicProtocolSpecsTest` `dieHardFork()` | IMPLICIT: `OpCodeGasSpec` |

### Gotham (Block 5,000,000) — ECIP-1017, ECIP-1039

| Spec | Feature | core-geth | Besu | Fukuii |
|------|---------|-----------|------|--------|
| ECIP-1017 | Emission schedule (5M era, 80% decay) | DIRECT: `ecip1017_test.go` `TestECIP1017EraRewardsIntegration`, `TestECIP1017EraCalculation`, `TestECIP1017RewardDecay`, `TestECIP1017RewardEventuallyZero`; `config_etc_test.go` `TestETCClassicECIP1017Config`, `TestETCMordorECIP1017Config` | DIRECT: `ClassicBlockProcessorTest` (26+ tests: era 0-N rewards, uncle rewards, Mordor 2M eras); `GenesisConfigClassicTest` `mordorEcip1017EraRounds()` | IMPLICIT: `ChainConfigValidationSpec` |
| ECIP-1039 | Monetary policy ratification | IMPLICIT: Covered by ECIP-1017 tests | IMPLICIT: Covered by ECIP-1017 tests | IMPLICIT: Covered by ECIP-1017 tests |

### Defuse Difficulty Bomb (Block 5,900,000) — ECIP-1041

| Spec | Feature | core-geth | Besu | Fukuii |
|------|---------|-----------|------|--------|
| ECIP-1041 | Remove difficulty bomb | DIRECT: `difficulty_etc_test.go` `TestDifficultyECIP1041BombRemoval` | DIRECT: `ClassicDifficultyCalculatorsTest` `bombRemoved*()` tests | IMPLICIT: `EthashDifficultyCalculatorSpec` |

### Atlantis / Byzantium (Block 8,772,000) — 8 EIPs

| EIP | Feature | core-geth | Besu | Fukuii |
|-----|---------|-----------|------|--------|
| EIP-100 | Difficulty adjustment (uncle-aware) | DIRECT: `difficulty_etc_test.go` `TestDifficultyETCAdjustmentDirection` | DIRECT: `ClassicDifficultyCalculatorsTest` `eip100*()` tests | IMPLICIT: `EthashDifficultyCalculatorSpec` |
| EIP-140 | REVERT opcode | DIRECT: `etc_fork_compliance_test.go` | DIRECT: `ClassicProtocolSpecsTest` `atlantisUsesSpuriousDragonGasCalculator()` | IMPLICIT: `OpCodeFunSpec` |
| EIP-196 | bn256 ECADD precompile | DIRECT: `contracts_etc_test.go` `TestETCPrecompilesPerFork`, `TestETCBn256GasRepricing` | INHERITED: Besu precompile tests | DIRECT: `PrecompiledContractsSpec` ECADD tests |
| EIP-197 | bn256 ECPAIRING precompile | DIRECT: `contracts_etc_test.go` `TestETCPrecompilesPerFork` | INHERITED: Besu precompile tests | DIRECT: `PrecompiledContractsSpec` ECPAIRING tests |
| EIP-198 | MODEXP precompile | DIRECT: `contracts_etc_test.go` `TestETCPrecompilesPerFork` | INHERITED: Besu precompile tests | DIRECT: `PrecompiledContractsSpec` MODEXP tests |
| EIP-211 | RETURNDATASIZE / RETURNDATACOPY | DIRECT: `etc_fork_compliance_test.go` | INHERITED | IMPLICIT: `OpCodeFunSpec` |
| EIP-214 | STATICCALL | DIRECT: `etc_fork_compliance_test.go` | INHERITED | DIRECT: `StaticCallOpcodeSpec` |
| EIP-658 | Transaction status codes in receipts | IMPLICIT: Block generation tests | INHERITED | IMPLICIT: Ledger tests |

### Agharta / Constantinople+Petersburg (Block 9,583,000) — 3 EIPs

| EIP | Feature | core-geth | Besu | Fukuii |
|-----|---------|-----------|------|--------|
| EIP-145 | Bitwise shifting (SHL/SHR/SAR) | DIRECT: `etc_fork_compliance_test.go` | INHERITED | DIRECT: `ShiftingOpCodeSpec` |
| EIP-1014 | CREATE2 | DIRECT: `etc_fork_compliance_test.go` | INHERITED | DIRECT: `CreateOpcodeSpec` |
| EIP-1052 | EXTCODEHASH | DIRECT: `etc_fork_compliance_test.go` | INHERITED | IMPLICIT: Call opcode tests |

### Phoenix / Istanbul (Block 10,500,839) — 6 EIPs

| EIP | Feature | core-geth | Besu | Fukuii |
|-----|---------|-----------|------|--------|
| EIP-152 | BLAKE2F precompile | DIRECT: `contracts_etc_test.go` precompile activation | DIRECT: `ClassicProtocolSpecsTest` `phoenixUsesIstanbulGasCalculator()` | DIRECT: `PrecompiledContractsSpec` BLAKE2F, `BlakeCompressionSpec` |
| EIP-1108 | bn256 gas repricing | DIRECT: `contracts_etc_test.go` `TestETCBn256GasRepricing` | IMPLICIT: Istanbul gas calculator | IMPLICIT: `ModExpEIP7883GasSpec` precompile gas |
| EIP-1344 | CHAINID opcode | DIRECT: `etc_fork_compliance_test.go` | IMPLICIT: Protocol spec tests | IMPLICIT: Config tests |
| EIP-1884 | SELFBALANCE + gas repricing | DIRECT: `etc_fork_compliance_test.go` | DIRECT: `ClassicProtocolSpecsTest` `phoenixUsesIstanbulGasCalculator()` | IMPLICIT: `OpCodeGasSpec` |
| EIP-2028 | Calldata gas reduction | IMPLICIT: Fork compliance | IMPLICIT: Istanbul gas calculator | IMPLICIT: `OpCodeGasSpec` |
| EIP-2200 | SSTORE gas metering (net gas) | DIRECT: `gas_table_test.go` `TestEIP2200` | IMPLICIT: Istanbul gas calculator | DIRECT: `SSTOREOpCodeGasPostConstantinopleSpec` |

### Thanos (Block 11,700,000) — ECIP-1099

| Spec | Feature | core-geth | Besu | Fukuii |
|------|---------|-----------|------|--------|
| ECIP-1099 | Etchash (epoch 30K->60K) | DIRECT: `difficulty_etc_test.go` `TestDifficultyECIP1099EpochLength`, `TestDifficultyECIP1099EpochCalculation`; `config_etc_test.go` `TestETCClassicEthashConfig` | DIRECT: `EtcHashTest` `testEcip1099EpochCalculator()`, `testEcip1099EpochCalculatorStartBlock()` | IMPLICIT: `EthashEpochBoundarySpec`, `EthashUtilsSpec` |

### Magneto / Berlin (Block 13,189,133) — 4 EIPs

| EIP | Feature | core-geth | Besu | Fukuii |
|-----|---------|-----------|------|--------|
| EIP-2565 | MODEXP gas repricing | DIRECT: `contracts_etc_test.go` `TestETCModExpEIP2565Repricing` | DIRECT: `ClassicProtocolSpecsTest` `magnetoUsesBerlinGasCalculator()` | DIRECT: `ModExpEIP7883GasSpec` |
| EIP-2929 | Cold/warm access gas costs | DIRECT: `etc_fork_compliance_test.go` | DIRECT: `ClassicProtocolSpecsTest` `magnetoUsesBerlinGasCalculator()` | DIRECT: `OpCodeGasSpecPostEip2929`, `CallOpcodesPostEip2929Spec` |
| EIP-2718 | Typed transaction envelope | DIRECT: `etc_fork_compliance_test.go` `TestETCTransactionTypes` | IMPLICIT: Protocol schedule tests | IMPLICIT: Domain transaction tests |
| EIP-2930 | Access list transactions (Type-1) | DIRECT: `etc_fork_compliance_test.go` `TestETCTransactionTypes` | IMPLICIT: Protocol schedule tests | IMPLICIT: `Eip3651Spec` access list interaction tests |

### Mystique / London (Block 14,525,000) — 2 EIPs

| EIP | Feature | core-geth | Besu | Fukuii |
|-----|---------|-----------|------|--------|
| EIP-3529 | Reduce SSTORE/SELFDESTRUCT refunds | DIRECT: `etc_fork_compliance_test.go` | DIRECT: `ClassicProtocolSpecsTest` `mystiqueUsesLondonGasCalculator()`, `mystiqueUsesLegacyFeeMarket()` | DIRECT: `Eip3529Spec` |
| EIP-3541 | Reject 0xEF-prefixed code | IMPLICIT: Fork compliance | IMPLICIT: Protocol spec tests | DIRECT: `Eip3541Spec` |

### Spiral / Shanghai (Block 19,250,000) — 4 EIPs

| EIP | Feature | core-geth | Besu | Fukuii |
|-----|---------|-----------|------|--------|
| EIP-3651 | Warm COINBASE | IMPLICIT: Fork compliance | DIRECT: `ClassicProtocolSpecsTest` `spiralUsesShanghaiGasCalculator()` | DIRECT: `Eip3651Spec` (9 tests: warm access, gas savings, EXTCODESIZE/EXTCODEHASH interaction) |
| EIP-3855 | PUSH0 opcode | DIRECT: `etc_fork_compliance_test.go` | DIRECT: `ClassicProtocolScheduleDeepTest` `spiralUsesEthereum2024Evm()` | DIRECT: `Push0Spec`, `PreOlympiaForkComplianceSpec` (PUSH0 gating tests) |
| EIP-3860 | Initcode size limit | DIRECT: `gas_table_test.go` `TestCreateGas` (EIP-3860 initcode cost) | IMPLICIT: Shanghai gas calculator | DIRECT: `Eip3860Spec` |
| EIP-6049 | SELFDESTRUCT deprecation | IMPLICIT: Fork compliance | IMPLICIT: Protocol spec tests | DIRECT: `Eip6049Spec` |

---

## Cross-Cutting Features

### ECBP-1100 (MESS — Modified Exponential Subjective Scoring)

| Feature | core-geth | Besu | Fukuii |
|---------|-----------|------|--------|
| MESS configuration | DIRECT: `config_etc_test.go` `TestETCClassicECBP1100Config` | DIRECT: `GenesisConfigClassicTest` `mordorEcbp1100Block()`, `mordorEcbp1100DeactivateBlock()` | DIRECT: `MESScorerSpec` (scorer config + artificial finality) |

### EIP-2124 (Fork ID Dissemination)

| Feature | core-geth | Besu | Fukuii |
|---------|-----------|------|--------|
| Fork ID computation (ETC) | DIRECT: `forkid/forkid_test.go` Classic cases (12 fork blocks) | DIRECT: `ForkIdsNetworkConfigTest` | DIRECT: `ForkIdSpec` `gatherForks for the etc chain correctly` |
| Fork ID computation (Mordor) | DIRECT: `forkid/forkid_test.go` Mordor cases (7 fork blocks) | DIRECT: `ForkIdsNetworkConfigTest` | DIRECT: `ForkIdSpec` `create correct ForkId for mordor blocks` |
| Fork ID validation | DIRECT: `forkid/forkid_test.go` `TestValidation` | INHERITED | DIRECT: `ForkIdValidatorSpec` |
| Fork ID RLP encoding | DIRECT: `forkid/forkid_test.go` `TestEncoding` | INHERITED | DIRECT: `ForkIdSpec` `be correctly encoded via rlp` |

### Chain Identity

| Feature | core-geth | Besu | Fukuii |
|---------|-----------|------|--------|
| ETC mainnet chain ID (61) | DIRECT: `config_etc_test.go` `TestETCClassicChainID` | DIRECT: `GenesisConfigClassicTest` `classicChainId()` | DIRECT: `EIP155BigIntChainIdSpec` |
| Mordor chain ID (63) | DIRECT: `config_etc_test.go` `TestETCMordorChainID` | DIRECT: `GenesisConfigClassicTest` `mordorChainId()` | DIRECT: `ForkIdSpec` (Mordor config) |
| ETC network ID (1) | DIRECT: `etc_fork_compliance_test.go` `TestETCClassicNetworkID` | IMPLICIT: Genesis config | IMPLICIT: Config tests |
| Mordor network ID (7) | DIRECT: `etc_fork_compliance_test.go` `TestETCMordorNetworkID` | IMPLICIT: Genesis config | IMPLICIT: Config tests |
| Genesis hashes | DIRECT: `config_etc_test.go` `TestGenesisHashes`, `etc_fork_compliance_test.go` `TestETCRequireBlockHashes` | DIRECT: Live tests `mordorGenesisHash()`, `etcMainnetGenesisHash()` | DIRECT: `GenesisBlockResponseSpec` |

### Mining & Proof-of-Work

| Feature | core-geth | Besu | Fukuii |
|---------|-----------|------|--------|
| Ethash algorithm (PoW computation) | DIRECT: `consensus/ethash/algorithm_test.go` | INHERITED: Besu Ethash implementation | DIRECT: `EthashUtilsSpec` (SuperSlow: full Ethash validation) |
| Ethash cache/dataset generation | DIRECT: `consensus/ethash/ethash_test.go` | INHERITED | IMPLICIT: `EthashUtilsSpec` |
| Ethash consensus validation | DIRECT: `consensus/ethash/consensus_test.go` | INHERITED | DIRECT: `EthashBlockHeaderValidatorSpec` |
| Block header PoW validation | IMPLICIT: Consensus tests | INHERITED | DIRECT: `EthashBlockHeaderValidatorSpec`, `RestrictedEthashBlockHeaderValidatorSpec` |
| Epoch calculation (30K standard) | DIRECT: `difficulty_etc_test.go` | DIRECT: `EtcHashTest` `testDefaultEpochCalculator()` | DIRECT: `EthashEpochBoundarySpec` |
| Epoch calculation (60K ECIP-1099) | DIRECT: `difficulty_etc_test.go` `TestDifficultyECIP1099EpochLength` | DIRECT: `EtcHashTest` `testEcip1099EpochCalculator()` | IMPLICIT: `EthashEpochBoundarySpec` |
| Difficulty minimum (131,072) | DIRECT: `difficulty_etc_test.go` `TestDifficultyETCMinimum` | DIRECT: `ClassicDifficultyCalculatorsTest` `bombPausedNeverBelowMinimum()` | IMPLICIT: `EthashDifficultyCalculatorSpec` |
| Difficulty bomb pause (ECIP-1010) | DIRECT: `difficulty_etc_test.go` `TestDifficultyECIP1010BombPause` | DIRECT: `ClassicDifficultyCalculatorsTest` `bombPaused*()` (4 tests) | IMPLICIT: `EthashDifficultyCalculatorSpec` |
| Difficulty bomb delay (Gotham) | IMPLICIT: `difficulty_etc_test.go` | DIRECT: `ClassicDifficultyCalculatorsTest` `bombDelayed*()` (3 tests) | IMPLICIT: `EthashDifficultyCalculatorSpec` |
| Difficulty bomb removal (ECIP-1041) | DIRECT: `difficulty_etc_test.go` `TestDifficultyECIP1041BombRemoval` | DIRECT: `ClassicDifficultyCalculatorsTest` `bombRemoved*()` tests | IMPLICIT: `EthashDifficultyCalculatorSpec` |
| EIP-100 uncle-aware difficulty | DIRECT: `difficulty_etc_test.go` | DIRECT: `ClassicDifficultyCalculatorsTest` `eip100*()` tests | IMPLICIT: `EthashDifficultyCalculatorSpec` |
| Classic vs Mordor difficulty configs | DIRECT: `difficulty_etc_test.go` `TestDifficultyClassicVsMordorConfigs` | IMPLICIT: Separate Classic/Mordor calculators | IMPLICIT: Per-network config tests |
| Block reward (ECIP-1017 eras) | DIRECT: `ecip1017_test.go` (4 tests), `etc_fork_compliance_test.go` `TestETCCoinbaseRewardAtGenesis` | DIRECT: `ClassicBlockProcessorTest` (26+ tests) | IMPLICIT: `ChainConfigValidationSpec` |
| Uncle reward (distance-based) | IMPLICIT: `ecip1017_test.go` block generation | DIRECT: `ClassicBlockProcessorTest` `era0OmmerRewardDistance1()`, `era0OmmerRewardDistance7()` | IMPLICIT: Consensus tests |

### Precompiled Contracts (Cumulative at Spiral)

| Address | Precompile | core-geth | Besu | Fukuii |
|---------|------------|-----------|------|--------|
| 0x01 | ECRECOVER | DIRECT: `contracts_etc_test.go` precompile activation | INHERITED | DIRECT: `PrecompiledContractsSpec` |
| 0x02 | SHA256 | DIRECT: `contracts_etc_test.go` precompile activation | INHERITED | DIRECT: `PrecompiledContractsSpec` |
| 0x03 | RIPEMD160 | DIRECT: `contracts_etc_test.go` precompile activation | INHERITED | DIRECT: `PrecompiledContractsSpec` |
| 0x04 | IDENTITY | DIRECT: `contracts_etc_test.go` precompile activation | INHERITED | DIRECT: `PrecompiledContractsSpec` |
| 0x05 | MODEXP | DIRECT: `contracts_etc_test.go` precompile activation + EIP-2565 repricing | INHERITED | DIRECT: `PrecompiledContractsSpec` |
| 0x06 | ECADD (bn256) | DIRECT: `contracts_etc_test.go` + gas repricing | INHERITED | DIRECT: `PrecompiledContractsSpec` |
| 0x07 | ECMUL (bn256) | DIRECT: `contracts_etc_test.go` + gas repricing | INHERITED | DIRECT: `PrecompiledContractsSpec` |
| 0x08 | ECPAIRING (bn256) | DIRECT: `contracts_etc_test.go` | INHERITED | DIRECT: `PrecompiledContractsSpec` |
| 0x09 | BLAKE2F | DIRECT: `contracts_etc_test.go` precompile activation | INHERITED | DIRECT: `PrecompiledContractsSpec`, `BlakeCompressionSpec` |

### Multi-Client Integration (Pre-Olympia)

| Feature | Status | Details |
|---------|--------|---------|
| Mordor genesis hash cross-validation | VERIFIED | All 3 clients produce identical genesis hash on Mordor (chain ID 63) |
| Mordor block sync | PARTIAL | core-geth: synced to head; Besu: synced to 175K; Fukuii: SNAP sync started |
| Live network tests (Mordor) | Besu only | `MordorLiveTest.java` — 9 tests (chain ID, genesis hash, PoW structure, difficulty, gas limit, ECIP-1099, ECIP-1017) |
| Live network tests (ETC mainnet) | Besu only | `EtcMainnetLiveTest.java` — 9+ tests (chain ID, genesis hash, Spiral activation, era rewards) |

---

## Coverage Gaps (Pre-Olympia)

| Gap | core-geth | Besu | Fukuii | Priority |
|-----|-----------|------|--------|----------|
| Live network tests | NONE | DIRECT | NONE | LOW (manual verification done) |
| ECIP-1039 dedicated test | NONE | NONE | NONE | LOW (ratification meta-ECIP) |
| EIP-8 devp2p forward compat | INHERITED | INHERITED | INHERITED | LOW (p2p layer, not consensus) |
| Uncle reward per-era tests | IMPLICIT | DIRECT | IMPLICIT | MEDIUM (Besu has best coverage) |

---

## Test Execution Commands

### core-geth
```bash
go test ./params/... -count=1 -timeout 2m -v           # Chain config + fork tests
go test ./core/... -count=1 -timeout 10m -v             # Fork compliance, ECIP-1017, block processing
go test ./consensus/... -count=1 -timeout 5m -v         # Ethash, difficulty, PoW
go test ./core/forkid/... -count=1 -timeout 1m -v       # Fork ID
go test ./core/vm/... -count=1 -timeout 5m -v           # EVM, precompiles, gas
```

### Besu
```bash
./gradlew :config:test --tests "*Classic*" -i            # Genesis config
./gradlew :ethereum:core:test --tests "*Classic*" -i     # Protocol specs, difficulty, block processor
./gradlew :ethereum:core:test --tests "*EtcHash*" -i     # Ethash/ECIP-1099
./gradlew :ethereum:core:test --tests "*live*" -i        # Live network tests (requires running node)
```

### Fukuii
```bash
sbt "testOnly *PreOlympiaForkComplianceSpec"             # Fork compliance
sbt "testOnly *ForkIdSpec"                               # Fork ID (EIP-2124)
sbt "testOnly *PrecompiledContractsSpec"                  # Precompiles
sbt "testOnly *EthashDifficultyCalculatorSpec"            # Difficulty
sbt "testOnly *MESScorerSpec"                             # ECBP-1100
sbt test                                                  # Full suite (~2,192 tests)
```

---

## Document History

| Date | Change |
|------|--------|
| 2026-03-05 | Initial creation — pre-Olympia matrix covering ECIP-1066 (Frontier through Spiral) |
