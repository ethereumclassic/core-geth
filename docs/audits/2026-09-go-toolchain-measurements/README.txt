Go toolchain audit, September 2026: the measurements
====================================================

These are the outputs behind docs/audits/2026-09-go-toolchain.md, published so
its figures can be checked by reading them rather than by repeating the scans.
No downloaded binary was executed to produce them.

binaries.tsv
    One row per binary measured: the repository, release tag and archive
    platform, the file inside the archive, its SHA-256 and size, and the Go
    version, GOOS, GOARCH, GOARM and CGO_ENABLED from its build information.

go-version-m/<tag>_<platform>.txt
    The output of `go version -m` for each binary: the Go version it was built
    with and every module compiled into it. 59 archives, from v1.12.0 to
    v1.13.0.

govulncheck-symbol/<tag>_<platform>.json.gz
govulncheck-module/<tag>_<platform>.json.gz
    govulncheck v1.7.0 in binary mode with -json, at symbol level and at
    module level, on the 41 archives of v1.12.20 to v1.13.0.

docker/
    The same for geth in the two v1.13.0 image tarballs, amd64 and arm64:
    `go version -m`, and govulncheck at module level.

Every scan used one snapshot of the Go vulnerability database: vuln.go.dev as
modified 2026-09-15T18:39:25Z, 4,445 entries, SHA-256
fe51e51645e6b1e946731cef048339972256c6417de91d4cf3985fc87b4aa89b.

The "db" field in each scan's config named the local copy of that snapshot on
the machine that ran it. It is replaced with https://vuln.go.dev, and the
"db_last_modified" field beside it identifies the snapshot. Nothing else in
the files is changed.
