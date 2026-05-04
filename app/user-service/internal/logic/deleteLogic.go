package logic

import (
	"context"

	sharedcommon "micro-shared/common"
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
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}

	ids := common.SplitIDs(in.Ids)
	if len(ids) == 0 {
		return common.FailAction("删除用户失败"), nil
	}
	if err := l.svcCtx.UserModel.DeleteBatch(l.ctx, ids); err != nil {
		msg := sharedcommon.Msg("删除", "用户")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("删除成功"), nil
}
