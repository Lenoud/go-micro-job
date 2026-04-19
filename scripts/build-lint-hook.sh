#!/usr/bin/env bash

# Pre-commit hook: build + lint all micro services.
# Usage:
#   ./scripts/build-lint-hook.sh          # build + lint
#   ./scripts/build-lint-hook.sh --build  # build only
#   ./scripts/build-lint-hook.sh --lint   # lint only
#
# Install as git pre-commit hook:
#   ln -sf ../../scripts/build-lint-hook.sh .git/hooks/pre-commit

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MICRO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# shellcheck source=/dev/null
source "$SCRIPT_DIR/colors.sh"

SERVICES=("user-service" "api-gateway")
FAILED=0

run_build() {
  for svc in "${SERVICES[@]}"; do
    print_info "[build] $svc ..."
    if (cd "$MICRO_DIR/app/$svc" && go build ./...); then
      print_success "[build] $svc ok"
    else
      print_error "[build] $svc FAILED"
      FAILED=1
    fi
  done
}

run_lint() {
  if ! command -v golangci-lint &>/dev/null; then
    print_warn "[lint] golangci-lint not found, skipping"
    return
  fi

  for svc in "${SERVICES[@]}"; do
    print_info "[lint] $svc ..."
    if (cd "$MICRO_DIR/app/$svc" && golangci-lint run ./...); then
      print_success "[lint] $svc ok"
    else
      print_error "[lint] $svc FAILED"
      FAILED=1
    fi
  done
}

case "${1:-all}" in
  --build) run_build ;;
  --lint)  run_lint ;;
  *)
    run_build
    run_lint
    ;;
esac

if [[ "$FAILED" -eq 1 ]]; then
  print_error "[hook] some checks failed"
  exit 1
fi

print_success "[hook] all checks passed"
