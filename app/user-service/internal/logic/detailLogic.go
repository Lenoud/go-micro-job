package logic

import (
	"context"

	"user-service/internal/common"
	"user-service/internal/svc"
	"user-service/user"
	sharedcommon "micro-shared/common"

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
	if !common.HasRole(in.GetAuth(), common.RoleAdmin, common.RoleRecruiter, common.RoleJobseeker) {
		return common.FailUserInfoForbidden("无权访问"), nil
	}

	targetUserID, ok := common.DetailTargetUserID(in.GetAuth(), in.GetId())
	if !ok {
		return common.FailUserInfoForbidden("无权访问"), nil
	}

	u, err := l.svcCtx.UserModel.FindOne(l.ctx, targetUserID)
	if err != nil {
		msg := sharedcommon.Msg("查询", "用户信息")
		common.LogErr(l.Logger, msg, err)
		return common.FailUserInfo(msg), nil
	}
	if u == nil {
		return common.FailUserInfo("用户不存在"), nil
	}
	return common.OkUserInfo("操作成功", common.UserModelToProto(u)), nil
}
