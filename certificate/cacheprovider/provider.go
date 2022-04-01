// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cacheprovider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
	"io/ioutil"
)

var logger = log.NewLogger("msx.certificate.cacheprovider")

type ProviderConfig struct {
	// UpstreamSource identifies the source whose contents will be cached
	UpstreamSource string
	// CertFile and KeyFile represent the full path to the TLS certificate and Key in pem format
	CertFile string `config:"default=server.crt"`
	KeyFile  string `config:"default=server.key"`
	// CaCertFile is the full path to the x509 certificate in pem format
	CaCertFile string `config:"default=ca.crt"`
}

// Provider implements the cert provider interface to return certs from an upstream provider, with caching
type Provider struct {
	cfg      ProviderConfig
	upstream certificate.Provider
}

func (p Provider) GetCertificate(ctx context.Context) (*tls.Certificate, error) {
	cert, err := p.upstream.GetCertificate(ctx)
	if err != nil {
		return nil, err
	}

	// Store cert and private key
	err = p.writeIdentity(cert)
	if err != nil {
		return nil, err
	}

	return cert, err
}

func (p Provider) writeIdentity(cert *tls.Certificate) error {
	keyBytes, err := certificate.RenderPemPrivateKey(cert.PrivateKey)
	if err != nil {
		return err
	}

	client, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return err
	}

	certBytes, err := certificate.RenderPemCertificate(client)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(p.cfg.KeyFile, keyBytes, 0644); err != nil {
		return err
	}

	if err = ioutil.WriteFile(p.cfg.CertFile, certBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (p Provider) GetCaCertificate(ctx context.Context) (*x509.Certificate, error) {
	cert, err := p.upstream.GetCaCertificate(ctx)
	if err != nil {
		return nil, err
	}

	// Store cert to disk
	err = p.writeAuthority(cert)
	if err != nil {
		return nil, err
	}

	return cert, err
}

func (p Provider) writeAuthority(ca *x509.Certificate) error {
	certBytes, err := certificate.RenderPemCertificate(ca)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(p.cfg.CaCertFile, certBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (p Provider) Renewable() bool {
	return p.upstream.Renewable()
}

type ProviderFactory struct{}

func (p ProviderFactory) Name() string {
	return "cache"
}

func (p ProviderFactory) New(ctx context.Context, configRoot string) (result certificate.Provider, err error) {
	var cfg ProviderConfig
	if err = config.MustFromContext(ctx).Populate(&cfg, configRoot); err != nil {
		return nil, errors.Wrapf(err, "Failed to load certificate source configuration for %q provider", p.Name())
	}

	var upstream certificate.Provider
	if upstream, err = certificate.NewProvider(ctx, cfg.UpstreamSource); err != nil {
		return nil, errors.Wrapf(err, "Failed to create provider for %q source", cfg.UpstreamSource)
	}

	return &Provider{
		cfg:      cfg,
		upstream: upstream,
	}, nil
}

func RegisterFactory(ctx context.Context) error {
	logger.WithContext(ctx).Info("Registering cache certificate provider")
	certificate.RegisterProviderFactory(new(ProviderFactory))
	return nil
}
