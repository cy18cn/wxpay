package wxpay

type Request struct {
	AppId    string `xml:"appid"`               // 必填：是，appid: 微信开放平台审核通过的应用APPID
	MchId    string `xml:"mch_id"`              // 必填：是，mch_id: 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"`           // 必填：是，nonce_str: 随机字符串，不长于32位
	Sign     string `xml:"sign"`                // 必填：是，sign: 签名
	SignType string `xml:"sign_type,omitempty"` // 必填：否，sign_type: 签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
}

type Response struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info,omitempty"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
}

type WxpayRequest interface {
	Payload() (string, error)
}

type WxpayResponse interface {
	FromXml(data string) interface{}
}
