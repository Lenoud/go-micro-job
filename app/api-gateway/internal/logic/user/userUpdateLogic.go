package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/userClient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UpdateUserReq) (resp *types.UpdateUserResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Update(l.ctx, &userclient.UpdateUserReq{
		Id:         req.Id,
		Username:   req.Username,
		Nickname:   req.Nickname,
		Mobile:     req.Mobile,
		Email:      req.Email,
		Role:       req.Role,
		Status:     req.Status,
		Password:   req.Password,
		PushEmail:  req.PushEmail,
		PushSwitch: req.PushSwitch,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc UserRpc.Update failed: %v", err)
		return &types.UpdateUserResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.UpdateUserResp{
		BaseResp: common.RpcBase(rpcResp),
	}, nil
}
