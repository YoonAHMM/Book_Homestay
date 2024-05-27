package main

import (
	"flag"
	"fmt"

	"Book_Homestay/app/travel/cmd/rpc/internal/config"
	homestayServer "Book_Homestay/app/travel/cmd/rpc/internal/server/homestay"
	homestaybussinessServer "Book_Homestay/app/travel/cmd/rpc/internal/server/homestaybussiness"
	"Book_Homestay/app/travel/cmd/rpc/internal/svc"
	"Book_Homestay/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/travel.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterHomestayServer(grpcServer, homestayServer.NewHomestayServer(ctx))
		pb.RegisterHomestaybussinessServer(grpcServer, homestaybussinessServer.NewHomestaybussinessServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
