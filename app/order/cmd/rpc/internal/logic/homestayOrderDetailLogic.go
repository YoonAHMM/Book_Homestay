package logic

import (
	"context"

	"Book_Homestay/app/order/cmd/rpc/internal/svc"
	"Book_Homestay/app/order/cmd/rpc/pb"
	"Book_Homestay/app/order/model"
	"Book_Homestay/common/errx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayOrderDetailLogic {
	return &HomestayOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿订单详情
func (l *HomestayOrderDetailLogic) HomestayOrderDetail(in *pb.HomestayOrderDetailReq) (*pb.HomestayOrderDetailResp, error) {
	homestayOrder, err := l.svcCtx.HomestayOrderModel.FindOneBySn(l.ctx,in.Sn)
	
	if err != nil && err != model.ErrNotFound {
		return nil, errx.NewErrCode(errx.DB_ERROR,err.Error())
	}

	var resp pb.HomestayOrder
	if homestayOrder != nil {
		_ = copier.Copy(&resp, homestayOrder)

		resp.CreateTime = homestayOrder.CreateTime.Unix()
		resp.LiveStartDate = homestayOrder.LiveStartDate.Unix()
		resp.LiveEndDate = homestayOrder.LiveEndDate.Unix()

	}

	return &pb.HomestayOrderDetailResp{
		HomestayOrder: &resp,
	}, nil
}
