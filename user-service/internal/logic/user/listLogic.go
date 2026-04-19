package userlogic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

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
	page := in.Page
	pageSize := in.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	list, total, err := l.svcCtx.UserModel.FindList(l.ctx, in.Keyword, page, pageSize)
	if err != nil {
		return common.FailUserList("查询用户列表失败"), nil
	}

	items := make([]*user.UserInfo, 0, len(list))
	for _, u := range list {
		items = append(items, common.UserModelToProto(u))
	}
	return common.SuccessUserList(&user.UserListData{
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}), nil
}
