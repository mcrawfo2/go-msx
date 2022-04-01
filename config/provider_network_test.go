// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
