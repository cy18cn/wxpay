package wxpay

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/viper"
)

type RefundResp struct {
	Response
	Error
	TradeType           string `xml:"trade_type"`                      //
	TransactionId       string `xml:"transaction_id"`                  // 微信订单号
	OutTradeNo          string `xml:"out_trade_no"`                    // 商户系统内部订单号
	OutRefundNo         string `xml:"out_refund_no"`                   // 商户退款订单号，商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	RefundId            string `xml:"refund_id"`                       // 微信退款单号
	RefundFee           int32  `xml:"refund_fee"`                      // 退款金额 单位为分,可以做部分退款
	SettlementRefundFee int32  `xml:"settlement_refund_fee,omitempty"` // 应结订单金额: 去掉非充值代金券退款金额后的退款金额，退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	TotalFee            int32  `xml:"total_fee"`                       // 订单总金额，单位为分，只能为整数
	SettlementTotalFee  int32  `xml:"settlement_total_fee,omitempty"`  // 应结订单金额: 去掉非充值代金券金额后的订单总金额，应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType             string `xml:"fee_type,omitempty"`              // 货币类型
	CashFee             int32  `xml:"cash_fee"`                        // 现金支付金额，单位为分，只能为整数
	CashFeeType         string `xml:"cash_fee_type,omitempty"`         // 现金支付货币类型
	CashRefundFee       int32  `xml:"cash_refund_fee,omitempty"`       // 现金退款金额
	CouponRefundFee     int32  `xml:"coupon_refund_fee,omitempty"`     // 代金券退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为现金
	CouponRefundCount   int32  `xml:"coupon_refund_count,omitempty"`   // 退款代金券使用数量
}

func NewResp(m map[string]string) *RefundResp {
	var (
		refundFee           int
		settlementRefundFee int
		totalFee            int
		settlementTotalFee  int
		cashFee             int
		cashRefundFee       int
		couponRefundFee     int
		couponRefundCount   int
	)

	refundFee, _ = strconv.Atoi(m["refund_fee"])
	settlementRefundFee, _ = strconv.Atoi(m["settlement_refund_fee"])
	totalFee, _ = strconv.Atoi(m["total_fee"])
	settlementTotalFee, _ = strconv.Atoi(m["settlement_total_fee"])
	cashFee, _ = strconv.Atoi(m["cash_fee"])
	cashRefundFee, _ = strconv.Atoi(m["cash_refund_fee"])
	couponRefundFee, _ = strconv.Atoi(m["coupon_refund_fee"])
	couponRefundCount, _ = strconv.Atoi(m["coupon_refund_count"])

	return &RefundResp{
		Response: Response{
			ReturnCode: m["return_code"],
			ReturnMsg:  m["return_msg"],
			AppId:      m["appid"],
			MchId:      m["mch_id"],
			DeviceInfo: m["device_info"],
			NonceStr:   m["nonce_str"],
			Sign:       m["sign"],
			ResultCode: m["result_code"],
		},
		Error: Error{
			ErrCode:    m["err_code"],
			ErrCodeDes: m["err_code_des"],
		},
		TradeType:           m["trade_type"],
		TransactionId:       m["transaction_id"],
		OutTradeNo:          m["out_trade_no"],
		OutRefundNo:         m["out_refund_no"],
		RefundId:            m["refund_id"],
		RefundFee:           int32(refundFee),
		SettlementRefundFee: int32(settlementRefundFee),
		TotalFee:            int32(totalFee),
		SettlementTotalFee:  int32(settlementTotalFee),
		FeeType:             m["fee_type"],
		CashFee:             int32(cashFee),
		CashFeeType:         m["cash_fee_type"],
		CashRefundFee:       int32(cashRefundFee),
		CouponRefundFee:     int32(couponRefundFee),
		CouponRefundCount:   int32(couponRefundCount),
	}
}

func (self *RefundResp) ToUrlValues() url.Values {
	ua := url.Values{}
	ua.Set("appid", self.AppId)
	ua.Set("mch_id", self.MchId)
	ua.Set("nonce_str", self.NonceStr)
	ua.Set("return_code", self.ReturnCode)
	ua.Set("return_msg", self.ReturnMsg)
	ua.Set("result_code", self.ResultCode)
	ua.Set("err_code", self.ErrCode)
	ua.Set("err_code_des", self.ErrCodeDes)
	ua.Set("device_info", self.DeviceInfo)

	ua.Set("trade_type", self.TradeType)
	ua.Set("transaction_id", self.TransactionId)
	ua.Set("out_trade_no", self.OutTradeNo)
	ua.Set("out_refund_no", self.OutRefundNo)
	ua.Set("refund_id", self.RefundId)
	ua.Set("refund_fee", strconv.Itoa(int(self.RefundFee)))
	ua.Set("settlement_refund_fee", strconv.Itoa(int(self.SettlementRefundFee)))
	ua.Set("total_fee", strconv.Itoa(int(self.TotalFee)))
	ua.Set("settlement_total_fee", strconv.Itoa(int(self.SettlementTotalFee)))
	ua.Set("fee_type", self.FeeType)
	ua.Set("cash_fee", strconv.Itoa(int(self.CashFee)))
	ua.Set("cash_fee_type", self.CashFeeType)
	ua.Set("cash_refund_fee", strconv.Itoa(int(self.CashRefundFee)))
	ua.Set("coupon_refund_fee", strconv.Itoa(int(self.CouponRefundFee)))
	ua.Set("coupon_refund_count", strconv.Itoa(int(self.CouponRefundCount)))

	return ua
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

func (self *RefundRequest) ToUrlValues() url.Values {
	ua := url.Values{}
	ua.Set("appid", self.AppId)
	ua.Set("mch_id", self.MchId)
	ua.Set("sign_type", self.SignType)
	ua.Set("nonce_str", self.NonceStr)

	ua.Set("transaction_id", self.TransactionId)
	ua.Set("out_trade_no", self.OutTradeNo)
	ua.Set("out_trade_no", self.OutRefundNo)
	ua.Set("total_fee", strconv.Itoa(int(self.TotalFee)))
	ua.Set("refund_fee", strconv.Itoa(int(self.RefundFee)))
	ua.Set("refund_fee_type", self.RefundFeeType)
	ua.Set("refund_desc", self.RefundDesc)
	ua.Set("refund_account", self.RefundAccount)
	ua.Set("notify_url", self.NotifyUrl)

	return ua
}

func Refund(req *RefundRequest, mchKey, mchId string) (*RefundResp, error) {
	var err error
	//tlsClient := tlsClients[mchId]
	if _, ok := tlsClients[mchId]; !ok {
		err = InitClient(req.MchId, viper.GetString("certPath"), viper.GetDuration("timeout"), viper.GetBool("isProd"))
		if err != nil {
			return nil, err
		}
	}

	tlsClient := tlsClients[mchId]

	if req.SignType == SIGN_MD5 {
		req.Sign = SignMD5(req.ToUrlValues(), mchKey)
	} else {
		req.Sign = SignHmacSha256(req.ToUrlValues(), mchKey)
	}

	var reply []byte
	reply, err = tlsClient.DoRequest(http.MethodPost, "/secapi/pay/refund", req)
	if err != nil {
		return nil, err
	}

	respMap := UnmarshalToMap(reply)
	if !VerifyResp(respMap["xml"].(map[string]interface{}), mchKey, respMap["sign"].(string), req.SignType) {
		return nil, errors.New("签名错误")
	}

	return NewResp(respMap["xml"].(map[string]string)), nil
}
