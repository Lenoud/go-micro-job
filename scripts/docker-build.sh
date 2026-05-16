#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MICRO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# shellcheck source=/dev/null
source "$SCRIPT_DIR/colors.sh"

SERVICES=("user-service" "department-service" "oplog-service" "api-gateway")
REGISTRY="${REGISTRY:-}"
TAG="${TAG:-latest}"

usage() {
  cat <<EOF
Usage: $(basename "$0") [options] [service...]

Build Docker images for micro services.

Services:
  $(IFS=$'\n'; echo "${SERVICES[*]}")

Options:
  --registry <url>   Docker registry prefix (default: none, env: REGISTRY)
  --tag <tag>        Image tag (default: latest, env: TAG)
  --no-cache         Pass --no-cache to docker build
  --push             Push images after building (requires --registry)
  -h, --help         Show this help

Examples:
  $(basename "$0")                              # Build all services
  $(basename "$0") user-service api-gateway     # Build specific services
  $(basename "$0") --registry ghcr.io/org --tag v1.0.0
  REGISTRY=localhost:5000 $(basename "$0") --push
EOF
}

NO_CACHE=""
PUSH=0
TARGET_SERVICES=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --registry)
      REGISTRY="$2"
      shift 2
      ;;
    --tag)
      TAG="$2"
      shift 2
      ;;
    --no-cache)
      NO_CACHE="--no-cache"
      shift
      ;;
    --push)
      PUSH=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    -*)
      print_error "Unknown option: $1"
      usage
      exit 1
      ;;
    *)
      TARGET_SERVICES+=("$1")
      shift
      ;;
  esac
done

# Default: build all services
if [[ ${#TARGET_SERVICES[@]} -eq 0 ]]; then
  TARGET_SERVICES=("${SERVICES[@]}")
fi

# Validate service names
for svc in "${TARGET_SERVICES[@]}"; do
  found=0
  for valid in "${SERVICES[@]}"; do
    if [[ "$svc" == "$valid" ]]; then
      found=1
      break
    fi
  done
  if [[ "$found" -eq 0 ]]; then
    print_error "Unknown service: $svc"
    print_info "Valid services: ${SERVICES[*]}"
    exit 1
  fi
done

cd "$MICRO_DIR"

print_info "[docker-build] Building ${#TARGET_SERVICES[@]} service(s) with tag=$TAG"

BUILD_FAILED=0
for svc in "${TARGET_SERVICES[@]}"; do
  if [[ -n "$REGISTRY" ]]; then
    IMAGE="${REGISTRY}/${svc}:${TAG}"
  else
    IMAGE="${svc}:${TAG}"
  fi

  print_info "[docker-build] Building $svc → $IMAGE ..."
  if docker build $NO_CACHE \
    --pull=false \
    -f "app/$svc/Dockerfile" \
    -t "$IMAGE" \
    .; then
    print_success "[docker-build] $svc → $IMAGE ✓"
  else
    print_error "[docker-build] $svc FAILED"
    BUILD_FAILED=1
  fi
done

if [[ "$BUILD_FAILED" -eq 1 ]]; then
  print_error "[docker-build] Some builds failed."
  exit 1
fi

# Push
if [[ "$PUSH" -eq 1 ]]; then
  if [[ -z "$REGISTRY" ]]; then
    print_error "[docker-build] --push requires --registry"
    exit 1
  fi
  for svc in "${TARGET_SERVICES[@]}"; do
    IMAGE="${REGISTRY}/${svc}:${TAG}"
    print_info "[docker-build] Pushing $IMAGE ..."
    docker push "$IMAGE"
    print_success "[docker-build] Pushed $IMAGE ✓"
  done
fi

print_success "[docker-build] All done."
