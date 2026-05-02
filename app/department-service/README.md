# department-service

部门管理微服务，提供部门的增删改查 gRPC 接口。

## 技术栈

| 组件 | 选型 | 说明 |
|------|------|------|
| 框架 | go-zero | gRPC 微服务 |
| ORM | **Ent** | Facebook 出品的 Go ORM，代码生成，类型安全 |
| 数据库 | MySQL | 表 `b_department` |
| 注册中心 | Etcd | 服务发现 |

## 为什么用 Ent

本服务作为 Ent ORM 的实验性试点，替代了之前手写 sqlx SQL 的方式：

- **之前**：每个 CRUD 方法手写 SQL 字符串，加字段要改 4 处，错误容易被 `_ =` 吞掉
- **现在**：只写 schema 定义（20 行），Ent 自动生成完整的 CRUD 代码

## 项目结构

```
department-service/
├── ent/
│   ├── schema/
│   │   └── department.go       # ← 唯一需要手写的文件
│   ├── client.go               # 自动生成：Ent client 入口
│   ├── department.go           # 自动生成：Department struct
│   ├── department_create.go    # 自动生成：创建操作
│   ├── department_query.go     # 自动生成：查询操作
│   ├── department_update.go    # 自动生成：更新操作
│   ├── department_delete.go    # 自动生成：删除操作
│   └── department/             # 自动生成：查询谓词（Where 条件）
├── internal/
│   ├── logic/                  # 业务逻辑，直接调用 Ent client
│   ├── common/                 # 响应构造、权限校验、日志
│   ├── svc/                    # 注入 Ent client
│   └── server/                 # gRPC server
├── department.proto            # gRPC 接口定义
└── department.go               # 启动入口
```

## 开发流程

### 新增字段

1. 编辑 `ent/schema/department.go`，在 `Fields()` 里加一行
2. 执行 `go generate ./ent`
3. 修改 `internal/logic/listLogic.go` 中的内联转换（如果需要返回新字段）
4. 完成

### 新建表

1. 在 `ent/schema/` 下新建 schema 文件
2. 执行 `go generate ./ent`
3. 在 `svc/serviceContext.go` 中已经有 `EntClient`，直接用
4. 写 logic 调用 `EntClient.Xxx.Create().Save()` 等

### 代码生成命令

```bash
# 生成 Ent 代码
cd /path/to/department-service
go generate ./ent

# 生成 proto 代码
protoc --go_out=. --go-grpc_out=. department.proto
```

## 数据库表结构

对应 MySQL 表 `b_department`：

```sql
CREATE TABLE b_department (
  id          int NOT NULL AUTO_INCREMENT,
  title       varchar(100) NOT NULL DEFAULT '',
  description varchar(500) DEFAULT '',
  parent_id   int DEFAULT '0',
  create_time datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
```

Ent schema 完全对应此表结构，**不使用 Ent 的自动迁移**，表结构由 SQL 脚本管理。

## 注意事项

- **不要调用 `client.Schema.Create()`**，Ent 只用于查询，不负责建表
- `ent/` 目录下的文件（除 `ent/schema/`）都是自动生成的，不要手动修改
- proto 的 id/parentId 类型是 `int64`，Ent 生成的 Go 类型是 `int`，需要显式 `int()` / `int64()` 转换
