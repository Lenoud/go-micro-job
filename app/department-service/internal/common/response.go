package common

import (
	"department-service/department"
	sharedcommon "micro-shared/common"
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
