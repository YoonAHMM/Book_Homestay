package homestay

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/common/calculate"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"Book_Homestay/app/travel/cmd/rpc/pb"
)

type BusinessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBusinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BusinessListLogic {
	return &BusinessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BusinessListLogic) BusinessList(req *types.BusinessListReq) (resp *types.BusinessListResp, err error) {
		homestaylistResp , err := l.svcCtx.Homestay_TravelRpc.BusinessList(l.ctx,&pb.BusinessListReq{
			LastId: req.LastId,
			PageSize: req.PageSize,
			Homestay_Business_Id: req.HomestayBusinessId,
		})

	if err != nil {
		return nil, err
	}

	var resp_list [] types.Homestay

	for _,homestay:= range homestaylistResp.Homestaylist{
		var Homestay types.Homestay
		_ = copier.Copy(&Homestay, homestay)

		Homestay.FoodPrice = calculate.Fen2Yuan(homestay.FoodPrice)
		Homestay.HomestayPrice = calculate.Fen2Yuan(homestay.HomestayPrice)
		Homestay.MarketHomestayPrice = calculate.Fen2Yuan(homestay.MarketHomestayPrice)

		resp_list = append(resp_list, Homestay)
	}
	
	return &types.BusinessListResp{
		List: resp_list,
	}, nil
}
