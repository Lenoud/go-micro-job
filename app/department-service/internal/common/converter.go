package common

import (
	"department-service/department"
	"department-service/internal/model"
)

func DepartmentModelToProto(d *model.Department) *department.DepartmentInfo {
	if d == nil {
		return nil
	}
	createTime := ""
	if d.CreateTime.Valid {
		createTime = d.CreateTime.Time.Format("2006-01-02 15:04:05")
	}
	return &department.DepartmentInfo{
		Id:          d.Id,
		Title:       d.Title,
		Description: d.Description,
		ParentId:    d.ParentId,
		CreateTime:  createTime,
	}
}
