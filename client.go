package wxpay

import (
	"errors"
	"net/http"
)

const (
	kSandboxURL    = "https://api.mch.weixin.qq.com/sandboxnew"
	kProductionURL = "https://api.mch.weixin.qq.com"

	NOT_FOUND_CERT_FILE  = "wxpay: not found cert file"
	NOT_FOUND_TLS_CLIENT = "wxpay: not found tls client"
)

type Client struct {
	ApiDomain    string // api.mch.weixin.qq.com
	IsProduction bool
	HttpClient   *http.Client
	CertPath     string
}

func New(isProd bool) *Client {
	client := &Client{}

	client.IsProduction = isProd
	client.HttpClient = http.DefaultClient

	if isProd {
		client.ApiDomain = kSandboxURL
	} else {
		client.ApiDomain = kProductionURL
	}

	return client
}

func NewTls(certPath string, isProd bool) (*Client, error) {

}
