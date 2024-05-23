package user

import (
	"context"

	"Book_Homestay/app/user/cmd/api/internal/svc"
	"Book_Homestay/app/user/cmd/api/internal/types"
	"Book_Homestay/app/user/cmd/rpc/pb"
	"Book_Homestay/common/vars"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	Resp,err:=l.svcCtx.UserRpc.Register(l.ctx,&pb.RegisterReq{
		AuthType: vars.UserAuthTypeSystem,
		AuthKey: req.Mobile,
		Password: req.Password,
		Mobile: req.Mobile,
	})

	if err != nil {
		return nil, err
	}

	var registerResp types.RegisterResp
	_ = copier.Copy(&registerResp, Resp)

	return &registerResp,nil
}
