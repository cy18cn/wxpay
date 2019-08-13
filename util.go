package wxpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

func Pkcs12ToPem(p12 []byte, password string) (cert tls.Certificate, err error) {
	var pemBytes []*pem.Block
	pemBytes, err = pkcs12.ToPEM(p12, password)

	if err != nil {
		return cert, err
	}

	var pemData []byte
	for _, b := range pemBytes {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err = tls.X509KeyPair(pemData, pemData)
	return cert, err
}

func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha256(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func SignHmacSha256(ua url.Values, key string) string {
	return strings.ToUpper(HmacSha256(BuildSignStr(ua), key))
}

func SignMD5(ua url.Values, key string) string {
	return strings.ToUpper(MD5(BuildSignStr(ua) + fmt.Sprintf("&key=%s", key)))
}

func BuildSignStr(ua url.Values) string {
	var noNil []string
	for k := range ua {
		if k == "sign" {
			continue
		}
		v := ua.Get(k)
		if len(v) > 0 {
			noNil = append(noNil, k)
		}
	}

	sort.Strings(noNil)

	var buf strings.Builder
	for _, k := range noNil {
		if buf.Len() > 0 {
			buf.WriteString("&")
		}

		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(ua.Get(k))
	}

	return buf.String()
}

func VerifySign(ua url.Values, key, sign, signType string) bool {
	if signType == SIGN_MD5 {
		return SignMD5(ua, key) == sign
	} else {
		return SignHmacSha256(ua, key) == sign
	}
}

func VerifyResp(resp map[string]interface{}, key, sign, signType string) bool {
	ua := url.Values{}
	for k, v := range resp {
		ua.Set(k, v.(string))
	}

	return VerifySign(ua, key, sign, signType)
}

func NonceStr() string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	rdn := rand.New(rand.NewSource(time.Now().UnixNano()))

	var buf strings.Builder
	for i := 0; i < 32; i++ {
		idx := rdn.Intn(len(chars) - 1)
		buf.WriteString(chars[idx : idx+1])
	}

	return buf.String()
}

func XmlDecoder() {
	b := []byte(`<xml>
<appid><![CDATA[wxxxx]]></appid>
<cash_fee><![CDATA[1]]></cash_fee>
<mch_id><![CDATA[xxx]]></mch_id>
<nonce_str><![CDATA[KptCSXZh1qBjK8wb]]></nonce_str>
<out_refund_no_0><![CDATA[1519535580802]]></out_refund_no_0>
<out_trade_no><![CDATA[1515845954891]]></out_trade_no>
<refund_account_0><![CDATA[REFUND_SOURCE_UNSETTLED_FUNDS]]></refund_account_0>
<refund_channel_0><![CDATA[ORIGINAL]]></refund_channel_0>
<refund_count>1</refund_count>
<refund_fee>1</refund_fee>
<refund_fee_0>1</refund_fee_0>
<refund_id_0><![CDATA[50000405502018022503594217469]]></refund_id_0>
<refund_recv_accout_0><![CDATA[支付用户的零钱]]></refund_recv_accout_0>
<refund_status_0><![CDATA[SUCCESS]]></refund_status_0>
<refund_success_time_0><![CDATA[2018-02-25 13:13:03]]></refund_success_time_0>
<result_code><![CDATA[SUCCESS]]></result_code>
<return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
<sign><![CDATA[E8A30F02296C6169860A92C2D52AD5A8]]></sign>
<total_fee><![CDATA[1]]></total_fee>
<transaction_id><![CDATA[4200000100201801133414066940]]></transaction_id>
</xml>`)
	decoder := xml.NewDecoder(bytes.NewBuffer(b))
	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}

		switch ele := t.(type) {
		case xml.StartElement:
			fmt.Println(xml.StartElement(ele).Name)
		case xml.CharData:
			fmt.Println(string(xml.CharData(ele)))
		case xml.Name:
			fmt.Println(xml.Name(ele).Local)
		}
	}
}
