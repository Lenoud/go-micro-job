package common

import (
	"user-service/internal/model"
	"user-service/user"
)

// UserModelToProto 将数据库 User 模型转为 proto UserInfo
func UserModelToProto(u *model.User) *user.UserInfo {
	if u == nil {
		return nil
	}
	createTime := ""
	if u.CreateTime.Valid {
		createTime = u.CreateTime.Time.Format("2006-01-02 15:04:05")
	}
	return &user.UserInfo{
		Id:         u.Id,
		Username:   u.Username,
		Password:   u.Password,
		Nickname:   u.Nickname,
		Mobile:     u.Mobile,
		Email:      u.Email,
		Role:       u.Role,
		Status:     u.Status,
		CreateTime: createTime,
	}
}
