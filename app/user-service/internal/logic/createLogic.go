package logic

import (
	"context"

	sharedcommon "micro-shared/common"
	"user-service/internal/common"
	"user-service/internal/model"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员创建用户
func (l *CreateLogic) Create(in *user.CreateUserReq) (*user.ActionResp, error) {
	if !common.HasRole(in.GetAuth(), common.RoleAdmin) {
		return common.FailActionForbidden("无权访问"), nil
	}

	if in.Username == "" || in.Password == "" {
		return common.FailAction("参数不完整"), nil
	}

	existing, err := l.svcCtx.UserModel.FindByUsername(l.ctx, in.Username)
	if err != nil {
		msg := sharedcommon.Msg("查询", "用户")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	if existing != nil {
		return common.FailActionDuplicate("用户名重复"), nil
	}

	md5Password := common.EncryptPassword(in.Password)

	role := in.Role
	if role == "" {
		role = common.RoleJobseeker
	}
	status := in.Status
	if status == "" {
		status = "0"
	}

	u := &model.User{
		Username: in.Username,
		Password: md5Password,
		Nickname: in.Nickname,
		Mobile:   in.Mobile,
		Email:    in.Email,
		Role:     role,
		Status:   status,
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, u)
	if err != nil {
		msg := sharedcommon.Msg("创建", "用户")
		common.LogErr(l.Logger, msg, err)
		return common.FailAction(msg), nil
	}
	return common.OkAction("创建成功"), nil
}
