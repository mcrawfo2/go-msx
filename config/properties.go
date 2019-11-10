package config

import (
	"context"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/pkg/errors"
)

type PropertiesFile struct {
	name string
	path string
}

func (f *PropertiesFile) Description() string {
	return fmt.Sprintf("%s: [%s]", f.name, f.path)
}

func (f *PropertiesFile) Load(ctx context.Context) (map[string]string, error) {
	logger.Infof("Loading properties config: %s", f.path)

	props, err := properties.LoadFile(f.path, properties.UTF8)
	if err != nil {
		return nil, err
	}

	settings := map[string]string{}
	for _, key := range props.Keys() {
		settings[NormalizeKey(key)], _ = props.Get(key)
	}

	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "Failed to load properties config")
	}

	return settings, nil
}

func NewPropertiesFile(name, path string) *PropertiesFile {
	return &PropertiesFile{
		name: name,
		path: path,
	}
}
