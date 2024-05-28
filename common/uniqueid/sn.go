package uniqueid

import (
	"Book_Homestay/common/Randx"
	"fmt"
	"time"
)

//生成sn单号
type SnPrefix string

const (
	SN_PREFIX_HOMESTAY_ORDER SnPrefix = "HSO" //民宿订单前缀 order/homestay_order
	SN_PREFIX_THIRD_PAYMENT  SnPrefix = "PMT" //第三方支付流水记录前缀 payment/third_payment
)

//生成单号
func GenSn(snPrefix SnPrefix) string {
	return fmt.Sprintf("%s%s%s", snPrefix, time.Now().Format("20060102150405"), Randx.Krand(8, Randx.KC_RAND_KIND_NUM))
}
