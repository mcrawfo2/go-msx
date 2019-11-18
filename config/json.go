package config

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	json "github.com/yosuke-furukawa/json5/encoding/json5"
	"io/ioutil"
)

type JSONFile struct {
	name string
	path string
}

func (f *JSONFile) Description() string {
	return fmt.Sprintf("%s: [%s]", f.name, f.path)
}

func (f *JSONFile) Load(ctx context.Context) (map[string]string, error) {
	logger.Infof("Loading JSON config: %s", f.path)

	encodedJSON, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}

	decodedJSON := map[string]interface{}{}
	if err := json.Unmarshal(encodedJSON, &decodedJSON); err != nil {
		return nil, err
	}

	settings, err := FlattenJSON(decodedJSON, "")
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "Failed to load JSON config")
	}

	return settings, nil

}

func FlattenJSON(input map[string]interface{}, namespace string) (map[string]string, error) {
	flattened := map[string]string{}

	for key, value := range input {
		var token string
		if namespace == "" {
			token = key
		} else {
			token = fmt.Sprintf("%s.%s", namespace, key)
		}

		if objectChild, ok := value.(map[string]interface{}); ok {
			settings, err := FlattenJSON(objectChild, token)
			if err != nil {
				return nil, err
			}

			for k, v := range settings {
				flattened[NormalizeKey(k)] = v
			}
		} else if arrayChild, ok := value.([]interface{}); ok {
			for k, v := range arrayChild {
				flattened[NormalizeKey(token) + fmt.Sprintf("[%d]", k)] = fmt.Sprintf("%v", v)
			}
		} else {
			flattened[NormalizeKey(token)] = fmt.Sprintf("%v", value)
		}
	}

	return flattened, nil
}

func NewJSONFile(name, path string) *JSONFile {
	return &JSONFile{
		name: name,
		path: path,
	}
}
