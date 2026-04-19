package logic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/model"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户注册
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.ActionResp, error) {
	if in.Username == "" || in.Password == "" || in.RePassword == "" {
		return common.FailAction("参数不完整"), nil
	}

	existing, err := l.svcCtx.UserModel.FindByUsername(l.ctx, in.Username)
	if err != nil {
		return common.FailAction("查询用户失败"), nil
	}
	if existing != nil {
		return common.FailAction("用户名重复"), nil
	}

	if in.Password != in.RePassword {
		return common.FailAction("密码不一致"), nil
	}

	md5Password := common.EncryptPassword(in.Password)
	token := common.GenerateToken(in.Username)

	u := &model.User{
		Username: in.Username,
		Password: md5Password,
		Nickname: in.Nickname,
		Mobile:   in.Mobile,
		Email:    in.Email,
		Role:     "1",
		Status:   "0",
		Token:    token,
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, u)
	if err != nil {
		return common.FailAction("注册用户失败"), nil
	}
	return common.OkAction("创建成功"), nil
}
