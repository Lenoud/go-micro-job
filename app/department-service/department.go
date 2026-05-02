package main

import (
	"flag"
	"fmt"

	"department-service/department"
	"department-service/internal/config"
	"department-service/internal/server"
	"department-service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/department.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		department.RegisterDepartmentServer(grpcServer, server.NewDepartmentServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	defer ctx.EntClient.Close()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
