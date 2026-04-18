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

func (l *UserUpdateLogic) UserUpdate(req *types.UpdateUserReq) (resp *types.ApiResponse, err error) {
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
