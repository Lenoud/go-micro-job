package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/userClient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdatePwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdatePwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdatePwdLogic {
	return &UserUpdatePwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdatePwdLogic) UserUpdatePwd(req *types.UpdatePwdReq) (resp *types.UpdatePwdResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.UpdatePwd(l.ctx, &userclient.UpdatePwdReq{
		UserId:      req.UserId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return &types.UpdatePwdResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.UpdatePwdResp{
		BaseResp: common.RpcBase(rpcResp),
	}, nil
}
