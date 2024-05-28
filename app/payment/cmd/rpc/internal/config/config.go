package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	DB struct {
		DataSource string
	}
	Cache cache.CacheConf
	KqPaymentUpdatePayStatusConf KqConfig
}

type KqConfig struct {
	Brokers []string
	Topic   string
}

type KqServerConfig struct {
	Address string `json:"AppId"`  //wechat mini appId
	Secret  string `json:"Secret"` //wechat mini secret
}
