package logic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetailLogic) Detail(in *user.IdReq) (*user.UserInfoResp, error) {
	if in.Id == "" {
		return common.FailUserInfo("用户ID不能为空"), nil
	}

	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return common.FailUserInfo("查询用户信息失败"), nil
	}
	if u == nil {
		return common.FailUserInfo("用户不存在"), nil
	}
	return common.OkUserInfo("操作成功", common.UserModelToProto(u)), nil
}
