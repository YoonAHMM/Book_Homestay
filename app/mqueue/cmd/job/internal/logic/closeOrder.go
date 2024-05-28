package logic

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	
	"Book_Homestay/app/mqueue/cmd/job/internal/svc"
	"Book_Homestay/app/mqueue/cmd/job/jobtype"
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/order/model"
	"Book_Homestay/common/errx"
)


var ErrCloseOrderFal ="close order fail"

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
		return errx.NewErrCode(errx.MQ_ERROR,ErrCloseOrderFal)
	}

	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(ctx, &order.HomestayOrderDetailReq{
		Sn: p.Sn,
	})

	if err != nil || resp.HomestayOrder == nil {
		return errx.NewErrCode(errx.MQ_ERROR,ErrCloseOrderFal)
	}
	
	if resp.HomestayOrder.TradeState == model.HomestayOrderTradeStateWaitPay {
		_, err := l.svcCtx.OrderRpc.UpdateHomestayOrderTradeState(ctx, &order.UpdateHomestayOrderTradeStateReq{
			Sn:         p.Sn,
			TradeState: model.HomestayOrderTradeStateCancel,
		})
		if err != nil {
			return errx.NewErrCode(errx.MQ_ERROR,ErrCloseOrderFal)
		}
	}

	return nil
}
