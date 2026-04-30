package logic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePwdLogic {
	return &UpdatePwdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改密码
func (l *UpdatePwdLogic) UpdatePwd(in *user.UpdatePwdReq) (*user.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleJobseeker, common.RoleRecruiter, common.RoleAdmin) {
		return common.FailAction("非法操作"), nil
	}

	targetUserID, allowed := common.ScopedMutationUserID(in.GetAuth(), in.GetUserId())
	if !allowed {
		return common.FailAction("非法操作"), nil
	}

	if in.OldPassword == "" || in.NewPassword == "" {
		return common.FailAction("参数不完整"), nil
	}

	u, err := l.svcCtx.UserModel.FindOne(l.ctx, targetUserID)
	if err != nil || u == nil {
		return common.FailAction("用户不存在"), nil
	}

	oldMd5Pwd := common.EncryptPassword(in.OldPassword)
	if u.Password != oldMd5Pwd {
		return common.FailAction("原密码错误"), nil
	}

	newHashedPwd := common.EncryptPassword(in.NewPassword)
	if err := l.svcCtx.UserModel.UpdatePassword(l.ctx, targetUserID, newHashedPwd); err != nil {
		return common.FailAction("修改密码失败"), nil
	}
	return common.OkAction("更新成功"), nil
}
