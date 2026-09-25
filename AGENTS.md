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

**This repository does not send pull requests to a separate upstream — it IS
the Ethereum Classic organization's own repository.** `origin` is
`github.com/ethereumclassic/core-geth`, no `upstream` remote is configured, and
`main` carries zero commits `origin/main` does not already have: this branch is
fully pushed to the organization's own remote, not a fork's local staging area.
Root wiring artifacts committed here land in this organization's own tree,
never in a separate maintainer's review queue. (Three auxiliary local-only
branches — `main-wired`, `pre-condense-backup`, `reference/modernize-go-1.26` —
do carry commits `origin` does not have; none of them is what this file
describes or what gets wired.) Recorded 2026-09-04 so a later session does not
re-ask.

## Branching

| Branch | What it is |
|---|---|
| `main` | the default branch — what workflows fire on, what pull requests target, and where releases are cut |
| `archive-etclabscore-2024-12` | the preserved pre-migration history, kept indefinitely |

- **`main` is cut from `archive-etclabscore-2024-12`** (`7ef3ecd7a`, 2024-12-16),
  the last commit before this repository was created on 2024-12-21. The annotated
  tag `archive/etclabscore-2024-12` marks the same commit.
- **`master` was deleted when `main` became the default.** Do not recreate it, and
  do not add a `master` trigger to a workflow: a trigger naming a branch that does
  not exist fires on nothing and reports nothing when it stops working. Anything
  described elsewhere as living "on `master`" is now on the archive branch alone.

Confirm which branch you are on before reading anything as current:

```bash
git rev-parse --abbrev-ref HEAD
```

## Toolchain

**Several places declare a Go version, and the three that build release
artifacts must name the same patch.**

| Where | Declares | Governs |
|---|---|---|
| `go.mod` | the language version | what the module compiles against |
| `.github/workflows/release-packages.yml` | `go-version`, an **exact patch** | the Linux x86_64, macOS and Windows release archives |
| `build/checksums.txt` | `# version:golang`, an **exact patch** | what `build/ci.go install -dlgo` downloads, which the ARM release archives build with |
| `Dockerfile`, `Dockerfile.alltools` | a `golang:` tag naming an **exact patch** | the container images |
| the other `.github/workflows/*.yml` | `go-version`, a **major** only | CI tests and checks; resolves to whatever patch the runner has |

Read the current values from those files rather than from this table.

**Change the release patch in all three release places together.** Each builds
part of one release's binaries, so a patch changed in only one of them ships a
release built by two toolchains.

`build/checksums.txt` also pins the linter release. `make lint` downloads that
exact version; it is never read from a locally installed copy.

**The Go module path is `github.com/ethereum/go-ethereum`, not a core-geth path.**
Every internal import uses it. That is deliberate downstream compatibility — do
not "fix" it, and do not be surprised when a package path does not match the
repository name.

## Commands

Every command below is defined in the `Makefile` or in `build/ci.go`. There is no
task runner other than `make`, and no command exists that is not listed here or
printed by `make help`.

**Cost note, since a tool that reads this file will attempt to run what it
lists:** `make all`, `make test` and every `test-coregeth*` target are long,
CPU-bound Go builds over a large tree — `make test` alone carries a 20-minute
timeout. Run one at a time. `CLAUDE.md` carries the fuller resource-discipline
note for Claude Code specifically; that guidance does not reach a tool reading
only this file.

```bash
make geth       # build cmd/geth into ./build/bin/geth
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

**Which workflow fires on what is not uniform, and the branch names are
mid-transition.** Verify against the workflow file rather than assuming a push
was tested:

| Fires on | Workflows |
|---|---|
| push to `main`, every pull request, dispatch | `test-linux.yml` — lint plus both suites |
| push to `main`, path-filtered to the docs | `docs-deploy.yml` |
| push to `main`, pull requests targeting `main`, dispatch | `evmc.yml` |
| dispatch only | the three `bench-*.yml` |
| every pull request, unqualified | `go-generate-check.yml` |
| pull requests targeting `main` touching `params/bootnode*`, plus a daily schedule | `audit-bootnodes.yml` |
| a `v*` tag | `release-packages.yml`, `docker-publish.yml` |

**`evmc.yml` had to move with the default branch, and the reason generalizes.**
The ruleset protecting the default branch requires this job's
`EVMC/EVM+EWASM State Tests` check. While the workflow was scoped to `master`
and `main` had become the default, that check could never report — so every pull
request was unmergeable, including the one that would fix the scoping. A
required check and the workflow producing it must name the same branch.

**The three `bench-*.yml` are dispatch-only, and that is deliberate.** Their push
trigger named `master` and was already dead once `main` became the default.
Re-aiming it would have spent up to six hours of runner time each — they carry
`timeout-minutes: 360` — on every push, for jobs that have never run here. Start
one when a benchmark is the question being asked. What was rejected is leaving a
trigger pointed at a branch that is going away, which reports nothing when it
stops working.

**`.travis.yml`, `circle.yml`, `appveyor.yml` and `Jenkinsfile` were removed from
`main`** on 2026-08-30 (`55ca851c2`, `100a0c6c7`) and are absent here. They
survive on the archive branch only — dead CI configs from before this
repository's migration, kept as history rather than as anything that runs.

## Structure

```
cmd/geth         the node binary; cmd/utils/flags.go defines the network flags
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
`ECBP1100FBlock`, which activates MESS. **`ECBP1100DeactivateFBlock` is 19,250,000
on Classic and 10,400,000 on Mordor, the default ECBP-1110 recommends: MESS is off
from those heights unless the operator passes `--mess`. `params/config_etc_test.go`
asserts both blocks, and a missing or changed block fails it.** ECBP-1100 is an Ethereum Classic
Best Practice, not a consensus rule — it changes which of two competing chains
this node prefers, never whether a block is valid. The Istanbul-equivalent set is
labeled `// ECIP-1088` in a comment rather than carried as its own field.

