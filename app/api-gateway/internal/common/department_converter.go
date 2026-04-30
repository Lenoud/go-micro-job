package common

import (
	"api-gateway/internal/types"

	departmentclient "department-service/departmentClient"
)

func ProtoToDepartmentInfo(d *departmentclient.DepartmentInfo) *types.DepartmentInfo {
	if d == nil {
		return nil
	}
	return &types.DepartmentInfo{
		Id:          d.Id,
		Title:       d.Title,
		Description: d.Description,
		ParentId:    d.ParentId,
		CreateTime:  d.CreateTime,
	}
}

func ProtoToDepartmentListData(d *departmentclient.DepartmentListData) *types.DepartmentListData {
	if d == nil {
		return nil
	}
	items := make([]types.DepartmentInfo, 0, len(d.List))
	for _, item := range d.List {
		if info := ProtoToDepartmentInfo(item); info != nil {
			items = append(items, *info)
		}
	}
	return &types.DepartmentListData{
		List:     items,
		Total:    d.Total,
		Page:     d.Page,
		PageSize: d.PageSize,
	}
}
