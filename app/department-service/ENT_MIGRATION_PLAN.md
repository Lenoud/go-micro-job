# department-service Ent 迁移 + 单体对齐计划（方案 B）

## 一、目标

两件事合并执行：

1. **Ent 迁移**：去掉 model 层手写 SQL，逻辑层直接用 Ent client 操作数据库
2. **单体对齐**：同步单体 api/ 近 20 次提交的改进（错误码、日志、响应模式）

**方案 B 原则**：不用兼容层，一步到位。proto 的 id 类型对齐数据库（int64），直接用 Ent 生成的 struct，删除手写的 model.Department 和 DepartmentModel 接口。

---

## 二、改动范围

### 为什么要改 proto

当前 proto 的 id 和 parentId 都是 `string`，但 MySQL 表里是 `int`。之前的方案为了不动 proto，加了一堆 `entToModel()`、`parseID()`、`sql1Result` 转换层。方案 B 直接改 proto 对齐数据库类型，从根源消除转换。

### 影响范围

| 改动 | 说明 |
|------|------|
| **department.proto** | id/parentId 从 string → int64 |
| **department/\*.pb.go** | 重新生成 |
| **internal/model/** | 删除 departmentmodel.go 和 helpers.go，整个 model 包清空 |
| **internal/common/converter.go** | 删除（对齐单体，转换内联到 logic） |
| **internal/common/response.go** | 加 LogErr、FailActionDuplicate、FailActionStateConflict |
| **internal/svc/serviceContext.go** | 注入 Ent client，不再注入 model |
| **internal/logic/ (4个文件)** | 直接用 Ent client，用 Msg+LogErr 模式 |
| **department.go** | 加 defer EntClient.Close() |
| **shared/common/response.go** | CodeSuccess 200→0，加 1001/1002 |
| **shared/common/msg.go** | 新增 Msg() |

### 跨服务影响（需注意但本次不改）

| 服务 | 影响 | 原因 |
|------|------|------|
| api-gateway | **会编译失败** | 引用了 `department.DepartmentInfo.Id`（原为 string，现为 int64） |
| web 前端 | JSON 中 id 从字符串变数字 | `"id":"1"` → `"id":1` |

api-gateway 和前端的适配是后续步骤，本次只改 department-service + shared。

---

## 三、单体 vs 微服务差异分析

### 3.1 错误码

**单体 api/internal/common/baseresponse.go：**
```go
const (
    CodeSuccess          int64 = 0
    CodeParam            int64 = 400
    CodeUnauthorized     int64 = 401
    CodeForbidden        int64 = 403
    CodeBizDuplicate     int64 = 1001
    CodeBizStateConflict int64 = 1002
)
```

**微服务 shared/common/response.go（当前）：**
```go
const (
    CodeSuccess      int64 = 200    // ← 应改为 0
    // ← 缺少 1001、1002
)
```

### 3.2 错误处理模式

**单体 api logic（以 department create 为例）：**
```go
if err != nil {
    msg := common.Msg("创建", "部门")
    common.LogErr(l.Logger, msg, err)
    return ..., nil
}
```

**微 service logic（当前）：**
```go
if err != nil {
    return common.FailAction("创建部门失败"), nil  // 硬编码，无日志
}
```

### 3.3 四个 logic 文件对照

| 文件 | 单体做法 | 微 service 当前 | 需改 |
|------|---------|----------------|------|
| createLogic | Msg+LogErr+Fail | 硬编码 FailAction | **改** |
| deleteLogic | Msg+LogErr+Fail | 硬编码 FailAction | **改** |
| listLogic | Msg+LogErr+Fail | 硬编码 FailDepartmentList | **改** |
| updateLogic | Msg+LogErr+Fail | 硬编码 FailAction | **改** |

---

## 四、执行步骤

### 步骤 1：对齐 shared/common

#### 1a. 修改 shared/common/response.go

```go
package common

import "time"

// ==================== 错误码（对齐单体）====================
const (
	CodeSuccess      int64 = 0    // 成功（从 200 改为 0）
	CodeParam        int64 = 400
	CodeUnauthorized int64 = 401
	CodeForbidden    int64 = 403
	CodeNotFound     int64 = 404
	CodeServer       int64 = 500

	// 1xxx — 业务规则冲突（对齐单体）
	CodeBizDuplicate     int64 = 1001 // 资源已存在
	CodeBizStateConflict int64 = 1002 // 状态不允许
)

func CurrentTimeMillis() int64 {
	return time.Now().UnixMilli()
}
```

#### 1b. 新增 shared/common/msg.go

```go
package common

// Msg 生成标准错误消息：op + target + "失败"
func Msg(op, target string) string {
	return op + target + "失败"
}
```

#### 1c. 编译验证

```bash
cd /Users/bobo/git_project/go-server-resume/micro/app/shared
go build ./...
```

- [ ] 编译通过
- [ ] shared/go.mod 没变（零外部依赖）

---

### 步骤 2：修改 proto

**文件：** `department.proto`

**当前：**
```protobuf
message DepartmentInfo {
  string id = 1;
  string title = 2;
  string description = 3;
  string parentId = 4;
  string createTime = 5;
}

message CreateDepartmentReq {
  string title = 1;
  string description = 2;
  string parentId = 3;
  DepartmentContext auth = 4;
}

message UpdateDepartmentReq {
  string id = 1;
  string title = 2;
  string description = 3;
  string parentId = 4;
  DepartmentContext auth = 5;
}

message DeleteReq {
  string ids = 1;
  DepartmentContext auth = 2;
}
```

**改为：**
```protobuf
message DepartmentInfo {
  int64 id = 1;           // string → int64，对齐数据库 int 类型
  string title = 2;
  string description = 3;
  int64 parentId = 4;     // string → int64
  string createTime = 5;
}

message CreateDepartmentReq {
  string title = 1;
  string description = 2;
  int64 parentId = 3;     // string → int64
  DepartmentContext auth = 4;
}

message UpdateDepartmentReq {
  int64 id = 1;           // string → int64
  string title = 2;
  string description = 3;
  int64 parentId = 4;     // string → int64
  DepartmentContext auth = 5;
}

// DeleteReq.ids 保持 string（逗号分隔），改动 api-gateway 范围太大，留后续
message DeleteReq {
  string ids = 1;
  DepartmentContext auth = 2;
}
```

#### 重新生成 proto

```bash
cd /Users/bobo/git_project/go-server-resume/micro/app/department-service
protoc --go_out=. --go-grpc_out=. department.proto
```

- [ ] `department/department.pb.go` 和 `department/department_grpc.pb.go` 已更新
- [ ] `DepartmentInfo.Id` 类型为 `int64`
- [ ] `CreateDepartmentReq.ParentId` 类型为 `int64`

---

### 步骤 3：初始化 Ent

```bash
cd /Users/bobo/git_project/go-server-resume/micro/app/department-service

go install entgo.io/ent/cmd/ent@latest
go run -mod=mod entgo.io/ent/cmd/ent init Department
go get entgo.io/ent@latest
go get github.com/go-sql-driver/mysql
```

- [ ] `ent/schema/department.go` 存在
- [ ] `go build ./...` 通过

---

### 步骤 4：编写 Ent Schema

**文件：** `ent/schema/department.go`

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type Department struct {
	ent.Schema
}

func (Department) Config() ent.Config {
	return ent.Config{
		Table: "b_department",
	}
}

func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			AutoIncrement(),
		field.String("title").
			Default(""),
		field.String("description").
			Default("").
			Optional(),
		field.Int("parent_id").
			Default(0).
			Optional(),
		field.Time("create_time").
			Default(time.Now).
			Optional(),
	}
}
```

**对照 MySQL 表结构：**
```
id          int NOT NULL AUTO_INCREMENT    → field.Int + AutoIncrement     ✅
title       varchar(100) NOT NULL DEFAULT '' → field.String + Default("")   ✅
description varchar(500) DEFAULT ''        → field.String + Optional       ✅
parent_id   int DEFAULT '0'                → field.Int + Default(0)        ✅
create_time datetime DEFAULT CURRENT_TIMESTAMP → field.Time + Default(now) ✅
```

- [ ] 5 个字段一一对应

---

### 步骤 5：生成 Ent 代码

```bash
go generate ./ent
```

- [ ] `ent/client.go`、`ent/department.go`、`ent/department_create.go` 等文件已生成
- [ ] `go build ./ent/...` 通过

---

### 步骤 6：删除 model 包

```bash
rm internal/model/departmentmodel.go
rm internal/model/helpers.go
```

Ent 已经生成了完整的 Department struct 和 CRUD 操作，手写的 model 层完全不需要了。

- [ ] `internal/model/` 目录已清空（或删除整个目录）
- [ ] 编译会报错（预期中，后续步骤修复）

---

### 步骤 7：删除 converter.go

**单体 api 已经删掉了独立的 converter 文件**，转换逻辑内联在 logic 里。微服务也照做。

```bash
rm internal/common/converter.go
```

- [ ] converter.go 已删除

---

### 步骤 8：修改 common/response.go

**文件：** `internal/common/response.go`

**在现有 import 中加 logx：**
```go
import (
	"department-service/department"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)
```

**在文件末尾追加：**
```go
// ==================== 业务规则冲突响应（对齐单体）====================

func FailActionDuplicate(msg string) *department.ActionResp {
	return &department.ActionResp{Code: sharedcommon.CodeBizDuplicate, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailActionStateConflict(msg string) *department.ActionResp {
	return &department.ActionResp{Code: sharedcommon.CodeBizStateConflict, Msg: msg, Timestamp: currentTimeMillis()}
}

// ==================== 日志记录（对齐单体 logerr.go）====================

func LogErr(lgr logx.Logger, msg string, err error) {
	lgr.Errorf("%s: %v", msg, err)
}
```

- [ ] 原有函数不动
- [ ] 新增 3 个函数

---

### 步骤 9：修改 ServiceContext

**文件：** `internal/svc/serviceContext.go`

**改为：**
```go
package svc

import (
	"department-service/ent"
	"department-service/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type ServiceContext struct {
	Config    config.Config
	EntClient *ent.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := ent.Open("mysql", c.MySQL.DataSource)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:    c,
		EntClient: client,
	}
}
```

**关键变化：**
- 不再有 `DepartmentModel` 接口
- 不再 import `model` 和 `sqlx`
- 只暴露 `EntClient`，logic 层直接用

---

### 步骤 10：修改 department.go（main）

**文件：** `department.go`

```go
// 在 defer s.Stop() 后面加一行：
defer s.Stop()
defer ctx.EntClient.Close()  // ← 新增
```

---

### 步骤 11：重写 4 个 logic 文件

每个 logic 文件两个变化：
1. 不再通过 `DepartmentModel` 接口，直接用 `l.svcCtx.EntClient.Department`
2. 错误处理用 `Msg + LogErr` 模式

#### 11a. createLogic.go

**改为：**
```go
package logic

import (
	"context"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/svc"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *department.CreateDepartmentReq) (*department.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}
	if in.GetTitle() == "" {
		return common.FailAction("部门名称不能为空"), nil
	}

	_, err := l.svcCtx.EntClient.Department.Create().
		SetTitle(in.GetTitle()).
		SetDescription(in.GetDescription()).
		SetParentID(int(in.GetParentId())).
		Save(l.ctx)
	if err != nil {
		msg := sharedcommon.Msg("创建", "部门")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("创建成功"), nil
}
```

**和当前版本的差异：**
- 删除了 `import "department-service/internal/model"`
- 新增了 `import sharedcommon "micro-shared/common"`
- 删除了 `data := &model.Department{...}` 中间变量
- 直接用 Ent 链式调用创建
- parentId 直接用 `in.GetParentId()`（int64），不用 string↔int 转换
- 错误处理改为 Msg+LogErr

#### 11b. deleteLogic.go

**改为：**
```go
package logic

import (
	"context"
	"fmt"

	"department-service/department"
	entdep "department-service/ent/department"
	"department-service/internal/common"
	"department-service/internal/svc"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogic) Delete(in *department.DeleteReq) (*department.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}
	idList := sharedcommon.SplitIDs(in.GetIds())
	if len(idList) == 0 {
		return common.FailAction("删除部门失败"), nil
	}

	// string ids → int ids
	intIDs := make([]int, 0, len(idList))
	for _, s := range idList {
		var n int
		if _, err := fmt.Sscanf(s, "%d", &n); err == nil && n > 0 {
			intIDs = append(intIDs, n)
		}
	}
	if len(intIDs) == 0 {
		return common.OkAction("删除成功"), nil
	}

	_, err := l.svcCtx.EntClient.Department.Delete().
		Where(entdep.IDIn(intIDs...)).
		Exec(l.ctx)
	if err != nil {
		msg := sharedcommon.Msg("删除", "部门")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("删除成功"), nil
}
```

**注意：** `DeleteReq.ids` 仍然是 string（逗号分隔），因为改 proto 的这个字段会影响 api-gateway，留后续处理。所以这里保留了一处 string→int 解析。这是唯一保留的"转换"，而且是业务需要（从请求参数解析），不是兼容层。

#### 11c. listLogic.go

**改为：**
```go
package logic

import (
	"context"

	"department-service/department"
	"department-service/ent"
	entdep "department-service/ent/department"
	"department-service/internal/common"
	"department-service/internal/svc"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *department.DepartmentListReq) (*department.DepartmentListResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleRecruiter, common.RoleAdmin) {
		return common.FailDepartmentListForbidden("无权访问"), nil
	}

	page := in.GetPage()
	pageSize := in.GetPageSize()
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 构建查询条件（Count 和 List 分别创建 query，避免复用）
	whereKeyword := in.GetKeyword() != ""

	// Count 查询
	countQuery := l.svcCtx.EntClient.Department.Query()
	if whereKeyword {
		countQuery = countQuery.Where(entdep.TitleContains(in.GetKeyword()))
	}
	total, err := countQuery.Count(l.ctx)
	if err != nil {
		msg := sharedcommon.Msg("查询", "部门列表")
		common.LogErr(l.Logger, msg, err)
		return common.FailDepartmentList(msg), nil
	}

	// List 查询（新建 query，不复用 countQuery）
	listQuery := l.svcCtx.EntClient.Department.Query()
	if whereKeyword {
		listQuery = listQuery.Where(entdep.TitleContains(in.GetKeyword()))
	}
	items, err := listQuery.
		Order(ent.Desc(entdep.FieldCreateTime)).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		All(l.ctx)
	if err != nil {
		msg := sharedcommon.Msg("查询", "部门列表")
		common.LogErr(l.Logger, msg, err)
		return common.FailDepartmentList(msg), nil
	}

	// 转换逻辑内联（对齐单体 api 做法已删除 converter.go）
	infos := make([]*department.DepartmentInfo, 0, len(items))
	for _, item := range items {
		createTime := ""
		if !item.CreateTime.IsZero() {
			createTime = item.CreateTime.Format("2006-01-02 15:04:05")
		}
		infos = append(infos, &department.DepartmentInfo{
			Id:          int64(item.ID),
			Title:       item.Title,
			Description: item.Description,
			ParentId:    int64(item.ParentID),
			CreateTime:  createTime,
		})
	}
	return common.OkDepartmentList(&department.DepartmentListData{
		List:     infos,
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
	}), nil
}
```

#### 11d. updateLogic.go

**改为：**
```go
package logic

import (
	"context"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/svc"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *department.UpdateDepartmentReq) (*department.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}
	if in.GetId() == 0 {
		return common.FailAction("部门不存在"), nil
	}

	existing, err := l.svcCtx.EntClient.Department.Get(l.ctx, int(in.GetId()))
	if err != nil || existing == nil {
		return common.FailAction("部门不存在"), nil
	}

	update := l.svcCtx.EntClient.Department.UpdateOneID(existing.ID)
	if in.GetTitle() != "" {
		update = update.SetTitle(in.GetTitle())
	}
	if in.GetDescription() != "" {
		update = update.SetDescription(in.GetDescription())
	}
	if in.GetParentId() != 0 {
		update = update.SetParentID(int(in.GetParentId()))
	}

	if err := update.Exec(l.ctx); err != nil {
		msg := sharedcommon.Msg("更新", "部门")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("更新成功"), nil
}
```

**关键变化：** 不再先查出来改字段再存回去，直接用 `UpdateOneID` 链式调用选择性更新。

---

### 步骤 12：全量编译

```bash
cd /Users/bobo/git_project/go-server-resume/micro/app/department-service
go build ./...
go vet ./...
```

- [ ] 编译通过
- [ ] 无 vet 警告

---

### 步骤 13：功能验证

```bash
go run department.go -f etc/department-local.yaml
```

#### Ent CRUD 验证

| # | 操作 | 预期 | 验证边界 |
|---|------|------|----------|
| 1 | 创建部门 title="测试A", parentId=0 | Code=**0**，成功 | 成功码从 200 变为 0 |
| 2 | 创建部门 title="" | Code=400 "部门名称不能为空" | 空值 |
| 3 | 查列表 keyword="" | 返回列表，id 为 int64 | id 类型验证 |
| 4 | 查列表 keyword="测试" | 只返回匹配项 | 模糊搜索 |
| 5 | 查列表 page=1, pageSize=1 | total > 1，list 1 条 | 分页 |
| 6 | 查列表 page=0 | 自动修正为 1 | 非法页码 |
| 7 | 更新部门 id=1, title="新名称" | 成功 | int64 id |
| 8 | 更新部门 id=0 | "部门不存在" | 零值边界 |
| 9 | 更新不存在的 id=999999 | "部门不存在" | Ent NotFound |
| 10 | 删除 ids="1" | 成功 | 单个删除 |
| 11 | 删除 ids="1,2" | 成功 | 批量删除 |
| 12 | 删除 ids="" | 成功（空值跳过） | 空值 |
| 13 | CreateTime 格式 | "2006-01-02 15:04:05" | time.Time 转换 |

#### 对齐验证

| # | 操作 | 预期 |
|---|------|------|
| 14 | 成功响应 Code | 0（不是 200） |
| 15 | 错误 Msg 格式 | "创建部门失败"/"查询部门列表失败" |
| 16 | 停 MySQL 触发错误 | 日志有 `"创建部门失败: <具体错误>"` |
| 17 | 正常操作后检查日志 | 无错误日志 |

**通过标准：** 17 个用例全部通过。

---

## 五、回退方案

```bash
cd /Users/bobo/git_project/go-server-resume/micro/app/department-service

# 恢复所有文件
git checkout department.proto
git checkout department/
git checkout internal/
git checkout svc/
git checkout department.go
git checkout go.mod go.sum
rm -rf ent/

# 恢复 shared
cd ../shared
git checkout common/response.go
rm -f common/msg.go

# 验证
cd ../department-service
go build ./...
```

---

## 六、已知风险

| 风险 | 影响 | 应对 |
|------|------|------|
| api-gateway 引用旧 proto 类型 | 编译失败 | 本次不改 gateway，单独做 gateway 适配 |
| 前端 JSON id 从 string 变 number | 可能报错 | 前端适配也是后续步骤 |
| DeleteReq.ids 保持 string | 有一处 string→int 解析 | 这是业务需要不是兼容层，改 ids 类型在 gateway 适配时一起做 |
| CodeSuccess 200→0 | user-service/oplog-service 成功码也跟着变 | 期望行为，和单体对齐 |

---

## 七、完整文件变更清单

```
micro/app/
├── shared/common/
│   ├── response.go                 # 修改：CodeSuccess 200→0，加 1001/1002
│   └── msg.go                      # 新增：Msg()
│   # shared/go.mod 不动
│
└── department-service/
    ├── department.proto            # 修改：id/parentId string → int64
    ├── department/                 # 重新生成：pb.go 文件
    ├── ent/                        # 新增：Ent schema + 生成代码（约 22 个文件）
    ├── internal/
    │   ├── model/                  # 删除：departmentmodel.go + helpers.go
    │   ├── common/
    │   │   ├── response.go         # 修改：加 LogErr + FailActionDuplicate/StateConflict
    │   │   ├── converter.go        # 删除（对齐单体，转换内联到 logic）
    │   │   └── roles.go            # 不动
    │   ├── logic/
    │   │   ├── createLogic.go      # 重写：直接用 Ent client
    │   │   ├── deleteLogic.go      # 重写：直接用 Ent client
    │   │   ├── listLogic.go        # 重写：直接用 Ent client
    │   │   └── updateLogic.go      # 重写：直接用 Ent client
    │   ├── server/
    │   │   └── departmentServer.go # 不动
    │   └── svc/
    │       └── serviceContext.go   # 重写：只注入 EntClient
    ├── go.mod                      # 修改：加 ent 依赖
    ├── go.sum                      # 修改
    └── department.go               # 修改：defer EntClient.Close()
```

**相比方案 A 省掉的东西：**
- ~~model.Department struct~~ — 直接用 Ent 生成的
- ~~DepartmentModel 接口~~ — logic 直接调 Ent client
- ~~entToModel() 转换函数~~ — 不需要
- ~~parseID() 转换函数~~ — proto 已经是 int64
- ~~sql1Result 适配器~~ — logic 不需要 sql.Result
- ~~NullTime 类型~~ — Ent 返回 time.Time，用 IsZero() 判断
