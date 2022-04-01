// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package certificate

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
)

// TLSConfig represents the configuration to be applied to a secure listener.
type TLSConfig struct {
	// Flags enables TLS
	Enabled bool `config:"default=false"`

	// InsecureSkipVerify disables cert validation
	InsecureSkipVerify bool `config:"default=false"`

	// MinVersion defines minimum supported TLS.  Should be one of:
	// tls10, tls11, tls12, tls13
	MinVersion string `config:"default=tls12"`

	// CertificateSource specifies the name of the certificate source binding
	CertificateSource string `config:"default=server"`

	// CipherSuites is a comma separated list of desired ciphersuites to use for secure connection
	// Default list is reasonable minimum as required by PSB
	CipherSuites []string `config:"default=TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305;TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256;TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384;TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA;TLS_RSA_WITH_AES_256_GCM_SHA384;TLS_RSA_WITH_AES_256_CBC_SHA"`

	// ServerName is used to verify the hostname on the returned
	// certificates unless InsecureSkipVerify is given. It is also included
	// in the client's handshake to support virtual hosting unless it is
	// an IP address.
	ServerName string `config:"default="`
}

func (cfg *TLSConfig) TlsConfig(ctx context.Context) (result *tls.Config, err error) {
	w, err := NewSource(ctx, cfg.CertificateSource)
	if err != nil {
		return nil, err
	}

	caCert, err := w.Provider().GetCaCertificate(ctx)
	if err != nil {
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCert)

	ciphers, err := cfg.Ciphers()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert ciphers")
	}

	result = &tls.Config{
		ClientAuth:           tls.VerifyClientCertIfGiven,
		ClientCAs:            caCertPool,
		MinVersion:           tlsVersions[cfg.MinVersion],
		CipherSuites:         ciphers,
		RootCAs:              caCertPool,
		GetCertificate:       w.TlsCertificate,
		GetClientCertificate: w.TlsClientCertificate,
		InsecureSkipVerify:   cfg.InsecureSkipVerify,
		ServerName:           cfg.ServerName,
	}

	return
}

// Ciphers returns the list of ciphers specified in the config mapped to internal values
func (cfg *TLSConfig) Ciphers() ([]uint16, error) {
	var suites []uint16

	for _, cipher := range cfg.CipherSuites {
		if v, ok := tlsCiphers[cipher]; ok {
			suites = append(suites, v)
		} else {
			return nil, fmt.Errorf("unsupported cipher %q", cipher)
		}
	}

	return suites, nil
}

// tlsCiphers maps the cipher names to the internal value
var tlsCiphers = map[string]uint16{
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

// tlsVersions maps the tls.MinVersion configuration to the internal value
var tlsVersions = map[string]uint16{
	"":      tls.VersionTLS10, // default in golang
	"tls10": tls.VersionTLS10,
	"tls11": tls.VersionTLS11,
	"tls12": tls.VersionTLS12,
	"tls13": tls.VersionTLS13,
}
