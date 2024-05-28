package thirdPayment

import (
	"context"
	"net/http"

	"Book_Homestay/app/payment/cmd/api/internal/svc"
	"Book_Homestay/app/payment/cmd/api/internal/types"
	"Book_Homestay/app/payment/cmd/rpc/payment"
	"Book_Homestay/app/payment/model"
	"Book_Homestay/common/errx"
	"Book_Homestay/common/uniqueid"

	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdPaymentWxPayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentWxPayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentWxPayCallbackLogic {
	return &ThirdPaymentWxPayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentWxPayCallbackLogic) ThirdPaymentWxPayCallback(rw http.ResponseWriter, req *http.Request) (resp *types.ThirdPaymentWxPayCallbackResp, err error) {
	

	//验证私钥
	_, err = svc.NewWxPayClientV3(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}

	//下载证书通过sha256得到密钥构建处理器
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(l.svcCtx.Config.WxPayConf.MchId)//从服务器下载证书
	handler := notify.NewNotifyHandler(l.svcCtx.Config.WxPayConf.APIv3Key, verifiers.NewSHA256WithRSAVerifier(certVisitor))
	

	transaction := new(payments.Transaction)

	//解析通知
	_, err = handler.ParseNotifyRequest(context.Background(), req, transaction)
	if err != nil {
		return nil, errx.NewErrCode(errx.WXMINIPAYCALLBACK_ERROR,err.Error())
	}

	returnCode := "SUCCESS"
	err = l.verifyAndUpdateState(transaction)
	if err != nil {
		returnCode = "FAIL"
	}

	return &types.ThirdPaymentWxPayCallbackResp{
		ReturnCode: returnCode,
	}, err
}


func (l *ThirdPaymentWxPayCallbackLogic) verifyAndUpdateState(notifyTrasaction *payments.Transaction) error {

	paymentResp, err := l.svcCtx.PaymentRpc.GetPaymentBySn(l.ctx, &payment.GetPaymentBySnReq{
		Sn: *notifyTrasaction.OutTradeNo,
	})

	if err != nil || paymentResp.PaymentDetail.Id == 0 {
		return err
	}

	//比对金额
	notifyPayTotal := *notifyTrasaction.Amount.PayerTotal
	if paymentResp.PaymentDetail.PayTotal != notifyPayTotal {
		return errx.NewErrCode(errx.WXMINIPAYCALLBACK_ERROR,"Order amount exception ")
	}


	payStatus := l.getPayStatusByWXPayTradeState(*notifyTrasaction.TradeState)

	if payStatus == model.ThirdPaymentPayTradeStateSuccess {

		if paymentResp.PaymentDetail.PayStatus != model.ThirdPaymentPayTradeStateWait {
			return nil
		}

		if _, err = l.svcCtx.PaymentRpc.UpdateTradeState(l.ctx, &payment.UpdateTradeStateReq{
			Sn:             *notifyTrasaction.OutTradeNo,
			TradeState:     *notifyTrasaction.TradeState,
			TransactionId:  *notifyTrasaction.TransactionId,
			TradeType:      *notifyTrasaction.TradeType,
			TradeStateDesc: *notifyTrasaction.TradeStateDesc,
			PayStatus:      l.getPayStatusByWXPayTradeState(*notifyTrasaction.TradeState),
		}); err != nil {
			return err
		}

	} else if payStatus == model.ThirdPaymentPayTradeStateWait {
	}

	return nil
}

func (l *ThirdPaymentWxPayCallbackLogic) getPayStatusByWXPayTradeState(wxPayTradeState string) int64 {

	switch wxPayTradeState {
	case uniqueid.SUCCESS: //支付成功
		return model.ThirdPaymentPayTradeStateSuccess
	case uniqueid.USERPAYING: //支付中
		return model.ThirdPaymentPayTradeStateWait
	case uniqueid.REFUND: //已退款
		return model.ThirdPaymentPayTradeStateWait
	default:
		return model.ThirdPaymentPayTradeStateFAIL
	}
}



