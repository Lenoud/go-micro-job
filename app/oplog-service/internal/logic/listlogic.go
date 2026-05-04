package logic

import (
	"context"

	sharedcommon "micro-shared/common"
	"oplog-service/internal/common"
	"oplog-service/internal/svc"
	"oplog-service/oplog"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *oplog.OpLogListReq) (*oplog.OpLogListResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailOpLogListForbidden("无权访问"), nil
	}

	page, pageSize := normalizePage(in.GetPage(), in.GetPageSize())
	total, err := l.svcCtx.OpLogModel.CountOpLogList(l.ctx)
	if err != nil {
		msg := sharedcommon.Msg("查询", "操作日志")
		common.LogErr(l.Logger, msg, err)
		return common.FailOpLogList(msg), nil
	}
	list, err := l.svcCtx.OpLogModel.FindOpLogList(l.ctx, page, pageSize)
	if err != nil {
		msg := sharedcommon.Msg("查询", "操作日志列表")
		common.LogErr(l.Logger, msg, err)
		return common.FailOpLogList(msg), nil
	}

	items := make([]*oplog.OpLogInfo, 0, len(list))
	for _, item := range list {
		items = append(items, common.OpLogModelToProto(item))
	}
	return common.OkOpLogList(&oplog.OpLogListData{
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}), nil
}
