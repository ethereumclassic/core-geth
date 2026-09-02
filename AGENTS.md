# core-geth

CoreGeth is an Ethereum protocol provider: a downstream of
[ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) that keeps chain
configuration data-driven rather than hard-coded, so one binary serves Ethereum
Classic, Mordor, the Ethereum Foundation network, MintMe and private chains.

[README.md](README.md) is the project introduction and the network support
matrix. This file is what an agent needs in order to work here without breaking
something.

**This client is consensus software for a live network.** A change that alters
which blocks a node accepts is a chain split, not a failing build. Everything
below about confirmation, verification and boundaries exists for that reason.

## Branching

The branch layout is mid-transition. Read it as a sequence, not a steady state.

| Branch | Now | After the transition |
|---|---|---|
| `main` | carries the modernization work, in progress | the default branch |
| `master` | the default branch, and the PR target `.github/CONTRIBUTING.md` names | retired |
| `archive-etclabscore-2024-12` | the preserved pre-migration history | kept as the archive branch |

- **`main` is cut from `archive-etclabscore-2024-12`** (`7ef3ecd7a`, 2024-12-16),
  the last commit before this repository was created on 2024-12-21. The annotated
  tag `archive/etclabscore-2024-12` marks the same commit.
- **`main` replaces `master` as the default when the modernization is written,
  tested, and ready to cut a release candidate — not before.** Until that point
  `master` is the default and is what contributors are directed to.
- **Changing the default branch is a repository setting and a deliberate act at
  that moment.** Do not change it incidentally, do not assume it has changed, and
  re-read this table rather than recalling which branch was default.

Confirm which branch you are on before reading anything as current:

```bash
git rev-parse --abbrev-ref HEAD
```

## Toolchain

**Three Go versions are declared in this tree and they do not agree.** This is
real, it is load-bearing, and it is one of the things a modernization pass has
to reconcile rather than assume:

| Where | Declares | Governs |
|---|---|---|
| `go.mod` | `go 1.21` | the language version the module compiles against |
| `.github/workflows/*.yml` | `go-version: '1.21'` | what CI builds and tests with |
| `build/checksums.txt` | `# version:golang 1.22.1` | what `build/ci.go install -dlgo` downloads for release builds |
| `Dockerfile`, `Dockerfile.alltools` | `golang:1.22-alpine` | the container build |

`build/checksums.txt` also pins `# version:golangci 1.55.2`. `make lint` downloads
that exact linter release; it is not read from any locally installed copy.

**The Go module path is `github.com/ethereum/go-ethereum`, not a core-geth path.**
Every internal import uses it. That is deliberate downstream compatibility — do
not "fix" it, and do not be surprised when a package path does not match the
repository name.

## Commands

Every command below is defined in the `Makefile` or in `build/ci.go`. There is no
task runner other than `make`, and no command exists that is not listed here or
printed by `make help`.

```bash
make geth            # build cmd/geth into ./build/bin/geth
make all             # build every executable
make test            # make all, then build/ci.go test -timeout 20m
make lint            # build/ci.go lint -> golangci-lint run --config .golangci.yml
make clean           # go clean -cache, remove build output
make help            # list the annotated targets
```

CoreGeth-specific suites, which is what proves this fork's configuration work:

```bash
make test-coregeth                      # features + clique consensus + condensed regression
make test-coregeth-features             # fork/feature/datatype equivalence
make test-coregeth-consensus            # clique consensus equivalence
make test-coregeth-chainspecs-coregeth  # CoreGeth JSON chainspec equivalence
make test-coregeth-regression-condensed # builds geth, imports simulated canonical chains
```

`make test-evmc` builds external EVM interpreters (`make hera`, `make evmone`)
and is not part of the default suite. `make tests-generate` **overwrites**
generated fixtures under `tests/testdata-etc/` — read the target before running
it.

### Submodules are required

`tests/testdata`, `tests/evm-benchmarks` and `tests/testdata-etc` are git
submodules. The test suites read from them and fail confusingly when they are
absent:

```bash
git submodule update --init --recursive
```

CI checks out with `submodules: recursive` for the test jobs and
`submodules: false` for lint.

## Continuous integration

