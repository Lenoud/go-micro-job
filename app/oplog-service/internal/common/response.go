package common

import (
	"oplog-service/oplog"

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

func OkAction(msg string) *oplog.ActionResp {
	return &oplog.ActionResp{Code: CodeSuccess, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailAction(msg string) *oplog.ActionResp {
	return &oplog.ActionResp{Code: CodeParam, Msg: msg, Timestamp: currentTimeMillis()}
}

func OkOpLogList(data *oplog.OpLogListData) *oplog.OpLogListResp {
	return &oplog.OpLogListResp{Code: CodeSuccess, Msg: "操作成功", Data: data, Timestamp: currentTimeMillis()}
}

func FailOpLogList(msg string) *oplog.OpLogListResp {
	return &oplog.OpLogListResp{Code: CodeParam, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailOpLogListForbidden(msg string) *oplog.OpLogListResp {
	return &oplog.OpLogListResp{Code: CodeForbidden, Msg: msg, Timestamp: currentTimeMillis()}
}
