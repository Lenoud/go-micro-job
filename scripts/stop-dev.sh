#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MICRO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
RUN_DIR="$MICRO_DIR/run"

# shellcheck source=/dev/null
source "$SCRIPT_DIR/dev-common.sh"

USER_SVC_PID_FILE="$RUN_DIR/user-service.pid"
GW_PID_FILE="$RUN_DIR/api-gateway.pid"

stop_pid_file "$GW_PID_FILE" "api-gateway"
stop_pid_file "$USER_SVC_PID_FILE" "user-service"

print_success "[stop-dev] done."
