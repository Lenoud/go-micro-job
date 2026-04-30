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

type DepartmentUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 管理员-更新部门
func NewDepartmentUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentUpdateLogic {
	return &DepartmentUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentUpdateLogic) DepartmentUpdate(req *types.UpdateDepartmentReq) (resp *types.DepartmentUpdateResp, err error) {
	auth, ok := common.DepartmentAuthFromContext(l.ctx)
	if !ok {
		return &types.DepartmentUpdateResp{BaseResp: common.FailBaseForbidden("无权访问")}, nil
	}

	rpcResp, err := l.svcCtx.DepartmentRpc.Update(l.ctx, &departmentclient.UpdateDepartmentReq{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		ParentId:    req.ParentId,
		Auth:        auth,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc DepartmentRpc.Update failed: %v", err)
		return &types.DepartmentUpdateResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.DepartmentUpdateResp{BaseResp: common.RpcBase(rpcResp)}, nil
}
