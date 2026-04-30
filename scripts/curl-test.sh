#!/usr/bin/env bash
# curl-test.sh — 测试 api-gateway (微服务模式) 的用户接口
set -euo pipefail

BASE="http://localhost:9000"
PASS="lb781023"

red()   { printf '\033[31m%s\033[0m\n' "$1"; }
green() { printf '\033[32m%s\033[0m\n' "$1"; }
bold()  { printf '\033[1m%s\033[0m\n' "$1"; }

fetch_token() {
  local endpoint="$1" username="$2"
  local raw
  raw=$(curl -s -X POST "${BASE}${endpoint}" \
    -H 'Content-Type: application/json' \
    -d "{\"username\":\"${username}\",\"password\":\"${PASS}\"}")
  local token
  token=$(echo "$raw" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('token',''))" 2>/dev/null || true)
  if [ -z "$token" ]; then
    red "  FAILED: $username @ $endpoint"
    echo "  Response: $raw"
    return 1
  fi
  echo "$token"
}

check() {
  local label="$1" expect="$2" actual
  actual="$3"
  if echo "$actual" | python3 -c "import sys,json; d=json.load(sys.stdin); sys.exit(0 if d.get('code')==$expect else 1)" 2>/dev/null; then
    green "  PASS: $label"
  else
    red "  FAIL: $label (expected code=$expect)"
  fi
}

bold "=== Login via api-gateway (gRPC -> user-service) ==="

bold "[Admin] admin / ${PASS}"
ADMIN_TOKEN=$(fetch_token "/api/user/login" "admin")

bold "[Jobseeker] 351719672@qq.com / ${PASS}"
SEEKER_TOKEN=$(fetch_token "/api/user/userLogin" "351719672@qq.com")

bold "[HR] skyrisai / ${PASS}"
HR_TOKEN=$(fetch_token "/api/user/userLogin" "skyrisai")

echo ""
bold "=== Auth-required endpoints ==="

bold "User list (admin)"
check "Admin lists users" 200 "$(curl -s "${BASE}/api/user/list?page=1&pageSize=5" -H "Authorization: Bearer ${ADMIN_TOKEN}")"

bold "User detail (admin)"
check "Admin views user detail" 200 "$(curl -s "${BASE}/api/user/detail?userId=1" -H "Authorization: Bearer ${ADMIN_TOKEN}")"

bold "Department list (admin)"
check "Admin lists departments" 200 "$(curl -s "${BASE}/api/department/list?page=1&pageSize=5" -H "Authorization: Bearer ${ADMIN_TOKEN}")"

echo ""
bold "=== Done ==="
