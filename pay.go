package wxpay

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/viper"
)

type UnifiedOrderResp struct {
	Response
	Error
	TradeType string `xml:"trade_type"`
	PrepayId  string `xml:"prepay_id"`
}

// UnifiedOrderRequest 微信统一下单
// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
type UnifiedOrderRequest struct {
	XMLName xml.Name `xml:"xml"`
	Request
	DeviceInfo     string `xml:"device_info,omitempty"`
	Body           string `xml:"body"`
	Detail         string `xml:"detail,omitempty"`
	Attach         string `xml:"attach,omitempty"`
	OutTradeNo     string `xml:"out_trade_no"`
	FeeType        string `xml:"fee_type,omitempty"`
	TotalFee       int32  `xml:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip,omitempty"`
	TimeStart      string `xml:"time_start,omitempty"`
	TimeExpire     string `xml:"time_expire,omitempty"`
	GoodsTag       string `xml:"goods_tag,omitempty"` // 订单优惠标记
	NotifyUrl      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
	LimitPay       string `xml:"limit_pay,omitempty"`
	Receipt        string `xml:"receipt,omitempty"`
}

func (self *UnifiedOrderRequest) Payload() (string, error) {
	if len(self.Sign) == 0 {
		return "", errors.New(MISSING_SIGN)
	}

	b, err := xml.Marshal(self)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (self *UnifiedOrderRequest) ToUrlValues() url.Values {
	ua := url.Values{}
	ua.Set("appid", self.AppId)
	ua.Set("mch_id", self.MchId)
	ua.Set("sign_type", self.SignType)
	ua.Set("nonce_str", self.NonceStr)

	ua.Set("device_info", self.DeviceInfo)
	ua.Set("body", self.Body)
	ua.Set("detail", self.Detail)
	ua.Set("attach", self.Attach)
	ua.Set("out_trade_no", self.OutTradeNo)
	ua.Set("fee_type", self.FeeType)
	ua.Set("total_fee", strconv.Itoa(int(self.TotalFee)))
	ua.Set("spbill_create_ip", self.SpbillCreateIP)
	ua.Set("time_start", self.TimeStart)
	ua.Set("time_expire", self.TimeExpire)
	ua.Set("goods_tag", self.GoodsTag)
	ua.Set("notify_url", self.NotifyUrl)
	ua.Set("trade_type", self.TradeType)
	ua.Set("limit_pay", self.LimitPay)
	ua.Set("receipt", self.Receipt)

	return ua
}

func UnifiedOrder(request *UnifiedOrderRequest, mchKey string) (*UnifiedOrderResp, error) {
	request.Sign = Sign(request.ToUrlValues(), mchKey)

	var err error
	if client == nil {
		client, err = InitClient(viper.GetBool("isProd"), "", "")
	}

	// var reply *UnifiedOrderResp 此处reply为nil
	// var reply UnifiedOrderResp 此处reply已初始化，属性值为各类型的默认值
	//client.DoRequest(http.MethodPost, "/pay/unifiedorder", request, &reply)

	// 此处reply为*UnifiedOrderResp
	// reply = &UnifiedOrderResp{}
	//client.DoRequest(http.MethodPost, "/pay/unifiedorder", request, reply)

	var reply UnifiedOrderResp
	client.DoRequest(http.MethodPost, "/pay/unifiedorder", request, &reply)

	if err != nil {
		return nil, err
	}

	return &reply, nil
}
