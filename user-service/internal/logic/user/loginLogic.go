package userlogic

import (
	"context"
	"time"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 管理员登录
func (l *LoginLogic) Login(in *user.LoginReq) (*user.UserInfoResp, error) {
	if in.Username == "" || in.Password == "" {
		return common.FailUserInfo("用户名或密码不能为空"), nil
	}

	password := common.EncryptPassword(in.Password)
	u, err := l.svcCtx.UserModel.FindAdminUser(l.ctx, in.Username, password)
	if err != nil {
		return common.FailUserInfo("查询管理员失败"), nil
	}
	if u == nil || u.Role != common.RoleAdmin {
		return common.FailUserInfo("用户名或密码错误"), nil
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"userId":   u.Id,
		"username": u.Username,
		"role":     u.Role,
		"exp":      now.Add(time.Duration(l.svcCtx.Config.JWT.AccessExpire) * time.Second).Unix(),
		"iat":      now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(l.svcCtx.Config.JWT.AccessSecret))
	if err != nil {
		return common.FailUserInfo("生成token失败"), nil
	}

	info := common.UserModelToProto(u)
	info.Token = tokenString
	return common.OkUserInfo("查询成功", info), nil
}
