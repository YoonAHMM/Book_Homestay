package homestayBussiness

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	

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
	bossReq,err:=l.svcCtx.Bussiness_TravelRpc.Homestaybussinessdetail(l.ctx,&pb.BussinessReq{	
		Id: req.Id,
	})

	if err!=nil {
		return nil,err
	}


	var boss types.HomestayBusinessBoss
	_ = copier.Copy(&boss,bossReq)
	

	return &types.HomestayBussinessDetailResp{
		Boss: boss,
	}, nil
}
