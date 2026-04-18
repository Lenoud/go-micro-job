package userlogic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/model"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员创建用户
func (l *CreateLogic) Create(in *user.CreateUserReq) (*user.ApiResponse, error) {
	if in.Username == "" || in.Password == "" {
		return common.Fail("参数不完整"), nil
	}

	existing, err := l.svcCtx.UserModel.FindByUsername(l.ctx, in.Username)
	if err != nil {
		return common.Fail("查询用户失败"), nil
	}
	if existing != nil {
		return common.Fail("用户名重复"), nil
	}

	md5Password := common.EncryptPassword(in.Password)
	token := common.GenerateToken(in.Username)

	role := in.Role
	if role == "" {
		role = "1"
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
		Token:    token,
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, u)
	if err != nil {
		return common.Fail("创建用户失败"), nil
	}
	return common.SuccessMsg("创建成功", nil), nil
}
