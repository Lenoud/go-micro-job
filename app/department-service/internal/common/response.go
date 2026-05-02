package common

import (
	"department-service/department"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	CodeSuccess   = sharedcommon.CodeSuccess
	CodeParam     = sharedcommon.CodeParam
	CodeForbidden = sharedcommon.CodeForbidden
)

func currentTimeMillis() int64 {
	return sharedcommon.CurrentTimeMillis()
}

func OkAction(msg string) *department.ActionResp {
	return &department.ActionResp{Code: CodeSuccess, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailAction(msg string) *department.ActionResp {
	return &department.ActionResp{Code: CodeParam, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailActionForbidden(msg string) *department.ActionResp {
	return &department.ActionResp{Code: CodeForbidden, Msg: msg, Timestamp: currentTimeMillis()}
}

func OkDepartmentList(data *department.DepartmentListData) *department.DepartmentListResp {
	return &department.DepartmentListResp{Code: CodeSuccess, Msg: "操作成功", Data: data, Timestamp: currentTimeMillis()}
}

func FailDepartmentList(msg string) *department.DepartmentListResp {
	return &department.DepartmentListResp{Code: CodeParam, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailDepartmentListForbidden(msg string) *department.DepartmentListResp {
	return &department.DepartmentListResp{Code: CodeForbidden, Msg: msg, Timestamp: currentTimeMillis()}
}

func SplitIDs(raw string) []string {
	return sharedcommon.SplitIDs(raw)
}

// ==================== 业务规则冲突响应（对齐单体 1001/1002）====================

func FailActionDuplicate(msg string) *department.ActionResp {
	return &department.ActionResp{Code: sharedcommon.CodeBizDuplicate, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailActionStateConflict(msg string) *department.ActionResp {
	return &department.ActionResp{Code: sharedcommon.CodeBizStateConflict, Msg: msg, Timestamp: currentTimeMillis()}
}

// ==================== 日志记录（对齐单体 api/internal/common/logerr.go）====================
// LogErr 放在这里而不是 shared，因为需要 logx 依赖，shared 保持零依赖。

// LogErr 记录错误日志，msg 通常以"失败"结尾
func LogErr(lgr logx.Logger, msg string, err error) {
	lgr.Errorf("%s: %v", msg, err)
}
