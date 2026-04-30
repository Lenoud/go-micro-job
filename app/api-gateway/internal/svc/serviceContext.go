package svc

import (
	"api-gateway/internal/config"

	departmentClient "department-service/departmentClient"
	userClient "user-service/userClient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	UserRpc       userClient.User
	DepartmentRpc departmentClient.Department
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UserRpc:       userClient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		DepartmentRpc: departmentClient.NewDepartment(zrpc.MustNewClient(c.DepartmentRpc)),
	}
}
