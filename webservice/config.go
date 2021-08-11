package webservice

import (
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"fmt"
	"strconv"
)

type WebServerConfig struct {
	Enabled      bool   `config:"default=false"`
	Host         string `config:"default=0.0.0.0"`
	Port         int    `config:"default=8080"`
	Tls          TLSConfig
	Cors         CorsConfig
	ContextPath  string `config:"default=/app"`
	StaticPath   string `config:"default=/www"`
	TraceEnabled bool   `config:"default=false"`
	DebugEnabled bool   `config:"default=false"`
}

type CorsConfig struct {
	Enabled              bool     `config:"default=true"`
	CustomAllowedHeaders []string `config:"default=${security.cors.allowedHeaders}"`
	CustomExposedHeaders []string `config:"default=${security.cors.exposedHeaders}"`
}

//TLSConfig represents the configuration to be applied to a secure listener.
type TLSConfig struct {
	// Flags TLS on or off for webserver
	Enabled bool `config:"default=false"`
	// MinVersion defines minimum supported TLS.  Should be one of:
	// tls10, tls11, tls12, tls13
	MinVersion string `config:"default=tls12"`
	//CertProvider defines the type of certprovider in use File is default
	CertificateSource string `config:"default=server"`
	//CaFile represents the full path to the CA to be used for validating client certs used in mTLS authentication
	CaFile string `config:"default=ca.pem"`
	// CipherSuites is a comma separated list of desired ciphersuites to use for secure connection
	// Default list is reasonable minimum as required by PSB
	CipherSuites []string `config:"default=TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305;TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256;TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384;TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA;TLS_RSA_WITH_AES_256_GCM_SHA384;TLS_RSA_WITH_AES_256_CBC_SHA"`
}

func (c WebServerConfig) Address() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (c WebServerConfig) Url() string {
	if c.Tls.Enabled {
		return "https://" + c.Address() + c.ContextPath
	}
	return "http://" + c.Address() + c.ContextPath
}

func NewWebServerConfig(cfg *config.Config) (*WebServerConfig, error) {
	var webServerConfig WebServerConfig
	if err := cfg.Populate(&webServerConfig, configRootWebServer); err != nil {
		return nil, err
	}
	return &webServerConfig, nil
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
