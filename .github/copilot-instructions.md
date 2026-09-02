# Copilot instructions for core-geth

This file is self-contained. It duplicates the project context in
[`AGENTS.md`](../AGENTS.md) because Copilot does not read that file on every
surface, and a pointer would leave the model with nothing on the surfaces that
do not. When either file changes, change both.

## What this repository is

CoreGeth is an Ethereum protocol provider: a downstream of
[ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) that keeps chain
configuration data-driven rather than hard-coded, so one binary serves Ethereum
Classic, Mordor, the Ethereum Foundation network, MintMe and private chains.

**This client is consensus software for a live network.** A change that alters
which blocks a node accepts is a chain split, not a failing build.

## Branching

The branch layout is mid-transition. Read it as a sequence, not a steady state.

| Branch | Now | After the transition |
|---|---|---|
| `main` | carries the modernization work, in progress | the default branch |
| `master` | the default branch, and the PR target `.github/CONTRIBUTING.md` names | retired |
| `archive-etclabscore-2024-12` | the preserved pre-migration history | kept as the archive branch |

- `main` is cut from `archive-etclabscore-2024-12` (`7ef3ecd7a`, 2024-12-16), the
  last commit before this repository was created on 2024-12-21.
- `main` replaces `master` as the default when the modernization is written,
  tested, and ready to cut a release candidate — not before.
- Changing the default branch is a repository setting and a deliberate act at that
  moment. Do not change it incidentally and do not assume it has changed.
- Confirm the branch before treating anything as current:
  `git rev-parse --abbrev-ref HEAD`.

## Toolchain

Three Go versions are declared in this tree and they do not agree. This is real
and load-bearing:

| Where | Declares | Governs |
|---|---|---|
| `go.mod` | the language version | what the module compiles against |
| `.github/workflows/*.yml` | `go-version`, a major only | what CI builds and tests with |
| `build/checksums.txt` | `# version:golang`, an exact patch | what `-dlgo` downloads |
| `Dockerfile`, `Dockerfile.alltools` | a `golang:` major tag | the container build |

Read the current values from those files. The major agrees across all four; the
patch does not, because two of them float. `make all` carries no `-dlgo`, so most
release archives use the runner's patch while only the ARM leg uses the pinned
one.

`build/checksums.txt` also pins the linter release. `make lint`
downloads that exact linter release rather than using a locally installed copy.

**The Go module path is `github.com/ethereum/go-ethereum`, not a core-geth path.**
Every internal import uses it. That is deliberate downstream compatibility — do
not "fix" it, and do not expect package paths to match the repository name.

## Commands

Every command below is defined in the `Makefile` or in `build/ci.go`. There is no
task runner other than `make`, and no command exists that is not listed here or
printed by `make help`.

```bash
make core-geth       # build cmd/core-geth into ./build/bin/core-geth
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
make test-coregeth-regression-condensed # builds core-geth, imports simulated chains
```

`make test-evmc` builds external EVM interpreters and is not part of the default
suite. `make tests-generate` **overwrites** generated fixtures under
`tests/testdata-etc/` — read the target before running it.

Test fixtures are git submodules and the consensus suites fail confusingly
without them:

```bash
git submodule update --init --recursive
```

## Continuous integration

GitHub Actions under `.github/workflows/` is what runs: `test-linux.yml` (lint
plus both test suites), `evmc.yml`, `go-generate-check.yml`, `bench-*.yml`,
`docs-deploy.yml`, `release-packages.yml`, `audit-bootnodes.yml`.

Triggers are not uniform and the branch names are mid-transition. `test-linux.yml`
fires on a push to `main`, every pull request, and dispatch. `docs-deploy.yml`
fires on `master` or `main`, path-filtered. `evmc.yml` and the three `bench-*.yml`
still fire on `master` only, so they stop silently when `master` retires. The
release and image workflows fire on a `v*` tag. Verify against the file.

`.travis.yml`, `circle.yml`, `appveyor.yml` and `Jenkinsfile` are historical CI
definitions, not what runs today. Leave them alone.

## Layout

