package common

import (
	"oplog-service/oplog"

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

// ==================== 业务规则冲突响应（对齐单体 1001/1002）====================

func FailActionDuplicate(msg string) *oplog.ActionResp {
	return &oplog.ActionResp{Code: sharedcommon.CodeBizDuplicate, Msg: msg, Timestamp: currentTimeMillis()}
}

func FailActionStateConflict(msg string) *oplog.ActionResp {
	return &oplog.ActionResp{Code: sharedcommon.CodeBizStateConflict, Msg: msg, Timestamp: currentTimeMillis()}
}

// ==================== 日志记录（对齐单体）====================

// LogErr 记录错误日志，msg 通常以"失败"结尾
func LogErr(lgr logx.Logger, msg string, err error) {
	lgr.Errorf("%s: %v", msg, err)
}
