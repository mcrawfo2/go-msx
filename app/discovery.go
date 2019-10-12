package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery/consulprovider"
	"github.com/pkg/errors"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, registerDiscoveryProviders)
	OnEvent(EventReady, PhaseBefore, registerServiceInstance)
	OnEvent(EventStop, PhaseBefore, deregisterServiceInstance)
}

func registerDiscoveryProviders(ctx context.Context) error {
	var cfg *config.Config
	if cfg = config.FromContext(ctx); cfg == nil {
		return errors.New("Failed to retrieve config from context")
	}

	logger.Info("Registering consul registration provider")
	registrationProvider, err := consulprovider.NewRegistrationProviderFromConfig(cfg)
	if err == consulprovider.ErrDisabled {
		logger.Error(err)
	} else if err != nil {
		return err
	} else if registrationProvider != nil {
		discovery.RegisterRegistrationProvider(registrationProvider)
	}

	logger.Info("Registering consul discovery provider")
	discoveryProvider, err := consulprovider.NewDiscoveryProviderFromConfig(cfg)
	if err == consulprovider.ErrDisabled {
		logger.Error(err)
	} else if err != nil {
		return err
	} else if discoveryProvider != nil {
		discovery.RegisterDiscoveryProvider(discoveryProvider)
	}

	return nil
}

func registerServiceInstance(ctx context.Context) error {
	if err := discovery.Register(ctx); err != nil && err != discovery.ErrRegistrationProviderNotDefined {
		return err
	}
	return nil
}

func deregisterServiceInstance(ctx context.Context) error {
	if err := discovery.Deregister(ctx); err != nil && err != discovery.ErrRegistrationProviderNotDefined {
		logger.Error(err)
	}
	return nil
}
