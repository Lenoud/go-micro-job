package userlogic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

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

// 批量删除用户
func (l *DeleteLogic) Delete(in *user.DeleteReq) (*user.ApiResponse, error) {
	ids := common.SplitIDs(in.Ids)
	if len(ids) == 0 {
		return common.Fail("删除用户失败"), nil
	}
	if err := l.svcCtx.UserModel.DeleteBatch(l.ctx, ids); err != nil {
		return common.Fail("删除用户失败"), nil
	}
	return common.SuccessMsg("删除成功", nil), nil
}
