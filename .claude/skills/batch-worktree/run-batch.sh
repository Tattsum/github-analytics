#!/usr/bin/env bash
# Run the github-analytics batch (-mode batch) inside a throwaway git worktree.
#
# Why a worktree: the batch is built and run from a *committed* snapshot of the
# code in an isolated checkout, so it is unaffected by in-progress edits in the
# main checkout or by other worktrees you have open. The Postgres DB is shared
# (the batch writes exactly one snapshot to it) — only the code is isolated.
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: run-batch.sh [--ref <git-ref>] [--keep] [-- <batch args...>]

  --ref <git-ref>  Commit/branch/tag to run the batch from (default: HEAD).
  --keep           Keep the worktree after the run (for debugging). Default: remove.
  -- <args>        Args passed straight through to the batch CLI. Examples:
                     -- -users user1,user2
                     -- -org myorg -team my-team -private

Runtime env (resolved automatically from the repo root if present):
  GITHUB_TOKEN   required; sourced from .envrc / .env when not already exported.
  DATABASE_URL   defaults to the local docker-compose Postgres.
EOF
}

REF="HEAD"
KEEP=0
BATCH_ARGS=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --ref) REF="${2:?--ref needs a value}"; shift 2 ;;
    --keep) KEEP=1; shift ;;
    -h | --help) usage; exit 0 ;;
    --) shift; BATCH_ARGS=("$@"); break ;;
    *) echo "unknown arg: $1" >&2; usage; exit 1 ;;
  esac
done

ORIGIN_ROOT="$(git rev-parse --show-toplevel)"
COMMIT="$(git -C "$ORIGIN_ROOT" rev-parse --short "$REF")"

# mktemp -d creates the dir, but `git worktree add` requires a non-existent path.
WT_DIR="$(mktemp -d "${TMPDIR:-/tmp}/gha-batch-XXXXXX")"
rmdir "$WT_DIR"

cleanup() {
  if [[ "$KEEP" -eq 1 ]]; then
    echo "worktree kept at: $WT_DIR"
    return
  fi
  git -C "$ORIGIN_ROOT" worktree remove --force "$WT_DIR" 2>/dev/null || true
  git -C "$ORIGIN_ROOT" worktree prune 2>/dev/null || true
}
trap cleanup EXIT

echo "Creating worktree at $WT_DIR (ref: $REF -> $COMMIT)"
git -C "$ORIGIN_ROOT" worktree add --detach "$WT_DIR" "$COMMIT" >/dev/null

# Worktrees do not carry gitignored files, so load secrets/config from the
# origin repo root. .env is KEY=VALUE; .envrc just exports GITHUB_TOKEN.
if [[ -f "$ORIGIN_ROOT/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$ORIGIN_ROOT/.env"
  set +a
fi
if [[ -f "$ORIGIN_ROOT/.envrc" ]]; then
  # shellcheck disable=SC1091
  source "$ORIGIN_ROOT/.envrc"
fi

if [[ -z "${GITHUB_TOKEN:-}" ]]; then
  echo "ERROR: GITHUB_TOKEN is not set (not found in env, .envrc, or .env)." >&2
  exit 1
fi

# Postgres is shared across worktrees by design; bring it up idempotently.
(cd "$ORIGIN_ROOT" && docker compose up -d --wait postgres >/dev/null)

DEFAULT_DB="postgres://github_analytics:github_analytics@localhost:5432/github_analytics?sslmode=disable"
echo "Running batch from worktree (commit $COMMIT)..."
(
  cd "$WT_DIR"
  DATABASE_URL="${DATABASE_URL:-$DEFAULT_DB}" \
    go run ./cmd/github-analytics -mode batch "${BATCH_ARGS[@]}"
)

echo "Batch finished (worktree will be cleaned up)."
