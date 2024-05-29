package homestaybussinesslogic

import (
	"context"

	"Book_Homestay/app/travel/cmd/rpc/internal/svc"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	"Book_Homestay/app/travel/model"
	"Book_Homestay/app/user/cmd/rpc/user"
	"Book_Homestay/common/errx"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GoodbossLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGoodbossLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodbossLogic {
	return &GoodbossLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GoodbossLogic) Goodboss(in *pb.GoodbossReq) (*pb.GoodbossResp, error) {
	Builder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(squirrel.Eq{
		"row_type":   model.HomestayActivityGoodBusiType,
		"row_status": model.HomestayActivityUpStatus,
	})

	homestayActivityList, err := l.svcCtx.HomestayActivityModel.FindPageListByPage(l.ctx, Builder, 0, 10, "data_id desc")
	if err != nil {
		return nil, errors.Wrapf(errx.NewErrCode(errx.DB_ERROR), "err : %v , in : %+v", err, in)
	}

	var resp_list []*pb.HomestayBusinessBoss
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
				var HomestayBusinessBoss pb.HomestayBusinessBoss
				_ = copier.Copy(&HomestayBusinessBoss, item)
				resp_list = append(resp_list, &HomestayBusinessBoss)
			}
		})
	}


	return &pb.GoodbossResp{
		Bosslist: resp_list,
	}, nil
}
