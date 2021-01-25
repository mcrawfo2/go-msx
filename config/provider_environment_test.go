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
