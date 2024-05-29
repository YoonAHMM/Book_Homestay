package svc

import (
	"Book_Homestay/app/payment/cmd/api/internal/config"
	"Book_Homestay/common/errx"
	"context"

	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

func NewWxPayClientV3(c config.Config) (*core.Client, error) {

	mchPrivateKey, err := utils.LoadPrivateKey(c.WxPayConf.PrivateKey)
	if err != nil {
		return nil, errors.Wrapf(errx.NewErrMsg("PrivateKey 失效"),"err :%v",err)
	}

	ctx := context.Background()
	
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(c.WxPayConf.MchId, c.WxPayConf.SerialNo, mchPrivateKey, c.WxPayConf.APIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, errors.Wrapf(errx.NewErrMsg("wx客户端创建失败"),"err : %v ",err)
	}

	return client, nil

}