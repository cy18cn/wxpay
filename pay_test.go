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
		Attach:         "{\"content\": \"test\"}",
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
				AppId:    "wxxxx",
				MchId:    "128633333",
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
		resp, err := UnifiedOrder(req, "Gd7FoZlNNqMhiolnfe4554dsfdsvPL")
		So(err, ShouldBeNil)
		So(resp.ReturnCode, ShouldEqual, "SUCCESS")
		So(resp.ResultCode, ShouldEqual, "SUCCESS")
		So(resp.PrepayId, ShouldNotBeNil)
	})
}

/**
type slice struct {
	array unsafe.Pointer
	len int
	cap int
}

slice可以修改array中的值，因为array是数组的指针
len和cap不能被修改，只是slice的一个copy
*/
func TestSlice(t *testing.T) {
	var s1 []int
	fmt.Printf("%p\n", &s1)
	AppendSlice(&s1)

	fmt.Println(s1)
	s2 := []int{3, 4, 5}
	fmt.Printf("%p\n", &s2)
	AppendSlice(&s2)
	fmt.Println(s2)
	ChangeSlice(s2)
	fmt.Println(s2)
}

func ChangeSlice(s []int) {
	s[0] -= 1
}

func AppendSlice(s *[]int) {
	*s = append(*s, 1)
	fmt.Printf("%p\n", &s)
}
