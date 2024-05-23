package user

import (
	"context"

	"Book_Homestay/app/user/cmd/api/internal/svc"
	"Book_Homestay/app/user/cmd/api/internal/types"
	"Book_Homestay/app/user/cmd/rpc/pb"
	"Book_Homestay/common/ctxdata"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	userid:=ctxdata.GetUidFromCtx(l.ctx)

	Resp,err:=l.svcCtx.UserRpc.GetUserInfo(l.ctx,&pb.GetUserInfoReq{
		Id: userid,
	})
	if err!=nil{
		return nil,err
	}

	var userInfo types.User
	_ = copier.Copy(&userInfo, Resp.User)
	
	return &types.UserInfoResp{
		UserInfo:userInfo,
	},nil
}
