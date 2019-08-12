package wxpay

import (
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUnmarshalToMap(t *testing.T) {
	convey.Convey("Failed unmarshal xml to map", t, func() {
		b := []byte(`<?xml version="1.0"?><xml><appid><![CDATA[wxxxx]]></appid><case><cash_fee><![CDATA[1]]></cash_fee></case><mch_id><![CDATA[xxx]]></mch_id><nonce_str><![CDATA[KptCSXZh1qBjK8wb]]></nonce_str><out_refund_no_0><![CDATA[1519535580802]]></out_refund_no_0><out_trade_no><![CDATA[1515845954891]]></out_trade_no><refund><refund_account_0><![CDATA[REFUND_SOURCE_UNSETTLED_FUNDS]]></refund_account_0><refund_channel_0><![CDATA[ORIGINAL]]></refund_channel_0><refund_count>1</refund_count><refund_fee>1</refund_fee><other_refund><refund_fee_0>1</refund_fee_0><refund_id_0><![CDATA[50000405502018022503594217469]]></refund_id_0><refund_recv_accout_0><![CDATA[支付用户的零钱]]></refund_recv_accout_0><refund_status_0><![CDATA[SUCCESS]]></refund_status_0></other_refund><refund_success_time_0><![CDATA[2018-02-25 13:13:03]]></refund_success_time_0></refund><result_code><![CDATA[SUCCESS]]></result_code><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg><sign><![CDATA[E8A30F02296C6169860A92C2D52AD5A8]]></sign><total_fee><![CDATA[1]]></total_fee><transaction_id><![CDATA[4200000100201801133414066940]]></transaction_id></xml>`)
		m := UnmarshalToMap(b)

		convey.So(m["appid"], assertions.ShouldEqual, "wxxxx")
		convey.So((m["case"].(map[string]interface{}))["cash_fee"], assertions.ShouldEqual, "1")
		convey.So(((m["refund"].(map[string]interface{}))["other_refund"].(map[string]interface{}))["refund_fee_0"], assertions.ShouldEqual, "1")
	})
}
