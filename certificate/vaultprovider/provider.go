// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vaultprovider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"github.com/pkg/errors"
	"time"
)

var logger = log.NewLogger("msx.certificate.vaultprovider")

// Provider implements a cert provider interface for Vault
type Provider struct {
	cfg ProviderConfig
}

// ProviderConfig defines the settings used to interact with Vault PKI
type ProviderConfig struct {
	Path     string        `config:"default="`          //Path is for using a non-default PKI provider path in vault
	Role     string        ``                           //Role is the Vault role with permissions to access PKI.
	TTL      time.Duration `config:"default=730h"`      //TTL sets requested time to live for certs. can't be longer than PKI's configured default
	CN       string        ``                           //CN for the certificate request, must be allowed by the role.
	AltNames []string      `config:"default=localhost"` //AltNames is a comma separated list of altnames for the certificate request, optional
	IPSans   []string      `config:"default=127.0.0.1"` //IPSans is a comma separated list of ipsans for the certificate request, optional
}

func (c ProviderConfig) IPSANS() []string {
	var results []string
	for _, v := range c.IPSans {
		if v != "" {
			results = append(results, v)
		}
	}
	return results
}

func (c ProviderConfig) SANS() []string {
	var results []string
	for _, v := range c.AltNames {
		if v != "" {
			results = append(results, v)
		}
	}
	return results
}

func (f Provider) GetCertificate(ctx context.Context) (*tls.Certificate, error) {
	request := vault.IssueCertificateRequest{
		CommonName: f.cfg.CN,
		Ttl:        f.cfg.TTL,
		AltNames:   f.cfg.SANS(),
		IpSans:     f.cfg.IPSANS(),
	}

	cert, _, err := vault.ConnectionFromContext(ctx).IssueCustomCertificate(ctx, f.cfg.Path, f.cfg.Role, request)
	return cert, err
}

func (f Provider) GetCaCertificate(ctx context.Context) (*x509.Certificate, error) {
	return vault.ConnectionFromContext(ctx).ReadCustomCaCertificate(ctx, f.cfg.Path)
}

func (f Provider) Renewable() bool {
	return true
}

type ProviderFactory struct{}

func (p ProviderFactory) Name() string {
	return "vault"
}

func (p ProviderFactory) New(ctx context.Context, configRoot string) (certificate.Provider, error) {
	var cfg ProviderConfig
	if err := config.MustFromContext(ctx).Populate(&cfg, configRoot); err != nil {
		return nil, errors.Wrapf(err, "Failed to load certificate source configuration for %q provider", p.Name())
	}

	return &Provider{
		cfg: cfg,
	}, nil
}

func RegisterFactory(ctx context.Context) error {
	vaultEnabled, _ := config.MustFromContext(ctx).BoolOr("spring.cloud.vault.enabled", false)
	if !vaultEnabled {
		logger.WithContext(ctx).Warn("Not registering vault certificate provider: vault is disabled")
		return nil
	}

	logger.WithContext(ctx).Info("Registering vault certificate provider")
	certificate.RegisterProviderFactory(new(ProviderFactory))
	return nil
}
