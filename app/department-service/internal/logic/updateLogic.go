package logic

import (
	"context"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/svc"

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
	if in.GetId() == "" {
		return common.FailAction("部门不存在"), nil
	}

	existing, err := l.svcCtx.DepartmentModel.FindOne(l.ctx, in.GetId())
	if err != nil || existing == nil {
		return common.FailAction("部门不存在"), nil
	}
	if in.GetTitle() != "" {
		existing.Title = in.GetTitle()
	}
	if in.GetDescription() != "" {
		existing.Description = in.GetDescription()
	}
	if in.GetParentId() != "" {
		existing.ParentId = in.GetParentId()
	}
	if err := l.svcCtx.DepartmentModel.Update(l.ctx, existing); err != nil {
		return common.FailAction("更新部门失败"), nil
	}
	return common.OkAction("更新成功"), nil
}
