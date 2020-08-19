package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/aliveprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/apilistprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/authprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/debugprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/envprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/healthprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/infoprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/loggersprovider"
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
		OnEvent(EventStart, PhaseBefore, registerDebugWebServices)
		OnEvent(EventStart, PhaseBefore, registerSwaggerWebService)
		OnEvent(EventStart, PhaseBefore, registerApiListWebService)
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
		aliveprovider.RegisterProvider(ctx),
		metricsprovider.RegisterProvider(ctx),
		envprovider.RegisterProvider(ctx),
		loggersprovider.RegisterProvider(ctx),
	}
	return err.Filter()
}

func registerDebugWebServices(ctx context.Context) error {
	debugEnabled, _ := config.FromContext(ctx).BoolOr("server.debug-enabled", false)
	if !debugEnabled {
		logger.Info("Debug endpoints disabled")
		return nil
	}

	logger.Info("Registering debug endpoints")
	return types.ErrorList{
		debugprovider.RegisterProvider(ctx),
	}.Filter()
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

func registerApiListWebService(ctx context.Context) error {
	logger.Info("Registering apilist documentation provider")
	if err := apilistprovider.RegisterProvider(ctx); err != nil {
		return err
	}

	return nil
}