GitHub Actions under `.github/workflows/` is what actually runs:
`test-linux.yml` (lint plus both test suites), `evmc.yml`, `go-generate-check.yml`,
`bench-*.yml`, `docs-deploy.yml`, `release-packages.yml`, `audit-bootnodes.yml`.

**`test-linux.yml` triggers on `push` to `master` only**, plus every pull request
and manual dispatch. A push to any other branch runs no workflow. Verify against
the workflow file rather than assuming a push was tested.

`.travis.yml`, `circle.yml`, `appveyor.yml` and `Jenkinsfile` are also present.
They are historical CI definitions, not what runs today. Leave them alone.

## Structure

```
cmd/geth              the node binary; cmd/utils/flags.go defines the network flags
params/               chain configuration - the core of what makes this a fork
params/config_classic.go   Ethereum Classic mainnet fork schedule
params/types/         the configuration interfaces that make chain config data-driven
core/, consensus/     block processing and consensus rules
eth/, p2p/, rpc/      networking and the JSON-RPC surface
tests/                consensus test harness plus the three fixture submodules
build/ci.go           the real build/test/lint entry point; the Makefile wraps it
docs/, mkdocs.yml     the documentation site source
```

## Chain configuration

Chain rules are data, not code branches. `params/config_classic.go` holds the
Ethereum Classic mainnet schedule as per-EIP activation blocks, and the equivalent
files hold the other supported networks.

**The Ethereum Classic fork schedule implemented here runs from Frontier through
Spiral**, Spiral being the head configuration at block 19,250,000 — `EIP3651FBlock`
(warm COINBASE), `EIP3855FBlock` (PUSH0), `EIP3860FBlock` (initcode metering) and
`EIP6049FBlock` (SELFDESTRUCT deprecation). `EIP4399FBlock` and `EIP4895FBlock` are
commented out with their reasons; Ethereum Classic is proof of work and does not
adopt them.

Alongside the EIP schedule the same struct sets the ECIP fields:
`ECIP1010PauseBlock`/`ECIP1010Length` (difficulty bomb defusal), `ECIP1017FBlock`/
`ECIP1017EraRounds` (monetary policy), `ECIP1099FBlock` (Etchash), and
`ECBP1100FBlock` with `ECBP1100DeactivateFBlock`, which switches the MESS
artificial-finality rule off at the Spiral block. The Istanbul-equivalent set is
labeled `// ECIP-1088` in a comment rather than carried as its own field.

Reading an activation block out of this file is the only reliable way to know
what is configured. Do not restate a fork schedule from memory, and do not infer
one network's rules from another's.

## Version

`params/version.go` is the single source: `1.12.21-unstable`, `VersionName`
`CoreGeth`. Release tooling reads it; nothing else should hard-code a version
string.

## Dependency updates

**`.github/dependabot.yml` exists and version updates are deliberately off**
(`open-pull-requests-limit: 0`). Recorded 2026-09-01.

- **Five ecosystems name something this repository actually holds** — `gomod`
  (`go.mod`, `go.sum`), `pip` (`requirements-mkdocs.txt`), `docker` (two
  Dockerfiles on `golang:1.22-alpine` and `alpine:latest`), `github-actions`
  (nine workflows) and `gitsubmodule` (three submodules).
- **The limit is zero because nobody is triaging a standing pull-request queue.**
  A queue nobody reads reports itself as a control while operating as noise.
- **Dependabot security updates are a repository setting with no key in that
  file.** Nothing written there turns them on or off, and a limit of zero does not
  withhold them — their pull requests are not subject to the limit and do not
  count toward it. Measured 2026-09-01: security updates are **enabled but
  paused**, vulnerability alerts are on, and five security pull requests are open
  against the `pip` surface. So this repository is not quiet, and the disabled
  config is not what makes it so either way.
- **What to re-check rather than re-read:** the repository-level security setting,
  which anyone can flip and which GitHub pauses on its own for an inactive
  repository, and the per-ecosystem support table, which grows. Neither is visible
  in the config file.
- **What would change the version-update decision:** someone taking ownership of
  the queue. Raising the limit brings a `cooldown:` block with it; while the limit
  is zero a cooldown would gate nothing and would read as a control that is
  operating.

Do not "fix" the disabled config into an active one. Its state is a decision.

