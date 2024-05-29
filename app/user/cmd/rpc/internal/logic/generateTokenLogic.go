package logic

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"Book_Homestay/app/user/cmd/rpc/internal/svc"
	"Book_Homestay/app/user/cmd/rpc/pb"
	"Book_Homestay/common/errx"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var ErrUsernamePwdError = errx.NewErrMsg("账号或密码不正确")

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	now := time.Now().Unix()

	accessExpire:=l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken,err:=l.makeJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret,now,accessExpire,in.UserId)

	if err!=nil {
		return nil, errors.Wrapf(errx.NewErrCode(errx.TOKEN_EXPIRE_ERROR), "getJwtToken err userId:%d , err:%v", in.UserId, err)
	}

	return &pb.GenerateTokenResp{
		AccessExpire: accessExpire,
		AccessToken: accessToken,
		RefreshAfter: now + accessExpire*2/3,
	}, nil
}


func(l*GenerateTokenLogic)makeJwtToken(secret string,time ,Expire ,id int64)(string,error){
	c := make(jwt.MapClaims)

	c["iat"]=time
	c["exp"]=time+Expire
	c["aud"]=id

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = c
	return token.SignedString([]byte(secret))
}