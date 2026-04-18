package userlogic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户更新自己的信息
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoReq) (*user.ApiResponse, error) {
	if in.Id == "" {
		return common.Fail("用户ID不能为空"), nil
	}

	existing, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil || existing == nil {
		return common.Fail("用户不存在"), nil
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
		return common.Fail("更新用户信息失败"), nil
	}
	return common.SuccessMsg("更新成功", nil), nil
}
