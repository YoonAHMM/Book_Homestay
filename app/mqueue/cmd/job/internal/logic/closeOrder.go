package logic

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"Book_Homestay/app/mqueue/cmd/job/internal/svc"
	"Book_Homestay/app/mqueue/cmd/job/jobtype"
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/order/model"
	"Book_Homestay/common/errx"
)


var ErrCloseOrderFal = errx.NewErrMsg("close order fail")

// CloseHomestayOrderHandler close no pay homestayOrder
type CloseHomestayOrderHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCloseHomestayOrderHandler(svcCtx *svc.ServiceContext) *CloseHomestayOrderHandler {
	return &CloseHomestayOrderHandler{
		svcCtx:svcCtx,
	}
}

// defer  close no pay homestayOrder  : if return err != nil , asynq will retry
func (l *CloseHomestayOrderHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {

	var p jobtype.DeferCloseHomestayOrderPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return errors.Wrapf(ErrCloseOrderFal, "closeHomestayOrderStateMqHandler payload err:%v, payLoad:%+v", err, t.Payload())
	}

	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(ctx, &order.HomestayOrderDetailReq{
		Sn: p.Sn,
	})

	if err != nil || resp.HomestayOrder == nil {
		return errors.Wrapf(ErrCloseOrderFal, "closeHomestayOrderStateMqHandler  get order fail or order no exists err:%v, sn:%s ,HomestayOrder : %+v", err, p.Sn, resp.HomestayOrder)
	}
	
	if resp.HomestayOrder.TradeState == model.HomestayOrderTradeStateWaitPay {
		_, err := l.svcCtx.OrderRpc.UpdateHomestayOrderTradeState(ctx, &order.UpdateHomestayOrderTradeStateReq{
			Sn:         p.Sn,
			TradeState: model.HomestayOrderTradeStateCancel,
		})
		if err != nil {
			return errors.Wrapf(ErrCloseOrderFal, "CloseHomestayOrderHandler close order fail  err:%v, sn:%s ", err, p.Sn)
		}
	}

	return nil
}
