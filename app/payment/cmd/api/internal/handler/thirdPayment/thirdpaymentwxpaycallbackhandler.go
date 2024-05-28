package thirdPayment

import (
	"net/http"

	"Book_Homestay/app/payment/cmd/api/internal/logic/thirdPayment"
	"Book_Homestay/app/payment/cmd/api/internal/svc"
	"Book_Homestay/app/payment/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ThirdPaymentWxPayCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdPaymentWxPayCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := thirdPayment.NewThirdPaymentWxPayCallbackLogic(r.Context(), svcCtx)
		resp, err := l.ThirdPaymentWxPayCallback(w,r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
