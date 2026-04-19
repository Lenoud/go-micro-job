package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/client/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailLogic) UserDetail(req *types.UserDetailReq) (resp *types.UserDetailResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Detail(l.ctx, &userclient.IdReq{
		Id: req.UserId,
	})
	if err != nil {
		return &types.UserDetailResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.UserDetailResp{
		BaseResp: common.RpcBase(rpcResp),
		Data:     common.ProtoToUserInfo(rpcResp.Data),
	}, nil
}
