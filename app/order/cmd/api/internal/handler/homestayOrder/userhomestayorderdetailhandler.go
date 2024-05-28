package homestayOrder

import (
	"net/http"

	"Book_Homestay/app/order/cmd/api/internal/logic/homestayOrder"
	"Book_Homestay/app/order/cmd/api/internal/svc"
	"Book_Homestay/app/order/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserHomestayOrderDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserHomestayOrderDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestayOrder.NewUserHomestayOrderDetailLogic(r.Context(), svcCtx)
		resp, err := l.UserHomestayOrderDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
