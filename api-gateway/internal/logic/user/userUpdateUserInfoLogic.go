package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/client/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateUserInfoLogic {
	return &UserUpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateUserInfoLogic) UserUpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.UpdateUserInfoResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.UpdateUserInfo(l.ctx, &userclient.UpdateUserInfoReq{
		Id:         req.Id,
		Nickname:   req.Nickname,
		Mobile:     req.Mobile,
		Email:      req.Email,
		PushEmail:  req.PushEmail,
		PushSwitch: req.PushSwitch,
	})
	if err != nil {
		return &types.UpdateUserInfoResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.UpdateUserInfoResp{
		BaseResp: common.RpcBase(rpcResp),
	}, nil
}
