package logic

import (
	"context"

	sharedcommon "micro-shared/common"
	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户更新自己的信息
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoReq) (*user.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleJobseeker, common.RoleRecruiter, common.RoleAdmin) {
		return common.FailAction("非法操作"), nil
	}

	targetUserID, allowed := common.ScopedMutationUserID(in.GetAuth(), in.GetId())
	if !allowed {
		return common.FailAction("非法操作"), nil
	}

	existing, err := l.svcCtx.UserModel.FindOne(l.ctx, targetUserID)
	if err != nil {
		msg := sharedcommon.Msg("查询", "用户")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	if existing == nil {
		return common.FailAction("用户不存在"), nil
	}

	if in.Nickname != "" {
		existing.Nickname = in.Nickname
	}
	if in.Mobile != "" {
		existing.Mobile = in.Mobile
	}
	if in.Email != "" {
		existing.Email = in.Email
	}

	if err := l.svcCtx.UserModel.Update(l.ctx, existing); err != nil {
		msg := sharedcommon.Msg("更新", "用户信息")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("更新成功"), nil
}
