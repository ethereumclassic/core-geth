# Releasing core-geth with verifiable artifacts

The release workflow lives at `.github/workflows/release-ghcr.yml` and is the only path that should produce verifiable release artifacts for GHCR.

## Trust model

Release verification is pinned to this workflow file and to the ref that ran it.

- Tagged releases must verify against `https://github.com/ethereumclassic/core-geth/.github/workflows/release-ghcr.yml@refs/tags/<tag>`.
- Manual dry runs must verify against `https://github.com/ethereumclassic/core-geth/.github/workflows/release-ghcr.yml@refs/heads/<branch>`.
- Protected release tags keep unreviewed refs from minting trusted release artifacts.
- Anything pushed outside this workflow — including images pushed by org admins — will fail verification because it will not have both the expected workflow identity and GitHub attestation provenance.

## Triggers

- Push a tag matching `v*` to build, push, sign, and attest the release image and release archives.
- Run `Release GHCR artifacts` with `dry_run=true` to build and attest without pushing, so reviewers can exercise the release path before cutting a tag.

Dry runs use `dry-run-<sha7>` in artifact names and never publish to GHCR.

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

```sh
gh attestation verify "oci://${IMAGE}@${DIGEST}" \
  --repo "${REPO}" \
  --cert-identity "${TAG_IDENTITY}" \
  --oidc-issuer "${OIDC_ISSUER}"
```

## Verify the container image signature

```sh
cosign verify "${IMAGE}@${DIGEST}" \
  --certificate-identity "${TAG_IDENTITY}" \
  --certificate-oidc-issuer "${OIDC_ISSUER}"
```

## Verify release binaries and checksums

Download the `release-binaries-<version>` artifact from the workflow run, then verify the archives and consolidated checksum file.

```sh
gh attestation verify "core-geth-linux-${TAG}.zip" \
  --repo "${REPO}" \
  --cert-identity "${TAG_IDENTITY}" \
  --oidc-issuer "${OIDC_ISSUER}"

gh attestation verify "core-geth-arm64-${TAG}.zip" \
  --repo "${REPO}" \
  --cert-identity "${TAG_IDENTITY}" \
  --oidc-issuer "${OIDC_ISSUER}"

gh attestation verify "SHA256SUMS-${TAG}.txt" \
  --repo "${REPO}" \
  --cert-identity "${TAG_IDENTITY}" \
  --oidc-issuer "${OIDC_ISSUER}"
```

For a dry run, switch the filenames to `dry-run-<sha7>` and replace `TAG_IDENTITY` with `BRANCH_IDENTITY`.
