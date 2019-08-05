package wxpay

type UnifiedOrderResp struct {
	Response
	Error
	TradeType string
	PrepayId  string
}

// UnifiedOrderRequest 微信统一下单
// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
type UnifiedOrderRequest struct {
	Request
	DeviceInfo     string `xml:"device_info,omitempty"`
	Body           string `xml:"body"`
	Detail         string `xml:"detail,omitempty"`
	Attach         string `xml:"attach,omitempty"`
	OutTradeNo     string `xml:"out_trade_no"`
	FeeType        string `xml:"fee_type,omitempty"`
	TotalFee       string `xml:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip,omitempty"`
	TimeStart      string `xml:"time_start,omitempty"`
	TimeExpire     string `xml:"time_expire,omitempty"`
	GoodsTag       string `xml:"goods_tag,omitempty"`
	NotifyUrl      string `xml:"notifyUrl"`
	TradeType      string `xml:"trade_type"`
	LimitPay       string `xml:"limit_pay,omitempty"`
	Receipt        string `xml:"receipt,omitempty"`
}
