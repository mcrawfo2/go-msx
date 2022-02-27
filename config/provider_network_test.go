package config

import (
	"context"
	"github.com/jackpal/gateway"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkProvider(t *testing.T) {
	defaultGateway, _ := gateway.DiscoverGateway()

	provider := NewNetworkProvider()

	entries, err := provider.Load(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, entries)

	settings := MapFromEntries(entries)
	assert.Equal(t, defaultGateway.String(), settings["network.outbound.gateway"])
}
