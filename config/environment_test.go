package config

import (
	"context"
	"os"
	"testing"
)

func TestEnvironmentLoad(t *testing.T) {
	os.Setenv("GLOBAL_TIMEOUT", "30")
	os.Setenv("GLOBAL_FREQUENCY", "0.5")
	os.Setenv("LOCAL_TIME_ZONE", "PST")
	os.Setenv("LOCAL_ENABLED", "true")

	p := NewEnvironment("env")

	expectedSettings := map[string]string{
		"global.timeout":   "30",
		"global.frequency": "0.5",
		"local.time.zone":  "PST",
		"local.enabled":    "true",
	}

	actualSettings, err := p.Load(context.Background())
	if err != nil {
		t.Error(err)
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
