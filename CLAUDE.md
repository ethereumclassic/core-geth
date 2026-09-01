@AGENTS.md

# Claude Code notes

`AGENTS.md` above is the project context and is imported, not summarized. What
follows is Claude-specific and adds to it; where the two could be read as
disagreeing, `AGENTS.md` wins.

## Before acting

- **Confirm the branch first.** `main` and `master` differ, and reading the wrong
  one returns a confident wrong answer with nothing to signal it:
  `git rev-parse --abbrev-ref HEAD`.
- **Read `params/` before saying anything about chain rules.** Fork schedules are
  data in this repository. A schedule recalled from memory is a guess, and here a
  guess about consensus is a chain split.
- **This is consensus software for a live network.** `AGENTS.md` lists what needs
  confirmation. That list is not advisory.

## Working here

- **One heavy task at a time.** `make all`, `make test` and `make test-coregeth`
  are long, CPU-bound Go builds over a large tree. Do not start a second one
  alongside the first, and do not run them in parallel with other heavy work.
- **Build before testing.** `make test` depends on `make all`; the CoreGeth
  regression target builds `core-geth` first.
- **Submodules before tests.** `git submodule update --init --recursive`, or the
  consensus suites fail in ways that look like real failures.
- **Small batches, then verify.** Change a few things, run the relevant target,
  confirm reality matches expectation before continuing.

## Model selection

Select by role, using the tier alias — an alias resolves to the newest model in
its tier, so it does not go stale the way a model name does.

| Alias | Use it for |
|---|---|
| `haiku` | typo fixes, mechanical edits, reading a file back |
| `sonnet` | the default: documentation, build and CI edits, review |
| `opus` | consensus rules, chain configuration, cryptography, p2p, state handling |

Switch with `/model haiku`, `/model sonnet`, `/model opus`. Current prices,
context windows and the model behind each alias are documented at
<https://platform.claude.com/docs/en/about-claude/models/overview>.

## Machine-local context

Anything specific to one contributor's machine — local checkout paths, personal
tool locations, private working notes — belongs in `CLAUDE.local.md` at this
repository's root, never in this file or in `AGENTS.md`. Both of those are
committed and travel to every clone.

Confirm `CLAUDE.local.md` is actually ignored before writing one, by effect
rather than by reading the patterns:

```bash
git check-ignore --no-index -q -- CLAUDE.local.md && echo ignored || echo "NOT ignored"
git check-ignore --no-index -q -- README.md && echo "check is broken" || echo "check discriminates"
```

The second line is the calibration. A check that cannot report "not ignored" is
not checking anything, and its all-clear means nothing.

## Response style

- Code and commands first; explain only what was asked about.
- Concise bullets over paragraphs. Tables for comparisons.
- No pleasantries, and do not repeat the prompt back.
- State what you verified and how. "Tests pass" without naming what ran is not a
  result.
