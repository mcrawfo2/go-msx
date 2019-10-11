package discovery

import "github.com/pkg/errors"

type RegistrationProvider interface {
	Register() error
}

type DiscoveryProvider interface {
	Discover(name string, healthyOnly bool, tags ...string) (ServiceInstances, error)
}

var (
	registrationProvider RegistrationProvider
	discoveryProvider    DiscoveryProvider

	ErrRegistrationProviderNotDefined = errors.New("Registration provider not registered")
	ErrDiscoveryProviderNotDefined    = errors.New("Discovery provider not registered")
)

func RegisterDiscoveryProvider(provider DiscoveryProvider) {
	discoveryProvider = provider
}

func Discover(name string, healthyOnly bool, tags ...string) (ServiceInstances, error) {
	if discoveryProvider == nil {
		return nil, ErrDiscoveryProviderNotDefined
	}

	return discoveryProvider.Discover(name, healthyOnly, tags...)
}

func RegisterRegistrationProvider(provider RegistrationProvider) {
	registrationProvider = provider
}

func Register() error {
	if registrationProvider == nil {
		return ErrRegistrationProviderNotDefined
	}
	return registrationProvider.Register()
}

func Deregister() error {
	return nil
}
