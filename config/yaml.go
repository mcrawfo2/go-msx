package config

import (
	"context"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"io/ioutil"
)

type YAMLFile struct {
	path string
}

func NewYAMLFile(path string) *YAMLFile {
	return &YAMLFile{
		path: path,
	}
}

func (f *YAMLFile) Load(ctx context.Context) (map[string]string, error) {
	logger.Infof("Loading YAML config: %s", f.path)

	encodedYAML, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}

	encodedJSON, err := yaml.YAMLToJSON(encodedYAML)
	if err != nil {
		return nil, err
	}

	decodedJSON := map[string]interface{}{}
	if err := json.Unmarshal(encodedJSON, &decodedJSON); err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "Failed to load YAML config")
	}

	return FlattenJSON(decodedJSON, "")
}
