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

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	Resp,err:=l.svcCtx.UserRpc.Login(l.ctx,&pb.LoginReq{
		AuthType: vars.UserAuthTypeSystem,
		AuthKey: req.Mobile,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	var loginResp types.LoginResp
	_ = copier.Copy(&resp, Resp)

	return &loginResp, nil
}
