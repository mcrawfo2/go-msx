package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"time"
)

const configRootHttpClient = "http.client"

type ClientConfig struct {
	Timeout         time.Duration `config:"default=30s"`
	IdleTimeout     time.Duration `config:"default=1s"`
	LocalCaFile     string        `config:"default="`
	CertFile        string        `config:"default="`
	KeyFile         string        `config:"default="`
	TlsInsecure     bool          `config:"default=true"`
	TlsMinVersion   string        `config:"default=tls10"`
	TlsCipherSuites []string      `config:"default="`
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

	ciphers, err := ParseCiphers(clientConfig.TlsCipherSuites)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: clientConfig.TlsInsecure,
		RootCAs:            rootCAs,
		Certificates:       clientCerts,
		MinVersion:         TLSLookup[clientConfig.TlsMinVersion],
		CipherSuites:       ciphers,
	}

	if len(clientCerts) > 0 {
		tlsConfig.BuildNameToCertificate()
	}

	return tlsConfig, nil
}

// ParseCiphers parse ciphersuites from the comma-separated string into
// recognized slice
func ParseCiphers(ciphers []string) ([]uint16, error) {
	var suites []uint16

	for _, cipher := range ciphers {
		if v, ok := CipherLookup[cipher]; ok {
			suites = append(suites, v)
		} else {
			return nil, fmt.Errorf("unsupported cipher %q", cipher)
		}
	}

	return suites, nil
}

// CipherLookup maps the cipher names to the internal value
var CipherLookup = map[string]uint16{
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":          tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":       tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":       tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256":       tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":          tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":          tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"TLS_RSA_WITH_AES_128_GCM_SHA256":               tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_RSA_WITH_AES_256_GCM_SHA384":               tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_RSA_WITH_AES_128_CBC_SHA256":               tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_RSA_WITH_AES_128_CBC_SHA":                  tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"TLS_RSA_WITH_AES_256_CBC_SHA":                  tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA":                 tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_RSA_WITH_RC4_128_SHA":                      tls.TLS_RSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_RSA_WITH_RC4_128_SHA":                tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":              tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256":   tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	"TLS_AES_128_GCM_SHA256":                        tls.TLS_AES_128_GCM_SHA256,
	"TLS_AES_256_GCM_SHA384":                        tls.TLS_AES_256_GCM_SHA384,
	"TLS_CHACHA20_POLY1305_SHA256":                  tls.TLS_CHACHA20_POLY1305_SHA256,
}

// TLSLookup maps the tls.MinVersion configuration to the internal value
var TLSLookup = map[string]uint16{
	"":      tls.VersionTLS10, // default in golang
	"tls10": tls.VersionTLS10,
	"tls11": tls.VersionTLS11,
	"tls12": tls.VersionTLS12,
	"tls13": tls.VersionTLS13,
}
