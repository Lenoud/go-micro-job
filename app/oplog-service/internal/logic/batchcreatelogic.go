package logic

import (
	"context"
	"time"

	"oplog-service/internal/common"
	"oplog-service/internal/model"
	"oplog-service/internal/svc"
	"oplog-service/oplog"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateLogic {
	return &BatchCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchCreate 批量写入操作日志。
// 安全说明：此接口无 OpLogContext 鉴权，仅由 api-gateway 在内网 gRPC 环境中调用，
// 不对外暴露。如需在非信任网络部署，应增加认证机制。
func (l *BatchCreateLogic) BatchCreate(in *oplog.BatchCreateReq) (*oplog.ActionResp, error) {
	if len(in.GetLogs()) == 0 {
		return common.OkAction("操作成功"), nil
	}

	logs := make([]*model.OpLog, 0, len(in.GetLogs()))
	now := time.Now().UnixMilli()
	for _, item := range in.GetLogs() {
		if item == nil {
			continue
		}
		success := item.GetSuccess()
		if success == "" {
			success = "1"
		}
		reTime := item.GetReTime()
		if reTime <= 0 {
			reTime = now
		}
		accessTime := item.GetAccessTime()
		if accessTime <= 0 {
			accessTime = now
		}
		logs = append(logs, &model.OpLog{
			RequestId:   item.GetRequestId(),
			UserId:      item.GetUserId(),
			ReIp:        item.GetReIp(),
			RequestTime: reTime,
			ReUa:        item.GetReUa(),
			ReUrl:       item.GetReUrl(),
			ReMethod:    item.GetReMethod(),
			ReContent:   item.GetReContent(),
			Success:     success,
			BizCode:     item.GetBizCode(),
			BizMsg:      item.GetBizMsg(),
			ResponseMs:  accessTime,
		})
	}
	if len(logs) == 0 {
		return common.OkAction("操作成功"), nil
	}
	if err := l.svcCtx.OpLogModel.BatchInsert(l.ctx, logs); err != nil {
		return common.FailAction("写入操作日志失败"), nil
	}
	return common.OkAction("操作成功"), nil
}
