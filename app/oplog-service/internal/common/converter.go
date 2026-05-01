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
	accessTime := ""
	if item.AccessTime > 0 {
		accessTime = strconv.FormatInt(item.AccessTime, 10)
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
		AccessTime:     accessTime,
		ReResponseTime: accessTime, // 别名：数据库无独立 re_response_time 列，复用 access_time
		ReUserAgent:    item.ReUa,  // 别名：数据库无独立 re_user_agent 列，复用 re_ua
	}
}

func formatUnixMilli(value int64) string {
	if value <= 0 {
		return ""
	}
	return time.UnixMilli(value).Format("2006-01-02 15:04:05")
}
