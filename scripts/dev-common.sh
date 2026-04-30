#!/usr/bin/env bash

# shellcheck source=/dev/null
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/colors.sh"

kill_cmd() {
  kill "$@"
}

sleep_cmd() {
  sleep "$@"
}

pgrep_cmd() {
  pgrep "$@"
}

lsof_cmd() {
  lsof "$@"
}

is_running() {
  local pid="$1"
  kill_cmd -0 "$pid" >/dev/null 2>&1
}

kill_tree() {
  local pid="$1"
  local child

  if [[ -z "${pid:-}" ]] || ! [[ "$pid" =~ ^[0-9]+$ ]]; then
    return
  fi

  while IFS= read -r child; do
    [[ -n "$child" ]] && kill_tree "$child"
  done < <(pgrep_cmd -P "$pid" 2>/dev/null || true)

  kill_cmd -TERM "$pid" >/dev/null 2>&1 || true
}

force_kill_tree() {
  local pid="$1"
  local child

  if [[ -z "${pid:-}" ]] || ! [[ "$pid" =~ ^[0-9]+$ ]]; then
    return
  fi

  while IFS= read -r child; do
    [[ -n "$child" ]] && force_kill_tree "$child"
  done < <(pgrep_cmd -P "$pid" 2>/dev/null || true)

  kill_cmd -KILL "$pid" >/dev/null 2>&1 || true
}

stop_pid() {
  local pid="$1"
  local name="${2:-process}"

  if [[ -n "$pid" ]] && is_running "$pid"; then
    print_warn "[dev] stopping existing $name (pid: $pid)..."
    kill_tree "$pid"
    sleep_cmd 1
    is_running "$pid" && force_kill_tree "$pid"
    wait "$pid" 2>/dev/null || true
  fi
}

stop_pid_file() {
  local pid_file="$1"
  local name="${2:-process}"
  local pid

  [[ -f "$pid_file" ]] || return 0

  pid="$(cat "$pid_file" 2>/dev/null || true)"
  stop_pid "$pid" "$name"

  rm -f "$pid_file"
}

find_port_pids() {
  local port="$1"

  lsof_cmd -tiTCP:"$port" -sTCP:LISTEN 2>/dev/null || true
}

stop_port_listener() {
  local port="$1"
  local name="${2:-port listener}"
  local pid

  while IFS= read -r pid; do
    [[ -n "$pid" ]] || continue
    stop_pid "$pid" "$name on port $port"
  done < <(find_port_pids "$port")
}

port_listening() {
  local port="$1"
  lsof_cmd -tiTCP:"$port" -sTCP:LISTEN 2>/dev/null | grep -q .
}

check_infra() {
	local missing=0

	if ! port_listening 3306; then
		print_error "[dev] mysql (3306) is not running"
		missing=1
	fi

  if ! port_listening 6379; then
    print_error "[dev] redis (6379) is not running"
    missing=1
  fi

  if ! port_listening 2379; then
    print_error "[dev] etcd (2379) is not running"
    missing=1
	fi

	if [[ "$missing" -eq 1 ]]; then
		print_error "[dev] start local mysql, redis, and etcd first; dev.sh no longer starts Docker infrastructure automatically"
		exit 1
	else
		print_success "[dev] infrastructure ready (mysql:3306, redis:6379, etcd:2379)"
	fi
}
