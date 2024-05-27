package homestay

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/common/Randx"
	"Book_Homestay/common/calculate"
	"Book_Homestay/common/errx"

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

	Builder := l.svcCtx.HomestayModel.SelectBuilder()
	count,err:=l.svcCtx.HomestayModel.FindCount(l.ctx,Builder,"id")
	if err != nil {
		return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
	}
	t,_:=Randx.GenerateRandomFixedRange(count,5)

	list, err := l.svcCtx.HomestayModel.FindPageListByIdDESC(l.ctx, l.svcCtx.HomestayModel.SelectBuilder(), t, 5)
	if err != nil {
		return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
	}
	
	var resp_list []types.Homestay
	for _, homestay := range list {

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
