package common

import "time"

// ==================== 错误码（对齐单体 api/internal/common/baseresponse.go）====================
const (
	CodeSuccess      int64 = 0   // 成功（单体已从 200 改为 0）
	CodeParam        int64 = 400 // 参数错误 / 通用业务失败
	CodeUnauthorized int64 = 401 // 未登录
	CodeForbidden    int64 = 403 // 无权限
	CodeNotFound     int64 = 404 // 资源不存在
	CodeServer       int64 = 500 // 服务器内部错误

	// 1xxx — 业务规则冲突（对齐单体）
	CodeBizDuplicate     int64 = 1001 // 资源已存在
	CodeBizStateConflict int64 = 1002 // 状态不允许
)

func CurrentTimeMillis() int64 {
	return time.Now().UnixMilli()
}
