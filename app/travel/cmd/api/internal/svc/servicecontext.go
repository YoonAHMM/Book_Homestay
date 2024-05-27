package svc

import (
	"Book_Homestay/app/travel/cmd/api/internal/config"
	
	"Book_Homestay/app/travel/cmd/rpc/client/homestay"
	"Book_Homestay/app/travel/cmd/rpc/client/homestaybussiness"
	"Book_Homestay/app/travel/cmd/rpc/client/homestaycomment"
	"Book_Homestay/app/travel/model"
	"Book_Homestay/app/user/cmd/rpc/user"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	//rpc
	UserRpc      user.UserZrpcClient
	Homestay_TravelRpc homestay.HomestayZrpcClient
	Bussiness_TravelRpc    homestaybussiness.Homestaybussiness
	Comment_TravelRpc      homestaycomment.Homestaycomment

	//model
	HomestayModel         model.HomestayModel
	HomestayActivityModel model.HomestayActivityModel
	HomestayBusinessModel model.HomestayBusinessModel
	HomestayCommentModel  model.HomestayCommentModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn:= sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config: c,

		UserRpc: user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
		Homestay_TravelRpc: homestay.NewHomestayZrpcClient(zrpc.MustNewClient(c.TravelRpcConf)),
		Bussiness_TravelRpc: homestaybussiness.NewHomestaybussiness(zrpc.MustNewClient(c.TravelRpcConf)),
		Comment_TravelRpc:   homestaycomment.NewHomestaycomment(zrpc.MustNewClient(c.TravelRpcConf)),

		HomestayModel: model.NewHomestayModel(sqlConn,c.Cache),
		HomestayActivityModel: model.NewHomestayActivityModel(sqlConn,c.Cache),
		HomestayBusinessModel: model.NewHomestayBusinessModel(sqlConn,c.Cache),
		HomestayCommentModel: model.NewHomestayCommentModel(sqlConn,c.Cache),
	}
}
