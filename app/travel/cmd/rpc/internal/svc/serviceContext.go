package svc

import (
	"Book_Homestay/app/travel/cmd/rpc/internal/config"
	"Book_Homestay/app/travel/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	Model model.HomestayModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn:= sqlx.NewMysql(c.DB.DataSource)
	
	return &ServiceContext{
		Config: c,
		Model: model.NewHomestayModel(sqlConn,c.Cache),
	}
}
