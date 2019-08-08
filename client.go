package wxpay

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	SANDBOX_URL = "https://api.mch.weixin.qq.com/sandboxnew"
	PROD_URL    = "https://api.mch.weixin.qq.com"

	NOT_FOUND_CERT_FILE  = "wxpay: not found cert file"
	NOT_FOUND_TLS_CLIENT = "wxpay: not found tls client"
	MISSING_SIGN         = "wxpay: missing required (sign)"

	RETURN_FAIL    = "FAIL"
	RETURN_SUCCESS = "SUCCESS"

	SIGN_MD5         = "MD5"
	SIGN_HMAC_SHA256 = "HMAC-SHA256"
)

var (
	lock      = &sync.Mutex{}
	client    *Client
	tlsClient *Client
)

type Client struct {
	ApiDomain    string // api.mch.weixin.qq.com
	IsProduction bool
	HttpClient   *http.Client
	CertPath     string
}

func newClient(isProd bool) *Client {
	client := &Client{}

	client.IsProduction = isProd
	client.HttpClient = http.DefaultClient

	if isProd {
		client.ApiDomain = SANDBOX_URL
	} else {
		client.ApiDomain = PROD_URL
	}

	return client
}

func newTlsClient(mchId, certPath string, isProd bool) (client *Client, err error) {
	if len(certPath) == 0 {
		return client, errors.New(NOT_FOUND_CERT_FILE)
	}

	var p12Cert []byte
	p12Cert, err = ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	var cert tls.Certificate
	cert, err = Pkcs12ToPem(p12Cert, mchId)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	transport := &http.Transport{
		TLSClientConfig:    config,
		DisableCompression: true,
	}

	client = &Client{}
	client.IsProduction = isProd

	if isProd {
		client.ApiDomain = SANDBOX_URL
	} else {
		client.ApiDomain = PROD_URL
	}
	client.HttpClient = &http.Client{
		Transport: transport,
	}

	return client, nil
}

func (self *Client) DoRequest(method, api string, wxpayReq WxpayRequest, reply interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
				return
			}

			if str, ok := r.(string); ok {
				err = errors.New(str)
			}
		}
	}()

	var payload string
	payload, err = wxpayReq.Payload()
	if err != nil {
		return err
	}

	var request *http.Request
	request, err = http.NewRequest(method, self.apiUrl(api), bytes.NewBufferString(payload))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/xml;charset=utf-8")
	request.Header.Set("Accept", "application/xml")

	var resp *http.Response
	resp, err = self.HttpClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return xml.Unmarshal(data, reply)
}

func (self *Client) apiUrl(api string) string {
	return fmt.Sprintf("%s%s", self.ApiDomain, api)
}

func InitClient(isProd bool, mchId, certPath string) (client *Client, err error) {
	lock.Lock()
	defer lock.Unlock()

	if len(certPath) == 0 {
		client = newClient(isProd)
		return
	}

	client, err = newTlsClient(mchId, certPath, isProd)
	return
}
