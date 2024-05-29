package logic

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"Book_Homestay/app/mqueue/cmd/job/jobtype"
	"Book_Homestay/app/order/cmd/rpc/internal/svc"
	"Book_Homestay/app/order/cmd/rpc/pb"
	"Book_Homestay/app/order/model"
	"Book_Homestay/app/travel/cmd/rpc/client/homestay"

	"Book_Homestay/common/Randx"
	"Book_Homestay/common/errx"
	"Book_Homestay/common/uniqueid"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const CloseOrderTimeMinutes = 30 
type CreateHomestayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿下订单
func (l *CreateHomestayOrderLogic) CreateHomestayOrder(in *pb.CreateHomestayOrderReq) (*pb.CreateHomestayOrderResp, error) {
	if in.LiveEndTime <= in.LiveStartTime {
		return nil, errors.Wrapf(errx.NewErrMsg("order illegal"),"in : %+v",in)
	}

	resp, err := l.svcCtx.Homestay_TravelRpc.HomestayDetail(l.ctx, &homestay.HomestayDetailReq{
		Id: in.HomestayId,
	})
	if err != nil {
		return nil, err
	}
	if resp.Homestay == nil {
		return nil, model.ErrNotFound
	}

	var cover string //Get the cover...
	if len(resp.Homestay.Banner) > 0 {
		cover = strings.Split(resp.Homestay.Banner, ",")[0]
	}

	order := new(model.HomestayOrder)
	order.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_HOMESTAY_ORDER)
	order.UserId = in.UserId
	order.HomestayId = in.HomestayId
	order.Title = resp.Homestay.Title
	order.SubTitle = resp.Homestay.SubTitle
	order.Cover = cover
	order.Info = resp.Homestay.Info
	order.PeopleNum = resp.Homestay.PeopleNum
	order.RowType = resp.Homestay.RowType
	order.HomestayPrice = resp.Homestay.HomestayPrice
	order.MarketHomestayPrice = resp.Homestay.MarketHomestayPrice
	order.HomestayBusinessId = resp.Homestay.HomestayBusinessId
	order.HomestayUserId = resp.Homestay.UserId
	order.LivePeopleNum = in.LivePeopleNum
	order.TradeState = model.HomestayOrderTradeStateWaitPay
	order.TradeCode = Randx.Krand(8, Randx.KC_RAND_KIND_ALL)
	order.Remark = in.Remark
	order.FoodInfo = resp.Homestay.FoodInfo
	order.FoodPrice = resp.Homestay.FoodPrice
	order.LiveStartDate = time.Unix(in.LiveStartTime, 0)
	order.LiveEndDate = time.Unix(in.LiveEndTime, 0)

	liveDays := int64(order.LiveEndDate.Sub(order.LiveStartDate).Seconds() / 86400) //Stayed a few days in total

	order.HomestayTotalPrice = int64(resp.Homestay.HomestayPrice * liveDays) //Calculate the total price of the B&B
	if in.IsFood {
		order.NeedFood = model.HomestayOrderNeedFoodYes
		//Calculate the total price of the meal.
		order.FoodTotalPrice = int64(resp.Homestay.FoodPrice * in.LivePeopleNum * liveDays)
	}

	order.OrderTotalPrice = order.HomestayTotalPrice + order.FoodTotalPrice //Calculate total order price.

	_, err = l.svcCtx.HomestayOrderModel.Insert(l.ctx,nil, order)
	if err != nil {
		return nil, errors.Wrapf(errx.NewErrCodeMsg(errx.DB_ERROR,"insert HomestayOrder fail"),"err:%v, in:%+v",err,in)
	}



	payload, err := json.Marshal(jobtype.DeferCloseHomestayOrderPayload{Sn: order.Sn})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("create defer close order task json Marshal fail err :%+v , sn : %s",err,order.Sn)
	}else{
		_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(jobtype.DeferCloseHomestayOrder, payload), asynq.ProcessIn(CloseOrderTimeMinutes * time.Minute))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("create defer close order task insert queue fail err :%+v , sn : %s",err,order.Sn)
		}
	}

	return &pb.CreateHomestayOrderResp{
		Sn: order.Sn,
	}, nil
}

