package homestay

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/app/travel/cmd/rpc/pb"


	"Book_Homestay/common/calculate"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

)

type HomestayListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomestayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayListLogic {
	return &HomestayListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayListLogic) HomestayList(req *types.HomestayListReq) (resp *types.HomestayListResp, err error) {

	homestaylistResp , err := l.svcCtx.Homestay_TravelRpc.HomestayList(l.ctx,&pb.HomestayListReq{
		Page: req.Page,
		PageSize: req.PageSize,
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

	return &types.HomestayListResp{
		List: resp_list,
	}, nil

}
