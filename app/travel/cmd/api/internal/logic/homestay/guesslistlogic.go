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

type GuessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGuessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GuessListLogic {
	return &GuessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GuessListLogic) GuessList(req *types.GuessListReq) (resp *types.GuessListResp, err error) {

	guesslistResp , err := l.svcCtx.Homestay_TravelRpc.GuessList(l.ctx,&pb.GuessListReq{
	})

	if err != nil {
		return nil, err
	}

	var resp_list [] types.Homestay

	for _,homestay:= range guesslistResp.Homestaylist{
		var Homestay types.Homestay
		_ = copier.Copy(&Homestay, homestay)

		Homestay.FoodPrice = calculate.Fen2Yuan(homestay.FoodPrice)
		Homestay.HomestayPrice = calculate.Fen2Yuan(homestay.HomestayPrice)
		Homestay.MarketHomestayPrice = calculate.Fen2Yuan(homestay.MarketHomestayPrice)

		resp_list = append(resp_list, Homestay)
	}
	
	
	return &types.GuessListResp{
		List: resp_list,
	}, nil

}
