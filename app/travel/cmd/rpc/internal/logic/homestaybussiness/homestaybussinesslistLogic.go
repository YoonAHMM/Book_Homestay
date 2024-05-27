package homestaybussinesslogic

import (
	"context"

	"Book_Homestay/app/travel/cmd/rpc/internal/svc"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	"Book_Homestay/common/errx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestaybussinesslistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestaybussinesslistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestaybussinesslistLogic {
	return &HomestaybussinesslistLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HomestaybussinesslistLogic) Homestaybussinesslist(in *pb.HomestaybussinesslistReq) (*pb.HomestaybussinesslistResp, error) {
	Builder := l.svcCtx.HomestayBusinessModel.SelectBuilder()
	list, err := l.svcCtx.HomestayBusinessModel.FindPageListByIdDESC(l.ctx, Builder, in.Lastid, in.Pagesize)
	if err != nil {
		return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
	}

	var resp_list []*pb.HomestayBusinessBoss
	if len(list) > 0 {
		for _, item := range list {
			var HomestayBusinessListInfo pb.HomestayBusinessBoss
			_ = copier.Copy(&HomestayBusinessListInfo, item)

			resp_list = append(resp_list, &HomestayBusinessListInfo)
		}
	}



	return &pb.HomestaybussinesslistResp{
		Bosslist: resp_list,
	}, nil
}
