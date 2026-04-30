package user

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/userClient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteLogic {
	return &UserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteLogic) UserDelete(req *types.DeleteReq) (resp *types.DeleteUserResp, err error) {
	auth, ok := common.AuthFromContext(l.ctx)
	if !ok {
		return &types.DeleteUserResp{BaseResp: common.FailBaseForbidden("无权访问")}, nil
	}

	rpcResp, err := l.svcCtx.UserRpc.Delete(l.ctx, &userclient.DeleteReq{
		Ids:  req.Ids,
		Auth: auth,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc UserRpc.Delete failed: %v", err)
		return &types.DeleteUserResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.DeleteUserResp{
		BaseResp: common.RpcBase(rpcResp),
	}, nil
}
