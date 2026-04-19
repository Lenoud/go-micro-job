package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/client/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (resp *types.UserRegisterResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Register(l.ctx, &userclient.RegisterReq{
		Username:   req.Username,
		Password:   req.Password,
		RePassword: req.RePassword,
		Nickname:   req.Nickname,
		Mobile:     req.Mobile,
		Email:      req.Email,
	})
	if err != nil {
		return &types.UserRegisterResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.UserRegisterResp{
		BaseResp: common.RpcBase(rpcResp),
	}, nil
}
