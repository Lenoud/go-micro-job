package logic

import (
	"context"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/svc"
	sharedcommon "micro-shared/common"

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

	_, err := l.svcCtx.EntClient.Department.Create().
		SetTitle(in.GetTitle()).
		SetDescription(in.GetDescription()).
		SetParentID(int(in.GetParentId())).
		Save(l.ctx)
	if err != nil {
		msg := sharedcommon.Msg("创建", "部门")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("创建成功"), nil
}
