package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/userClient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginReq) (resp *types.LoginResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &userclient.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc UserRpc.Login failed: %v", err)
		return &types.LoginResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.LoginResp{
		BaseResp: common.RpcBase(rpcResp),
		Data:     common.ProtoToUserInfo(rpcResp.Data),
	}, nil
}
