package homestaybussinesslogic

import (
	"context"

	"Book_Homestay/app/travel/cmd/rpc/internal/svc"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	"Book_Homestay/app/travel/model"
	"Book_Homestay/app/user/cmd/rpc/user"
	"Book_Homestay/common/errx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestaybussinessdetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestaybussinessdetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestaybussinessdetailLogic {
	return &HomestaybussinessdetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HomestaybussinessdetailLogic) Homestaybussinessdetail(in *pb.BussinessReq) (*pb.BussinessResp, error) {
	homestayBusiness, err := l.svcCtx.HomestayBusinessModel.FindOne(l.ctx,in.Id)
	if err != nil && err != model.ErrNotFound {
		return nil,errx.NewErrCode(errx.DB_ERROR,err.Error())
	}

	var HomestayBusinessBoss pb.HomestayBusinessBoss
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


	return &pb.BussinessResp{
		Boss: &HomestayBusinessBoss,
	}, nil
}
