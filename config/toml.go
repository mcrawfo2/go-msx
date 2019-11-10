package config

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"io/ioutil"
)

type TOMLFile struct {
	name string
	path string
}

func (f *TOMLFile) Description() string {
	return fmt.Sprintf("%s: [%s]", f.name, f.path)
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

func NewTOMLFile(name, path string) *TOMLFile {
	return &TOMLFile{
		name: name,
		path: path,
	}
}
