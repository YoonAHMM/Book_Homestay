package logic

import (
	"context"
	"encoding/json"

	"Book_Homestay/app/mqueue/cmd/job/jobtype"
	"Book_Homestay/app/order/cmd/rpc/internal/svc"
	"Book_Homestay/app/order/cmd/rpc/pb"
	"Book_Homestay/app/order/model"
	"Book_Homestay/common/errx"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHomestayOrderTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateHomestayOrderTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHomestayOrderTradeStateLogic {
	return &UpdateHomestayOrderTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新民宿订单状态
func (l *UpdateHomestayOrderTradeStateLogic) UpdateHomestayOrderTradeState(in *pb.UpdateHomestayOrderTradeStateReq) (*pb.UpdateHomestayOrderTradeStateResp, error) {


	homestayOrder, err := l.svcCtx.HomestayOrderModel.FindOneBySn(l.ctx,in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	if homestayOrder == nil {
		return nil, model.ErrNotFound
	}

	if homestayOrder.TradeState == in.TradeState {
		return &pb.UpdateHomestayOrderTradeStateResp{}, nil
	}

	
	if err := l.verifyOrderTradeState(in.TradeState, homestayOrder.TradeState); err != nil {
		return nil, errors.WithMessagef(err, " , in : %+v", in)
	}

	
	homestayOrder.TradeState = in.TradeState
	if err := l.svcCtx.HomestayOrderModel.UpdateWithVersion(l.ctx,nil, homestayOrder); err != nil {
		return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
	}

	if in.TradeState == model.HomestayOrderTradeStateWaitUse {
		payload, err := json.Marshal(jobtype.PaySuccessNotifyUserPayload{Order: homestayOrder})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("pay success notify user task json Marshal fail, err :%+v , sn : %s",err,homestayOrder.Sn)
		}else{
			_, err := l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(jobtype.MsgPaySuccessNotifyUser, payload))
			if err != nil {
				logx.WithContext(l.ctx).Errorf("pay success notify user  insert queue fail err :%+v , sn : %s",err,homestayOrder.Sn)
			}
		}
	}


	return &pb.UpdateHomestayOrderTradeStateResp{
		Id:              homestayOrder.Id,
		UserId:          homestayOrder.UserId,
		Sn:              homestayOrder.Sn,
		TradeCode:       homestayOrder.TradeCode,
		Title:           homestayOrder.Title,
		LiveStartDate:   homestayOrder.LiveStartDate.Unix(),
		LiveEndDate:     homestayOrder.LiveEndDate.Unix(),
		OrderTotalPrice: homestayOrder.OrderTotalPrice,
	}, nil
}


// Update homestay order status
func (l *UpdateHomestayOrderTradeStateLogic) verifyOrderTradeState(newTradeState, oldTradeState int64) error {
	if newTradeState == model.HomestayOrderTradeStateWaitPay {
		return errx.NewErrCode(errx.ORDER_ERROR,"Changing this status is not supported")
	}

	if newTradeState == model.HomestayOrderTradeStateCancel {

		if oldTradeState != model.HomestayOrderTradeStateWaitPay {
			return  errx.NewErrCode(errx.ORDER_ERROR,"只有待支付的订单才能被取消")
		}

	} else if newTradeState == model.HomestayOrderTradeStateWaitUse {
		if oldTradeState != model.HomestayOrderTradeStateWaitPay {
			return errx.NewErrCode(errx.ORDER_ERROR,"Only orders pending payment can change this status")
		}

	} else if newTradeState == model.HomestayOrderTradeStateUsed {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errx.NewErrCode(errx.ORDER_ERROR,"Only unused orders can be changed to this status")
			
		}
	} else if newTradeState == model.HomestayOrderTradeStateRefund {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errx.NewErrCode(errx.ORDER_ERROR,"Only unused orders can be changed to this status")

		}
	} else if newTradeState == model.HomestayOrderTradeStateExpire {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errx.NewErrCode(errx.ORDER_ERROR,"Only unused orders can be changed to this status")
		}
	}

	return nil
}
