// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEnvironmentProvider(t *testing.T) {
	err := os.Setenv("TESTING", "true")
	assert.NoError(t, err)
	provider := NewEnvironmentProvider("environment")

	entries, err := provider.Load(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, entries)

	settings := MapFromEntries(entries)
	assert.Equal(t, "true", settings["TESTING"])
}
