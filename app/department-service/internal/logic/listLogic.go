package logic

import (
	"context"

	"department-service/department"
	"department-service/ent"
	entdep "department-service/ent/department"
	"department-service/internal/common"
	"department-service/internal/svc"
	sharedcommon "micro-shared/common"

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

	// 构建查询条件（Count 和 List 分别创建 query，避免复用）
	whereKeyword := in.GetKeyword() != ""

	// Count 查询
	countQuery := l.svcCtx.EntClient.Department.Query()
	if whereKeyword {
		countQuery = countQuery.Where(entdep.TitleContains(in.GetKeyword()))
	}
	total, err := countQuery.Count(l.ctx)
	if err != nil {
		msg := sharedcommon.Msg("查询", "部门列表")
		common.LogErr(l.Logger, msg, err)
		return common.FailDepartmentList(msg), nil
	}

	// List 查询（新建 query，不复用 countQuery）
	listQuery := l.svcCtx.EntClient.Department.Query()
	if whereKeyword {
		listQuery = listQuery.Where(entdep.TitleContains(in.GetKeyword()))
	}
	items, err := listQuery.
		Order(ent.Desc(entdep.FieldCreateTime)).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		All(l.ctx)
	if err != nil {
		msg := sharedcommon.Msg("查询", "部门列表")
		common.LogErr(l.Logger, msg, err)
		return common.FailDepartmentList(msg), nil
	}

	// 转换逻辑内联（对齐单体 api 做法已删除 converter.go）
	infos := make([]*department.DepartmentInfo, 0, len(items))
	for _, item := range items {
		createTime := ""
		if !item.CreateTime.IsZero() {
			createTime = item.CreateTime.Format("2006-01-02 15:04:05")
		}
		infos = append(infos, &department.DepartmentInfo{
			Id:          int64(item.ID),
			Title:       item.Title,
			Description: item.Description,
			ParentId:    int64(item.ParentID),
			CreateTime:  createTime,
		})
	}
	return common.OkDepartmentList(&department.DepartmentListData{
		List:     infos,
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
	}), nil
}
