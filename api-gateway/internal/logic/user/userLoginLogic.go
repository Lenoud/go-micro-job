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

func (l *UserLoginLogic) UserLogin(req *types.LoginReq) (resp *types.ApiResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &userclient.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return &types.ApiResponse{Code: -1, Msg: "rpc调用失败", Timestamp: time.Now().UnixMilli()}, nil
	}

	var data interface{}
	if rpcResp.Data != "" {
		_ = json.Unmarshal([]byte(rpcResp.Data), &data)
	}
	return &types.ApiResponse{
		Code:      rpcResp.Code,
		Msg:       rpcResp.Msg,
		Data:      data,
		Timestamp: time.Now().UnixNano() / 1e6,
	}, nil
}
