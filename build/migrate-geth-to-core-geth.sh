#!/usr/bin/env bash
#
# migrate-geth-to-core-geth.sh
#
# Core-Geth's executable was renamed from `geth` to `core-geth` so it no longer
# collides with go-ethereum's `geth` on the same machine.
#
# THIS SCRIPT MOVES NO DATA, AND IT DOES NOT NEED TO.
#
# The datadir instance directory (`<datadir>/geth/`) and the IPC socket
# (`geth.ipc`) are set from constants in the client, not from the executable's
# filename, so they are unchanged by the rename. Your chaindata, nodekey and
# nodes database stay exactly where they are. There is no resync.
#
# What actually needs migrating is every place that *invokes* the binary:
# systemd units, Docker commands, shell scripts, cron entries, monitoring checks.
# That is what this script finds and, with --apply, rewrites.
#
# Usage:
#   ./migrate-geth-to-core-geth.sh                  # report only (default)
#   ./migrate-geth-to-core-geth.sh --apply          # rewrite, with .bak backups
#   ./migrate-geth-to-core-geth.sh --scan-dir DIR   # add a directory to scan
#
set -euo pipefail

MODE=report
SCAN_DIRS=()
FOUND_ANY=0
NEEDS_MANUAL=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --apply)    MODE=apply; shift ;;
    --scan-dir) SCAN_DIRS+=("$2"); shift 2 ;;
    -h|--help)  sed -n '2,28p' "$0"; exit 0 ;;
    *) echo "unknown argument: $1" >&2; exit 2 ;;
  esac
done

[[ ${#SCAN_DIRS[@]} -eq 0 ]] && SCAN_DIRS=("$HOME" /etc /opt /srv /usr/local/bin)

say()  { printf '%s\n' "$*"; }
warn() { printf 'WARN: %s\n' "$*" >&2; }

# A `geth` on disk may be go-ethereum's or Core-Geth's. Core-Geth prints "Core"
# as the first line of `version`; go-ethereum prints "Geth". Anything else is
# unknown and is never touched.
identify_binary() {
  local bin="$1" first
  [[ -x "$bin" ]] || { echo unknown; return; }
  if ! first=$("$bin" version 2>/dev/null | head -1); then echo unknown; return; fi
  case "$first" in
    Core*) echo core-geth ;;
    Geth*) echo go-ethereum ;;
    *)     echo unknown ;;
  esac
}

say "=== 1. Binaries named 'geth' on PATH ==="
IFS=':' read -r -a path_dirs <<< "$PATH"
for d in "${path_dirs[@]}"; do
  [[ -n "$d" && -f "$d/geth" ]] || continue
  kind=$(identify_binary "$d/geth")
  case "$kind" in
    core-geth)
      say "  $d/geth  ->  Core-Geth. Rename or replace with 'core-geth'."
      FOUND_ANY=1
      if [[ $MODE == apply ]]; then
        if mv -n "$d/geth" "$d/core-geth" 2>/dev/null; then
          say "    renamed to $d/core-geth"
        else
          warn "could not rename $d/geth (permissions, or core-geth already exists)"
        fi
      fi
      ;;
    go-ethereum)
      say "  $d/geth  ->  go-ethereum. LEAVE ALONE; this is the collision you are removing."
      ;;
    *)
      warn "  $d/geth  ->  could not identify; skipping"
      ;;
  esac
done

# Rewriting is POSITIONAL, not path-shaped. An absolute path ending in /geth is
# ambiguous -- /usr/local/bin/geth is the executable, /var/lib/etc/geth is very
# likely a datadir -- so a rule that rewrites any such path will happily strand a
# node's chaindata. Only three unambiguous positions are rewritten:
#
#   A. the first token after systemd ExecStart= / ExecStop= / ExecReload=
#   B. a path under a bin/ or sbin/ directory
#   C. a bare `geth` word (no slash), which cannot be a path
#
# Anything else that mentions geth is REPORTED for manual review, never rewritten.
SEP=$'[[:space:]"\']'

build_rewrite() {
  REWRITE=(
    -E
    -e "s#^([[:space:]]*Exec[A-Za-z]*=[-@]*)([^[:space:]]*/)?geth([[:space:]]|\$)#\\1\\2core-geth\\3#"
    -e "s#(^|$SEP|=)(/[^[:space:]]*/s?bin/)geth($SEP|\$)#\\1\\2core-geth\\3#g"
    -e "s#(^|$SEP|=)geth($SEP|\$)#\\1core-geth\\2#g"
  )
}
build_rewrite

