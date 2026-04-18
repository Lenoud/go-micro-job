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

func (l *UserUserLoginLogic) UserUserLogin(req *types.UserLoginReq) (resp *types.ApiResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpc.UserLogin(l.ctx, &userclient.UserLoginReq{
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
