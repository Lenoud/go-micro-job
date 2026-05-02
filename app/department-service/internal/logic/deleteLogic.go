package logic

import (
	"context"
	"fmt"

	"department-service/department"
	entdep "department-service/ent/department"
	"department-service/internal/common"
	"department-service/internal/svc"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogic) Delete(in *department.DeleteReq) (*department.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}
	idList := sharedcommon.SplitIDs(in.GetIds())
	if len(idList) == 0 {
		return common.FailAction("删除部门失败"), nil
	}

	// string ids → int ids
	intIDs := make([]int, 0, len(idList))
	for _, s := range idList {
		var n int
		if _, err := fmt.Sscanf(s, "%d", &n); err == nil && n > 0 {
			intIDs = append(intIDs, n)
		}
	}
	if len(intIDs) == 0 {
		return common.OkAction("删除成功"), nil
	}

	_, err := l.svcCtx.EntClient.Department.Delete().
		Where(entdep.IDIn(intIDs...)).
		Exec(l.ctx)
	if err != nil {
		msg := sharedcommon.Msg("删除", "部门")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("删除成功"), nil
}
