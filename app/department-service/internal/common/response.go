package common

import (
	"strings"
	"time"

	"department-service/department"
)

const (
	CodeSuccess   int64 = 200
	CodeParam     int64 = 400
	CodeForbidden int64 = 403
)

func currentTimeMillis() int64 {
	return time.Now().UnixMilli()
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
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
