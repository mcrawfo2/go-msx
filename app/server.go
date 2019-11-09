package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/healthprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/infoprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/jwtprovider"
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
		OnEvent(EventStart, PhaseBefore, registerJwtSecurityProvider)
		OnEvent(EventStart, PhaseBefore, registerHealthWebService)
		OnEvent(EventStart, PhaseBefore, registerInfoWebService)
		OnEvent(EventStart, PhaseBefore, registerSwaggerWebService)
		OnEvent(EventStart, PhaseAfter, webservice.Start)
		OnEvent(EventStop, PhaseBefore, webservice.Stop)
	}

	return nil
}

func registerJwtSecurityProvider(ctx context.Context) error {
	logger.Info("Registering JWT web security provider")
	return jwtprovider.RegisterSecurityProvider(ctx)
}

func registerHealthWebService(ctx context.Context) error {
	logger.Info("Registering health web service")
	return healthprovider.RegisterHealthProvider(ctx)
}

func registerInfoWebService(ctx context.Context) error {
	logger.Info("Registering info web service")
	return infoprovider.RegisterInfoProvider(ctx)
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
