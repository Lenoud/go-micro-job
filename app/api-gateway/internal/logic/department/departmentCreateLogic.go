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

type DepartmentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 管理员-创建部门
func NewDepartmentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentCreateLogic {
	return &DepartmentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentCreateLogic) DepartmentCreate(req *types.CreateDepartmentReq) (resp *types.DepartmentCreateResp, err error) {
	auth, ok := common.DepartmentAuthFromContext(l.ctx)
	if !ok {
		return &types.DepartmentCreateResp{BaseResp: common.FailBaseForbidden("无权访问")}, nil
	}

	rpcResp, err := l.svcCtx.DepartmentRpc.Create(l.ctx, &departmentclient.CreateDepartmentReq{
		Title:       req.Title,
		Description: req.Description,
		ParentId:    req.ParentId,
		Auth:        auth,
	})
	if err != nil {
		logx.Errorf("[gateway] rpc DepartmentRpc.Create failed: %v", err)
		return &types.DepartmentCreateResp{BaseResp: common.FailBaseMsg("rpc调用失败")}, nil
	}

	return &types.DepartmentCreateResp{BaseResp: common.RpcBase(rpcResp)}, nil
}
