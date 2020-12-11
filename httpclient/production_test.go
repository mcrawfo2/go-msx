package httpclient

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestProductionHttpClientFactory_AddClientConfigurationFunc(t *testing.T) {
	cfg := config.NewConfig()
	_ = cfg.Load(context.Background())
	ctx := config.ContextWithConfig(context.Background(), cfg)
	factory, _ := NewProductionHttpClientFactory(ctx)

	clientFunc := func(c *http.Client) {}
	factory.AddClientConfigurationFunc(clientFunc)
	assert.Len(t, factory.configurer.ClientFuncs, 1)
}

func TestProductionHttpClientFactory_AddTransportConfigurationFunc(t *testing.T) {
	cfg := config.NewConfig()
	_ = cfg.Load(context.Background())
	ctx := config.ContextWithConfig(context.Background(), cfg)
	factory, _ := NewProductionHttpClientFactory(ctx)

	transportFunc := func(c *http.Transport) {}
	factory.AddTransportConfigurationFunc(transportFunc)
	assert.Len(t, factory.configurer.TransportFuncs, 1)
}
