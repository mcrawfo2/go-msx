package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/healthprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/infoprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/swaggerprovider"
)

func init() {
	OnEvent(EventStart, PhaseBefore, registerHealthWebService)
	OnEvent(EventStart, PhaseBefore, registerInfoWebService)
	OnEvent(EventStart, PhaseBefore, registerSwaggerWebService)
	OnEvent(EventStart, PhaseAfter, webservice.Start)
	OnEvent(EventStop, PhaseBefore, webservice.Stop)
}

func registerHealthWebService(context.Context) error {
	logger.Info("Registering health web service")
	healthprovider.RegisterHealthProvider()
	return nil
}


func registerInfoWebService(context.Context) error {
	logger.Info("Registering info web service")
	infoprovider.RegisterInfoProvider()
	return nil
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
