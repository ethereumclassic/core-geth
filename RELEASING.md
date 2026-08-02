# Releasing core-geth with verifiable artifacts

The release workflow lives at `.github/workflows/release-ghcr.yml`. It is the only path that produces core-geth release artifacts: it builds the release archives, publishes and signs the container image, attests everything, and drafts the GitHub Release.

## Trust model

Release verification is pinned to this workflow file and to the ref that ran it.

- Tagged releases must verify against `https://github.com/ethereumclassic/core-geth/.github/workflows/release-ghcr.yml@refs/tags/<tag>`.
- Manual dry runs must verify against `...@refs/heads/<branch>`; pull request dry runs against `...@refs/pull/<n>/merge`.
- Protected release tags keep unreviewed refs from minting trusted release artifacts.
- Anything pushed outside this workflow — including images pushed by org admins — will fail verification because it will not have both the expected workflow identity and GitHub attestation provenance.

## Triggers

- Push a tag matching `v*` to build, push, sign, and attest the release image and release archives, then draft the GitHub Release.
- Run `Release GHCR artifacts` with `dry_run=true` to build and attest without pushing, so reviewers can exercise the release path before cutting a tag.
- Opening or updating a pull request that touches the release path itself (`release-ghcr.yml`, `Dockerfile`, `build/archive-signing.sh`, `build/ci.go`) runs the same dry run automatically, so a broken release surfaces on the PR instead of at tag time.

Dry runs use `dry-run-<sha7>` in artifact names, never publish to GHCR, and never
create a release. Pull requests from forks build but skip attestation, because a
fork's token carries no OIDC identity to sign with.

## What a release contains

Five platforms, each with a `geth` archive and an `alltools` archive plus a
`.sha256` sidecar:

| archive prefix | target |
| --- | --- |
| `core-geth-linux-` | linux/amd64 |
| `core-geth-arm64-` | linux/arm64 |
| `core-geth-osx-amd64-` | darwin/amd64 |
| `core-geth-osx-arm64-` | darwin/arm64 |
| `core-geth-win64-` | windows/amd64 |

Plus a consolidated `SHA256SUMS-<version>.txt` and the multi-arch container
image `ghcr.io/ethereumclassic/core-geth`.

The macOS archives carry their architecture in the name. The older bare `osx`
name silently changed meaning when GitHub moved `macos-latest` to Apple Silicon,
so an `osx` archive could be either an Intel or an arm64 binary depending on when
it was cut.

The 32-bit `arm`, `arm5`, `arm6`, and `arm7` targets are not built. `arm` was
already announced as deprecated in `build/archive-signing.sh` for removal at
v1.12.7.

`core-geth-osx-amd64-` is built with blst's portable C path
(`-D__BLST_NO_ASM__`), because blst's x86_64 assembly no longer assembles under
the Clang on current macOS runner images. BLS and KZG operations are slower in
that one archive. Every other platform, including `osx-arm64`, uses blst's
assembly as normal.

## How the image is built

Each platform is built on a native runner — `linux/amd64` on `ubuntu-latest`,
`linux/arm64` on `ubuntu-24.04-arm` — and the per-platform digests are merged
into a manifest list with `docker buildx imagetools create`. Nothing is built
under QEMU emulation; a cgo build of go-ethereum for arm64 under emulation takes
hours and risks the job timeout.

The consequence for verification is that a real release has **one index digest**
covering both platforms. That index digest is what gets signed and attested.

Dry runs cannot produce a manifest list, because there is no registry to merge
digests in. Instead each platform is exported as an OCI archive, uploaded as the
`image-oci-<platform>` workflow artifact, and attested **as a file**. Verify a
dry-run image with the file recipe below, not the `oci://` recipe.

## GHCR package visibility

After the first successful push, set `ghcr.io/ethereumclassic/core-geth` to **public** in the GitHub Packages UI. GHCR creates new packages as private by default, and reviewers pulling a verified image will otherwise get an authentication error that looks like a workflow problem.

## Verification inputs

```sh
REPO=ethereumclassic/core-geth
WORKFLOW=.github/workflows/release-ghcr.yml
TAG=v1.12.34
BRANCH=master
OIDC_ISSUER=https://token.actions.githubusercontent.com
TAG_IDENTITY="https://github.com/${REPO}/${WORKFLOW}@refs/tags/${TAG}"
BRANCH_IDENTITY="https://github.com/${REPO}/${WORKFLOW}@refs/heads/${BRANCH}"
IMAGE=ghcr.io/ethereumclassic/core-geth
DIGEST=sha256:<manifest-digest>
```

Use the tag identity for real releases and the branch identity for manual dry runs.

## Verify the container image provenance

Use `DIGEST` from the manifest list, not a per-platform digest:

```sh
DIGEST="$(docker buildx imagetools inspect "${IMAGE}:${TAG}" \
  --format '{{json .Manifest}}' | jq -r '.digest')"

gh attestation verify "oci://${IMAGE}@${DIGEST}" \
  --repo "${REPO}" \
  --cert-identity "${TAG_IDENTITY}" \
  --cert-oidc-issuer "${OIDC_ISSUER}"
```

## Verify the container image signature

The index is signed with `--recursive`, so both the index digest and each
per-platform manifest digest carry a signature.

```sh
cosign verify "${IMAGE}@${DIGEST}" \
  --certificate-identity "${TAG_IDENTITY}" \
  --certificate-oidc-issuer "${OIDC_ISSUER}"
```

## Verify a dry-run image

Download the `image-oci-linux-amd64` / `image-oci-linux-arm64` artifacts from
the run, then verify the archives as files against the branch identity:

```sh
gh attestation verify "core-geth-image-linux-amd64-dry-run-<sha7>.tar" \
  --repo "${REPO}" \
  --cert-identity "${BRANCH_IDENTITY}" \
  --cert-oidc-issuer "${OIDC_ISSUER}"
```

To inspect what was actually built, load the archive locally:

```sh
# Docker 25+ with the containerd image store:
docker load -i "core-geth-image-linux-amd64-dry-run-<sha7>.tar"

# Otherwise:
skopeo copy \
  "oci-archive:core-geth-image-linux-amd64-dry-run-<sha7>.tar" \
  "docker-daemon:core-geth:dry-run"
```

## Verify release binaries and checksums

Download the assets from the GitHub Release (or, for a dry run, the
`release-binaries-<version>` workflow artifact), then verify every archive and
the consolidated checksum file:

```sh
gh release download "${TAG}" --repo "${REPO}"

for f in core-geth-*-"${TAG}".zip "SHA256SUMS-${TAG}.txt"; do
  gh attestation verify "$f" \
    --repo "${REPO}" \
    --cert-identity "${TAG_IDENTITY}" \
    --cert-oidc-issuer "${OIDC_ISSUER}" || echo "FAILED: $f"
done
```

`--cert-identity` is what makes this meaningful. Without it, `gh attestation
verify` accepts any attestation from the repository, including one minted by a
different workflow.

Once the checksum file is verified, the archives can be checked against it
directly:

```sh
sha256sum -c "SHA256SUMS-${TAG}.txt"
```

For a dry run, switch the filenames to `dry-run-<sha7>` and replace `TAG_IDENTITY` with `BRANCH_IDENTITY`.
