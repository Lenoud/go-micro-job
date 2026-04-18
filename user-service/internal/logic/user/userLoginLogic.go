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

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 前台用户登录
func (l *UserLoginLogic) UserLogin(in *user.UserLoginReq) (*user.ApiResponse, error) {
	if in.Username == "" || in.Password == "" {
		return common.Fail("用户名或密码不能为空"), nil
	}

	password := common.EncryptPassword(in.Password)
	u, err := l.svcCtx.UserModel.FindNormalUser(l.ctx, in.Username, password)
	if err != nil {
		return common.Fail("查询用户失败"), nil
	}
	if u == nil || (u.Role != common.RoleJobseeker && u.Role != common.RoleRecruiter) {
		return common.Fail("用户名或密码错误"), nil
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
		return common.Fail("生成token失败"), nil
	}

	respData := map[string]interface{}{
		"id":         u.Id,
		"username":   u.Username,
		"nickname":   u.Nickname,
		"mobile":     u.Mobile,
		"email":      u.Email,
		"role":       u.Role,
		"status":     u.Status,
		"token":      tokenString,
		"createTime": u.CreateTime,
		"pushEmail":  u.PushEmail,
		"pushSwitch": u.PushSwitch,
		"avatar":     "",
	}
	return common.SuccessMsg("查询成功", respData), nil
}
