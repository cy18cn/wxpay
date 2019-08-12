package wxpay

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type PayNotification struct {
	Response
	Error
	OpenId        string `xml:"openid"`             // required，用户标识
	IsSubscribe   string `xml:"is_subscribe"`       // required，用户是否关注公众账号，Y-关注，N-未关注
	TradeType     string `xml:"trade_type"`         // required, 交易类型APP...
	BankType      string `xml:"bank_type"`          // required, CMC, 银行类型，采用字符串类型的银行标识
	TotalFee      int32  `xml:"total_fee"`          // required, 订单总金额，单位为分
	FeeType       string `xml:"fee_type,omitempty"` // optional, 货币类型，符合ISO4217标准的三位字母代码，默认人民币：CNY
	CashFee       int32  `xml:"cash_fee"`           // required, 现金支付金额订单现金支付金额
	CashFeeType   string `xml:"cash_fee_type"`      // optional, 货币类型
	CouponFee     int32  `xml:"coupon_fee"`         // optional, 代金券或立减优惠金额<=订单总金额，订单总金额-代金券或立减优惠金额=现金支付金额
	CouponCount   int32  `xml:"coupon_count"`       // optional, 代金券或立减优惠使用数量
	TransactionId string `xml:"transaction_id"`     // required, 微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`       // required, 商户订单号
	Attach        string `xml:"attach"`             // optional, 商家数据包，原样返回
	TimeEnd       string `xml:"time_end"`           // required, 支付完成时间
	CouponId0     string `xml:"coupon_id_0"`        // optional, 第一张优惠券id
	CouponFee0    int32  `xml:"coupon_fee_0"`       // optional, 第一张优惠券优惠金额
}

func (self *PayNotification) ToUrlValues() url.Values {
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

	ua.Set("openid", self.OpenId)
	ua.Set("is_subscribe", self.IsSubscribe)
	ua.Set("trade_type", self.TradeType)
	ua.Set("bank_type", self.BankType)
	ua.Set("total_fee", strconv.Itoa(int(self.TotalFee)))
	ua.Set("fee_type", self.FeeType)
	ua.Set("cash_fee", strconv.Itoa(int(self.CashFee)))
	ua.Set("cash_fee_type", self.CashFeeType)
	ua.Set("coupon_fee", strconv.Itoa(int(self.CouponFee)))
	ua.Set("coupon_count", strconv.Itoa(int(self.CouponCount)))
	ua.Set("transaction_id", self.TransactionId)
	ua.Set("out_trade_no", self.OutTradeNo)
	ua.Set("attach", self.Attach)
	ua.Set("TimeEnd", self.TimeEnd)

	return ua
}

type NotifyResp struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

func AckNotification(w http.ResponseWriter) error {
	resp := NotifyResp{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}

	b, err := xml.Marshal(&resp)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return nil
}

func GetNotification(req *http.Request) (*PayNotification, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var nofi *PayNotification
	err = xml.Unmarshal(body, nofi)
	if err != nil {
		return nil, err
	}

	return nofi, nil
}
