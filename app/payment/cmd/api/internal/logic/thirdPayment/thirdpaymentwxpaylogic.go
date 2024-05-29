package thirdPayment

import (
	"context"

	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/payment/cmd/api/internal/svc"
	"Book_Homestay/app/payment/cmd/api/internal/types"
	"Book_Homestay/app/payment/cmd/rpc/payment"
	"Book_Homestay/app/payment/model"
	"Book_Homestay/app/user/cmd/rpc/user"
	"Book_Homestay/common/ctxdata"
	"Book_Homestay/common/errx"
	"Book_Homestay/common/vars"

	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdPaymentwxPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentwxPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentwxPayLogic {
	return &ThirdPaymentwxPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentwxPayLogic) ThirdPaymentwxPay(req *types.ThirdPaymentWxPayReq) (resp *types.ThirdPaymentWxPayResp, err error) {
	
	var totalPrice int64   // Total amount paid for current order(cent)
	var description string // Current Payment Description.

	switch req.ServiceType {
	case model.ThirdPaymentServiceTypeHomestayOrder:

		homestayTotalPrice, homestayDescription, err := l.getPayHomestayPriceDescription(req.OrderSn)
		if err != nil {
			return nil, err
		}
		totalPrice = homestayTotalPrice
		description = homestayDescription

	default:
		return nil, errors.Wrap(errx.NewErrMsg("Payment for this business type is not supported"),"")
	}

	
	wechatPrepayRsp, err := l.createWxPrePayOrder(req.ServiceType, req.OrderSn, totalPrice, description)
	if err != nil {
		return nil, err
	}

	return &types.ThirdPaymentWxPayResp{
		Appid:     l.svcCtx.Config.WxMiniConf.AppId,
		NonceStr:  *wechatPrepayRsp.NonceStr,
		PaySign:   *wechatPrepayRsp.PaySign,
		Package:   *wechatPrepayRsp.Package,
		Timestamp: *wechatPrepayRsp.TimeStamp,
		SignType:  *wechatPrepayRsp.SignType,
	}, nil
}

// wx下单
func (l *ThirdPaymentwxPayLogic) createWxPrePayOrder(serviceType, orderSn string, totalPrice int64, description string) (*jsapi.PrepayWithRequestPaymentResponse, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)

	userResp, err := l.svcCtx.UserRpc.GetUserAuthByUserId(l.ctx, &user.GetUserAuthByUserIdReq{
		UserId:   userId,
		AuthType: vars.UserAuthTypeSmallWX,
	})

	if err != nil {
		return nil, err
	}
	if userResp.UserAuth == nil || userResp.UserAuth.Id == 0 {
		return nil, model.ErrNotFound
	}
	openId := userResp.UserAuth.AuthKey

	
	createPaymentResp, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		UserId:      userId,
		PayModel:    model.ThirdPaymentPayModelWechatPay,
		PayTotal:    totalPrice,
		OrderSn:     orderSn,
		ServiceType: serviceType,
	})

	if err != nil || createPaymentResp.Sn == "" {
		return nil, errors.Wrapf(errx.NewErrCodeMsg(errx.WXMINIPAY_ERROR,"支付流水无法创建"),"err : %v , ordersn :%d",err,openId)
	}


	wxPayClient, err := svc.NewWxPayClientV3(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}

	jsApiSvc := jsapi.JsapiApiService{Client: wxPayClient}

	//wx支付下单
	resp, _, err := jsApiSvc.PrepayWithRequestPayment(l.ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(l.svcCtx.Config.WxMiniConf.AppId),
			Mchid:       core.String(l.svcCtx.Config.WxPayConf.MchId),
			Description: core.String(description),
			OutTradeNo:  core.String(createPaymentResp.Sn),
			Attach:      core.String(description),
			NotifyUrl:   core.String(l.svcCtx.Config.WxPayConf.NotifyUrl),
			Amount: &jsapi.Amount{
				Total: core.Int64(totalPrice),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(openId),
			},
		},
	)
	if err != nil {
		return nil,errors.Wrapf(errx.NewErrCode(errx.WXMINIPAY_ERROR),"支付下单错误，err :%v",err)
	}

	return resp, nil

}



func (l *ThirdPaymentwxPayLogic) getPayHomestayPriceDescription(orderSn string) (int64, string, error) {

	description := "homestay pay"

	// get user openid
	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{
		Sn: orderSn,
	})
	if err != nil {
		return 0, description, err
	}
	if resp.HomestayOrder == nil || resp.HomestayOrder.Id == 0 {
		return 0, description, err
	}

	return resp.HomestayOrder.OrderTotalPrice, description, nil
}
