// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import (
	"cto-github.cisco.com/NFV-BU/go-msx/build/npm"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"fmt"
	"gopkg.in/pipe.v2"
	"path/filepath"
	"regexp"
)

func init() {
	AddTarget("generate-asyncapi-schema", "Generates patched AsyncApi schema", GenerateAsyncApiSchema)
	AddTarget("generate-asyncapi-entities", "Generates AsyncApi entities", GenerateAsyncApiEntities)
	AddTarget("install-asyncapi-ui", "Installs AsyncAPI Studio package", InstallAsyncApiStudio)
}

const (
	AsyncApiVersion      = "2.4.0"
	AsyncApiSchemaUrl    = "https://raw.githubusercontent.com/asyncapi/spec-json-schemas/3d5d18e00fe5775964d9cd77aa5ea968555bf7ff/schemas/" + AsyncApiVersion + ".json"
	AsyncApiPatchFile    = "schema/asyncapi/asyncapi-" + AsyncApiVersion + "-patch.json"
	AsyncApiSchemaFile   = "schema/asyncapi/asyncapi-" + AsyncApiVersion + ".json"
	AsyncApiEntitiesFile = "schema/asyncapi/entities.go"
)

func GenerateAsyncApiSchema(_ []string) error {
	schemaPatchFile, _ := filepath.Abs(AsyncApiPatchFile)
	schemaFile, _ := filepath.Abs(AsyncApiSchemaFile)

	schemaTempFile, err := downloadTemp(AsyncApiSchemaUrl, "schema.*.json")
	if err != nil {
		return err
	}

	// 2. Generate a patched schema
	jsonCliPipe := pipe.Line(
		exec.ReadUrl(AsyncApiSchemaUrl),
		exec.ExecSimple(
			"docker", "run", "-i",
			"-v", fmt.Sprintf("%s:%s", schemaTempFile, "/schema.json"),
			"-v", fmt.Sprintf("%s:%s", schemaPatchFile, "/patch.json"),
			"swaggest/json-cli", "json-cli",
			"apply", "/patch.json", "/schema.json",
		),
		pipe.WriteFile(schemaFile, 0644),
	)

	return exec.ExecutePipes(jsonCliPipe)
}

var reJsonSchema = regexp.MustCompile(`"jsonschema"`)
var pkgJsonSchema = []byte(`"github.com/swaggest/jsonschema-go"`)

func GenerateAsyncApiEntities(_ []string) error {
	schemaFile, _ := filepath.Abs(AsyncApiSchemaFile)
	entitiesFile, _ := filepath.Abs(AsyncApiEntitiesFile)

	topOfFile := true

	jsonCliPipe := pipe.Line(
		pipe.ReadFile(schemaFile),
		exec.ExecSimple(
			"docker", "run", "-i",
			"swaggest/json-cli", "json-cli",
			"gen-go", "-",
			"--package-name", "asyncapi",
			"--root-name", "Spec",
			"--with-zero-values",
			"--fluent-setters",
		),
		exec.ExecSimple("grep", "-v", "Deprecated"),
		pipe.Filter(func(line []byte) bool {
			topOfFile = len(line) == 0 && topOfFile
			return !topOfFile
		}),
		pipe.Replace(func(line []byte) []byte {
			return reJsonSchema.ReplaceAll(line, pkgJsonSchema)
		}),
		pipe.WriteFile(entitiesFile, 0644),
	)

	return exec.ExecutePipes(jsonCliPipe)
}

func InstallAsyncApiStudio(_ []string) error {
	return npm.InstallNodePackageContents(
		BuildConfig.Msx.Platform.AsyncApi.Artifact,
		BuildConfig.Msx.Platform.AsyncApi.Version,
		"package/build",
		filepath.Join(BuildConfig.OutputStaticPath(), "asyncapi", "studio"))
}
