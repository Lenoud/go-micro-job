package user

import (
	"context"
	"encoding/json"
	"time"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	userclient "user-service/client/user"

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
		return &types.UserLoginResp{BaseResp: types.BaseResp{Code: -1, Msg: "rpc调用失败", Timestamp: time.Now().UnixMilli()}}, nil
	}

	var data types.UserInfo
	if rpcResp.Data != "" {
		_ = json.Unmarshal([]byte(rpcResp.Data), &data)
	}
	return &types.UserLoginResp{
		BaseResp: types.BaseResp{Code: rpcResp.Code, Msg: rpcResp.Msg, Timestamp: time.Now().UnixMilli()},
		Data:     &data,
	}, nil
}
