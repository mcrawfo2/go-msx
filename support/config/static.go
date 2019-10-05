package config

import (
	"context"
	"github.com/pkg/errors"
)

type Static struct {
	settings map[string]string
}

func NewStatic(settings map[string]string) *Static {
	return &Static{
		settings: settings,
	}
}

func (s *Static) Load(ctx context.Context) (map[string]string, error) {
	logger.Info("Loading static config")

	settings := map[string]string{}

	for key, value := range s.settings {
		settings[key] = value
	}

	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "Failed to load static config")
	}

	return settings, nil
}

func (s *Static) Set(key, val string) {
	s.settings[key] = val
}