**Actions in the workflows are pinned to mutable tags** (`actions/checkout@v2`,
`actions/setup-go@v5`, `softprops/action-gh-release@v1` and others). A tag is a
mutable reference; a commit SHA is not. Repointing them is a workflow change and
needs confirmation before it is made.

## Facts that mislead if you do not know them

- **`git.diff` is tracked, not a stray file.** It is a 6.0 MB conflicted diff
  committed by accident in a 2021 upstream merge (`d4a8d7365`) and still present
  in `master`. Nothing reads it, and `.dockerignore` does not exclude it, so it is
  carried into the Docker build context. Removing it is a deliberate repository
  change, not tidying — ask first.
- **`swarm/` is legacy.** It is retained; it is not the active networking stack.
- **The `sync-parity-chainspecs` target is marked deprecated in the `Makefile`
  itself.** Parity configuration support is not maintained past the Istanbul fork.
- **`AUTHORS` is generated, not written.** `build/update-license.go` produces it
  from `git shortlog -s -n -e`, canonicalized through `.mailmap`. Nothing runs it
  — no `make` target, no CI job — so the file is a stale snapshot. **Never edit it
  by hand**; the header line is a constant in the generator and a hand edit is
  reverted on the next run.
- **That generator also rewrites the license header of every source file** from a
  template hardcoded to `The go-ethereum Authors`, and no file attributed to
  `The core-geth Authors` is in its skip list. Running it as it stands deletes
  that attribution. Do not run it without deciding first what it should assert.
- **`SECURITY.md` is upstream's and points at the Ethereum Foundation** — reports
  are directed to `bounty@ethereum.org` under the Foundation's PGP key, and the
  audit links go to go-ethereum. It is not this project's policy. Do not follow it
  and do not cite it as the reporting path.
- **`geth version-check` queries go-ethereum's vulnerability feed**
  (`cmd/geth/misccmd.go`) and prints `No vulnerabilities found` when nothing
  matches. That feed does not track this client, so a clean result from it says
  nothing about this client.

## Boundaries

### Ask first

- **Any push, to any remote.** This is a public repository of the Ethereum
  Classic organization. Nothing leaves the machine without explicit confirmation.
- **Any commit.** Including a commit that only touches documentation.
- **Anything under `.github/workflows/`.** These run in the organization's CI with
  the organization's secrets.
- **Any change to a chain configuration** — an activation block, a fork schedule,
  a genesis allocation, a bootnode list, a checkpoint hash. These decide which
  chain a node follows.
- **Any dependency change** — `go.mod`, `go.sum`, `requirements-mkdocs.txt`, a
  Dockerfile base image, a submodule pointer, a pinned version in
  `build/checksums.txt`.
- **Regenerating test fixtures** (`make tests-generate` and its sub-targets).
- **Removing `git.diff`.**
- **Opening, editing or closing a pull request or issue.**

### Never

- **Change repository settings, branch protection, rulesets or Actions
  permissions.** Report drift; do not correct it.
- **Commit a key, keystore, credential or token in any form**, encrypted or not.
- **Use `git add .` or `git add -A`.** Stage named paths.
- **Add, change, or recommend changing `COPYING`, `COPYING.LESSER`, or any license
  header.** Licensing is a legal question before it is a technical one.
- **Claim a test suite passed without having run it.** Name what ran.

## Conventions

- **Commit messages are prefixed with the package or area they modify**, per
  `.github/CONTRIBUTING.md` — for example `eth, rpc: make trace configs optional`.
  Present tense, lower case after the prefix.
- **Commit messages are public and permanent.** No internal notes, no
  characterization of any person or organization, no machine-local paths. If it
  cannot be verified from the repository or a published source, it does not belong
  in one.
- **Code is `gofmt`-formatted and documented per Go commentary conventions.**
  `make lint` is the gate; it runs the linter set enabled in `.golangci.yml`
  (`goimports`, `govet`, `staticcheck`, `unused`, `misspell` and others), not the
  golangci-lint defaults.
- **American English**, in code, comments, commits and documentation.
- **Comments carry reasoning, not description.** `params/config_classic.go`
  annotates activation blocks with the fork they belong to; match that.
- **Verify by effect, and calibrate the check so it can fail.** A check that
  cannot report a negative proves nothing. This applies to gitignore coverage
  (`git check-ignore --no-index -q -- <path>`, never `-v` as the condition), to
  test results, and to any claim that something works.
