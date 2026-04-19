package common

import (
	"strings"

	"user-service/user"
)

// ---- 错误码常量 ----

const (
	CodeSuccess int64 = 200
	CodeParam   int64 = 400
)

// ---- 类型化响应 helpers ----

func SuccessAction() *user.ActionResp {
	return &user.ActionResp{Code: CodeSuccess, Msg: "操作成功"}
}

func SuccessActionMsg(msg string) *user.ActionResp {
	return &user.ActionResp{Code: CodeSuccess, Msg: msg}
}

func FailAction(msg string) *user.ActionResp {
	return &user.ActionResp{Code: CodeParam, Msg: msg}
}

func SuccessUserInfo(msg string, info *user.UserInfo) *user.UserInfoResp {
	return &user.UserInfoResp{Code: CodeSuccess, Msg: msg, Data: info}
}

func FailUserInfo(msg string) *user.UserInfoResp {
	return &user.UserInfoResp{Code: CodeParam, Msg: msg}
}

func SuccessUserList(data *user.UserListData) *user.UserListResp {
	return &user.UserListResp{Code: CodeSuccess, Msg: "操作成功", Data: data}
}

func FailUserList(msg string) *user.UserListResp {
	return &user.UserListResp{Code: CodeParam, Msg: msg}
}

// ---- 工具函数 ----

func SplitIDs(ids string) []string {
	if ids == "" {
		return nil
	}
	parts := make([]string, 0)
	for _, id := range splitComma(ids) {
		id = trimSpace(id)
		if id != "" {
			parts = append(parts, id)
		}
	}
	return parts
}

func splitComma(s string) []string {
	result := make([]string, 0)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}
