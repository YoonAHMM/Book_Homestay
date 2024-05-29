package logic

import (
	"Book_Homestay/app/user/cmd/rpc/internal/svc"
	"Book_Homestay/app/user/cmd/rpc/pb"
	"Book_Homestay/app/user/cmd/rpc/user"
	"Book_Homestay/app/user/model"
	"Book_Homestay/common/Randx"
	"Book_Homestay/common/calculate"
	"Book_Homestay/common/errx"
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	u, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx,in.Mobile)
	if err != nil && err != model.ErrNotFound {
		return nil,errors.Wrapf(errx.NewErrCode(errx.DB_ERROR), "mobile:%s,err:%v", in.Mobile, err)
	}
	if u != nil {
		return nil, errors.Wrapf(errx.NewErrMsg("user has been registered"), "Register user exists mobile:%s,err:%v", in.Mobile, err)
	}

	var userId int64
	if err := l.svcCtx.UserModel.Trans(l.ctx,func(ctx context.Context,session sqlx.Session) error {
		user := new(model.User)
		user.Mobile = in.Mobile

		if len(user.Nickname) == 0 {
			user.Nickname = Randx.RandName()
		}
		if len(in.Password) > 0 {
			user.Password = calculate.Md5ByString(in.Password)
		}

		insertResult, err := l.svcCtx.UserModel.Insert(ctx,session, user)
		if err != nil {
			return errors.Wrapf(errx.NewErrCode(errx.DB_ERROR), "Register db user Insert err:%v,user:%+v", err, user)
		}
		lastId, err := insertResult.LastInsertId()
		if err != nil {
			return errors.Wrapf(errx.NewErrCode(errx.DB_ERROR), "Register db user insertResult.LastInsertId err:%v,user:%+v", err, user)
		}
		userId = lastId

		userAuth := new(model.UserAuth)
		userAuth.UserId = lastId
		userAuth.AuthKey = in.AuthKey
		userAuth.AuthType = in.AuthType
		if _, err := l.svcCtx.UserAuthModel.Insert(ctx,session,userAuth); err != nil {
			return  errors.Wrapf(errx.NewErrCode(errx.DB_ERROR), "Register db user_auth Insert err:%v,userAuth:%v", err, userAuth)
		}
		return nil
	}); err != nil {
		return nil, err
	}


	g :=NewGenerateTokenLogic(l.ctx,l.svcCtx)
	tokenResp,err:=g.GenerateToken(&user.GenerateTokenReq{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.Wrapf(errx.NewErrCode(errx.TOKEN_GENERATE_ERROR),"GenerateToken userId : %d", userId)
	}

	return &user.RegisterResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
