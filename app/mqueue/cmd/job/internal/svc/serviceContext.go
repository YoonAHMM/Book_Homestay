package svc

import (
	"github.com/hibiken/asynq"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/zeromicro/go-zero/zrpc"
	"Book_Homestay/app/mqueue/cmd/job/internal/config"
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/user/cmd/rpc/user"
)

type ServiceContext struct {
	Config config.Config
	AsynqServer *asynq.Server
	MiniProgram *miniprogram.MiniProgram

	OrderRpc order.Order
	UsercenterRpc user.UserZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		AsynqServer:newAsynqServer(c),
		MiniProgram:newMiniprogramClient(c),
		OrderRpc:order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc:user.NewUserZrpcClient(zrpc.MustNewClient(c.UsercenterRpcConf)),
	}
}



