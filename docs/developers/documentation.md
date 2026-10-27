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

## Logo

`docs/img/logo.svg` is the master for the site's logo, and the header shows it directly. PNG
renderings are made from it and never edited by hand. The favicon is rendered with
[Inkscape](https://inkscape.org/), from the project root:

```
inkscape docs/img/logo.svg --export-type=png --export-width=64 --export-filename=docs/img/favicon.png
```

