package logic

import (
	"context"

	"Book_Homestay/app/user/cmd/rpc/internal/svc"
	"Book_Homestay/app/user/cmd/rpc/pb"
	"Book_Homestay/app/user/cmd/rpc/user"
	"Book_Homestay/app/user/model"
	"Book_Homestay/common/errx"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	
	u, err := l.svcCtx.UserModel.FindOne(l.ctx,in.Id)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errx.NewErrCodeMsg(errx.DB_ERROR,"get userinfo fail"),"err : %v , in : %+v", err, in)
	}
	if u == nil {
		return nil, model.ErrNotFound
	}

	var resp user.User
	_ = copier.Copy(&resp, u)

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil

}
