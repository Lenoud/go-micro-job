package common

import (
	sharedcommon "micro-shared/common"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

// ==================== 错误码 ====================
const (
	CodeSuccess      = sharedcommon.CodeSuccess
	CodeParam        = sharedcommon.CodeParam
	CodeUnauthorized = sharedcommon.CodeUnauthorized
	CodeForbidden    = sharedcommon.CodeForbidden
	CodeNotFound     = sharedcommon.CodeNotFound
	CodeServer       = sharedcommon.CodeServer
)

func currentTimeMillis() int64 {
	return sharedcommon.CurrentTimeMillis()
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
	return sharedcommon.SplitIDs(raw)
}

// ==================== 业务规则冲突响应（对齐单体 1001/1002）====================

func FailActionDuplicate(msg string) *user.ActionResp {
	return &user.ActionResp{Code: sharedcommon.CodeBizDuplicate, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailActionStateConflict(msg string) *user.ActionResp {
	return &user.ActionResp{Code: sharedcommon.CodeBizStateConflict, Msg: msg, Timestamp: currentTimeMillis()}
}

// ==================== 日志记录（对齐单体）====================

// LogErr 记录错误日志，msg 通常以"失败"结尾
func LogErr(lgr logx.Logger, msg string, err error) {
	lgr.Errorf("%s: %v", msg, err)
}
