package svc

import (
	"Book_Homestay/app/order/cmd/rpc/internal/config"
	"Book_Homestay/app/order/model"
	"Book_Homestay/app/travel/cmd/rpc/client/homestay"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	AsynqClient *asynq.Client

	Homestay_TravelRpc  homestay.HomestayZrpcClient

	HomestayOrderModel model.HomestayOrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		AsynqClient:asynq.NewClient(asynq.RedisClientOpt{Addr: c.Redis.Host, Password: c.Redis.Pass}),

		Homestay_TravelRpc: homestay.NewHomestayZrpcClient(zrpc.MustNewClient(c.TravelRpcConf)),

		HomestayOrderModel: model.NewHomestayOrderModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
