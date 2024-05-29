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

type GetUserAuthByAuthKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAuthByAuthKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAuthByAuthKeyLogic {
	return &GetUserAuthByAuthKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAuthByAuthKeyLogic) GetUserAuthByAuthKey(in *pb.GetUserAuthByAuthKeyReq) (*pb.GetUserAuthByAuthKeyResp, error) {
	userAuth, err := l.svcCtx.UserAuthModel.FindOneByAuthTypeAuthKey(l.ctx,in.AuthType, in.AuthKey)

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errx.NewErrCodeMsg(errx.DB_ERROR,"get user auth fail"), "err : %v , in : %+v", err, in)
	}

	var resp user.UserAuth
	_ = copier.Copy(&resp, userAuth)

	return &pb.GetUserAuthByAuthKeyResp{
		UserAuth: &resp,
	}, nil
	
}

