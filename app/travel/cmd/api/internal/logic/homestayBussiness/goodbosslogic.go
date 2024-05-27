package homestayBussiness

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/app/travel/model"
	"Book_Homestay/app/user/cmd/rpc/user"
	"Book_Homestay/common/errx"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GoodBossLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGoodBossLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodBossLogic {
	return &GoodBossLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoodBossLogic) GoodBoss(req *types.GoodBossReq) (resp *types.GoodBossResp, err error) {
	whereBuilder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(squirrel.Eq{
		"row_type":   model.HomestayActivityGoodBusiType,
		"row_status": model.HomestayActivityUpStatus,
	})

	homestayActivityList, err := l.svcCtx.HomestayActivityModel.FindPageListByPage(l.ctx, whereBuilder, 0, 10, "data_id desc")
	if err != nil {
		return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
	}

	var resp_list []types.HomestayBusinessBoss
	if len(homestayActivityList) > 0 {

		mr.MapReduceVoid(func(source chan<- interface{}) {
			for _, homestayActivity := range homestayActivityList {
				source <- homestayActivity.DataId
			}
		}, func(item interface{}, writer mr.Writer[*user.User], cancel func(error)) {
			id := item.(int64)

			userResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoReq{
				Id: id,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("GoodListLogic GoodList fail userId : %d ,err:%v", id, err)
				return
			}
			if userResp.User != nil && userResp.User.Id > 0 {
				writer.Write(userResp.User)
			}
		}, func(pipe <-chan *user.User, cancel func(error)) {

			for item := range pipe {
				var HomestayBusinessBoss types.HomestayBusinessBoss
				_ = copier.Copy(&HomestayBusinessBoss, item)

				// compute star todo
				resp_list = append(resp_list, HomestayBusinessBoss)
			}
		})
	}

	return &types.GoodBossResp{
		List: resp_list,
	}, nil
}
