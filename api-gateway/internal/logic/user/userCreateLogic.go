package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/client/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCreateLogic) UserCreate(req *types.CreateUserReq) (resp *types.CreateUserResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Create(l.ctx, &userclient.CreateUserReq{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Mobile:   req.Mobile,
		Email:    req.Email,
		Role:     req.Role,
		Status:   req.Status,
	})
	if err != nil {
		return &types.CreateUserResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.CreateUserResp{
		BaseResp: common.RpcBase(rpcResp),
	}, nil
}
