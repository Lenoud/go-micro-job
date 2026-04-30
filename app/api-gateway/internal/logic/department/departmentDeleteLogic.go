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

type DepartmentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 管理员-删除部门
func NewDepartmentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentDeleteLogic {
	return &DepartmentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentDeleteLogic) DepartmentDelete(req *types.DeleteReq) (resp *types.DepartmentDeleteResp, err error) {
	auth, ok := common.DepartmentAuthFromContext(l.ctx)
	if !ok {
		return &types.DepartmentDeleteResp{BaseResp: common.FailBaseForbidden("无权访问")}, nil
	}

	rpcResp, err := l.svcCtx.DepartmentRpc.Delete(l.ctx, &departmentclient.DeleteReq{
		Ids:  req.Ids,
		Auth: auth,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc DepartmentRpc.Delete failed: %v", err)
		return &types.DepartmentDeleteResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.DepartmentDeleteResp{BaseResp: common.RpcBase(rpcResp)}, nil
}
