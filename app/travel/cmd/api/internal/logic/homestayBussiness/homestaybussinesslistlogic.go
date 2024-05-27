package homestayBussiness

import (
	"context"

	"Book_Homestay/app/travel/cmd/api/internal/svc"
	"Book_Homestay/app/travel/cmd/api/internal/types"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBussinessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomestayBussinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBussinessListLogic {
	return &HomestayBussinessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayBussinessListLogic) HomestayBussinessList(req *types.HomestayBussinessListReq) (resp *types.HomestayBussinessListResp, err error) {
	bossReq,err:=l.svcCtx.Bussiness_TravelRpc.Homestaybussinesslist(l.ctx,&pb.HomestaybussinesslistReq{	
		Lastid: req.LastId,
		Pagesize: req.PageSize,
	})

	if err!=nil {
		return nil,err
	}

	var resp_list [] types.HomestayBusinessListInfo

	for _,Boss:= range bossReq.Bosslist{

		var boss types.HomestayBusinessListInfo
		_ = copier.Copy(&boss,Boss)

		resp_list = append(resp_list, boss)
	}
	
	return &types.HomestayBussinessListResp{
		List: resp_list,
	}, nil

}
