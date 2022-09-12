package asyncapiprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/spf13/cobra"
	"io/ioutil"
)

const (
	defaultSpecFileName = "${fs.sources}/api/asyncapi.yaml"
)

func GenerateAsyncApiSpecCommand(ctx context.Context, _ []string) (err error) {
	specProvider, err := NewRegistrySpecProvider(ctx)
	if err != nil {
		return err
	}

	yamlBytes, err := specProvider.Spec()
	if err != nil {
		logger.WithContext(ctx).WithError(err).Fatal("Failed to render AsyncApi spec as YAML")
	}

	outputFileName, _ := config.FromContext(ctx).StringOr("cli.flag.output", defaultSpecFileName)
	err = ioutil.WriteFile(outputFileName, yamlBytes, 0664)
	if err != nil {
		logger.WithContext(ctx).WithError(err).Fatalf("Failed to save AsyncApi spec to %q", outputFileName)
	}

	logger.WithContext(ctx).Infof("Saved AsyncApi spec to %q", outputFileName)
	return nil
}

func CustomizeAsyncApiSpecCommand(cmd *cobra.Command) {
	cmd.Flags().String("output", defaultSpecFileName, "Specify the output file for the AsyncApi specification")
}
