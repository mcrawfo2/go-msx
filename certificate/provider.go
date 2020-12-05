//go:generate mockery --inpackage --name=Provider --structname=mockProvider --filename mock_provider_test.go
//go:generate mockery --inpackage --name=ProviderFactory --structname=mockProviderFactory --filename mock_provider_factory_test.go
package certificate

import (
	"context"
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/pkg/errors"
)

var ErrNoSuchProvider = errors.New("No such provider registered")

type Provider interface {
	GetCertificate(ctx context.Context) (*tls.Certificate, error)
	Renewable() bool
}

type ProviderFactory interface {
	Name() string
	New(ctx context.Context, configRoot string) (Provider, error)
}

var factories = map[string]ProviderFactory{}

func RegisterProviderFactory(factory ProviderFactory) {
	factories[factory.Name()] = factory
}

func NewProvider(ctx context.Context, sourceName string) (Provider, error) {
	configRoot := "certificate.source." + sourceName
	configProvider := configRoot + ".provider"

	// Select which provider is configured
	providerName, err := config.MustFromContext(ctx).StringOr(configProvider, "file")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read source provider configuration")
	}

	logger.
		WithContext(ctx).
		Infof("Creating certificate source %q using provider %q", sourceName, providerName)

	// Find the provider factory
	providerFactory, ok := factories[providerName]
	if !ok {
		return nil, errors.Wrapf(ErrNoSuchProvider, "Provider %q not found", providerName)
	}

	// Create a new provider
	provider, err := providerFactory.New(ctx, configRoot)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create certificate source %q", sourceName)
	}

	return provider, nil
}
