package httpclient

import (
	"crypto/tls"
	"net/http"
)

type ProductionHttpClientFactory struct {
	tlsConfig *tls.Config
}

var tlsConfig = &tls.Config{
	InsecureSkipVerify: true,
	ClientAuth:         tls.VerifyClientCertIfGiven,
}

func NewProductionHttpClientFactory() *ProductionHttpClientFactory {
	tlsConfig.BuildNameToCertificate()
	return &ProductionHttpClientFactory{
		tlsConfig: tlsConfig,
	}
}

func (f *ProductionHttpClientFactory) NewHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
}
