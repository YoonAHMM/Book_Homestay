package logic

import (
	"context"
	"encoding/json"
	"time"

	"Book_Homestay/app/payment/cmd/rpc/internal/svc"
	"Book_Homestay/app/payment/cmd/rpc/pb"
	"Book_Homestay/app/payment/model"
	"Book_Homestay/common/errx"
	"Book_Homestay/common/kqueue"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateLogic {
	return &UpdateTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新交易状态
func (l *UpdateTradeStateLogic) UpdateTradeState(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {
	payment,err:=l.svcCtx.ThirdPaymentModel.FindOneBySn(l.ctx,in.Sn)
	if err!=nil{
		return nil,errors.Wrapf(errx.NewErrCode(errx.DB_ERROR),"err:%v",err)
	}

	if payment==nil{
		return nil,model.ErrNotFound
	}

	//无法修改为待支付状态，支付成功或失败状态只能由待支付状态转移，只有支付成功才能退款
	if in.PayStatus==model.ThirdPaymentPayTradeStateFAIL||in.PayStatus==model.ThirdPaymentPayTradeStateSuccess{
		if payment.PayStatus != model.ThirdPaymentPayTradeStateWait {
			return &pb.UpdateTradeStateResp{}, nil
		}
	}else if in.PayStatus == model.ThirdPaymentPayTradeStateRefund {
		if payment.PayStatus != model.ThirdPaymentPayTradeStateSuccess {
			return nil,errors.Wrapf(errx.NewErrMsg("Only orders with successful payment can be refunded"),"in :%+v",in)
		}
	} else {
		return nil,errors.Wrapf(errx.NewErrMsg("This status is not currently supported"),"in :%+v",in)
	}


	payment.TradeState = in.TradeState
	payment.TransactionId = in.TransactionId
	payment.TradeType = in.TradeType
	payment.TradeStateDesc = in.TradeStateDesc
	payment.PayStatus = in.PayStatus
	payment.PayTime = time.Unix(in.PayTime, 0)
	if err := l.svcCtx.ThirdPaymentModel.UpdateWithVersion(l.ctx,nil, payment); err != nil {
		return nil, errors.Wrapf(errx.NewErrCode(errx.DB_ERROR),"err:%v",err)
	}

	if err:=l.pubKqPaySuccess(in.Sn,in.PayStatus);err != nil{
		logx.WithContext(l.ctx).Errorf("l.pubKqPaySuccess : %+v",err)
	}

	return &pb.UpdateTradeStateResp{}, nil
}

func (l *UpdateTradeStateLogic) pubKqPaySuccess(orderSn string,payStatus int64) error{

	m := kqueue.ThirdPaymentUpdatePayStatusNotifyMessage{
		OrderSn:  orderSn ,
		PayStatus: payStatus,
	}

	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return  l.svcCtx.KqueuePaymentUpdatePayStatusClient.Push(string(body))
}
