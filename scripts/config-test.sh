#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MICRO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# shellcheck source=/dev/null
source "$SCRIPT_DIR/colors.sh"

FAILED=0

require_text() {
  local file="$1"
  local text="$2"

  if ! grep -Fq "$text" "$file"; then
    print_error "[config] missing '$text' in ${file#$MICRO_DIR/}"
    FAILED=1
  fi
}

check_rest_config() {
  local file="$1"

  require_text "$file" "Mode:"
  require_text "$file" "Timeout:"
  require_text "$file" "MaxConns:"
  require_text "$file" "MaxBytes:"
  require_text "$file" "CpuThreshold:"
  require_text "$file" "Log:"
  require_text "$file" "DevServer:"
  require_text "$file" "Middlewares:"
  require_text "$file" "Breaker: true"
  require_text "$file" "Timeout: true"
  require_text "$file" "Recover: true"
  require_text "$file" "Prometheus: true"
  require_text "$file" "UserRpc:"
  require_text "$file" "DepartmentRpc:"
  require_text "$file" "NonBlock: true"
  require_text "$file" "KeepaliveTime:"
}

check_rpc_config() {
  local file="$1"

  require_text "$file" "Mode:"
  require_text "$file" "Timeout:"
  require_text "$file" "CpuThreshold:"
  require_text "$file" "Health: true"
  require_text "$file" "Log:"
  require_text "$file" "DevServer:"
  require_text "$file" "Middlewares:"
  require_text "$file" "Recover: true"
  require_text "$file" "Stat: true"
  require_text "$file" "Prometheus: true"
  require_text "$file" "Breaker: true"
}

check_rest_config "$MICRO_DIR/app/api-gateway/etc/apigateway-local.yaml"
check_rest_config "$MICRO_DIR/app/api-gateway/etc/apigateway.yaml"
check_rpc_config "$MICRO_DIR/app/user-service/etc/user-local.yaml"
check_rpc_config "$MICRO_DIR/app/user-service/etc/user.yaml"
check_rpc_config "$MICRO_DIR/app/department-service/etc/department-local.yaml"
check_rpc_config "$MICRO_DIR/app/department-service/etc/department.yaml"

if [[ "$FAILED" -eq 1 ]]; then
  print_error "[config] explicit config check failed"
  exit 1
fi

print_success "[config] explicit config check passed"
