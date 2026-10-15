"""MkDocs build hook: stamp the client version and build commit into the site footer.

docs.coregeth.com is rebuilt from `main` on every merge, so the values here are read at
build time: the version is `params.VersionWithMeta` from `params/version.go`, and the commit
is the deployed `HEAD`. Together they say which client state the live docs describe -- the
development version heading to the next release, not the last tag.

Wired through `hooks:` in `mkdocs.yml`. The stamp is appended to `config.copyright`, which
Material renders as the site footer, so it appears once site-wide rather than per page (the
per-page "last updated" comes from the git-revision-date-localized plugin).
"""

import json
import os
import re
import subprocess
from pathlib import Path

_REPO = Path(__file__).resolve().parent
_VERSION_GO = _REPO / "params" / "version.go"
_COMMIT_URL = "https://github.com/ethereumclassic/core-geth/commit/{sha}"


def _client_version():
    """Reproduce params.VersionWithMeta by reading params/version.go.

    VersionWithMeta is "<major>.<minor>.<patch>" with "-<meta>" appended whenever meta is
    non-empty -- which it always is (the invariant in version.go forbids an empty meta), so a
    stable release reads "1.13.1-stable" and the current dev tree reads "1.13.1-unstable".
    """
    try:
        src = _VERSION_GO.read_text()
    except OSError:
        return None

    def const(name):
        m = re.search(rf"\b{name}\b\s*=\s*\"?([0-9A-Za-z]+)\"?", src)
        return m.group(1) if m else None

    major, minor, patch = const("VersionMajor"), const("VersionMinor"), const("VersionPatch")
    meta = const("VersionMeta")
    if None in (major, minor, patch):
        return None
    version = f"{major}.{minor}.{patch}"
    if meta:
        version += f"-{meta}"
    return version


def _commit():
    """The deployed commit: HEAD locally, or $GITHUB_SHA in CI. Returns (full, short) or (None, None)."""
    sha = os.environ.get("GITHUB_SHA")
    if not sha:
        try:
            sha = subprocess.check_output(
                ["git", "rev-parse", "HEAD"], cwd=_REPO, text=True, stderr=subprocess.DEVNULL
            ).strip()
        except (OSError, subprocess.CalledProcessError):
            return None, None
    if not re.fullmatch(r"[0-9a-f]{40}", sha or ""):
        return None, None
    return sha, sha[:8]


def on_config(config, **kwargs):
    version = _client_version()
    sha, short = _commit()

    parts = []
    if version:
        parts.append(f"Core-Geth {version}")
    if short:
        parts.append(f'built from <a href="{_COMMIT_URL.format(sha=sha)}">{short}</a>')
    if not parts:
        return config

    stamp = "These docs are current for " + ", ".join(parts) + "."

    site_url = (config.get("site_url") or "").rstrip("/")
    releases_url = f"{site_url}/releases/" if site_url else "releases/"
    stamp += (
        " Static, version-specific docs for each release are listed on the "
        f'<a href="{releases_url}">Releases page</a>.'
    )

    existing = config.get("copyright") or ""
    config["copyright"] = f"{existing}<br>{stamp}" if existing else stamp
    return config


# --- schema.org FAQPage structured data for the FAQ -------------------------
# Emits JSON-LD built from the page's "## question" sections, so search engines
# read the Q&A as structured data. Scoped to the FAQ page. Rich-result display
# is restricted by Google to some site classes, but the semantic signal applies
# regardless.

FAQ_SRC = "about/faq.md"


def _md_to_text(md):
    """Minimal Markdown -> plain text for JSON-LD question/answer strings."""
    md = re.sub(r"\[([^\]]+)\]\([^)]+\)", r"\1", md)  # [text](url) -> text
    md = re.sub(r"`([^`]+)`", r"\1", md)              # `code` -> code
    md = re.sub(r"\*\*([^*]+)\*\*", r"\1", md)        # **bold** -> bold
    md = re.sub(r"\*([^*]+)\*", r"\1", md)            # *italic* -> italic
    md = re.sub(r"\s+", " ", md)
    return md.strip()


def _faq_jsonld(markdown):
    """Build a schema.org FAQPage object from '## <question>' sections."""
    sections = re.split(r"^## +", markdown, flags=re.MULTILINE)[1:]  # drop the intro
    qa = []
    for sec in sections:
        head, _, body = sec.partition("\n")
        question, answer = _md_to_text(head), _md_to_text(body)
        if question and answer:
            qa.append({
                "@type": "Question",
                "name": question,
                "acceptedAnswer": {"@type": "Answer", "text": answer},
            })
    if not qa:
        return None
    return json.dumps(
        {"@context": "https://schema.org", "@type": "FAQPage", "mainEntity": qa},
        ensure_ascii=False,
    )


def on_page_markdown(markdown, page, config, files, **kwargs):
    if page.file.src_path.replace("\\", "/") == FAQ_SRC:
        data = _faq_jsonld(markdown)
        if data:
            page.meta["faq_jsonld"] = data
    return markdown


def on_post_page(output, page, config, **kwargs):
    data = page.meta.get("faq_jsonld")
    if data and "</head>" in output:
        # Escape '<' so the JSON can never close the <script> element early.
        data = data.replace("<", "\\u003c")
        script = f'<script type="application/ld+json">{data}</script>\n'
        output = output.replace("</head>", script + "</head>", 1)
    return output
