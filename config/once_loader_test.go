package config

import (
	"context"
	"testing"
)

func TestOnceLoader(t *testing.T) {
	static := NewStatic("once", map[string]string{"enabled": "true"})
	loader := NewOnceLoader(static)

	settings, err := loader.Load(context.Background())
	if err != nil {
		t.Error(err)
	}

	if settings["enabled"] != "true" {
		t.Errorf("Enabled was '%s', expected 'true'", settings["enabled"])
	}

	static.Set("enabled", "false")

	settings, err = loader.Load(context.Background())
	if err != nil {
		t.Error(err)
	}

	if settings["enabled"] != "true" {
		t.Errorf("Enabled was '%s', expected 'true'", settings["enabled"])
	}
}
