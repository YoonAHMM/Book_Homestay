package homestaylogic

import (
	"context"

	"Book_Homestay/app/travel/cmd/rpc/internal/svc"
	"Book_Homestay/app/travel/cmd/rpc/pb"

	"Book_Homestay/common/errx"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type BusinessListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBusinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BusinessListLogic {
	return &BusinessListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BusinessListLogic) BusinessList(in *pb.BusinessListReq) (*pb.BusinessListResp, error) {
	
	Builder := l.svcCtx.HomestayModel.SelectBuilder().Where(squirrel.Eq{"homestay_business_id": in.Homestay_Business_Id})
	list, err := l.svcCtx.HomestayModel.FindPageListByIdDESC(l.ctx, Builder, in.LastId, in.PageSize)
	if err != nil {
		return nil, errors.Wrapf(errx.NewErrCode(errx.DB_ERROR), "err : %v , in : %+v", err, in)
	}

	var resp_list [] *pb.Homestay
	for _, homestay := range list {

		var Homestay  pb.Homestay
		_ = copier.Copy(&Homestay, homestay)

		Homestay.FoodPrice = homestay.FoodPrice
		Homestay.HomestayPrice = homestay.HomestayPrice
		Homestay.MarketHomestayPrice = homestay.MarketHomestayPrice

		resp_list = append(resp_list, &Homestay)
	}


	return &pb.BusinessListResp{
		Homestaylist: resp_list,
	}, nil
}
