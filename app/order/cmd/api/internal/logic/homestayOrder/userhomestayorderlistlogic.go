package homestayOrder

import (
	"context"

	"Book_Homestay/app/order/cmd/api/internal/svc"
	"Book_Homestay/app/order/cmd/api/internal/types"
	"Book_Homestay/app/order/cmd/rpc/order"
	"Book_Homestay/common/calculate"
	"Book_Homestay/common/ctxdata"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserHomestayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderListLogic {
	return &UserHomestayOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserHomestayOrderListLogic) UserHomestayOrderList(req *types.UserHomestayOrderListReq) (resp *types.UserHomestayOrderListResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx) //get login user id

	orderresp, err := l.svcCtx.OrderRpc.UserHomestayOrderList(l.ctx, &order.UserHomestayOrderListReq{
		UserId:      userId,
		TraderState: req.TradeState,
		PageSize:    req.PageSize,
		LastId:      req.LastId,
	})
	if err != nil {
		return nil, err
	}

	var typesUserHomestayOrderList []types.UserHomestayOrderListView

	if len(orderresp.List) > 0 {

		for _, homestayOrder := range orderresp.List {

			var typeHomestayOrder types.UserHomestayOrderListView
			_ = copier.Copy(&typeHomestayOrder, homestayOrder)

			typeHomestayOrder.OrderTotalPrice = calculate.Fen2Yuan(homestayOrder.OrderTotalPrice)

			typesUserHomestayOrderList = append(typesUserHomestayOrderList, typeHomestayOrder)
		}
	}

	return &types.UserHomestayOrderListResp{
		List: typesUserHomestayOrderList,
	}, nil

}