Reading an activation block out of this file is the only reliable way to know
what is configured. Do not restate a fork schedule from memory, and do not infer
one network's rules from another's.

## Version

`params/version.go` is the single source, and `VersionName` is `CoreGeth`.
Release tooling reads it; nothing else should hard-code a version string —
including this file, which quoted one and went stale at the first bump. Read the
constants there rather than repeating them.

`VersionMeta` carries the release stage and advances `unstable` → `RC1`, `RC2`,
… → `stable`, set when each tag is cut. It must never be empty: the archive and
version helpers branch on `!= "stable"`, so an empty value yields a malformed
`1.13.0--<commit>` rather than a clean one.

## Dependency updates

**`.github/dependabot.yml` exists and version updates are deliberately off**
(`open-pull-requests-limit: 0`). Recorded 2026-09-01.

- **Five ecosystems name something this repository actually holds** — `gomod`
  (`go.mod`, `go.sum`), `pip` (`requirements-mkdocs.txt`), `docker` (two
  Dockerfiles), `github-actions` (every workflow under `.github/workflows/`) and
  `gitsubmodule` (the submodules in `.gitmodules`).
- **A key existing is not the same as the surface being covered.** `gomod`,
  `pip` and `github-actions` support Dependabot security updates; `docker` and
  `gitsubmodule` do not, at any setting. For those two, nothing here and no
  repository toggle delivers a patch — their coverage, if any, is external to
  Dependabot entirely.
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

**Every `uses:` reference in the workflows is pinned to a full commit SHA**, with
the human-readable version in a trailing comment. A tag is a mutable reference
and a commit is not, so the comment is a label and the SHA is the contract.

Two things follow. Bumping one means resolving the new SHA, reading the diff at
that commit, and updating both the SHA and its comment — the pin defends against
a compromised upstream commit, and skipping the diff gives that up. And the
comment can drift from the SHA without anything failing, so trust the SHA.

## Facts that mislead if you do not know them

- **A Windows build writes `zz_resources_windows_amd64.syso` into each `cmd/`
  package** and removes it afterwards. It is the executable's version information
  and manifest, generated by `internal/build/winres.go`, and `.gitignore` covers
  it. A leftover one after an interrupted build is harmless and is overwritten.
- **`swarm/` was removed from `main`** on 2026-08-30 (`55ca851c2`); it no longer
  exists in this tree. It survives on the archive branch only.
- **The `sync-parity-chainspecs` target was removed from the `Makefile`.** It
  invoked a script this repository does not contain, so it could never run, and
  Parity configuration support is not maintained past the Istanbul fork. Do not
  restore it without restoring the script it needs.
- **`AUTHORS` is generated, not written.** `build/update-license.go` produces it
  from `git shortlog -s -n -e`, canonicalized through `.mailmap`. Nothing runs it
  — no `make` target, no CI job — so the file is a stale snapshot. **Never edit it
  by hand**; the header line is a constant in the generator and a hand edit is
  reverted on the next run.
- **That generator also rewrites the license header of every source file** from a
  template hardcoded to `The go-ethereum Authors`, and no file attributed to
  `The core-geth Authors` is in its skip list. Running it as it stands deletes
  that attribution. Do not run it without deciding first what it should assert.
- **`SECURITY.md` is this project's own policy and is the reporting path.** It was
  upstream's until this release series, directing reports to the Ethereum
  Foundation under the Foundation's PGP key; it now routes them to this
  repository's private advisories or to `security@ethereumclassic.com`. Cite it.
- **`geth version-check` queries go-ethereum's vulnerability feed** and prints
  `No vulnerabilities found` when nothing matches. That feed does not track this
  client, so a clean result from it says nothing about this client.

  **It used to report the opposite error, and the guard against that is load
  bearing.** The advisories' patterns begin `Geth/` and are unanchored, and this
  client identifies as `Core-Geth/` — which contains `Geth`. An unanchored search
  therefore matched a substring of our own name and reported a go-ethereum
  advisory against a release carrying the fix: measured, GETH-2024-01 at severity
  High against `Core-Geth/v1.13.0`. `cmd/geth/version_check.go` now requires the
  match to begin at position 0. Do not "simplify" that back to `MatchString`, and
  do not anchor by prepending `^` to the pattern — alternation binds loosely, so
  the anchor would apply to only the first branch.

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
