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

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.UserListReq) (resp *types.UserListResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.List(l.ctx, &userclient.UserListReq{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return &types.UserListResp{BaseResp: types.BaseResp{Code: -1, Msg: "rpc调用失败", Timestamp: time.Now().UnixMilli()}}, nil
	}

	var data types.UserListData
	if rpcResp.Data != "" {
		_ = json.Unmarshal([]byte(rpcResp.Data), &data)
	}
	return &types.UserListResp{
		BaseResp: types.BaseResp{Code: rpcResp.Code, Msg: rpcResp.Msg, Timestamp: time.Now().UnixMilli()},
		Data:     &data,
	}, nil
}
