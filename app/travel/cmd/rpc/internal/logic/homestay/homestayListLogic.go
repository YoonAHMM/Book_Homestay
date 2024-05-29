package homestaylogic

import (
	"context"

	"Book_Homestay/app/travel/cmd/rpc/internal/svc"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	"Book_Homestay/app/travel/model"
	"Book_Homestay/common/errx"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type HomestayListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayListLogic {
	return &HomestayListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HomestayListLogic) HomestayList(in *pb.HomestayListReq) (*pb.HomestayListResp, error) {
	Builder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(squirrel.Eq{
		"row_type":   model.HomestayActivityPreferredType,
		"row_status": model.HomestayActivityUpStatus,
	})
	homestayActivityList, err := l.svcCtx.HomestayActivityModel.FindPageListByPage(l.ctx, Builder, in.Page, in.PageSize, "data_id desc")

	if err != nil {
		return nil, errors.Wrapf(errx.NewErrCode(errx.DB_ERROR), "err : %v , in : %+v", err, in)
	}

	var resp_list []*pb.Homestay
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
				var tyHomestay pb.Homestay
				_ = copier.Copy(&tyHomestay, homestay)

				tyHomestay.FoodPrice = homestay.FoodPrice
				tyHomestay.HomestayPrice = homestay.HomestayPrice
				tyHomestay.MarketHomestayPrice = homestay.MarketHomestayPrice

				resp_list = append(resp_list, &tyHomestay)
			}
		})
	}


	return &pb.HomestayListResp{
		Homestaylist: resp_list,
	}, nil
}
