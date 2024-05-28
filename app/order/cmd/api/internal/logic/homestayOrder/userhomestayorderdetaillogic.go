package homestayOrder

import (
	"context"

	"Book_Homestay/app/order/cmd/api/internal/svc"
	"Book_Homestay/app/order/cmd/api/internal/types"
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/order/model"
	"Book_Homestay/app/payment/cmd/rpc/payment"
	"Book_Homestay/common/calculate"
	"Book_Homestay/common/ctxdata"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserHomestayOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderDetailLogic {
	return &UserHomestayOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserHomestayOrderDetailLogic) UserHomestayOrderDetail(req *types.UserHomestayOrderDetailReq) (resp *types.UserHomestayOrderDetailResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)

	orderresp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{
		Sn: req.Sn,
	})
	if err != nil {
		return nil, err
	}

	var typesOrderDetail types.UserHomestayOrderDetailResp
	if orderresp.HomestayOrder != nil && orderresp.HomestayOrder.UserId == userId {

		copier.Copy(&typesOrderDetail, orderresp.HomestayOrder)

		//重置价格.
		typesOrderDetail.OrderTotalPrice = calculate.Fen2Yuan(orderresp.HomestayOrder.OrderTotalPrice)
		typesOrderDetail.FoodTotalPrice = calculate.Fen2Yuan(orderresp.HomestayOrder.FoodTotalPrice)
		typesOrderDetail.HomestayTotalPrice = calculate.Fen2Yuan(orderresp.HomestayOrder.HomestayTotalPrice)
		typesOrderDetail.HomestayPrice = calculate.Fen2Yuan(orderresp.HomestayOrder.HomestayPrice)
		typesOrderDetail.FoodPrice = calculate.Fen2Yuan(orderresp.HomestayOrder.FoodPrice)
		typesOrderDetail.MarketHomestayPrice = calculate.Fen2Yuan(orderresp.HomestayOrder.MarketHomestayPrice)

		//支付信息.
		if typesOrderDetail.TradeState != model.HomestayOrderTradeStateCancel && typesOrderDetail.TradeState != model.HomestayOrderTradeStateWaitPay {
			paymentResp, err := l.svcCtx.PaymentRpc.GetPaymentSuccessRefundByOrderSn(l.ctx, &payment.GetPaymentSuccessRefundByOrderSnReq{
				OrderSn: orderresp.HomestayOrder.Sn,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("Failed to get order payment information err : %v , orderSn:%s", err, orderresp.HomestayOrder.Sn)
			}

			if paymentResp.PaymentDetail != nil {
				typesOrderDetail.PayTime = paymentResp.PaymentDetail.PayTime
				typesOrderDetail.PayType = paymentResp.PaymentDetail.PayMode
			}
		}

		return &typesOrderDetail, nil
	}

	return nil, nil
}
