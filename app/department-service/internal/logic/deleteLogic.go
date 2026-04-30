package logic

import (
	"context"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogic) Delete(in *department.DeleteReq) (*department.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}
	if len(common.SplitIDs(in.GetIds())) == 0 {
		return common.FailAction("删除部门失败"), nil
	}

	if err := l.svcCtx.DepartmentModel.Delete(l.ctx, in.GetIds()); err != nil {
		return common.FailAction("删除部门失败"), nil
	}
	return common.OkAction("删除成功"), nil
}
