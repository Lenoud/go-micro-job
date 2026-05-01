// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package oplog

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	oplogclient "oplog-service/oplogclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登录日志列表
func NewLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogListLogic {
	return &LoginLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogListLogic) LoginLogList(req *types.LoginLogListReq) (resp *types.LoginLogListResp, err error) {
	auth, ok := common.OpLogAuthFromContext(l.ctx)
	if !ok {
		return &types.LoginLogListResp{BaseResp: common.FailBaseForbidden("无权访问")}, nil
	}

	rpcResp, err := l.svcCtx.OpLogRpc.LoginList(l.ctx, &oplogclient.OpLogListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Auth:     auth,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc OpLogRpc.LoginList failed: %v", err)
		return &types.LoginLogListResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.LoginLogListResp{
		BaseResp: common.RpcBase(rpcResp),
		Data:     common.ProtoToOpLogListData(rpcResp.Data),
	}, nil
}
