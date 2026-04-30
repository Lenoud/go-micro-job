package logic

import (
	"context"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/model"
	"department-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *department.CreateDepartmentReq) (*department.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}
	if in.GetTitle() == "" {
		return common.FailAction("部门名称不能为空"), nil
	}

	data := &model.Department{
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		ParentId:    in.GetParentId(),
	}
	if _, err := l.svcCtx.DepartmentModel.Insert(l.ctx, data); err != nil {
		return common.FailAction("创建部门失败"), nil
	}
	return common.OkAction("创建成功"), nil
}
