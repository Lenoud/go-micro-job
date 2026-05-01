#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MICRO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# shellcheck source=/dev/null
source "$SCRIPT_DIR/dev-common.sh"

LOG_DIR="$MICRO_DIR/logs"
RUN_DIR="$MICRO_DIR/run"
export GOCACHE="${GOCACHE:-/tmp/go-server-resume-micro-gocache}"
USER_SVC_LOG="$LOG_DIR/user-service.log"
DEPARTMENT_SVC_LOG="$LOG_DIR/department-service.log"
OPLOG_SVC_LOG="$LOG_DIR/oplog-service.log"
GW_LOG="$LOG_DIR/api-gateway.log"
USER_SVC_PID_FILE="$RUN_DIR/user-service.pid"
DEPARTMENT_SVC_PID_FILE="$RUN_DIR/department-service.pid"
OPLOG_SVC_PID_FILE="$RUN_DIR/oplog-service.pid"
GW_PID_FILE="$RUN_DIR/api-gateway.pid"

CLEANED_UP=0

cleanup() {
  local exit_code="${1:-0}"

  if [[ "$CLEANED_UP" -eq 1 ]]; then
    return
  fi
  CLEANED_UP=1

  trap - EXIT INT TERM

  printf '\n'
  print_warn "[dev] stopping processes..."
  stop_pid_file "$GW_PID_FILE" "api-gateway"
  stop_pid_file "$OPLOG_SVC_PID_FILE" "oplog-service"
  stop_pid_file "$DEPARTMENT_SVC_PID_FILE" "department-service"
  stop_pid_file "$USER_SVC_PID_FILE" "user-service"
  print_success "[dev] cleanup complete."

  exit "$exit_code"
}

trap 'cleanup $?' EXIT
trap 'exit 130' INT TERM

mkdir -p "$LOG_DIR" "$RUN_DIR"
stop_pid_file "$GW_PID_FILE" "api-gateway"
stop_pid_file "$OPLOG_SVC_PID_FILE" "oplog-service"
stop_pid_file "$DEPARTMENT_SVC_PID_FILE" "department-service"
stop_pid_file "$USER_SVC_PID_FILE" "user-service"
stop_port_listener 9103 "oplog-service"
stop_port_listener 9102 "department-service"
stop_port_listener 9101 "user-service"
stop_port_listener 9000 "api-gateway"

: > "$USER_SVC_LOG"
: > "$DEPARTMENT_SVC_LOG"
: > "$OPLOG_SVC_LOG"
: > "$GW_LOG"

# 检查并启动基础设施
check_infra

# 确保 micro_job 数据库和表存在（本地 mysql 3306）
print_info "[dev] ensuring database micro_job exists..."
mysql --protocol=TCP -h127.0.0.1 -P3306 -uroot -proot123 --default-character-set=utf8mb4 \
  -e "CREATE DATABASE IF NOT EXISTS micro_job CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null \
  || print_warn "[dev] could not create database, may already exist or mysql not accessible"

print_info "[dev] importing schema into micro_job..."
USER_SQL="$MICRO_DIR/sql/user.sql"
if [[ -f "$USER_SQL" ]]; then
  mysql --protocol=TCP -h127.0.0.1 -P3306 -uroot -proot123 --default-character-set=utf8mb4 < "$USER_SQL" 2>/dev/null \
    || print_warn "[dev] schema import failed or already imported"
fi
DEPARTMENT_SQL="$MICRO_DIR/sql/department.sql"
if [[ -f "$DEPARTMENT_SQL" ]]; then
  mysql --protocol=TCP -h127.0.0.1 -P3306 -uroot -proot123 --default-character-set=utf8mb4 < "$DEPARTMENT_SQL" 2>/dev/null \
    || print_warn "[dev] department schema import failed or already imported"
fi
OPLOG_SQL="$MICRO_DIR/sql/oplog.sql"
if [[ -f "$OPLOG_SQL" ]]; then
  mysql --protocol=TCP -h127.0.0.1 -P3306 -uroot -proot123 --default-character-set=utf8mb4 < "$OPLOG_SQL" 2>/dev/null \
    || print_warn "[dev] oplog schema import failed or already imported"
fi

# 并行构建
BUILD_SERVICES=("user-service" "department-service" "oplog-service" "api-gateway")
BUILD_PIDS=()
BUILD_FAILED=0

