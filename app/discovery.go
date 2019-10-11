package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery/consulprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, registerDiscoveryProviders)
	OnEvent(EventReady, PhaseBefore, registerServiceInstance)
	OnEvent(EventStop, PhaseBefore, deregisterServiceInstance)
}

func registerDiscoveryProviders() {
	logger.Info("Registering consul registration provider")
	registrationProvider, err := consulprovider.NewRegistrationProviderFromConfig(Config())
	if err == consulprovider.ErrDisabled {
		logger.Error(err)
	} else if err != nil {
		Shutdown()
		logger.Error(err)
	} else if registrationProvider != nil {
		discovery.RegisterRegistrationProvider(registrationProvider)
	}

	logger.Info("Registering consul discovery provider")
	discoveryProvider, err := consulprovider.NewDiscoveryProviderFromConfig(Config())
	if err == consulprovider.ErrDisabled {
		logger.Error(err)
	} else if err != nil {
		Shutdown()
		logger.Error(err)
	} else if discoveryProvider != nil {
		discovery.RegisterDiscoveryProvider(discoveryProvider)
	}
}

func registerServiceInstance() {
	if err := discovery.Register(Context()); err != nil && err != discovery.ErrRegistrationProviderNotDefined {
		Shutdown()
		logger.Error(err)
	}
}

func deregisterServiceInstance() {
	if err := discovery.Deregister(Context()); err != nil && err != discovery.ErrRegistrationProviderNotDefined {
		logger.Error(err)
	}
}