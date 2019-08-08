package wxpay

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/pem"
	"fmt"
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

func Sign(ua url.Values, key string) string {
	var noNil []string
	for k := range ua {
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

	if ua.Get("sign_type") == SIGN_MD5 {
		return strings.ToUpper(MD5(buf.String() + fmt.Sprintf("&key=%s", key)))
	} else {
		return strings.ToUpper(HmacSha256(buf.String(), key))
	}
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
