---
hide:
  - toc        # Hide table of contents
description: "How Core-Geth is versioned: the version components and release-stage metadata in params/version.go, and how archive and release names are formed."
---

# Versioning

`ethereumclassic/core-geth` uses [Semantic Versioning](https://semver.org). The API definition that would demand increments to the major version is basically nil;
it can be expected that a major version bump would be accompanied by an entirely new repository and name.

`VersionMeta` (`params/version.go`) carries the release stage. It advances `unstable` →
`RC1`, `RC2`, ... → `stable` as a release is prepared, must never be empty, and is set once
per tag; `geth version` prints it appended to the semantic version, for example
`1.13.0-RC2` for a release candidate. See [Publishing a release](create-new-release.md)
for when each stage is set.

!!! note "See also"

    You can find some historical discussions on versioning at the following links.

    - https://github.com/etclabscore/core-geth/pull/29#issuecomment-588977383
    - https://github.com/etclabscore/multi-geth-fork/issues/153
    - https://github.com/etclabscore/core-geth/pull/30#issuecomment-591979271
    - https://github.com/etclabscore/core-geth/issues/83
