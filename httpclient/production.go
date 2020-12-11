package httpclient

import (
	"context"
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"net/http"
)

type ProductionHttpClientFactory struct {
	tlsConfig    *tls.Config
	clientConfig *ClientConfig
	configurer   ClientConfigurer
}

func (f *ProductionHttpClientFactory) NewHttpClient() *http.Client {
	return f.NewHttpClientWithConfigurer(context.Background(), nil)
}

func (f *ProductionHttpClientFactory) NewHttpClientWithConfigurer(ctx context.Context, configurer Configurer) *http.Client {
	contextConfigurer := ConfigurerFromContext(ctx)

	var tlsConfig = &tls.Config{
		InsecureSkipVerify: f.tlsConfig.InsecureSkipVerify,
		RootCAs:            f.tlsConfig.RootCAs,
		Certificates:       f.tlsConfig.Certificates[:],
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		IdleConnTimeout: f.clientConfig.IdleTimeout,
		Proxy:           http.ProxyFromEnvironment,
	}

	f.configurer.HttpTransport(transport)
	if contextConfigurer != nil {
		contextConfigurer.HttpTransport(transport)
	}
	if configurer != nil {
		configurer.HttpTransport(transport)
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   f.clientConfig.Timeout,
	}

	f.configurer.HttpClient(client)
	if contextConfigurer != nil {
		contextConfigurer.HttpClient(client)
	}
	if configurer != nil {
		configurer.HttpClient(client)
	}

	return client
}

func (f *ProductionHttpClientFactory) AddClientConfigurationFunc(fn ClientConfigurationFunc) {
	f.configurer.ClientFuncs = append(f.configurer.ClientFuncs, fn)
}

func (f *ProductionHttpClientFactory) AddTransportConfigurationFunc(fn TransportConfigurationFunc) {
	f.configurer.TransportFuncs = append(f.configurer.TransportFuncs, fn)
}

func NewProductionHttpClientFactoryFromConfig(cfg *config.Config) (*ProductionHttpClientFactory, error) {
	clientConfig, err := NewClientConfig(cfg)
	if err != nil {
		return nil, err
	}

	tlsConfig, err := NewTlsConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return &ProductionHttpClientFactory{
		tlsConfig:    tlsConfig,
		clientConfig: clientConfig,
	}, nil
}

func NewProductionHttpClientFactory(ctx context.Context) (*ProductionHttpClientFactory, error) {
	return NewProductionHttpClientFactoryFromConfig(config.MustFromContext(ctx))
}
