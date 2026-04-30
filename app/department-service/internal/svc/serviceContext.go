package svc

import (
	"department-service/internal/config"
	"department-service/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	DepartmentModel model.DepartmentModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:          c,
		DepartmentModel: model.NewDepartmentModel(conn),
	}
}
