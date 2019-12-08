package config

import (
	"context"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type Environment struct {
	name string
}

func (e *Environment) Description() string {
	return e.name
}

func (e *Environment) Load(ctx context.Context) (map[string]string, error) {
	logger.Info("Loading environment config")

	settings := map[string]string{}

	for _, line := range os.Environ() {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := NormalizeKey(parts[0])
		if len(key) == 0 || key[0] == '.' {
			continue
		}
		settings[key] = parts[1]
	}

	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "Failed to load environment config")
	}

	return settings, nil
}

func NewEnvironment(name string) *Environment {
	return &Environment{
		name: name,
	}
}
