package logic

import (
	"context"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/svc"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *department.UpdateDepartmentReq) (*department.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}
	if in.GetId() == 0 {
		return common.FailAction("部门不存在"), nil
	}

	existing, err := l.svcCtx.EntClient.Department.Get(l.ctx, int(in.GetId()))
	if err != nil || existing == nil {
		return common.FailAction("部门不存在"), nil
	}

	update := l.svcCtx.EntClient.Department.UpdateOneID(existing.ID)
	if in.GetTitle() != "" {
		update = update.SetTitle(in.GetTitle())
	}
	if in.GetDescription() != "" {
		update = update.SetDescription(in.GetDescription())
	}
	if in.GetParentId() != 0 {
		update = update.SetParentID(int(in.GetParentId()))
	}

	if err := update.Exec(l.ctx); err != nil {
		msg := sharedcommon.Msg("更新", "部门")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("更新成功"), nil
}
