---
description: "How the Core-Geth documentation is written and built: the Markdown sources under /docs and the MkDocs site that publishes them."
---

# Documentation

Project documentation lives in `/docs` and is written in [Markdown](https://daringfireball.net/projects/markdown/syntax).  

For web-based access, these files are passed through a static site generator [MkDocs](https://www.mkdocs.org/),
specifically [MkDocs Material](https://squidfunk.github.io/mkdocs-material/) and served via Github Pages.

## Development

You can run a live-reloading web server in imitation of the production generator. To do so:

- Ensure that your python environment is using python3 and its package manager pip3. 
  You can then install the required `mkdocs` executable and its dependencies using:
  ```
  python -m pip install -r requirements-mkdocs.txt
  ```
- Run `mkdocs serve` from the project root. A convenience Make command is likewise provided as `make mkdocs-serve`.
- Open `http://localhost:8000` in a web browser.
- Write some docs!

## Writing

Keep pages clean, crisp, clear and accurate. A step says what to do and what to expect. Leave out what a reader does not
act on, and say a thing once, on one page, then link to it.

## Code blocks


A reader copies a block with its copy button and pastes it into a shell, so write each block to run
as it stands. Leave out the `$` or `>` prompt. Put expected output in a sentence or a block of its
own, not in a `#` comment: stock zsh, the macOS default, does not treat a pasted `#` as a comment,
and the line fails. Where one line depends on the one before it, such as a checksum and the install
that follows it, chain them with `&&` so a pasted block stops at the failure.

A page or tab for one platform names that platform's files, with a note under the block for its variants, so a reader
on that platform can run it as it stands. Commands look up the latest release instead of naming one, because a fixed
version goes stale at the next release.

## Logo

`docs/img/logo.svg` is the master for the site's logo, and the header shows it directly. PNG
renderings are made from it and never edited by hand. The favicon is rendered with
[Inkscape](https://inkscape.org/), from the project root:

```
inkscape docs/img/logo.svg --export-type=png --export-width=64 --export-filename=docs/img/favicon.png
```

