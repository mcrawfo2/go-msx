package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestProviderSources(t *testing.T) {
	commandDir, _ := os.Getwd()
	parentDir := filepath.Dir(commandDir)
	p := NewSourcesProvider("Sources")
	expected := ProviderEntries{
		NewEntry(p, "fs.roots.command", commandDir),
		NewEntry(p, "fs.sources", parentDir),
	}
	entries, err := p.Load(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expected, entries)
}
