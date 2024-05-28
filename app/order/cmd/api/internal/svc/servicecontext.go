package svc

import (
	"Book_Homestay/app/order/cmd/api/internal/config"
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/payment/cmd/rpc/payment"
	"Book_Homestay/app/travel/cmd/rpc/client/homestay"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc   order.Order
	PaymentRpc payment.Payment
	Homestay_TravelRpc  homestay.HomestayZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:   order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		PaymentRpc: payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		Homestay_TravelRpc: homestay.NewHomestayZrpcClient(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
