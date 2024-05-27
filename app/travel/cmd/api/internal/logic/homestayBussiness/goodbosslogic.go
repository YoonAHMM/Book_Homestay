package homestayBussiness

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/app/travel/cmd/rpc/pb"

	
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	
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
	goodbossReq,err:=l.svcCtx.Bussiness_TravelRpc.Goodboss(l.ctx,&pb.GoodbossReq{	
	})

	if err!=nil {
		return nil,err
	}

	var resp_list [] types.HomestayBusinessBoss

	for _,goodboss:= range goodbossReq.Bosslist{

		var boss types.HomestayBusinessBoss
		_ = copier.Copy(&boss, goodboss)

		resp_list = append(resp_list, boss)
	}
	
	return &types.GoodBossResp{
		List: resp_list,
	}, nil
}
