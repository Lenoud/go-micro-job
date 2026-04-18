#!/usr/bin/env bash

supports_color() {
  [[ -t 1 && -z "${NO_COLOR:-}" ]]
}

print_with_color() {
  local color="$1"
  local message="$2"

  if supports_color; then
    printf '\033[%sm%s\033[0m\n' "$color" "$message"
  else
    printf '%s\n' "$message"
  fi
}

print_info() {
  print_with_color "1;36" "$1"
}

print_success() {
  print_with_color "1;32" "$1"
}

print_warn() {
  print_with_color "1;33" "$1"
}

print_error() {
  print_with_color "1;31" "$1"
}
