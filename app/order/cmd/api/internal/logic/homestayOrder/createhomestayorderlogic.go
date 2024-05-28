package homestayOrder

import (
	"context"

	"Book_Homestay/app/order/cmd/api/internal/svc"
	"Book_Homestay/app/order/cmd/api/internal/types"
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/app/travel/cmd/rpc/pb"
	"Book_Homestay/app/travel/model"
	"Book_Homestay/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateHomestayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateHomestayOrderLogic) CreateHomestayOrder(req *types.CreateHomestayOrderReq) (resp *types.CreateHomestayOrderResp, err error) {
	homestayResp , err:=l.svcCtx.Homestay_TravelRpc.HomestayDetail(l.ctx,&pb.HomestayDetailReq{
		Id: req.HomestayId,
	})
	if err != nil{
		return nil, err
	}
	if homestayResp.Homestay == nil || homestayResp.Homestay .Id == 0{
		return nil,model.ErrNotFound
	}
	userId := ctxdata.GetUidFromCtx(l.ctx)

	order_resp, err := l.svcCtx.OrderRpc.CreateHomestayOrder(l.ctx, &order.CreateHomestayOrderReq{
		HomestayId:    req.HomestayId,
		IsFood:        req.IsFood,
		LiveStartTime: req.LiveStartTime,
		LiveEndTime:   req.LiveEndTime,
		UserId:        userId,
		LivePeopleNum: req.LivePeopleNum,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateHomestayOrderResp{
		OrderSn: order_resp.Sn,
	}, nil
}