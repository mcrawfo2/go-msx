package config

import (
	"context"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"io/ioutil"
)

type TOMLFile struct {
	path string
}

func NewTOMLFile(path string) *TOMLFile {
	return &TOMLFile{
		path: path,
	}
}

func (f *TOMLFile) Load(ctx context.Context) (map[string]string, error) {
	logger.Infof("Loading TOML config: %s", f.path)

	data, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}

	out := make(map[string]interface{})
	if _, err := toml.Decode(string(data), &out); err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "Failed to load TOML config")
	}

	return FlattenJSON(out, "")
}
