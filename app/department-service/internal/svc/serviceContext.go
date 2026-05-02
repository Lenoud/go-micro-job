package svc

import (
	"department-service/ent"
	"department-service/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type ServiceContext struct {
	Config    config.Config
	EntClient *ent.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := ent.Open("mysql", c.MySQL.DataSource)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:    c,
		EntClient: client,
	}
}
