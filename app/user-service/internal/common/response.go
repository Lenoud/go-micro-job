package common

import (
	"strings"
	"time"

	"user-service/user"
)

// ==================== 错误码 ====================
const (
	CodeSuccess      int64 = 200
	CodeParam        int64 = 400
	CodeUnauthorized int64 = 401
	CodeForbidden    int64 = 403
	CodeNotFound     int64 = 404
	CodeServer       int64 = 500
)

func currentTimeMillis() int64 {
	return time.Now().UnixMilli()
}

// ==================== ActionResp helpers ====================

func OkAction(msg string) *user.ActionResp {
	return &user.ActionResp{Code: CodeSuccess, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailAction(msg string) *user.ActionResp {
	return &user.ActionResp{Code: CodeParam, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailActionForbidden(msg string) *user.ActionResp {
	return &user.ActionResp{Code: CodeForbidden, Msg: msg, Timestamp: currentTimeMillis()}
}

// ==================== UserInfoResp helpers ====================

func OkUserInfo(msg string, info *user.UserInfo) *user.UserInfoResp {
	return &user.UserInfoResp{Code: CodeSuccess, Msg: msg, Data: info, Timestamp: currentTimeMillis()}
}

func FailUserInfo(msg string) *user.UserInfoResp {
	return &user.UserInfoResp{Code: CodeParam, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailUserInfoForbidden(msg string) *user.UserInfoResp {
	return &user.UserInfoResp{Code: CodeForbidden, Msg: msg, Timestamp: currentTimeMillis()}
}

// ==================== UserListResp helpers ====================

func OkUserList(data *user.UserListData) *user.UserListResp {
	return &user.UserListResp{Code: CodeSuccess, Msg: "操作成功", Data: data, Timestamp: currentTimeMillis()}
}

func FailUserList(msg string) *user.UserListResp {
	return &user.UserListResp{Code: CodeParam, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailUserListForbidden(msg string) *user.UserListResp {
	return &user.UserListResp{Code: CodeForbidden, Msg: msg, Timestamp: currentTimeMillis()}
}

// ==================== 工具函数 ====================

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
