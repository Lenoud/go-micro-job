package user

import (
	"context"
	"time"

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
		return &types.UserRegisterResp{BaseResp: types.BaseResp{Code: -1, Msg: "rpc调用失败", Timestamp: time.Now().UnixMilli()}}, nil
	}

	return &types.UserRegisterResp{
		BaseResp: types.BaseResp{Code: rpcResp.Code, Msg: rpcResp.Msg, Timestamp: time.Now().UnixMilli()},
	}, nil
}
