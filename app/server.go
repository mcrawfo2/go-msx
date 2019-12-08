package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/authprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/envprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/healthprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/infoprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/metricsprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/swaggerprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, withConfig(registerRegistrations))
}

func registerRegistrations(cfg *config.Config) error {
	serverEnabled, err := cfg.BoolOr("server.enabled", true)
	if err != nil {
		return err
	}

	if serverEnabled {
		OnEvent(EventStart, PhaseBefore, registerAuthenticationProvider)
		OnEvent(EventStart, PhaseBefore, registerAdminWebServices)
		OnEvent(EventStart, PhaseBefore, registerSwaggerWebService)
		OnEvent(EventStart, PhaseAfter, webservice.Start)
		OnEvent(EventStop, PhaseBefore, webservice.Stop)
	}

	return nil
}

func registerAuthenticationProvider(ctx context.Context) error {
	logger.Info("Registering resource path glob security provider")
	return authprovider.RegisterAuthenticationProvider(ctx)
}

func registerAdminWebServices(ctx context.Context) error {
	logger.Info("Registering admin endpoints")
	err := types.ErrorList{
		adminprovider.RegisterProvider(ctx),
		healthprovider.RegisterProvider(ctx),
		infoprovider.RegisterProvider(ctx),
		metricsprovider.RegisterProvider(ctx),
		envprovider.RegisterProvider(ctx),
	}
	return err.Filter()
}

func registerSwaggerWebService(ctx context.Context) error {
	logger.Info("Registering swagger documentation provider")
	if err := swaggerprovider.RegisterSwaggerProvider(ctx); err != nil && err != swaggerprovider.ErrDisabled {
		return err
	} else if err == swaggerprovider.ErrDisabled {
		logger.Info("Swagger documentation provider disabled")
	}

	return nil
}
