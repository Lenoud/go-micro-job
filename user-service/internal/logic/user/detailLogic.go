package userlogic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetailLogic) Detail(in *user.IdReq) (*user.ApiResponse, error) {
	if in.Id == "" {
		return common.Fail("用户ID不能为空"), nil
	}

	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return common.Fail("查询用户信息失败"), nil
	}
	if u == nil {
		return common.Fail("用户不存在"), nil
	}
	return common.Success(u), nil
}
