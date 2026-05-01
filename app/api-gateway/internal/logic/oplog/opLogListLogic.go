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

type OpLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 操作日志列表
func NewOpLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpLogListLogic {
	return &OpLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpLogListLogic) OpLogList(req *types.OpLogListReq) (resp *types.OpLogListResp, err error) {
	auth, ok := common.OpLogAuthFromContext(l.ctx)
	if !ok {
		return &types.OpLogListResp{BaseResp: common.FailBaseForbidden("无权访问")}, nil
	}

	rpcResp, err := l.svcCtx.OpLogRpc.List(l.ctx, &oplogclient.OpLogListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Auth:     auth,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc OpLogRpc.List failed: %v", err)
		return &types.OpLogListResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.OpLogListResp{
		BaseResp: common.RpcBase(rpcResp),
		Data:     common.ProtoToOpLogListData(rpcResp.Data),
	}, nil
}