for svc in "${BUILD_SERVICES[@]}"; do
  print_info "[dev] building $svc..."
  (
    cd "$MICRO_DIR/app/$svc"
    go build ./...
  ) &
  BUILD_PIDS+=("$!")
done

for i in "${!BUILD_SERVICES[@]}"; do
  svc="${BUILD_SERVICES[$i]}"
  pid="${BUILD_PIDS[$i]}"
  if wait "$pid"; then
    print_success "[dev] build $svc ok"
  else
    print_error "[dev] build $svc failed"
    BUILD_FAILED=1
  fi
done

if [[ "$BUILD_FAILED" -eq 1 ]]; then
  exit 1
fi

# 启动 user-service (gRPC)
print_info "[dev] starting user-service (gRPC :9101)..."
(
  cd "$MICRO_DIR/app/user-service"
  exec go run user.go -f etc/user-local.yaml
) >> "$USER_SVC_LOG" 2>&1 &
USER_SVC_PID=$!
echo "$USER_SVC_PID" > "$USER_SVC_PID_FILE"

# 启动 department-service (gRPC)
print_info "[dev] starting department-service (gRPC :9102)..."
(
  cd "$MICRO_DIR/app/department-service"
  exec go run department.go -f etc/department-local.yaml
) >> "$DEPARTMENT_SVC_LOG" 2>&1 &
DEPARTMENT_SVC_PID=$!
echo "$DEPARTMENT_SVC_PID" > "$DEPARTMENT_SVC_PID_FILE"

# 启动 oplog-service (gRPC)
print_info "[dev] starting oplog-service (gRPC :9103)..."
(
  cd "$MICRO_DIR/app/oplog-service"
  exec go run oplog.go -f etc/oplog.yaml
) >> "$OPLOG_SVC_LOG" 2>&1 &
OPLOG_SVC_PID=$!
echo "$OPLOG_SVC_PID" > "$OPLOG_SVC_PID_FILE"

# 等待 RPC 服务就绪，避免 gateway 启动时找不到依赖
wait_for_port 9101 "user-service" 30
wait_for_port 9102 "department-service" 30
wait_for_port 9103 "oplog-service" 30

# 启动 api-gateway (REST)
print_info "[dev] starting api-gateway (REST :9000)..."
(
  cd "$MICRO_DIR/app/api-gateway"
  exec go run apigateway.go -f etc/apigateway-local.yaml
) >> "$GW_LOG" 2>&1 &
GW_PID=$!
echo "$GW_PID" > "$GW_PID_FILE"

print_success "[dev] user-service log:      $USER_SVC_LOG"
print_success "[dev] department-service log: $DEPARTMENT_SVC_LOG"
print_success "[dev] oplog-service log:      $OPLOG_SVC_LOG"
print_success "[dev] api-gateway log:        $GW_LOG"
print_warn "[dev] Ctrl+C will stop all processes."
printf '\n'

EXITED_STATUS=0
while true; do
  if ! is_running "$USER_SVC_PID"; then
    wait "$USER_SVC_PID" 2>/dev/null || EXITED_STATUS=$?
    break
  fi

  if ! is_running "$DEPARTMENT_SVC_PID"; then
    wait "$DEPARTMENT_SVC_PID" 2>/dev/null || EXITED_STATUS=$?
    break
  fi

  if ! is_running "$OPLOG_SVC_PID"; then
    wait "$OPLOG_SVC_PID" 2>/dev/null || EXITED_STATUS=$?
    break
  fi

  if ! is_running "$GW_PID"; then
    wait "$GW_PID" 2>/dev/null || EXITED_STATUS=$?
    break
  fi

  sleep 1
done

if ! is_running "$USER_SVC_PID"; then
  print_error "[dev] user-service exited. Check $USER_SVC_LOG"
fi

if ! is_running "$DEPARTMENT_SVC_PID"; then
  print_error "[dev] department-service exited. Check $DEPARTMENT_SVC_LOG"
fi

if ! is_running "$OPLOG_SVC_PID"; then
  print_error "[dev] oplog-service exited. Check $OPLOG_SVC_LOG"
fi

if ! is_running "$GW_PID"; then
  print_error "[dev] api-gateway exited. Check $GW_LOG"
fi

exit "$EXITED_STATUS"
