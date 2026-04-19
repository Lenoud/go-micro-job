package userlogic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除用户
func (l *DeleteLogic) Delete(in *user.DeleteReq) (*user.ActionResp, error) {
	ids := common.SplitIDs(in.Ids)
	if len(ids) == 0 {
		return common.FailAction("删除用户失败"), nil
	}
	if err := l.svcCtx.UserModel.DeleteBatch(l.ctx, ids); err != nil {
		return common.FailAction("删除用户失败"), nil
	}
	return common.OkAction("删除成功"), nil
}
