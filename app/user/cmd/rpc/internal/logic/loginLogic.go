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

	"github.com/pkg/errors"
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
		return nil,errors.Wrapf(errx.NewErrMsg("authtype illegal")," in : %+v", in)
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
		return nil, errors.Wrapf(errx.NewErrMsg("账号或密码错误"),"账号或密码错误")
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
