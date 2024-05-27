package homestay

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/app/travel/model"
	"Book_Homestay/common/errx"

	"Book_Homestay/common/calculate"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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
	Builder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(squirrel.Eq{
		"row_type":   model.HomestayActivityPreferredType,
		"row_status": model.HomestayActivityUpStatus,
	})
	homestayActivityList, err := l.svcCtx.HomestayActivityModel.FindPageListByPage(l.ctx, Builder, req.Page, req.PageSize, "data_id desc")

	if err != nil {
		return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
	}

	var resp_list []types.Homestay
	if len(homestayActivityList) > 0 { // mapreduce example
		mr.MapReduceVoid(func(source chan<- interface{}) {
			for _, homestayActivity := range homestayActivityList {
				source <- homestayActivity.DataId
			}
		}, func(item interface{}, writer mr.Writer[*model.Homestay], cancel func(error)) {
			id := item.(int64)

			homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, id)
			if err != nil && err != model.ErrNotFound {
				logx.WithContext(l.ctx).Errorf("ActivityHomestayListLogic ActivityHomestayList 获取活动数据失败 id : %d ,err : %v", id, err)
				return
			}
			writer.Write(homestay)
		}, func(pipe <-chan *model.Homestay, cancel func(error)) {

			for homestay := range pipe {
				var tyHomestay types.Homestay
				_ = copier.Copy(&tyHomestay, homestay)

				tyHomestay.FoodPrice = calculate.Fen2Yuan(homestay.FoodPrice)
				tyHomestay.HomestayPrice = calculate.Fen2Yuan(homestay.HomestayPrice)
				tyHomestay.MarketHomestayPrice = calculate.Fen2Yuan(homestay.MarketHomestayPrice)

				resp_list = append(resp_list, tyHomestay)
			}
		})
	}

	return &types.HomestayListResp{
		List: resp_list,
	}, nil

}
