package logic

import (
	"context"

	"department-service/department"
	"department-service/internal/common"
	"department-service/internal/svc"

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

func (l *ListLogic) List(in *department.DepartmentListReq) (*department.DepartmentListResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleRecruiter, common.RoleAdmin) {
		return common.FailDepartmentListForbidden("无权访问"), nil
	}

	page := in.GetPage()
	pageSize := in.GetPageSize()
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	list, total, err := l.svcCtx.DepartmentModel.FindList(l.ctx, in.GetKeyword(), page, pageSize)
	if err != nil {
		return common.FailDepartmentList("查询部门列表失败"), nil
	}

	items := make([]*department.DepartmentInfo, 0, len(list))
	for _, item := range list {
		items = append(items, common.DepartmentModelToProto(item))
	}
	return common.OkDepartmentList(&department.DepartmentListData{
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}), nil
}
