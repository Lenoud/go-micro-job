#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# shellcheck source=/dev/null
source "$SCRIPT_DIR/dev-common.sh"

assert_eq() {
  local want="$1"
  local got="$2"
  local label="$3"

  if [[ "$want" != "$got" ]]; then
    printf 'FAIL %s: want %s, got %s\n' "$label" "$want" "$got" >&2
    exit 1
  fi
}

test_wait_for_port_succeeds_when_port_becomes_ready() {
  local checks=0
  local sleeps=0

  port_listening() {
    local _port="$1"
    checks=$((checks + 1))
    if [[ "$checks" -ge 3 ]]; then
      return 0
    fi
    return 1
  }

  sleep_cmd() {
    sleeps=$((sleeps + 1))
  }

  wait_for_port 9101 "user-service" 5

  assert_eq "3" "$checks" "checks before ready"
  assert_eq "2" "$sleeps" "sleeps before ready"
}

test_wait_for_port_fails_after_timeout() {
  local checks=0
  local sleeps=0

  port_listening() {
    local _port="$1"
    checks=$((checks + 1))
    return 1
  }

  sleep_cmd() {
    sleeps=$((sleeps + 1))
  }

  if wait_for_port 9102 "department-service" 2; then
    printf 'FAIL timeout: expected wait_for_port to fail\n' >&2
    exit 1
  fi

  assert_eq "2" "$checks" "checks before timeout"
  assert_eq "2" "$sleeps" "sleeps before timeout"
}

test_wait_for_port_succeeds_when_port_becomes_ready
test_wait_for_port_fails_after_timeout

printf 'dev-common tests passed\n'
