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
	// DB: re_time = 请求时间(unix ms), access_time = 响应耗时(ms)
	// 前端: accessTime = "访问时间"(何时), reResponseTime = "耗时(ms)"(多久)
	responseTimeMs := ""
	if item.AccessTime > 0 {
		responseTimeMs = strconv.FormatInt(item.AccessTime, 10)
	}
	return &oplog.OpLogInfo{
		Id:             item.Id,
		RequestId:      item.RequestId,
		UserId:         item.UserId,
		ReIp:           item.ReIp,
		ReTime:         formatUnixMilli(item.ReTime),
		ReUa:           item.ReUa,
		ReUrl:          item.ReUrl,
		ReMethod:       item.ReMethod,
		ReContent:      item.ReContent,
		Success:        item.Success,
		BizCode:        item.BizCode,
		BizMsg:         item.BizMsg,
		AccessTime:     formatUnixMilli(item.ReTime), // 访问时间 = 请求发生的时间
		ReResponseTime: responseTimeMs,                // 耗时(ms) = 响应耗时
		ReUserAgent:    item.ReUa,                      // 别名：复用 re_ua
	}
}

func formatUnixMilli(value int64) string {
	if value <= 0 {
		return ""
	}
	return time.UnixMilli(value).Format("2006-01-02 15:04:05")
}
