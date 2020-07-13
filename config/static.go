package config

import (
	"context"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	yaml3 "gopkg.in/yaml.v3"
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

func NewStaticFromMap(name string, values map[string]interface{}) (*Static, error) {
	valuesBytes, err := yaml3.Marshal(values)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := yaml.YAMLToJSON(valuesBytes)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &values)
	if err != nil {
		return nil, err
	}

	settings, err := FlattenJSON(values, "")
	if err != nil {
		return nil, err
	}

	return &Static{
		name:     name,
		settings: settings,
	}, nil
}
