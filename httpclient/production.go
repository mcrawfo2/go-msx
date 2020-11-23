package httpclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

const configRootHttpClient = "http.client"

type ClientConfigurationFunc func(c *http.Client)
type TransportConfigurationFunc func(c *http.Transport)

func TlsInsecure(insecure bool) TransportConfigurationFunc {
	return func(c *http.Transport) {
		c.TLSClientConfig.InsecureSkipVerify = insecure
	}
}

type ClientConfig struct {
	Timeout     time.Duration `config:"default=30s"`
	IdleTimeout time.Duration `config:"default=1s"`
	TlsInsecure bool          `config:"default=true"`
	LocalCaFile string        `config:"default="`
	CertFile    string        `config:"default="`
	KeyFile     string        `config:"default="`
}

func getRootCAs(cfg *ClientConfig) (*x509.CertPool, error) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		logger.Warn("System certificate pool empty")
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	if cfg.LocalCaFile == "" {
		return rootCAs, nil
	}

	certs, err := ioutil.ReadFile(cfg.LocalCaFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to append %q to RootCAs", cfg.LocalCaFile)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		logger.Warn("No certs appended, using system certs only")
	} else {
		logger.Infof("Added certificates from %q to RootCAs", cfg.LocalCaFile)
	}

	return rootCAs, nil
}

func getClientCerts(cfg *ClientConfig) ([]tls.Certificate, error) {
	if cfg.CertFile == "" && cfg.KeyFile == "" {
		logger.Warn("TLS client certificate not specified.")
		return nil, nil
	} else if cfg.CertFile == "" || cfg.KeyFile == "" {
		return nil, errors.New("Must specify both TLS client cert and key files")
	}

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, err
	}
	logger.Infof("Loaded client certificate from %q", cfg.CertFile)
	return []tls.Certificate{cert}, nil
}

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
		RootCAs: f.tlsConfig.RootCAs,
		Certificates: f.tlsConfig.Certificates[:],
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
		Timeout: f.clientConfig.Timeout,
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
	var clientConfig ClientConfig
	if err := cfg.Populate(&clientConfig, configRootHttpClient); err != nil {
		return nil, err
	}

	rootCAs, err := getRootCAs(&clientConfig)
	if err != nil {
		return nil, err
	}

	clientCerts, err := getClientCerts(&clientConfig)
	if err != nil {
		return nil, err
	}

	var tlsConfig = &tls.Config{
		InsecureSkipVerify: clientConfig.TlsInsecure,
		RootCAs:            rootCAs,
		Certificates:       clientCerts,
	}

	if len(clientCerts) > 0 {
		tlsConfig.BuildNameToCertificate()
	}

	return &ProductionHttpClientFactory{
		tlsConfig:    tlsConfig,
		clientConfig: &clientConfig,
	}, nil
}

func NewProductionHttpClientFactory(ctx context.Context) (*ProductionHttpClientFactory, error) {
	return NewProductionHttpClientFactoryFromConfig(config.MustFromContext(ctx))
}

type Configurer interface {
	HttpClient(*http.Client)
	HttpTransport(*http.Transport)
}

type ClientConfigurer struct {
	ClientFuncs []ClientConfigurationFunc
	TransportFuncs []TransportConfigurationFunc
}

func (c ClientConfigurer) HttpClient(client *http.Client) {
	for _, configure := range c.ClientFuncs {
		configure(client)
	}
}

func (c ClientConfigurer) HttpTransport(transport *http.Transport) {
	for _, configure := range c.TransportFuncs {
		configure(transport)
	}
}

type CompositeConfigurer struct {
	Service Configurer
	Endpoint Configurer
}

func (c CompositeConfigurer) HttpClient(client *http.Client) {
	c.Service.HttpClient(client)
	c.Endpoint.HttpClient(client)
}

func (c CompositeConfigurer) HttpTransport(transport *http.Transport) {
	c.Service.HttpTransport(transport)
	c.Endpoint.HttpTransport(transport)
}
