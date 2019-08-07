package wxpay

import (
	"crypto/tls"
	"encoding/pem"

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
