package svc

import (
	"Book_Homestay/app/travel/cmd/rpc/internal/config"
	"Book_Homestay/app/travel/model"
	"Book_Homestay/app/user/cmd/rpc/user"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	//model
	HomestayModel         model.HomestayModel
	HomestayActivityModel model.HomestayActivityModel
	HomestayBusinessModel model.HomestayBusinessModel
	HomestayCommentModel  model.HomestayCommentModel

	UserRpc      user.UserZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn:= sqlx.NewMysql(c.DB.DataSource)
	
	return &ServiceContext{
		Config: c,
		
		HomestayModel: model.NewHomestayModel(sqlConn,c.Cache),
		HomestayActivityModel: model.NewHomestayActivityModel(sqlConn,c.Cache),
		HomestayBusinessModel: model.NewHomestayBusinessModel(sqlConn,c.Cache),
		HomestayCommentModel: model.NewHomestayCommentModel(sqlConn,c.Cache),

		UserRpc: user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
