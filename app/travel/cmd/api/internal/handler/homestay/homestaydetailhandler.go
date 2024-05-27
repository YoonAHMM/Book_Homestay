package homestay

import (
	"net/http"

	"Book_Homestay/app/travel/cmd/api/internal/logic/homestay"
	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func HomestayDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HomestayDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestay.NewHomestayDetailLogic(r.Context(), svcCtx)
		resp, err := l.HomestayDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
