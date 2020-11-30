package fileprovider

import (
	"context"
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
)

var logger = log.NewLogger("msx.certificate.fileprovider")

type ProviderConfig struct {
	// CertFile and KeyFile represent the full path to the TLS certificate and Key in pem format
	CertFile string `config:"default=server.crt"`
	KeyFile  string `config:"default=server.key"`
}

// Provider implements the cert provider interface to return certs from a file
type Provider struct {
	cfg ProviderConfig
}

func (f Provider) GetCertificate(ctx context.Context) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(f.cfg.CertFile, f.cfg.KeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load certificate/key from file provider")
	}

	return &cert, nil
}

func (f Provider) Renewable() bool {
	return false
}

type ProviderFactory struct{}

func (p ProviderFactory) Name() string {
	return "file"
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
	logger.Info("Registering file certificate provider")

	certificate.RegisterProviderFactory(new(ProviderFactory))
	return nil
}
