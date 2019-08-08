package wxpay

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnifiedOrderRequest_Payload(t *testing.T) {
	req := &UnifiedOrderRequest{
		Request: Request{
			AppId:    "wxa26bd20b8854ba2a",
			MchId:    "1286242401",
			NonceStr: NonceStr(),
			SignType: "MD5",
			Sign:     "9A0A8659F005D6984697E2CA0A9CF3B7",
		},
		Body:           "test product",
		NotifyUrl:      "http://www.test.com",
		TradeType:      "APP",
		SpbillCreateIP: "202.105.107.18",
		TotalFee:       101,
		OutTradeNo:     "test-111111125",
	}

	str, err := req.Payload()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str)
}

func TestUnifiedOrder(t *testing.T) {
	Convey("wxpay: unified order is failed!", t, func() {
		req := &UnifiedOrderRequest{
			Request: Request{
				AppId:    "wx00858d13475e52d",
				MchId:    "12862345411",
				NonceStr: strconv.Itoa(int(time.Now().UnixNano())),
				SignType: "MD5",
				Sign:     "9A0A8659F005D6984697E2CA0A9CF3B7",
			},
			Body:           "共享停车订单支付",
			NotifyUrl:      "http://www.test.com",
			TradeType:      "APP",
			SpbillCreateIP: "202.105.107.18",
			TotalFee:       10,
			OutTradeNo:     "test-111111125",
			Attach:         "{\"content\": \"test\"}",
		}
		resp, err := UnifiedOrder(req, "Gd7FoZlNN165456434354dsadwsvPL")
		So(err, ShouldBeNil)
		So(resp.ReturnCode, ShouldEqual, "SUCCESS")
		So(resp.ResultCode, ShouldEqual, "SUCCESS")
		So(resp.PrepayId, ShouldNotBeNil)
	})
}
