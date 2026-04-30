#!/usr/bin/env bash
# 对比测试：单体 API (9100) vs 微服务 Gateway (9000) 的 user 模块一致性
# 用法: bash scripts/compare-test.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/colors.sh"

MONO="http://localhost:9100"
MICRO="http://localhost:9000"

PASS=0
FAIL=0

compare() {
  local label="$1" mono_out="$2" micro_out="$3"

  local mono_code micro_code mono_msg micro_msg
  mono_code=$(echo "$mono_out" | python3 -c "import sys,json; print(json.load(sys.stdin).get('code',''))" 2>/dev/null || echo "ERR")
  micro_code=$(echo "$micro_out" | python3 -c "import sys,json; print(json.load(sys.stdin).get('code',''))" 2>/dev/null || echo "ERR")
  mono_msg=$(echo "$mono_out" | python3 -c "import sys,json; print(json.load(sys.stdin).get('msg',''))" 2>/dev/null || echo "")
  micro_msg=$(echo "$micro_out" | python3 -c "import sys,json; print(json.load(sys.stdin).get('msg',''))" 2>/dev/null || echo "")

  if [ "$mono_code" = "$micro_code" ] && [ "$mono_msg" = "$micro_msg" ]; then
    print_success "  OK   [$label] code=$mono_code msg=$mono_msg"
    PASS=$((PASS + 1))
  else
    print_error "  FAIL [$label] mono: code=$mono_code msg=$mono_msg | micro: code=$micro_code msg=$micro_msg"
    FAIL=$((FAIL + 1))
  fi
}

# 检查服务是否在线
check_service() {
  local name="$1" url="$2"
  if ! curl -s -o /dev/null -w '' "$url" 2>/dev/null; then
    print_error "[$name] 服务未启动 ($url)"
    exit 1
  fi
  print_success "[$name] 服务在线"
}

echo "========================================"
echo " 单体 vs 微服务 User 模块对比测试"
echo "========================================"
echo ""

check_service "单体 API" "$MONO/api/job/list?page=1"
check_service "微服务 GW" "$MICRO/api/user/login"
echo ""

