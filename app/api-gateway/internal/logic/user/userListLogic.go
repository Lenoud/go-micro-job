package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/userClient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.UserListReq) (resp *types.UserListResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.List(l.ctx, &userclient.UserListReq{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc UserRpc.List failed: %v", err)
		return &types.UserListResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.UserListResp{
		BaseResp: common.RpcBase(rpcResp),
		Data:     common.ProtoToUserListData(rpcResp.Data),
	}, nil
}
