package logic

import (
	"context"
	"strings"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"
	sharedcommon "micro-shared/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员更新用户
func (l *UpdateLogic) Update(in *user.UpdateUserReq) (*user.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}

	if in.Id == "" {
		return common.FailAction("用户ID不能为空"), nil
	}

	existing, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		msg := sharedcommon.Msg("查询", "用户")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	if existing == nil {
		return common.FailAction("用户不存在"), nil
	}

	if in.Username != "" {
		existing.Username = in.Username
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
	if in.Role != "" {
		existing.Role = in.Role
	}
	if in.Status != "" {
		existing.Status = in.Status
	}
	if strings.TrimSpace(in.Password) != "" {
		if err := l.svcCtx.UserModel.UpdatePassword(l.ctx, existing.Id, common.EncryptPassword(in.Password)); err != nil {
			msg := sharedcommon.Msg("重置", "密码")
			common.LogErr(l.Logger, msg, err)
			return common.FailAction(msg), nil
		}
	}

	if err := l.svcCtx.UserModel.Update(l.ctx, existing); err != nil {
		msg := sharedcommon.Msg("更新", "用户")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("更新成功"), nil
}
