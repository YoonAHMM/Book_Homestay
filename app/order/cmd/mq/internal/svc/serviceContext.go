package svc

import (
	"Book_Homestay/app/order/cmd/mq/internal/config"
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/user/cmd/rpc/user"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc      order.Order
	UserRpc user.UserZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UserRpc:      user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
