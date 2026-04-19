# go-micro-job

招聘系统微服务架构，基于 go-zero 框架，从 [go-server-resume](https://github.com/Lenoud/go-server-resume) 单体后端拆分。

## 服务一览

| 服务 | 类型 | 端口 | 说明 |
|------|------|------|------|
| api-gateway | REST API | 9000 | API 网关，JWT 鉴权，路由转发到 gRPC 服务 |
| user-service | gRPC | 9101 | 用户模块，注册/登录/CRUD |

## 依赖

- Go 1.25.7
- go-zero v1.10.1
- etcd v3.5（服务注册/发现）
- Redis 7（缓存）
- MySQL 8（共享父仓库 `go_job_mysql` 容器，使用独立库 `micro_job`）

## 目录结构

```
├── app/
│   ├── api-gateway/       # REST API 网关
│   │   ├── gateway.api    # API 定义（goctl 源文件）
│   │   ├── apigateway.go  # 入口
│   │   ├── etc/           # 配置文件
│   │   └── internal/
│   │       ├── handler/   # HTTP handler（goctl 生成）
│   │       ├── logic/     # 业务逻辑
│   │       ├── svc/       # ServiceContext
│   │       └── types/     # 请求/响应类型
│   └── user-service/      # gRPC 用户服务
│       ├── user.proto     # Protobuf 定义
│       ├── user.go        # 入口
│       ├── etc/           # 配置文件
│       └── internal/
│           ├── logic/     # 业务逻辑
│           ├── model/     # 数据库 model
│           ├── server/    # gRPC server
│           ├── svc/       # ServiceContext
│           └── common/    # 工具函数
├── scripts/               # 开发脚本
├── go.work                # Go workspace
└── docker-compose.yml     # 基础设施编排（Redis + etcd）
```

## 快速开始

### 前置条件

1. MySQL 容器运行中（来自父仓库 `docker compose up -d mysql`）
2. Redis 和 etcd 启动：

```bash
docker compose up -d redis etcd
```

### 启动开发环境

```bash
scripts/dev.sh
```

自动完成：构建 → 启动 user-service → 启动 api-gateway → 日志输出到 `logs/`。

停止：`Ctrl+C` 或 `scripts/stop-dev.sh`

### 测试

```bash
scripts/curl-test.sh
```

自动获取三种角色 token 并验证接口。

### 手动测试

```bash
# 登录拿 token
TOKEN=$(curl -s -X POST 'http://localhost:9000/api/user/login' \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"lb781023"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

# 用户列表
curl -s 'http://localhost:9000/api/user/list?page=1&pageSize=5' \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
```

## API 接口

### 无鉴权

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/user/login | 管理员登录（role=3） |
| POST | /api/user/userLogin | 前台登录（role=1/2） |
| POST | /api/user/userRegister | 用户注册 |

### 需鉴权（JWT）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/user/list | 用户列表 |
| GET | /api/user/detail | 用户详情 |
| POST | /api/user/create | 创建用户 |
| POST | /api/user/update | 更新用户 |
| POST | /api/user/delete | 批量删除用户 |
| POST | /api/user/updateUserInfo | 用户更新自己信息 |
| POST | /api/user/updatePwd | 修改密码 |

## 代码生成

### api-gateway（goctl）

修改 `gateway.api` 后重新生成：

```bash
cd app/api-gateway && goctl api go --api gateway.api --dir . --style goZero
```

### user-service（protoc + goctl）

修改 `user.proto` 后重新生成：

```bash
cd app/user-service && goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style goZero
```

## Docker 部署

```bash
docker compose up -d
```

包含 Redis、etcd、user-service、api-gateway 四个容器。
MySQL 使用独立库 `micro_job`，依赖父仓库的 MySQL 容器。
