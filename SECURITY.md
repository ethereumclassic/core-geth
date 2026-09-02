# Security Policy

## Reporting a vulnerability

**Do not open a public issue for a security vulnerability.**

Report it privately through GitHub's private vulnerability reporting on this
repository:

<https://github.com/ethereumclassic/core-geth/security/advisories/new>

That channel is enabled and reaches this project's maintainers. It needs a GitHub
account and nothing else — no key exchange and no mailing list.

Include what you can: the affected version or commit, the network, what an
attacker gains, and a reproduction if you have one. A report without a
reproduction is still worth sending.

### What to expect

Reports are acknowledged and triaged privately. Where a fix is warranted it is
prepared under a GitHub security advisory and disclosed once a release carrying
it is available. Reporters are credited unless they ask not to be.

## Scope

This repository is Core-Geth, the Ethereum Classic execution client. Consensus
rules, chain configuration, peer-to-peer networking, the JSON-RPC surface, key
handling and the node binaries are all in scope.

Vulnerabilities in other Ethereum Classic software — other clients, explorers,
wallets, bridges or infrastructure — are outside this repository. Report those to
the project that maintains them.

## Core-Geth is a derivative of go-ethereum

A defect here may be inherited from upstream rather than specific to this client.
Where it affects code shared with
[go-ethereum](https://github.com/ethereum/go-ethereum), it is worth reporting
there as well: <https://github.com/ethereum/go-ethereum/security/policy>.

**Neither project forwards reports to the other.** Reporting to both is welcome
and is the right move for a shared defect.

## `geth version-check` checks go-ethereum, not this client

The `version-check` command queries go-ethereum's vulnerability feed at
`geth.ethereum.org` and matches advisories against the running version string.

**That feed does not track Core-Geth**, and Core-Geth's version numbering has
diverged from go-ethereum's, so its output is not meaningful here in either
direction: a `No vulnerabilities found` result says nothing about this client,
and a reported match may describe a defect this client never carried.

Do not treat that command as a clean bill of health for Core-Geth.
