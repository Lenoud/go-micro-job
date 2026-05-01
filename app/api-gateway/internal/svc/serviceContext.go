package svc

import (
	"api-gateway/internal/config"

	departmentClient "department-service/departmentClient"
	oplogClient "oplog-service/oplogclient"
	userClient "user-service/userClient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	UserRpc       userClient.User
	DepartmentRpc departmentClient.Department
	OpLogRpc      oplogClient.OpLog
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UserRpc:       userClient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		DepartmentRpc: departmentClient.NewDepartment(zrpc.MustNewClient(c.DepartmentRpc)),
		OpLogRpc:      oplogClient.NewOpLog(zrpc.MustNewClient(c.OpLogRpc)),
	}
}
