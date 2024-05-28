package svc

import (
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/payment/cmd/api/internal/config"
	"Book_Homestay/app/payment/cmd/rpc/payment"
	"Book_Homestay/app/user/cmd/rpc/user"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	WxPayClient  *core.Client

	PaymentRpc    payment.Payment
	OrderRpc      order.Order
	UserRpc       user.UserZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		PaymentRpc:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UserRpc:      user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
