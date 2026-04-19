package common

import (
	"api-gateway/internal/types"
	userclient "user-service/userClient"
)

// ProtoToUserInfo 将 proto UserInfo 转为 API types.UserInfo
func ProtoToUserInfo(u *userclient.UserInfo) *types.UserInfo {
	if u == nil {
		return nil
	}
	return &types.UserInfo{
		Id:         u.Id,
		Username:   u.Username,
		Nickname:   u.Nickname,
		Mobile:     u.Mobile,
		Email:      u.Email,
		Role:       u.Role,
		Status:     u.Status,
		Token:      u.Token,
		CreateTime: u.CreateTime,
		PushEmail:  u.PushEmail,
		PushSwitch: u.PushSwitch,
		Avatar:     u.Avatar,
	}
}

// ProtoToUserListData 将 proto UserListData 转为 API types.UserListData
func ProtoToUserListData(d *userclient.UserListData) *types.UserListData {
	if d == nil {
		return nil
	}
	items := make([]types.UserInfo, 0, len(d.List))
	for _, u := range d.List {
		if info := ProtoToUserInfo(u); info != nil {
			items = append(items, *info)
		}
	}
	return &types.UserListData{
		List:     items,
		Total:    d.Total,
		Page:     d.Page,
		PageSize: d.PageSize,
	}
}