```
cmd/core-geth         the node binary; cmd/utils/flags.go defines the network flags
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
Ethereum Classic mainnet schedule as per-EIP activation blocks.

**The Ethereum Classic fork schedule implemented here runs from Frontier through
Spiral**, Spiral being the head configuration at block 19,250,000 —
`EIP3651FBlock`, `EIP3855FBlock`, `EIP3860FBlock` and `EIP6049FBlock`.
`EIP4399FBlock` and `EIP4895FBlock` are commented out with their reasons;
Ethereum Classic is proof of work and does not adopt them. The same struct sets
`ECIP1010PauseBlock`/`ECIP1010Length`, `ECIP1017FBlock`/`ECIP1017EraRounds`,
`ECIP1099FBlock` (Etchash), and `ECBP1100FBlock` with `ECBP1100DeactivateFBlock`,
which switches the MESS artificial-finality rule off at the Spiral block.

Read activation blocks out of the file. Do not restate a fork schedule from
memory, and do not infer one network's rules from another's.

`params/version.go` is the single source of the version: `1.12.21-unstable`.

## Dependency updates

`.github/dependabot.yml` exists and version updates are deliberately off
(`open-pull-requests-limit: 0`), recorded 2026-09-01. Five ecosystems name
something this repository holds: `gomod`, `pip`, `docker`, `github-actions` and
`gitsubmodule`. Dependabot security updates are a repository setting with no key
in that file, and a limit of zero does not withhold them. Do not turn the
disabled config into an active one — its state is a decision.

Every `uses:` reference is pinned to a full commit SHA, with the readable version
in a trailing comment. Bumping one means resolving the new SHA, reading the diff
at that commit, and updating both — and it is a workflow change, so it needs
confirmation. The comment can drift from the SHA; trust the SHA.

## Facts that mislead if you do not know them

- **`swarm/` is legacy**, retained but not the active networking stack.
- **`sync-parity-chainspecs` is marked deprecated in the `Makefile` itself.**
- **`AUTHORS` is generated, not written**, by `build/update-license.go` from
  `git shortlog` via `.mailmap`. Nothing runs it, so it is stale. Never hand-edit
  it — the header is a constant in the generator and an edit is reverted on the
  next run. That generator also rewrites every source file's license header from a
  template hardcoded to `The go-ethereum Authors`, and no file attributed to
  `The core-geth Authors` is in its skip list, so running it deletes that
  attribution.
- **`SECURITY.md` is upstream's and points at the Ethereum Foundation** —
  `bounty@ethereum.org`, the Foundation's PGP key, go-ethereum audit links. It is
  not this project's policy; do not cite it as the reporting path.
- **`core-geth version-check` queries go-ethereum's vulnerability feed** and prints
  `No vulnerabilities found` when nothing matches. That feed does not track this
  client.

## Boundaries

Ask before: pushing to any remote; committing anything, documentation included;
touching `.github/workflows/`; changing any chain configuration — activation
block, fork schedule, genesis allocation, bootnode list, checkpoint hash;
changing any dependency — `go.mod`, `go.sum`, `requirements-mkdocs.txt`, a
Dockerfile base image, a submodule pointer, a pin in `build/checksums.txt`;
regenerating test fixtures; opening, editing or closing a
pull request or issue.

Never: change repository settings, branch protection, rulesets or Actions
permissions; commit a key, keystore, credential or token in any form; use
`git add .` or `git add -A`; add, change or recommend changing `COPYING`,
`COPYING.LESSER` or any license header; claim a test suite passed without having
run it.

## Conventions

- Commit messages are prefixed with the package or area they modify, per
  `.github/CONTRIBUTING.md` — for example `eth, rpc: make trace configs optional`.
- Commit messages are public and permanent. No internal notes, no
  characterization of any person or organization, no machine-local paths.
- Code is `gofmt`-formatted and documented per Go commentary conventions.
  `make lint` is the gate and runs the linter set enabled in `.golangci.yml`, not
  the golangci-lint defaults.
- American English, in code, comments, commits and documentation.
- Comments carry reasoning, not description.
- Verify by effect, and calibrate the check so it can fail. A check that cannot
  report a negative proves nothing.
