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

type HomestayDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailLogic {
	return &HomestayDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayDetailLogic) HomestayDetail(req *types.HomestayDetailReq) (resp *types.HomestayDetailResp, err error) {

	homestay , err := l.svcCtx.Homestay_TravelRpc.HomestayDetail(l.ctx,&pb.HomestayDetailReq{
		Id: req.Id,
	})
	
	if err != nil {
		return nil, err
	}

	var Homestay types.Homestay
	_ = copier.Copy(&Homestay, homestay)

	Homestay.FoodPrice = calculate.Fen2Yuan(homestay.Homestay.FoodPrice)
	Homestay.HomestayPrice = calculate.Fen2Yuan(homestay.Homestay.HomestayPrice)
	Homestay.MarketHomestayPrice = calculate.Fen2Yuan(homestay.Homestay.MarketHomestayPrice)
	
	return &types.HomestayDetailResp{
		Homestay: Homestay,
	}, nil
}
