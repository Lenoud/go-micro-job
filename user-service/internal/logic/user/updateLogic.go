package userlogic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员更新用户
func (l *UpdateLogic) Update(in *user.UpdateUserReq) (*user.ApiResponse, error) {
	if in.Id == "" {
		return common.Fail("用户ID不能为空"), nil
	}

	existing, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil || existing == nil {
		return common.Fail("用户不存在"), nil
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

	if err := l.svcCtx.UserModel.Update(l.ctx, existing); err != nil {
		return common.Fail("更新用户失败"), nil
	}
	return common.SuccessMsg("更新成功", nil), nil
}
