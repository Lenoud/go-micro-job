package common

import (
	"api-gateway/internal/types"

	oplogclient "oplog-service/oplogclient"
)

func ProtoToOpLogInfo(p *oplogclient.OpLogInfo) *types.OpLogInfo {
	if p == nil {
		return nil
	}
	return &types.OpLogInfo{
		Id:             p.Id,
		RequestId:      p.RequestId,
		UserId:         p.UserId,
		ReIp:           p.ReIp,
		ReTime:         p.ReTime,
		ReUa:           p.ReUa,
		ReUrl:          p.ReUrl,
		ReMethod:       p.ReMethod,
		ReContent:      p.ReContent,
		Success:        p.Success,
		BizCode:        p.BizCode,
		BizMsg:         p.BizMsg,
		ReResponseTime: p.ReResponseTime,
		ReUserAgent:    p.ReUserAgent,
	}
}

func ProtoToOpLogListData(d *oplogclient.OpLogListData) *types.OpLogListData {
	if d == nil {
		return nil
	}
	items := make([]types.OpLogInfo, 0, len(d.List))
	for _, item := range d.List {
		if info := ProtoToOpLogInfo(item); info != nil {
			items = append(items, *info)
		}
	}
	return &types.OpLogListData{
		List:     items,
		Total:    d.Total,
		Page:     d.Page,
		PageSize: d.PageSize,
	}
}
