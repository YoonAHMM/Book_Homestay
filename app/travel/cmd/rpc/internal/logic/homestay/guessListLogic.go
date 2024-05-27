package homestaylogic

import (
	"context"

	"Book_Homestay/app/travel/cmd/rpc/internal/svc"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	"Book_Homestay/common/Randx"
	"Book_Homestay/common/errx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GuessListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGuessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GuessListLogic {
	return &GuessListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GuessListLogic) GuessList(in *pb.GuessListReq) (*pb.GuessListResp, error) {
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
	
	var resp_list []*pb.Homestay

	for _, homestay := range list {

		var Homestay pb.Homestay
		_ = copier.Copy(&Homestay, homestay)

		Homestay.FoodPrice = homestay.FoodPrice
		Homestay.HomestayPrice = homestay.HomestayPrice
		Homestay.MarketHomestayPrice = homestay.MarketHomestayPrice

		resp_list = append(resp_list, &Homestay)
	}

	return &pb.GuessListResp{
		Homestaylist: resp_list,
	}, nil
}
