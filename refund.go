package wxpay

import (
	"encoding/xml"
	"errors"
	"github.com/spf13/viper"
)

type RefundResp struct {
	Response
	Error
	TradeType           string                 `xml:"trade_type"`                      //
	TransactionId       string                 `xml:"transaction_id"`                  // 微信订单号
	OutTradeNo          string                 `xml:"out_trade_no"`                    // 商户系统内部订单号
	OutRefundNo         string                 `xml:"out_refund_no"`                   // 商户退款订单号，商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	RefundId            string                 `xml:"refund_id"`                       // 微信退款单号
	RefundFee           int32                  `xml:"refund_fee"`                      // 退款金额 单位为分,可以做部分退款
	SettlementRefundFee int32                  `xml:"settlement_refund_fee,omitempty"` // 应结订单金额: 去掉非充值代金券退款金额后的退款金额，退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	TotalFee            int32                  `xml:"total_fee"`                       // 订单总金额，单位为分，只能为整数
	SettlementTotalFee  int32                  `xml:"settlement_total_fee,omitempty"`  // 应结订单金额: 去掉非充值代金券金额后的订单总金额，应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType             string                 `xml:"fee_type,omitempty"`              // 货币类型
	CashFee             int32                  `xml:"cash_fee"`                        // 现金支付金额，单位为分，只能为整数
	CashFeeType         string                 `xml:"cash_fee_type,omitempty"`         // 现金支付货币类型
	CashRefundFee       int32                  `xml:"cash_refund_fee,omitempty"`       // 现金退款金额
	CouponRefundFee     int32                  `xml:"coupon_refund_fee,omitempty"`     // 代金券退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为现金
	CouponRefundCount   int32                  `xml:"coupon_refund_count,omitempty"`   // 退款代金券使用数量
	Coupon              map[string]interface{} `xml:",any"`
}

type RefundRequest struct {
	XMLName xml.Name `xml:"xml"`
	Request
	TransactionId string `xml:"transaction_id,omitempty"`  // 必填：out_trade_no二选一，微信生成的订单号，在支付通知中有返回
	OutTradeNo    string `xml:"out_trade_no,omitempty"`    // 必填：transaction_id二选一，商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	OutRefundNo   string `xml:"out_refund_no"`             // 必填：是，商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔
	TotalFee      int32  `xml:"total_fee"`                 // 必填：是，订单总金额，单位为分，只能为整数
	RefundFee     int32  `xml:"refund_fee"`                // 必填：是，退款金额
	RefundFeeType string `xml:"refund_fee_type,omitempty"` // 必填：否，退款货币类型，需与支付一致，或者不填
	RefundDesc    string `xml:"refund_desc,omitempty"`     // 必填：否，退款原因， 若商户传入，会在下发给用户的退款消息中体现退款原因
	RefundAccount string `xml:"refund_account,omitempty"`  // 必填：否，仅针对老资金流商户使用 REFUND_SOURCE_UNSETTLED_FUNDS---未结算资金退款（默认使用未结算资金退款） REFUND_SOURCE_RECHARGE_FUNDS---可用余额退款
	NotifyUrl     string `xml:"notify_url,omitempty"`      // 必填：否，异步接收微信支付退款结果通知的回调地址，通知URL必须为外网可访问的url，不允许带参数 如果参数中传了notify_url，则商户平台上配置的回调地址将不会生效
}

func (self *RefundRequest) Payload() (string, error) {
	if len(self.Sign) == 0 {
		return "", errors.New(MISSING_SIGN)
	}

	if len(self.TransactionId) == 0 && len(self.OutTradeNo) == 0 {
		return "", errors.New("refund: transaction_id or out_trade_no cannot be empty")
	}

	b, err := xml.Marshal(self)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func Refund(req *RefundRequest, mchKey, mchId string) (*RefundResp, error) {
	var err error
	//tlsClient := tlsClients[mchId]
	if _, ok := tlsClients[mchId]; !ok {
		err = InitClient(viper.GetBool("isProd"), req.MchId, viper.GetString("certPath"))
		if err != nil {
			return nil, err
		}
	}

	tlsClient := tlsClients[mchId]

	if req.SignType == SIGN_MD5 {
		//req.Sign = SignMD5()
	}

}

type RefundCoupon struct {
}
