package wxpay

type Request struct {
	AppId    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
	SignType string `xml:"sign_type,omitempty"`
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
