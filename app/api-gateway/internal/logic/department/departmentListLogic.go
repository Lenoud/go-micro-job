// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package department

import (
	"context"

	"api-gateway/internal/common"
	"api-gateway/internal/svc"
	"api-gateway/internal/types"

	departmentclient "department-service/departmentClient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 部门列表
func NewDepartmentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentListLogic {
	return &DepartmentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentListLogic) DepartmentList(req *types.DepartmentListReq) (resp *types.DepartmentListResp, err error) {
	auth, ok := common.DepartmentAuthFromContext(l.ctx)
	if !ok {
		return &types.DepartmentListResp{BaseResp: common.FailBaseForbidden("无权访问")}, nil
	}

	rpcResp, err := l.svcCtx.DepartmentRpc.List(l.ctx, &departmentclient.DepartmentListReq{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
		Auth:     auth,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc DepartmentRpc.List failed: %v", err)
		return &types.DepartmentListResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.DepartmentListResp{
		BaseResp: common.RpcBase(rpcResp),
		Data:     common.ProtoToDepartmentListData(rpcResp.Data),
	}, nil
}