# ===== 1. 登录 =====
print_info "=== 1. 管理员登录 ==="
r1=$(curl -s -X POST "$MONO/api/user/login" -H 'Content-Type: application/json' -d '{"username":"admin","password":"lb781023"}')
r2=$(curl -s -X POST "$MICRO/api/user/login" -H 'Content-Type: application/json' -d '{"username":"admin","password":"lb781023"}')
compare "admin login" "$r1" "$r2"
T1=$(echo "$r1" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")
T2=$(echo "$r2" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

print_info "=== 2. HR 登录 ==="
r1=$(curl -s -X POST "$MONO/api/user/userLogin" -H 'Content-Type: application/json' -d '{"username":"skyrisai","password":"lb781023"}')
r2=$(curl -s -X POST "$MICRO/api/user/userLogin" -H 'Content-Type: application/json' -d '{"username":"skyrisai","password":"lb781023"}')
compare "HR login" "$r1" "$r2"
HR1=$(echo "$r1" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")
HR2=$(echo "$r2" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

print_info "=== 3. 求职者登录 ==="
r1=$(curl -s -X POST "$MONO/api/user/userLogin" -H 'Content-Type: application/json' -d '{"username":"351719672@qq.com","password":"lb781023"}')
r2=$(curl -s -X POST "$MICRO/api/user/userLogin" -H 'Content-Type: application/json' -d '{"username":"351719672@qq.com","password":"lb781023"}')
compare "求职者 login" "$r1" "$r2"

# ===== 4. 用户列表 =====
print_info "=== 4. 用户列表 ==="
r1=$(curl -s "$MONO/api/user/list?page=1&pageSize=10" -H "Authorization: Bearer $T1")
r2=$(curl -s "$MICRO/api/user/list?page=1&pageSize=10" -H "Authorization: Bearer $T2")
COUNT=$(echo "$r1" | python3 -c "import sys,json; print(len(json.load(sys.stdin)['data']['list']))")
for i in $(seq 0 $((COUNT - 1))); do
  ID=$(echo "$r1" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['list'][$i]['id'])")
  compare "list[$i] id=$ID" "$r1" "$r2"
  # 仅对比一次，因为列表整体一致即可
  if [ "$i" -eq 0 ]; then
    FIELDS_OK=true
    for field in id username nickname mobile email role status; do
      V1=$(echo "$r1" | python3 -c "import sys,json; d=json.load(sys.stdin)['data']['list'][$i]; print(d.get('$field',''))" 2>/dev/null)
      V2=$(echo "$r2" | python3 -c "import sys,json; d=json.load(sys.stdin)['data']['list'][$i]; print(d.get('$field',''))" 2>/dev/null)
      if [ "$V1" != "$V2" ]; then
        print_error "  DIFF list[$i].$field mono=$V1 micro=$V2"
        FAIL=$((FAIL + 1))
        FIELDS_OK=false
      fi
    done
    if $FIELDS_OK; then
      print_success "  OK   list[$i] 所有字段一致"
      PASS=$((PASS + 1))
    fi
  fi
done

# ===== 5. 用户详情 =====
print_info "=== 5. 用户详情 ==="
for uid in 1 2 3 4 5; do
  r1=$(curl -s "$MONO/api/user/detail?userId=$uid" -H "Authorization: Bearer $T1")
  r2=$(curl -s "$MICRO/api/user/detail?userId=$uid" -H "Authorization: Bearer $T2")
  compare "detail id=$uid" "$r1" "$r2"
  for field in username nickname mobile email role status; do
    V1=$(echo "$r1" | python3 -c "import sys,json; d=json.load(sys.stdin).get('data',{}); print(d.get('$field',''))" 2>/dev/null)
    V2=$(echo "$r2" | python3 -c "import sys,json; d=json.load(sys.stdin).get('data',{}); print(d.get('$field',''))" 2>/dev/null)
    if [ "$V1" != "$V2" ]; then
      print_error "  DIFF detail[$uid].$field mono=$V1 micro=$V2"
      FAIL=$((FAIL + 1))
    fi
  done
done

# ===== 6. 管理员 update =====
print_info "=== 6. 管理员 update ==="
r1=$(curl -s -X POST "$MONO/api/user/update" -H "Authorization: Bearer $T1" -H 'Content-Type: application/json' \
  -d '{"id":"3","nickname":"对比测试昵称"}')
r2=$(curl -s -X POST "$MICRO/api/user/update" -H "Authorization: Bearer $T2" -H 'Content-Type: application/json' \
  -d '{"id":"3","nickname":"对比测试昵称"}')
compare "update id=3" "$r1" "$r2"

# 验证更新结果
r1=$(curl -s "$MONO/api/user/detail?userId=3" -H "Authorization: Bearer $T1")
r2=$(curl -s "$MICRO/api/user/detail?userId=3" -H "Authorization: Bearer $T2")
compare "验证 update" "$r1" "$r2"

# 恢复
curl -s -X POST "$MONO/api/user/update" -H "Authorization: Bearer $T1" -H 'Content-Type: application/json' \
  -d '{"id":"3","nickname":"张求职"}' > /dev/null
curl -s -X POST "$MICRO/api/user/update" -H "Authorization: Bearer $T2" -H 'Content-Type: application/json' \
  -d '{"id":"3","nickname":"张求职"}' > /dev/null

# ===== 7. updateUserInfo =====
print_info "=== 7. updateUserInfo ==="
r1=$(curl -s -X POST "$MONO/api/user/updateUserInfo" -H "Authorization: Bearer $HR1" -H 'Content-Type: application/json' \
  -d '{"id":"2","nickname":"对比HR昵称"}')
r2=$(curl -s -X POST "$MICRO/api/user/updateUserInfo" -H "Authorization: Bearer $HR2" -H 'Content-Type: application/json' \
  -d '{"id":"2","nickname":"对比HR昵称"}')
compare "updateUserInfo" "$r1" "$r2"

r1=$(curl -s "$MONO/api/user/detail?userId=2" -H "Authorization: Bearer $T1")
r2=$(curl -s "$MICRO/api/user/detail?userId=2" -H "Authorization: Bearer $T2")
compare "验证 updateUserInfo" "$r1" "$r2"

# 恢复
curl -s -X POST "$MONO/api/user/updateUserInfo" -H "Authorization: Bearer $HR1" -H 'Content-Type: application/json' \
  -d '{"id":"2","nickname":"玛咖HR-刘经理"}' > /dev/null
curl -s -X POST "$MICRO/api/user/updateUserInfo" -H "Authorization: Bearer $HR2" -H 'Content-Type: application/json' \
  -d '{"id":"2","nickname":"玛咖HR-刘经理"}' > /dev/null

# ===== 8. updatePwd =====
print_info "=== 8. updatePwd ==="
r1=$(curl -s -X POST "$MONO/api/user/updatePwd" -H "Authorization: Bearer $T1" -H 'Content-Type: application/json' \
  -d '{"userId":"3","oldPassword":"lb781023","newPassword":"test1234"}')
r2=$(curl -s -X POST "$MICRO/api/user/updatePwd" -H "Authorization: Bearer $T2" -H 'Content-Type: application/json' \
  -d '{"userId":"3","oldPassword":"lb781023","newPassword":"test1234"}')
compare "updatePwd" "$r1" "$r2"

# 验证新密码登录
r1=$(curl -s -X POST "$MONO/api/user/userLogin" -H 'Content-Type: application/json' -d '{"username":"351719672@qq.com","password":"test1234"}')
r2=$(curl -s -X POST "$MICRO/api/user/userLogin" -H 'Content-Type: application/json' -d '{"username":"351719672@qq.com","password":"test1234"}')
compare "新密码登录" "$r1" "$r2"

# 恢复密码
curl -s -X POST "$MONO/api/user/updatePwd" -H "Authorization: Bearer $T1" -H 'Content-Type: application/json' \
  -d '{"userId":"3","oldPassword":"test1234","newPassword":"lb781023"}' > /dev/null
curl -s -X POST "$MICRO/api/user/updatePwd" -H "Authorization: Bearer $T2" -H 'Content-Type: application/json' \
  -d '{"userId":"3","oldPassword":"test1234","newPassword":"lb781023"}' > /dev/null

# ===== 9. create + delete =====
print_info "=== 9. 创建用户 ==="
r1=$(curl -s -X POST "$MONO/api/user/create" -H "Authorization: Bearer $T1" -H 'Content-Type: application/json' \
  -d '{"username":"_cmp_test_user_","password":"test123456","nickname":"对比测试","role":"1"}')
r2=$(curl -s -X POST "$MICRO/api/user/create" -H "Authorization: Bearer $T2" -H 'Content-Type: application/json' \
  -d '{"username":"_cmp_test_user_","password":"test123456","nickname":"对比测试","role":"1"}')
compare "create user" "$r1" "$r2"

print_info "=== 10. 删除测试用户 ==="
IDS_Mono=$(curl -s "$MONO/api/user/list?page=1&pageSize=20" -H "Authorization: Bearer $T1" \
  | python3 -c "import sys,json; ls=json.load(sys.stdin)['data']['list']; print(','.join(u['id'] for u in ls if u['username']=='_cmp_test_user_'))" 2>/dev/null || echo "")
IDS_Micro=$(curl -s "$MICRO/api/user/list?page=1&pageSize=20" -H "Authorization: Bearer $T2" \
  | python3 -c "import sys,json; ls=json.load(sys.stdin)['data']['list']; print(','.join(u['id'] for u in ls if u['username']=='_cmp_test_user_'))" 2>/dev/null || echo "")
if [ -n "$IDS_Mono" ] && [ -n "$IDS_Micro" ]; then
  r1=$(curl -s -X POST "$MONO/api/user/delete" -H "Authorization: Bearer $T1" -H 'Content-Type: application/json' -d "{\"ids\":\"$IDS_Mono\"}")
  r2=$(curl -s -X POST "$MICRO/api/user/delete" -H "Authorization: Bearer $T2" -H 'Content-Type: application/json' -d "{\"ids\":\"$IDS_Micro\"}")
  compare "delete user" "$r1" "$r2"
else
  print_warn "  SKIP 测试用户未找到"
fi

# ===== 11. 错误场景 =====
print_info "=== 11. 错误场景 ==="
r1=$(curl -s -X POST "$MONO/api/user/login" -H 'Content-Type: application/json' -d '{"username":"admin","password":"wrong"}')
r2=$(curl -s -X POST "$MICRO/api/user/login" -H 'Content-Type: application/json' -d '{"username":"admin","password":"wrong"}')
compare "密码错误" "$r1" "$r2"

r1=$(curl -s -X POST "$MONO/api/user/userLogin" -H 'Content-Type: application/json' -d '{"username":"admin","password":"lb781023"}')
r2=$(curl -s -X POST "$MICRO/api/user/userLogin" -H 'Content-Type: application/json' -d '{"username":"admin","password":"lb781023"}')
compare "管理员用userLogin被拒" "$r1" "$r2"

# ===== 总结 =====
echo ""
echo "========================================"
TOTAL=$((PASS + FAIL))
if [ "$FAIL" -eq 0 ]; then
  print_success "ALL $TOTAL TESTS PASSED"
else
  print_error "$FAIL / $TOTAL TESTS FAILED"
  exit 1
fi
