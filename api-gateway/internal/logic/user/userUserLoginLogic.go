package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/userClient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUserLoginLogic {
	return &UserUserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUserLoginLogic) UserUserLogin(req *types.UserLoginReq) (resp *types.UserLoginResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.UserLogin(l.ctx, &userclient.UserLoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return &types.UserLoginResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.UserLoginResp{
		BaseResp: common.RpcBase(rpcResp),
		Data:     common.ProtoToUserInfo(rpcResp.Data),
	}, nil
}
