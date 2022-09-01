package swaggerprovider

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

func CustomizeCommand(cmd *cobra.Command) {
	cmd.Args = cobra.ExactArgs(1)
	cmd.Use += " target.{json|yml}"
}

func SaveSpec(ctx context.Context, args []string) (err error) {
	filename := args[0]
	format := filepath.Ext(filename)

	swagger := provider.GetSpecDocument()

	data, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return err
	}

	switch format {
	case ".yml", ".yaml":
		data, err = jsonToYaml(data)
	case ".json":
	default:
		return errors.Errorf("Unknown file format: %q", format)
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func jsonToYaml(specJsonBytes []byte) (specYamlBytes []byte, err error) {
	var specYaml = yaml.MapSlice{}
	err = yaml.Unmarshal(specJsonBytes, &specYaml)
	if err != nil {
		return nil, err
	}

	specYamlBytes, err = yaml.Marshal(specYaml)
	if err != nil {
		return nil, err
	}

	return specYamlBytes, nil
}
