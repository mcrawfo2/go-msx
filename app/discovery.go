package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery/consulprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, registerDiscoveryProviders)
	OnEvent(EventReady, PhaseBefore, registerServiceInstance)
}

func registerDiscoveryProviders() {
	registrationProvider, err := consulprovider.NewRegistrationProviderFromConfig(Config())
	if err != nil {
		Shutdown()
		logger.Error(err)
	}

	discoveryProvider, err := consulprovider.NewDiscoveryProviderFromConfig(Config())
	if err != nil {
		Shutdown()
		logger.Error(err)
	}

	discovery.RegisterRegistrationProvider(registrationProvider)
	discovery.RegisterDiscoveryProvider(discoveryProvider)
}

func registerServiceInstance() {
	if err := discovery.Register(); err != nil {
		Shutdown()
		logger.Error(err)
	}
}
