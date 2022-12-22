// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/spf13/cobra"
	"io/ioutil"
)

const (
	defaultSpecFileName = "${fs.sources}/api/openapi.yaml"
)

var specProvider SpecProvider

type SpecProvider interface {
	RenderYamlSpec() ([]byte, error)
}

func GenerateOpenApiSpecCommand(ctx context.Context, _ []string) error {
	if specProvider == nil {
		logger.WithContext(ctx).Fatal("No OpenApi specification provider registered")
	}

	yamlBytes, err := specProvider.RenderYamlSpec()
	if err != nil {
		logger.WithContext(ctx).WithError(err).Fatal("Failed to render OpenApi spec as YAML")
	}

	outputFileName, _ := config.FromContext(ctx).StringOr("cli.flag.output", defaultSpecFileName)
	err = ioutil.WriteFile(outputFileName, yamlBytes, 0664)
	if err != nil {
		logger.WithContext(ctx).WithError(err).Fatalf("Failed to save OpenApi spec to %q", outputFileName)
	}

	logger.WithContext(ctx).Infof("Saved OpenApi spec to %q", outputFileName)
	return nil
}

func CustomizeOpenApiSpecCommand(cmd *cobra.Command) {
	cmd.Args = cobra.NoArgs
	cmd.Flags().String("output", defaultSpecFileName, "Specify the output file for the OpenApi specification")
}
