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
	Role     string        ``                           //Role is the Vault role with permissions to access PKI.
	TTL      time.Duration `config:"default=730h"`      //TTL sets requested time to live for certs. can't be longer than PKI's configured default
	CN       string        ``                           //CN for the certificate request, must be allowed by the role.
	AltNames []string      `config:"default=localhost"` //AltNames is a comma separated list of altnames for the certificate request, optional
	IPSans   []string      `config:"default=127.0.0.1"` //IPSans is a comma separated list of ipsans for the certificate request, optional
}

func (f Provider) GetCertificate(ctx context.Context) (*tls.Certificate, error) {
	request := vault.IssueCertificateRequest{
		CommonName: f.cfg.CN,
		Ttl:        f.cfg.TTL,
		AltNames:   f.cfg.AltNames,
		IpSans:     f.cfg.IPSans,
	}

	return vault.ConnectionFromContext(ctx).IssueCertificate(ctx, f.cfg.Role, request)
}

func (f Provider) GetCaCertificate(ctx context.Context) (*x509.Certificate, error) {
	return vault.ConnectionFromContext(ctx).ReadCaCertificate(ctx)
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
		logger.Warn("Not registering vault certificate provider: vault is disabled")
		return nil
	}

	logger.Info("Registering vault certificate provider")
	certificate.RegisterProviderFactory(new(ProviderFactory))
	return nil
}
