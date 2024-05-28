package kq

import (
	"context"
	"encoding/json"
	"Book_Homestay/app/order/cmd/mq/internal/svc"
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/order/model"
	paymentModel "Book_Homestay/app/payment/model"
	"Book_Homestay/common/kqueue"
	"Book_Homestay/common/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

/**
	Listening to the payment flow status change notification message queue
*/
type PaymentUpdateStatusMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentUpdateStatusMq(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentUpdateStatusMq {
	return &PaymentUpdateStatusMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//消费
func (l *PaymentUpdateStatusMq) Consume(_, val string) error {

	var message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage

	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Error("PaymentUpdateStatusMq->Consume Unmarshal err : %v , val : %s", err, val)
		return err
	}

	if err := l.execService(message); err != nil {
		logx.WithContext(l.ctx).Error("PaymentUpdateStatusMq->execService  err : %v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *PaymentUpdateStatusMq) execService(message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage) error {

	orderTradeState := l.getOrderTradeStateByPaymentTradeState(message.PayStatus)
	if orderTradeState != -99 {
		//update homestay order state
		_, err := l.svcCtx.OrderRpc.UpdateHomestayOrderTradeState(l.ctx, &order.UpdateHomestayOrderTradeStateReq{
			Sn:         message.OrderSn,
			TradeState: orderTradeState,
		})
		if err != nil {
			return errx.NewErrCode(errx.MQ_ERROR,err.Error())
		}
	}

	return nil
}

//将支付提供商的状态代码转换为内部民宿订单状态
func (l *PaymentUpdateStatusMq) getOrderTradeStateByPaymentTradeState(thirdPaymentPayStatus int64) int64 {

	switch thirdPaymentPayStatus {
	case paymentModel.ThirdPaymentPayTradeStateSuccess:
		return model.HomestayOrderTradeStateWaitUse
	case paymentModel.ThirdPaymentPayTradeStateRefund:
		return model.HomestayOrderTradeStateRefund
	default:
		return -99
	}

}
