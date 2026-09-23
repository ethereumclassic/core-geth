---
hide:
  - toc        # Hide table of contents
title: Publishing a Release
---

# Publishing a release

Releases are cut from **`main`** in
[`ethereumclassic/core-geth`](https://github.com/ethereumclassic/core-geth). The
`master` branch is gone and the previous `etclabscore` remote is not this project's
to publish to.

`v1.13.0-rc1` is used as the worked example below.

## Before the tag

- [ ] **`main` is green.** `make lint`, `make test` and `make test-coregeth` all pass.
      The release artifacts are produced by CI, so a red pipeline means no artifacts
      rather than a bad release.
- [ ] **Set the release stage.** Edit `params/version.go` so `VersionMeta` names the
      stage being cut: `unstable` → `RC1`, `RC2`, … → `stable`. It must never be
      empty: the archive and version helpers branch on `!= "stable"`, so an empty
      value produces a malformed `1.13.0--<commit>` name rather than a clean one.
      Check `gofmt` after editing; the constant block realigns.
- [ ] **Commit and push it**, then let CI finish before tagging. The tag should point
      at a commit CI has already judged.

## Cut the tag

```shell
$ git tag -a v1.13.0-rc1 -m 'Core-Geth v1.13.0-rc1'
$ git push origin v1.13.0-rc1
```

Creating a `v*` tag is restricted to repository admins by a ruleset, so this push
reports a bypass. That restriction is what the build attestation rests on. An
attestation binds an archive to the run that produced it, but only the tag rule
decides who may start such a run.

**A release tag is immutable.** Do not move or delete one after artifacts exist:
the attestations are written to a public append-only log, so a reused tag name ends
up bound to two different artifacts with no way to tell which is which. If a tag is
wrong, cut the next one.

## What the tag triggers

| Workflow | Produces |
| --- | --- |
| `release-packages.yml` | 18 archives: 9 platforms × (`geth`, `alltools`), each with a `.sha256` and a build attestation. A single publish job uploads all of them to one **draft** release once every build has finished |
| `docker-publish.yml` | multi-architecture images for `linux/amd64` and `linux/arm64`, pushed to GHCR, then the images saved as four `.tar.gz` files, each with a `.sha256` and a build attestation, attached to the draft release. It waits for that draft, and stops if more than one release carries the tag |

A tag whose name contains a hyphen is treated as a prerelease: the GitHub release is
marked as one, and the container image does **not** take `:latest`. Only a full
release such as `v1.13.0` takes the moving tag.

## Before publishing the draft

- [ ] **All 18 archives and the four image tarballs are attached**, and each `.sha256` matches its file.
- [ ] **The archive names carry the tag**, not a commit SHA. A bare SHA means the
      tag was not present in the build checkout.
- [ ] **The attestation verifies against a downloaded archive**, which is a
      different claim from the workflow step having gone green:

    ```shell
    $ gh attestation verify core-geth-linux-v1.13.0-rc1.zip --repo ethereumclassic/core-geth
    ```

- [ ] **Load an image tarball.** Its attestation verifies, it loads under the published name, and the
      config digest its `manifest.json` names equals the image ID the release run's save job logged
      for that architecture. Compare those two, not `docker image inspect`'s ID: a Docker using the
      containerd image store reports the manifest digest there instead:

    ```shell
    $ gh attestation verify core-geth-docker-amd64-<tag>.tar.gz --repo ethereumclassic/core-geth
    $ docker load -i core-geth-docker-amd64-<tag>.tar.gz
    $ docker run --rm ghcr.io/ethereumclassic/core-geth:<tag> version
    $ tar -xOzf core-geth-docker-amd64-<tag>.tar.gz manifest.json   # "Config": "blobs/sha256/<digest>"
    ```

- [ ] **Spot-check a binary.** It should report the version you set and a commit
      date, its glibc floor should be the one the build targets rather than the
      runner's, and its Go version should be the patch `build/checksums.txt` pins, the
      same in every archive:

    ```shell
    $ ./geth version
    $ objdump -T geth | grep -o 'GLIBC_[0-9.]*' | sort -uV | tail -1
    $ go version geth
    ```

- [ ] **Every archive carries `COPYING`** beside its binaries:

    ```shell
    $ unzip -l core-geth-linux-<tag>.zip
    ```

- [ ] **The Windows and macOS metadata checks passed.** The release run's *Verify the
      Windows version information* and *Verify the macOS minimum version* steps are
      green. On Windows, the Details tab of `geth.exe`'s Properties shows Core-Geth, the
      version and Ethereum Classic DAO LLC.
- [ ] **The published image index carries its annotations**: description, source,
      vendor and documentation:

    ```shell
    $ docker buildx imagetools inspect --raw ghcr.io/ethereumclassic/core-geth:<tag>
    ```

- [ ] **Regenerate the `--help` dump.** From a local `make geth` at this tag,
      replace the dump in
      [`run-cli.md`](../getting-started/run-cli.md#command-line-options) with
      this command's output in full, never by hand. `HOME='~'` keeps the
      builder's real home directory out of the printed defaults.

      **It is deliberately not kept current between releases.** The dump names a
      version and a commit, so it has to describe a released binary; a flag added or
      changed on `main` reaches this page at the tag and not before:

    ```shell
    $ HOME='~' build/bin/geth --help
    ```

- [ ] **Regenerate the JSON-RPC module pages** (`docs/JSON-RPC-API/modules/*.md`),
      from the same checkout, never by hand. They carry the same version stamp as
      the `--help` dump and lag `main` for the same reason:

    ```shell
    $ make docs-generate
    ```

- [ ] **Test the example code on [Adding a network](add-network.md).** Its five Go
      files are not in the repository: add them under `params/` in the same
      checkout, run the page's test command from the repository root, then remove
      them:

    ```shell
    $ go test ./params/ -run 'TestGenesisHashABC|TestABCBootnodes|ExampleDefaultABCGenesisBlock' -count=1 -v
    ```

- [ ] **A second person has reviewed** the notes, the artifact list and the
      fingerprints.
- [ ] **For a full release, create the archive branch.** The migration guide tells
      readers to clone `archive-release-<version>`, and points them at the releases
      page when that clone fails. Creating the branch is what makes the documented
      build path work. Not applicable to a release candidate; the branch is cut with
      the full release.

## After publishing

- [ ] **Publish the documentation withheld for this release.** A page drafted for an
      unreleased version is held back two ways and both have to be undone: remove it
      from `exclude_docs` in `mkdocs.yml`, and add its entry back to `nav`. Neither
      alone is enough. Removing only the exclusion publishes a page nothing links to,
      and removing only the navigation entry publishes it anyway, because a page
      absent from `nav` is still built, still in the sitemap and still in the site
      search. Pushing that to `main` is what deploys it.
- [ ] Set `VersionMeta` to the next stage and push that commit.
- [ ] For a first-ever image push, confirm the GHCR package visibility. A registry
      creates new packages private by default, so nothing is pullable until it is
      made public, and making it public is the moment every tag it carries starts
      serving real traffic. Check what `:latest` points at before flipping it.
