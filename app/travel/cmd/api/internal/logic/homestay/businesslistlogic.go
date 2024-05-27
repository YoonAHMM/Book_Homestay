package homestay

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/common/calculate"
	"Book_Homestay/common/errx"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
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
	
	Builder := l.svcCtx.HomestayModel.SelectBuilder().Where(squirrel.Eq{"homestay_business_id": req.HomestayBusinessId})
	list, err := l.svcCtx.HomestayModel.FindPageListByIdDESC(l.ctx, Builder, req.LastId, req.PageSize)
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

	return &types.BusinessListResp{
		List: resp_list,
	}, nil
}
