package config

import (
	"context"
	"github.com/pkg/errors"
)

type Static struct {
	name     string
	settings map[string]string
}

func (s *Static) Description() string {
	return s.name
}

func (s *Static) Load(ctx context.Context) (map[string]string, error) {
	logger.Infof("Loading %s config", s.name)

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

func NewStatic(name string, settings map[string]string) *Static {
	return &Static{
		name:     name,
		settings: settings,
	}
}
