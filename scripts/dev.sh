#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MICRO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# shellcheck source=/dev/null
source "$SCRIPT_DIR/dev-common.sh"

LOG_DIR="$MICRO_DIR/logs"
RUN_DIR="$MICRO_DIR/run"
USER_SVC_LOG="$LOG_DIR/user-service.log"
GW_LOG="$LOG_DIR/api-gateway.log"
USER_SVC_PID_FILE="$RUN_DIR/user-service.pid"
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
  stop_pid_file "$USER_SVC_PID_FILE" "user-service"
  print_success "[dev] cleanup complete."

  exit "$exit_code"
}

trap 'cleanup $?' EXIT
trap 'exit 130' INT TERM

mkdir -p "$LOG_DIR" "$RUN_DIR"
stop_pid_file "$GW_PID_FILE" "api-gateway"
stop_pid_file "$USER_SVC_PID_FILE" "user-service"
stop_port_listener 9101 "user-service"
stop_port_listener 9000 "api-gateway"

: > "$USER_SVC_LOG"
: > "$GW_LOG"

# 检查并启动基础设施
check_infra

# 确保 micro_job 数据库和表存在（复用根目录 mysql 3306）
print_info "[dev] ensuring database micro_job exists..."
docker exec go_job_mysql mysql -uroot -proot123 --default-character-set=utf8mb4 \
  -e "CREATE DATABASE IF NOT EXISTS micro_job CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null \
  || print_warn "[dev] could not create database, may already exist or mysql not accessible via docker"

print_info "[dev] importing schema into micro_job..."
USER_SQL="$MICRO_DIR/sql/user.sql"
if [[ -f "$USER_SQL" ]]; then
  docker exec -i go_job_mysql mysql -uroot -proot123 --default-character-set=utf8mb4 < "$USER_SQL" 2>/dev/null \
    || print_warn "[dev] schema import failed or already imported"
fi

# 构建
print_info "[dev] building user-service..."
(
  cd "$MICRO_DIR/user-service"
  go build ./...
)

print_info "[dev] building api-gateway..."
(
  cd "$MICRO_DIR/api-gateway"
  go build ./...
)

# 启动 user-service (gRPC)
print_info "[dev] starting user-service (gRPC :9101)..."
(
  cd "$MICRO_DIR/user-service"
  exec go run user.go -f etc/user-local.yaml
) >> "$USER_SVC_LOG" 2>&1 &
USER_SVC_PID=$!
echo "$USER_SVC_PID" > "$USER_SVC_PID_FILE"

# 等待 user-service 注册到 etcd
print_info "[dev] waiting for user-service to register on etcd..."
sleep_cmd 3

# 启动 api-gateway (REST)
print_info "[dev] starting api-gateway (REST :9000)..."
(
  cd "$MICRO_DIR/api-gateway"
  exec go run apigateway.go -f etc/apigateway-local.yaml
) >> "$GW_LOG" 2>&1 &
GW_PID=$!
echo "$GW_PID" > "$GW_PID_FILE"

print_success "[dev] user-service log: $USER_SVC_LOG"
print_success "[dev] api-gateway log:  $GW_LOG"
print_warn "[dev] Ctrl+C will stop both processes."
printf '\n'

EXITED_STATUS=0
while true; do
  if ! is_running "$USER_SVC_PID"; then
    wait "$USER_SVC_PID" 2>/dev/null || EXITED_STATUS=$?
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

if ! is_running "$GW_PID"; then
  print_error "[dev] api-gateway exited. Check $GW_LOG"
fi

exit "$EXITED_STATUS"