# Lines mentioning geth that none of the three rules will touch.
UNSAFE_RE="(^|$SEP|=)/[^[:space:]]*/geth($SEP|\$)"

scan_file() {
  local f="$1" hits manual
  hits=$(grep -nE "(^|$SEP|=)(/[^[:space:]]*/)?geth($SEP|\$)" "$f" 2>/dev/null || true)
  [[ -n "$hits" ]] || return 0
  FOUND_ANY=1
  say "  $f"
  printf '%s\n' "$hits" | sed 's/^/      /' | head -6

  # Report ambiguous absolute paths that are NOT under bin/ or sbin/.
  manual=$(grep -nE "$UNSAFE_RE" "$f" 2>/dev/null \
           | grep -vE "/s?bin/geth($SEP|\$)" \
           | grep -vE "^[0-9]+:[[:space:]]*Exec[A-Za-z]*=" || true)
  if [[ -n "$manual" ]]; then
    NEEDS_MANUAL=1
    say "      NOT REWRITTEN -- verify by hand whether these are executables or data paths:"
    printf '%s\n' "$manual" | sed 's/^/        /' | head -4
  fi

  if [[ $MODE == apply ]]; then
    cp -p "$f" "$f.bak"
    sed "${REWRITE[@]}" -i "$f"
    say "      rewritten (backup: $f.bak)"
  fi
}

say ""
say "=== 2. systemd units invoking geth ==="
for unit_dir in /etc/systemd/system /lib/systemd/system /usr/lib/systemd/system "$HOME/.config/systemd/user"; do
  [[ -d "$unit_dir" ]] || continue
  while IFS= read -r -d '' f; do scan_file "$f"; done \
    < <(grep -rlZ --include='*.service' -E '(^|[[:space:]/=])geth([[:space:]]|$)' "$unit_dir" 2>/dev/null || true)
done

say ""
say "=== 3. scripts, compose files and cron entries invoking geth ==="
for root in "${SCAN_DIRS[@]}"; do
  [[ -d "$root" ]] || continue
  while IFS= read -r -d '' f; do scan_file "$f"; done \
    < <(grep -rlZ --binary-files=without-match \
          --include='*.sh' --include='*.bash' --include='*.service' --include='*.yml' --include='*.yaml' \
          --include='*.conf' --include='*.env' --include='crontab' \
          -E '(^|[[:space:]/=])geth([[:space:]]|$)' "$root" 2>/dev/null || true)
done

if crontab -l >/dev/null 2>&1; then
  if crontab -l 2>/dev/null | grep -qE '(^|[[:space:]/])geth([[:space:]]|$)'; then
    FOUND_ANY=1; NEEDS_MANUAL=1
    say ""
    say "=== 4. user crontab references geth -- EDIT BY HAND ==="
    crontab -l 2>/dev/null | grep -nE '(^|[[:space:]/])geth([[:space:]]|$)' | sed 's/^/      /'
    say "  This script never edits a crontab. Run 'crontab -e' and change the"
    say "  executable token only."
  fi
fi

say ""
say "=== Data check: nothing here should have changed ==="
say "  Your datadir layout is unaffected by the rename:"
say "    <datadir>/geth/chaindata   <- still 'geth', do NOT rename"
say "    <datadir>/geth/nodekey     <- still 'geth', do NOT rename"
say "    <datadir>/geth.ipc         <- still 'geth.ipc'"
say "  If you renamed any of those, move them back; the client looks for them"
say "  under 'geth' regardless of what the executable is called."

say ""
if [[ $FOUND_ANY -eq 0 ]]; then
  say "RESULT: no references found. Nothing to migrate."
elif [[ $MODE == report ]]; then
  say "RESULT: references found (report only). Re-run with --apply to rewrite."
  exit 1
else
  say "RESULT: rewrite complete. Backups are alongside each file as *.bak."
  say "  systemctl daemon-reload   # if any unit changed"
  [[ $NEEDS_MANUAL -eq 1 ]] && say "  One or more items still need manual editing (see above)."
fi
