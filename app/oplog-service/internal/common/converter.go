package common

import (
	"strconv"
	"time"

	"oplog-service/internal/model"
	"oplog-service/oplog"
)

func OpLogModelToProto(item *model.OpLog) *oplog.OpLogInfo {
	if item == nil {
		return nil
	}
	responseMs := ""
	if item.ResponseMs > 0 {
		responseMs = strconv.FormatInt(item.ResponseMs, 10)
	}
	return &oplog.OpLogInfo{
		Id:             item.Id,
		RequestId:      item.RequestId,
		UserId:         item.UserId,
		ReIp:           item.ReIp,
		ReTime:         formatUnixMilli(item.RequestTime),
		ReUa:           item.ReUa,
		ReUrl:          item.ReUrl,
		ReMethod:       item.ReMethod,
		ReContent:      item.ReContent,
		Success:        item.Success,
		BizCode:        item.BizCode,
		BizMsg:         item.BizMsg,
		ReResponseTime: responseMs,
		ReUserAgent:    item.ReUa,
	}
}

func formatUnixMilli(value int64) string {
	if value <= 0 {
		return ""
	}
	return time.UnixMilli(value).Format("2006-01-02 15:04:05")
}
