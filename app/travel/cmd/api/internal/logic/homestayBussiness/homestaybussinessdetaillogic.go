package homestayBussiness

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/app/travel/model"
	"Book_Homestay/app/user/cmd/rpc/user"
	"Book_Homestay/common/errx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBussinessDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomestayBussinessDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBussinessDetailLogic {
	return &HomestayBussinessDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayBussinessDetailLogic) HomestayBussinessDetail(req *types.HomestayBussinessDetailReq) (resp *types.HomestayBussinessDetailResp, err error) {
	homestayBusiness, err := l.svcCtx.HomestayBusinessModel.FindOne(l.ctx,req.Id)
	if err != nil && err != model.ErrNotFound {
		return nil,errx.NewErrCode(errx.DB_ERROR,err.Error())
	}

	var HomestayBusinessBoss types.HomestayBusinessBoss
	if homestayBusiness != nil {

		userResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoReq{
			Id: homestayBusiness.UserId,
		})
		if err != nil {
			return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
		}
		if userResp.User != nil && userResp.User.Id > 0 {
			_ = copier.Copy(&HomestayBusinessBoss, userResp.User)
		}
	}

	return &types.HomestayBussinessDetailResp{
		Boss: HomestayBusinessBoss,
	}, nil
}
