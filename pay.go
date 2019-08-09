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
	TradeType string `xml:"trade_type"` //
	PrepayId  string `xml:"prepay_id"`  // 微信生成的预支付回话标识，用于后续接口调用中使用，该值有效期为2小时
}

func (self *UnifiedOrderResp) ToUrlValues() url.Values {
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
	ua.Set("prepay_id", self.PrepayId)

	return ua
}

// UnifiedOrderRequest 微信统一下单
// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
type UnifiedOrderRequest struct {
	XMLName xml.Name `xml:"xml"`
	Request
	DeviceInfo     string `xml:"device_info,omitempty"`      // 必填：否，device_info：终端设备号(门店号或收银设备ID)，默认请传"WEB"
	Body           string `xml:"body"`                       // 必填：是，body: APP——需传入应用市场上的APP名字-实际商品名称，天天爱消除-游戏充值。
	Detail         string `xml:"detail,omitempty"`           // 必填：否，detail：商品详细描述，对于使用单品优惠的商户，该字段必须按照规范上传
	Attach         string `xml:"attach,omitempty"`           // 必填：否，attach：附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	OutTradeNo     string `xml:"out_trade_no"`               // 必填：是，out_trade_no：商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*且在同一个商户号下唯一
	FeeType        string `xml:"fee_type,omitempty"`         // 必填：否，fee_type：符合ISO 4217标准的三位字母代码，默认人民币：CNY
	TotalFee       int32  `xml:"total_fee"`                  // 必填：是，total_fee：订单总金额，单位为分
	SpbillCreateIP string `xml:"spbill_create_ip,omitempty"` // 必填：否，spbill_create_ip：支持IPV4和IPV6两种格式的IP地址。调用微信支付API的机器IP
	TimeStart      string `xml:"time_start,omitempty"`       // 必填：否，time_start：订单生成时间，格式为yyyyMMddHHmmss
	TimeExpire     string `xml:"time_expire,omitempty"`      // 必填：否，time_expire：订单失效时间，格式为yyyyMMddHHmmss
	GoodsTag       string `xml:"goods_tag,omitempty"`        // 必填：否，goods_tag：订单优惠标记
	NotifyUrl      string `xml:"notify_url"`                 // 必填：是，notify_url：接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	TradeType      string `xml:"trade_type"`                 // 必填：是，trade_type：支付类型（APP,JSAPI,NATIVE）
	LimitPay       string `xml:"limit_pay,omitempty"`        // 必填：否，limit_pay：no_credit--指定不能使用信用卡支付
	Receipt        string `xml:"receipt,omitempty"`          // 必填：否，receipt：Y，传入Y时，支付成功消息和支付详情页将出现开票入口
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
	if request.SignType == SIGN_MD5 {
		request.Sign = SignMD5(request.ToUrlValues(), mchKey)
	} else {
		request.Sign = SignHmacSha256(request.ToUrlValues(), mchKey)
	}

	var err error
	if client == nil {
		err = InitClient(viper.GetBool("isProd"), "", "")
		if err != nil {
			return nil, err
		}
	}

	// var reply *UnifiedOrderResp 此处reply为nil
	// var reply UnifiedOrderResp 此处reply已初始化，属性值为各类型的默认值
	//client.DoRequest(http.MethodPost, "/pay/unifiedorder", request, &reply)

	// 此处reply为*UnifiedOrderResp
	// reply = &UnifiedOrderResp{}
	//client.DoRequest(http.MethodPost, "/pay/unifiedorder", request, reply)

	var reply []byte
	reply, err = client.DoRequest(http.MethodPost, "/pay/unifiedorder", request)
	if err != nil {
		return nil, err
	}

	var resp UnifiedOrderResp
	xml.Unmarshal(reply, &resp)
	if err != nil {
		return nil, err
	}

	if !VerifySign(resp.ToUrlValues(), mchKey, resp.Sign, request.SignType) {
		return nil, errors.New("Result签名错误")
	}

	return &resp, nil
}
