# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

招聘系统的微服务版本，从父仓库 `go-server-resume` 拆分而来。两个 go module 通过 `go.work` 管理：

- **api-gateway**（`app/api-gateway/`）: REST 网关，端口 9000，接收 HTTP 请求，转发 gRPC 调用
- **user-service**（`app/user-service/`）: gRPC 服务，端口 9101，含用户模块全部业务逻辑和数据库访问

- **Go 版本**: 1.25.7
- **框架**: go-zero v1.10.1
- **数据库**: MySQL `micro_job`（复用父仓库 Docker 容器 3306 端口，`root:root123`）
- **服务发现**: etcd（2379）
- **缓存**: Redis（6379，密码 `redis123`）

## Build & Run

```bash
# 一键启动（检查基础设施、建库建表、编译、依次启动 user-service → api-gateway）
scripts/dev.sh

# 停止
scripts/stop-dev.sh

# 单独编译
cd app/api-gateway && go build ./...
cd app/user-service && go build ./...
```

基础设施（redis、etcd）由 `dev.sh` 自动通过 `docker compose` 启动。MySQL 依赖父仓库的 `go_job_mysql` 容器。

## Architecture

```
micro/
├── go.work                    # Go workspace: api-gateway + user-service
├── docker-compose.yml         # 基础设施 + 生产部署
├── scripts/                   # dev.sh, curl-test.sh, stop-dev.sh
├── app/
│   ├── api-gateway/           # REST 网关 (go-zero rest)
│   │   ├── gateway.api        # API 定义（唯一源文件）
│   │   ├── apigateway.go      # 入口
│   │   ├── etc/apigateway.yaml    # 生产配置
│   │   ├── etc/apigateway-local.yaml  # 本地开发配置
│   │   └── internal/
│   │       ├── config/        # Config 结构体
│   │       ├── handler/user/  # HTTP handler（goctl 生成）
│   │       ├── logic/user/    # 业务逻辑（调用 gRPC client）
│   │       ├── svc/           # ServiceContext，注入 UserRpc client
│   │       └── types/         # 请求/响应类型（goctl 生成）
│   └── user-service/          # gRPC 服务 (go-zero zrpc)
│       ├── user.proto         # Protobuf 定义（唯一源文件）
│       ├── user.go            # 入口
│       ├── etc/user.yaml      # 生产配置
│       ├── etc/user-local.yaml # 本地开发配置
│       ├── client/user/       # gRPC client（goctl 生成，api-gateway 引用）
│       ├── user/              # protobuf 生成代码（.pb.go）
│       └── internal/
│           ├── config/        # Config 结构体
│           ├── server/user/   # gRPC server（goctl 生成）
│           ├── logic/user/    # 核心业务逻辑
│           ├── model/         # 数据库 model + helpers
│           ├── svc/           # ServiceContext，注入 UserModel
│           └── common/        # 响应构造、密码加密、工具函数
```

**调用链**: HTTP → api-gateway handler → logic → gRPC client → user-service server → logic → model → MySQL

**api-gateway 是薄层**: handler 仅转发请求，logic 做 JSON 序列化/反序列化（gRPC `ApiResponse.data` 是 JSON string），所有业务逻辑在 user-service 中。

## 代码生成

### api-gateway（goctl api）

`gateway.api` 是唯一源文件。修改 API 定义后：

```bash
cd app/api-gateway && goctl api go --api gateway.api --dir . --style goZero
```

goctl 不覆盖已有 logic 文件，但会覆盖 handler/types/routes/config/svc。生成后需恢复 `config.go`、`serviceContext.go`、YAML 配置中的自定义内容。

### user-service（goctl rpc）

`user.proto` 是唯一源文件。修改 protobuf 定义后：

```bash
cd app/user-service && goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style goZero
```

同样不覆盖已有 logic 文件，但会覆盖 server/client/config/svc。生成后需恢复自定义内容。

### 新增接口的标准工作流

1. 修改 `user.proto` 添加 message + rpc 方法
2. `goctl rpc protoc` 重新生成 user-service 代码
3. 恢复 config.go / serviceContext.go 的自定义内容
4. 在 user-service `internal/logic/user/` 中实现新 logic 方法
5. 修改 `gateway.api` 添加对应 HTTP 端点
6. `goctl api go` 重新生成 api-gateway 代码
7. 恢复 config.go / serviceContext.go 的自定义内容
8. 在 api-gateway `internal/logic/user/` 中实现转发逻辑

## Key Patterns

### 统一响应
- user-service logic 返回 `*user.ApiResponse`，用 `common.Success(data)` / `common.Fail("msg")` / `common.SuccessMsg("msg", data)` / `common.SuccessPage(list, total, page, pageSize)`
- gRPC `ApiResponse.data` 字段是 JSON string，api-gateway logic 层负责 `json.Unmarshal` 转为 `interface{}`
- HTTP 响应格式：`{"code": 200, "msg": "success", "data": ..., "timestamp": ...}`

### JWT 鉴权
- Secret: 配置文件 `JWT.AccessSecret`
- Claims: `userId` / `username` / `role`
- api-gateway 通过 `rest.WithJwt()` 注册鉴权路由，JWT 验证在网关层完成

### 业务规则
- 密码加密: `MD5(明文 + "abcd1234")`
- 批量删除: ids 为逗号分隔字符串
- 用户角色: 1=求职者, 2=HR, 3=管理员
- 登录接口区分: `/login` 仅管理员(role=3)，`/userLogin` 仅求职者/HR(role=1,2)

### Model 层
- 使用 `IFNULL(col, default) AS col` 处理 NULL，定义 `xxxFields` 常量复用字段列表
- `NullTime` / `NullString` 类型处理可空字段
- 分页: `FindList` 返回 `(list, total, error)`
- sqlx scan 要求 SQL 列与 struct `db:` tag 一一对应

## curl 测试

```bash
# 自动获取三种角色 token 并验证
scripts/curl-test.sh
```

测试账号（与父仓库共享 MySQL 数据）：

| 角色 | username | password | 登录接口 |
|------|----------|----------|----------|
| 管理员(3) | admin | lb781023 | /api/user/login |
| 求职者(1) | 351719672@qq.com | lb781023 | /api/user/userLogin |
| HR(2) | skyrisai | lb781023 | /api/user/userLogin |

手动获取 token：
```bash
# 管理员
TOKEN=$(curl -s -X POST 'http://localhost:9000/api/user/login' \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"lb781023"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

# 使用 token
curl -s 'http://localhost:9000/api/user/list?page=1&pageSize=5' \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
```

## 添加新微服务的模式

参照 user-service 的模式：
1. 创建新的 protobuf 定义文件
2. 用 goctl rpc 生成 gRPC 服务骨架
3. 实现业务逻辑（logic → model）
4. 在 api-gateway 的 `gateway.api` 中添加 HTTP 端点
5. 在 api-gateway 的 `config.go` 中添加新 RPC client 配置
6. 在 api-gateway 的 `serviceContext.go` 中注入新 RPC client
7. 在 `docker-compose.yml` 中添加新服务
8. 更新 `dev.sh` 启动脚本
