package logic

import (
	"context"

	"Book_Homestay/app/user/cmd/rpc/internal/svc"
	"Book_Homestay/app/user/cmd/rpc/pb"
	"Book_Homestay/app/user/cmd/rpc/user"
	"Book_Homestay/app/user/model"
	"Book_Homestay/common/calculate"
	"Book_Homestay/common/errx"
	"Book_Homestay/common/vars"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	if in.AuthType != vars.UserAuthTypeSystem{
		return nil,errx.NewErrCode(errx.LOGIN_ERROR,"UserAuthTypeSystem wrong")
	}

	mobile:=in.AuthKey

	u, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx,mobile)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	if u == nil {
		return nil, model.ErrNotFound
	}

	if calculate.Md5ByString(in.Password) != u.Password {
		return nil, errx.NewErrCode(errx.LOGIN_ERROR,"Password wrong")
	}
	
	g:=NewGenerateTokenLogic(l.ctx,l.svcCtx)

	tokenresp,err:=g.GenerateToken(&pb.GenerateTokenReq{
		UserId: u.Id,
	})
	if err != nil{
		return nil,err
	}

	return &user.LoginResp{
		AccessToken:  tokenresp.AccessToken,
		AccessExpire: tokenresp.AccessExpire,
		RefreshAfter: tokenresp.RefreshAfter,
	}, nil
}
