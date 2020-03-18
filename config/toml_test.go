package config

import (
	"context"
	"testing"
)

func TestTOMLFile(t *testing.T) {
	configFile := "test/config.toml"

	tomlProvider := NewTOMLFile("TOMLFile", configFile, FileContentReader(configFile))
	actualSettings, err := tomlProvider.Load(context.Background())
	if err != nil {
		t.Error(err)
	}

	expectedSettings := map[string]string{
		"global.timeout":   "30",
		"global.frequency": "0.5",
		"local.time.zone":  "PST",
		"local.enabled":    "true",
	}

	for key, expected := range expectedSettings {
		actual, ok := actualSettings[key]

		if !ok {
			t.Errorf("Key '%s' not in settings", key)
		}

		if actual != expected {
			t.Errorf("Setting '%s' was '%s', expected '%s'", key, actual, expected)
		}
	}
}
