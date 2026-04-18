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

type UserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailLogic) UserDetail(req *types.UserDetailReq) (resp *types.ApiResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Detail(l.ctx, &userclient.IdReq{
		Id: req.UserId,
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
