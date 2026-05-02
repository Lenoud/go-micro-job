package logic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *user.UserListReq) (*user.UserListResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailUserListForbidden("无权访问"), nil
	}

	page := in.GetPage()
	pageSize := in.GetPageSize()
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	list, total, err := l.svcCtx.UserModel.FindList(l.ctx, in.GetKeyword(), page, pageSize)
	if err != nil {
		msg := sharedcommon.Msg("查询", "用户列表")
		common.LogErr(l.Logger, msg, err)
		return common.FailUserList(msg), nil
	}

	items := make([]*user.UserInfo, 0, len(list))
	for _, u := range list {
		items = append(items, common.UserModelToProto(u))
	}
	return common.OkUserList(&user.UserListData{
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}), nil
}
