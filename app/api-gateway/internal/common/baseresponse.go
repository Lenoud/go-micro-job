package common

import (
	"api-gateway/internal/types"
	"time"
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

// OkBase returns a BaseResp with Code=200 and the given message.
func OkBase(msg string) types.BaseResp {
	return types.BaseResp{Code: CodeSuccess, Msg: msg, Timestamp: currentTimeMillis()}
}

// FailBase returns a BaseResp with the given code and message.
func FailBase(code int64, msg string) types.BaseResp {
	return types.BaseResp{Code: code, Msg: msg, Timestamp: currentTimeMillis()}
}

// FailBaseMsg returns a BaseResp with Code=400.
func FailBaseMsg(msg string) types.BaseResp {
	return FailBase(CodeParam, msg)
}

// FailBaseForbidden returns a BaseResp with Code=403.
func FailBaseForbidden(msg string) types.BaseResp {
	return FailBase(CodeForbidden, msg)
}

// RpcBase extracts code/msg/timestamp from a gRPC response into a BaseResp.
// For gRPC responses, the timestamp comes from user-service (no local override).
func RpcBase(resp interface {
	GetCode() int64
	GetMsg() string
	GetTimestamp() int64
}) types.BaseResp {
	return types.BaseResp{Code: resp.GetCode(), Msg: resp.GetMsg(), Timestamp: resp.GetTimestamp()}
}
