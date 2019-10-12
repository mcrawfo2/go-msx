package config

import (
	"context"
	"github.com/magiconair/properties"
	"github.com/pkg/errors"
)

type PropertiesFile struct {
	path string
}

func NewPropertiesFile(path string) *PropertiesFile {
	return &PropertiesFile{
		path: path,
	}
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
