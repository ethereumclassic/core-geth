---
title: Installation
description: "Download, verify and install a Core-Geth release on Linux, macOS, Windows or Docker, and check an archive against its published checksum."
---

!!! danger "Security advisory: the v1.12.x line"
    Core-Geth moved to
    [`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth) in
    December 2024. Archives and images under the previous `etclabscore` namespace are
    not built from this source and do not carry the fixes released here.

    **Affected:** all `v1.12.x` releases. **Fixed:** `v1.13.0`.

    | Identifier | Severity | What it allows |
    | --- | --- | --- |
    | [CVE-2026-22862](../audits/2026-03-security-audit.md#cve-2026-22862-ecies-elliptic-curve-integrated-encryption-scheme-decrypt-ciphertext-length-undercheck) | High | Any peer crashes the node during the RLPx handshake. **Exploited against Ethereum Classic bootnodes, 18 March 2026.** |
    | [CVE-2026-26315](../audits/2026-03-security-audit.md#cve-2026-26315-ecies-generateshared-accepts-unvalidated-public-key) | High | Repeated handshakes leak bits of the node's own P2P key |
    | [CVE-2026-26314](../audits/2026-03-security-audit.md#cve-2026-26314-secp256k1-isoncurve-field-boundary-bypass) | High | secp256k1 coordinate field boundary bypass |
    | [CVE-2025-24883](../audits/2026-03-security-audit.md#cve-2025-24883-off-curve-public-key-in-unmarshalpubkey) | High | Off-curve public key accepted by `UnmarshalPubkey` |
    | [CVE-2026-26313](../audits/2026-03-security-audit.md#cve-2026-26313-p2p-rlp-recursive-length-prefix-item-count-memory-exhaustion) | High | One crafted p2p message exhausts the node's memory |
    | [CVE-2026-22868](../audits/2026-03-security-audit.md#cve-2026-22868-kzg-kate-zaverucha-goldberg-blob-proof-verification-dos) | Medium | KZG proof verification denial of service |
    | [GraphQL query depth](../audits/2026-03-security-audit.md#graphql-query-depth-dos) | Medium | Unbounded GraphQL query nesting; no CVE identifier assigned |

    Every `v1.12.x` archive measured was also built on a Go version that is no longer
    supported, most recently [Go 1.21 and Go 1.22](../audits/2026-09-go-toolchain.md),
    whose support ended in August 2024 and February 2025.

    Later releases in the line backport some of these; the
    [March 2026 audit](../audits/2026-03-security-audit.md) records which, and the
    [migration guide](../tutorials/v1.13.0-migration.md) covers the upgrade.

## Choose your route

Pick the one row that matches your machine and follow only that section. You do not
need the others.

| Your machine | Route |
| --- | --- |
| Linux, x86_64: most servers and desktops | [Linux, x86_64](#linux-x86_64) |
| Linux, 64-bit ARM: Raspberry Pi on a 64-bit OS, AWS Graviton, Ampere | [Linux, ARM](#linux-arm) |
| Linux, 32-bit ARM: older Raspberry Pi and embedded boards | [Linux, ARM](#linux-arm) |
| macOS, Apple Silicon | [macOS](#macos) |
| macOS, Intel | [macOS](#macos) |
| Windows, x86_64 | [Windows](#windows) |
| Anything, in a container | [Docker](#docker) |
| Anything, compiled yourself | [Build from source](#build-from-source) |

On Linux and macOS, `uname -sm` prints the operating system and machine architecture,
which is what the rows above distinguish.

**The executable is named `geth`.** `core-geth` appears throughout as the project,
the repository, the container image and the archive filenames. There is no
executable by that name.

Every archive route ends with a `geth` on your `PATH`. Archives are attached to each
tagged release on the
[releases page](https://github.com/ethereumclassic/core-geth/releases), and every one
of them is a `.zip` containing the binary and `COPYING`, the license, at the top level,
published alongside a `.sha256` file.

## Linux, x86_64

```shell
$ VERSION=v1.13.0
$ BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
$ curl -LO $BASE/core-geth-linux-$VERSION.zip
$ curl -LO $BASE/core-geth-linux-$VERSION.zip.sha256
$ sha256sum -c core-geth-linux-$VERSION.zip.sha256
$ unzip core-geth-linux-$VERSION.zip
$ sudo install -m 0755 geth /usr/local/bin/geth
$ geth version
```

`sha256sum -c` must print `OK`. If it prints `FAILED`, stop: delete the file and
download it again rather than running it.

**This archive needs glibc 2.17 or newer**: RHEL 7, CentOS 7 and anything more
recent. `ldd --version` prints what you have. On an older distribution
the binary will not start; [build from source](#build-from-source) there instead.

Before running it on a node that holds value, also
[verify where the archive came from](#verify-where-the-archive-came-from).

## Linux, ARM

Four ARM archives are published. `uname -m` tells you which one you need:

| `uname -m` prints | Download | Needs glibc |
| --- | --- | --- |
| `aarch64` or `arm64` | `core-geth-arm64-<tag>.zip` | 2.17 or newer |
| `armv7l` | `core-geth-arm7-<tag>.zip` | 2.28 or newer |
| `armv6l` | `core-geth-arm6-<tag>.zip` | 2.28 or newer |
| `armv5tel` and other `armv5` | `core-geth-arm5-<tag>.zip` | 2.28 or newer |

**A 64-bit capable board does not mean a 64-bit archive.** Raspberry Pi OS reports
`aarch64` in its 64-bit edition and `armv7l` in its 32-bit edition on the same
hardware. Follow what `uname -m` actually prints.

Then, substituting your architecture for `arm64` below:

```shell
$ VERSION=v1.13.0
$ ARCH=arm64
$ BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
$ curl -LO $BASE/core-geth-$ARCH-$VERSION.zip
$ curl -LO $BASE/core-geth-$ARCH-$VERSION.zip.sha256
$ sha256sum -c core-geth-$ARCH-$VERSION.zip.sha256
$ unzip core-geth-$ARCH-$VERSION.zip
$ sudo install -m 0755 geth /usr/local/bin/geth
$ geth version
```

A `core-geth-arm-<tag>.zip` is published as well. It is byte-identical to the `arm5`
archive and exists only so that existing scripts referring to it keep working: use
`arm5`.

Before running it on a node that holds value, also
[verify where the archive came from](#verify-where-the-archive-came-from).

## macOS

Two archives, chosen by processor:

| `uname -m` prints | Download |
| --- | --- |
| `arm64`: Apple Silicon, M1 and later | `core-geth-osx-arm64-<tag>.zip` |
| `x86_64`: Intel | `core-geth-osx-<tag>.zip` |

macOS has no `sha256sum`; use `shasum -a 256 -c`, which reads the same file.

```shell
$ VERSION=v1.13.0
$ PLATFORM=osx-arm64            # or: osx
$ BASE=https://github.com/ethereumclassic/core-geth/releases/download/$VERSION
$ curl -LO $BASE/core-geth-$PLATFORM-$VERSION.zip
$ curl -LO $BASE/core-geth-$PLATFORM-$VERSION.zip.sha256
$ shasum -a 256 -c core-geth-$PLATFORM-$VERSION.zip.sha256
$ unzip core-geth-$PLATFORM-$VERSION.zip
$ sudo install -m 0755 geth /usr/local/bin/geth
$ geth version
```

!!! tip "If macOS refuses to run it"
    These binaries are not code-signed, so Gatekeeper blocks one that carries the
    quarantine attribute with *"cannot be opened because the developer cannot be
    verified"*. Browsers set that attribute; `curl` does not, so the commands above
    avoid it. If you downloaded through a browser, clear it:

    ```shell
    $ xattr -d com.apple.quarantine geth
    ```

    Once macOS has blocked it, **Open Anyway** in System Settings, under Privacy & Security,
    also lets it run ([Apple's instructions](https://support.apple.com/en-us/102445)).

Before running it on a node that holds value, also
[verify where the archive came from](#verify-where-the-archive-came-from).

## Windows

Download `core-geth-win64-<tag>.zip` and its `.sha256`. In PowerShell:

```powershell
> $VERSION = "v1.13.0"
> $BASE = "https://github.com/ethereumclassic/core-geth/releases/download/$VERSION"
> Invoke-WebRequest "$BASE/core-geth-win64-$VERSION.zip" -OutFile "core-geth-win64-$VERSION.zip"
> Invoke-WebRequest "$BASE/core-geth-win64-$VERSION.zip.sha256" -OutFile "core-geth-win64-$VERSION.zip.sha256"
> (Get-FileHash "core-geth-win64-$VERSION.zip" -Algorithm SHA256).Hash
> Get-Content "core-geth-win64-$VERSION.zip.sha256"
> Expand-Archive "core-geth-win64-$VERSION.zip" -DestinationPath .
> .\geth.exe version
```

**PowerShell has no equivalent of `sha256sum -c`, so compare the two lines yourself.**
`Get-FileHash` prints uppercase hex and the `.sha256` file holds lowercase followed by
the filename; the letters must match, the case does not.

The archive contains `geth.exe` and `COPYING`, the license. Move `geth.exe` wherever you keep command-line tools
and add that directory to `PATH` if you want to run it by name.

!!! tip "If Windows warns about the publisher"
    These binaries are not code-signed, so Windows shows their publisher as unknown. That
    appears in the firewall alert on the first start, and in the *"Windows protected your
    PC"* screen if Microsoft Defender SmartScreen stops a downloaded `geth.exe`. Check the
    archive first, against its checksum above and
    [where it came from](#verify-where-the-archive-came-from). Then allow the firewall
    access, and in the SmartScreen screen select **More info**, then **Run anyway**. A
    browser marks what it downloads, and removing the mark from the archive before you
    extract it avoids the SmartScreen screen:

    ```powershell
    > Unblock-File "core-geth-win64-$VERSION.zip"
    ```

Before running it on a node that holds value, also
[verify where the archive came from](#verify-where-the-archive-came-from).

## Verify where the archive came from

**A checksum proves the download was not corrupted. It does not prove what built the
file**, because whoever publishes the archive publishes the checksum beside it. From
v1.13.0 each archive also carries a build attestation, which records the workflow,
commit and ref that produced it:

```shell
$ gh attestation verify core-geth-linux-v1.13.0.zip --repo ethereumclassic/core-geth
```

That needs the [GitHub CLI](https://cli.github.com/). It succeeds only for an artifact
built by this repository's release workflow. An archive from anywhere else fails, as
does one whose bytes have changed. No signing key is involved: the attestation is minted
against a short-lived certificate issued to the workflow run itself, so there is no key
to steal or rotate.

Worth doing once on any binary you are about to run on a node that holds value.

## Docker

Images are published to the GitHub Container Registry for each tagged release,
built for `linux/amd64` and `linux/arm64`:

```shell
$ docker pull ghcr.io/ethereumclassic/core-geth:latest
```

Tags mirror the release tags. `ghcr.io/ethereumclassic/core-geth:v1.13.0` is a
specific release; `latest` follows the most recent non-prerelease. The image
built from `Dockerfile.alltools`, containing the full tool set rather than the
node alone, is published under the same name with an `alltools-` prefix, as
`ghcr.io/ethereumclassic/core-geth:alltools-latest`.

!!! warning "Images under the previous namespace are not these images"
    Images published as `etclabscore/core-geth` on Docker Hub are not built from
    this source and receive nothing released here.

!!! note "If the pull returns `unauthorized`"
    The images are published but not yet public. Until they are, load the image from the
    release or build it locally, as below. It runs the same way.

### Load the image from the release

Each release also carries its images as files, one per architecture:
`core-geth-docker-amd64-<tag>.tar.gz` and `core-geth-docker-arm64-<tag>.tar.gz`, and
`core-geth-alltools-docker-amd64-<tag>.tar.gz` and `core-geth-alltools-docker-arm64-<tag>.tar.gz`
for the full tool set. Each is the published image saved from the registry, with a `.sha256` and
a build attestation. Verify the one for your machine and load it:

```shell
$ gh attestation verify core-geth-docker-amd64-v1.13.0.tar.gz --repo ethereumclassic/core-geth
$ docker load -i core-geth-docker-amd64-v1.13.0.tar.gz
```

The loaded image keeps its published name, `ghcr.io/ethereumclassic/core-geth:v1.13.0`, so the
`docker run` command below works unchanged.

### Build an image from the release archive

This packages the binary the release attested, rather than compiling one. Download
`core-geth-linux-<tag>.zip` for `linux/amd64` or `core-geth-arm64-<tag>.zip` for `linux/arm64`
from the release, [verify it](#verify-where-the-archive-came-from), and unpack it next to this
`Dockerfile`:

```dockerfile
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
COPY geth /usr/local/bin/geth
EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["geth"]
```

```shell
$ unzip core-geth-linux-v1.13.0.zip
$ docker build -t core-geth:v1.13.0 .
$ docker run --rm core-geth:v1.13.0 version
```

Use a glibc-based image such as Debian or Ubuntu. The release binaries link against glibc, so an
Alpine base cannot run them.

### Build an image from source

The `Dockerfile` produces an image containing `geth`, and `Dockerfile.alltools` one containing
the full tool set:

```shell
$ git clone https://github.com/ethereumclassic/core-geth.git
$ cd core-geth
$ docker build -t core-geth:local .
```

### Run the container

Run the image the same way whichever method produced it. Its entry point is the `geth` binary,
so flags are passed straight to the node. `--name` below labels the running container and is
yours to choose:

```shell
$ docker run -d \
    --name core-geth \
    -v $LOCAL_DATADIR:/root \
    -p 30303:30303 -p 30303:30303/udp \
    -p 127.0.0.1:8545:8545 \
    ghcr.io/ethereumclassic/core-geth:latest \
    --classic \
    --http --http.addr 0.0.0.0 --http.port 8545
```

If you built the image yourself, use your own tag (`core-geth:v1.13.0` or `core-geth:local`,
above) in place of `ghcr.io/ethereumclassic/core-geth:latest`.

That maps the devp2p port over both TCP and UDP, keeps chain data in
`$LOCAL_DATADIR` on the host so it survives the container, and reaches the
JSON-RPC endpoint from the host and nowhere else.

!!! warning "Both halves of that RPC line are deliberate"
    `--http` alone binds the listener to `localhost` **inside the container**, while
    Docker forwards a published port to the container's *external* interface, so
    `-p 8545:8545` with the default address publishes a port that reaches nothing, and
    the RPC appears dead with no error explaining why.

    `--http.addr 0.0.0.0` binds every interface inside the container's own network
    namespace, which is isolated. `-p 127.0.0.1:8545:8545` is what keeps it off the
    host's public interfaces. **Publishing as `-p 8545:8545` instead exposes your RPC
    on every interface of the host, which on a VPS is the open internet.**

    If you do not need RPC from the host at all, drop both the `-p 127.0.0.1:8545:8545`
    and the `--http*` flags.

The image also exposes 8546 for the WebSocket endpoint, which needs
`--ws --ws.addr 0.0.0.0` to be served.

## Build from source

Building needs **Go 1.26** and a C compiler. This is also the route for a platform no
archive covers, or a distribution older than the glibc floors above.

```shell
$ git clone https://github.com/ethereumclassic/core-geth.git
$ cd core-geth
$ make geth
$ ./build/bin/geth version
```

**Do not clone with `--recursive`.** The three submodules hold consensus test fixtures,
they are several gigabytes, and `make geth` does not read them. They are needed only to
run the test suites.

`make all` builds every tool rather than the node alone.
[Build from source](../developers/build-from-source.md) covers the full set of targets,
cross-compilation, and building inside Docker.

## What else is in a release

Each platform publishes two archives:

| Archive | Contains |
| --- | --- |
| `core-geth-<platform>-<tag>.zip` | the `geth` node binary, and `COPYING` |
| `core-geth-alltools-<platform>-<tag>.zip` | `geth` plus the other tools built from this source, and `COPYING` |

Every route above uses the first. The `alltools` archive adds every other executable
built from this source (`abidump`, `abigen`, `ancient-store-mem`, `bootnode`,
`clef`, `devp2p`, `echainspec`, `era`, `ethkey`, `evm`, `p2psim` and `rlpdump`) and
installs exactly the same way. Read the archive itself rather than this list if you
need to be certain: the release ships whatever `make all` produced.

Once `geth` is installed, [Running a node](run-a-node.md) takes it to a synced node that keeps
running, with a guide for each platform and the flags for each kind of node.
