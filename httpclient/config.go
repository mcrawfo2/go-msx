package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/pkg/errors"
	"io/ioutil"
	"time"
)

const configRootHttpClient = "http.client"

type ClientConfig struct {
	Timeout     time.Duration `config:"default=30s"`
	IdleTimeout time.Duration `config:"default=1s"`
	TlsInsecure bool          `config:"default=true"`
	LocalCaFile string        `config:"default="`
	CertFile    string        `config:"default="`
	KeyFile     string        `config:"default="`
}

func NewClientConfig(cfg *config.Config) (*ClientConfig, error) {
	var clientConfig ClientConfig
	if err := cfg.Populate(&clientConfig, configRootHttpClient); err != nil {
		return nil, err
	}

	return &clientConfig, nil
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

func NewTlsConfig(clientConfig *ClientConfig) (*tls.Config, error) {
	rootCAs, err := getRootCAs(clientConfig)
	if err != nil {
		return nil, err
	}

	clientCerts, err := getClientCerts(clientConfig)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: clientConfig.TlsInsecure,
		RootCAs:            rootCAs,
		Certificates:       clientCerts,
	}

	if len(clientCerts) > 0 {
		tlsConfig.BuildNameToCertificate()
	}

	return tlsConfig, nil
}
