package wxpay

import "net/http"

type Client struct {
	ApiDomain    string // api.mch.weixin.qq.com
	IsProduction bool
	NotifyUrl    string
	HttpClient   *http.Client
}
