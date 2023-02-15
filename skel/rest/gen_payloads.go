// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema/openapi"
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/payloads"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"fmt"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"github.com/mcrawfo2/go-jsonschema/pkg/generator"
	"path"
)

const (
	ExtraPropertiesMsxAction           = "x-msx-action"
	ExtraPropertiesMsxActions          = "x-msx-actions"
	ExtraPropertiesMsxInjectedProperty = "x-msx-injected-property"
)

type DomainPayloadsGenerator struct {
	Domain  string
	Folder  string
	Actions types.ComparableSlice[string]

	// Go payload generator from jsonschema-go
	Generator *generator.Generator

	Spec Spec

	*text.GoFile
}

func (g DomainPayloadsGenerator) createPayloadSnippets() error {
	packageName := g.GoFile.Inflector.Inflect(skel.InflectionLowerPlural)
	generatedSchemas := make(types.StringSet)

	for _, payload := range g.Spec.Payloads.ForActions(g.Actions...) {
		schemaName := openapi.SchemaRefName(payload.Schema.SchemaReference)
		if generatedSchemas.Contains(schemaName) {
			continue
		}
		generatedSchemas.Add(schemaName)

		payloadSchema := g.Spec.GetJsonSchema(payload.Schema)

		logger.Infof("  🥽 Analyzing schema %s", schemaName)

		// Grab all the transitive schemas
		collectedSchema, err := g.Spec.CollectJsonSchema(schemaName, *payloadSchema)
		if err != nil {
			return err
		}

		// Skip types that are resolved outside of this module (error, paging, uuid, time, etc)
		optionalType := payloads.GoJsonSchemaForSchema(&collectedSchema).Type()
		if optionalType.IsPresent() {
			logger.Infof("  ⏭️ Skipping external schema %s", optionalType.Value())
			continue
		}

		// Rotate schema so desired schema is at root
		resolvedSchema := collectedSchema
		for resolvedSchema.Ref != nil {
			refName := path.Base(*resolvedSchema.Ref)
			resolvedSchema = *collectedSchema.Definitions[refName].TypeObject
		}
		resolvedSchema.Definitions = collectedSchema.Definitions
		collectedSchema = resolvedSchema

		// Convert to JSON
		schemaBytes, err := json.Marshal(collectedSchema)
		if err != nil {
			return err
		}

		// Add the spec to the generator
		err = g.Generator.AddSource(generator.Source{
			PackageName: packageName,
			RootType:    schemaName,
			Folder:      g.Folder,
			Data:        schemaBytes,
		})
		if err != nil {
			return err
		}
	}

	// Convert the returned file to templates
	for _, file := range g.Generator.Files() {
		logger.Infof("  📇 Collecting %s", file.FileName)
		for _, decl := range file.Package.Decls {
			namedDecl := decl.(codegen.Named)

			err := g.GoFile.AddNewDecl(
				"Payloads",
				namedDecl.GetName(),
				decl,
				file.Package.Imports)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g DomainPayloadsGenerator) Generate() error {
	return types.ErrorList{
		g.createPayloadSnippets(),
	}.Filter()
}

func (g DomainPayloadsGenerator) Filename() string {
	target := path.Join(g.Folder, "payloads_lowersingular.go")
	return g.GoFile.Inflector.Inflect(target)
}

func NewDomainPayloadsGenerator(spec Spec) ComponentGenerator {
	inflector := skel.NewInflector(generatorConfig.Domain)

	return DomainPayloadsGenerator{
		Domain:  generatorConfig.Domain,
		Folder:  generatorConfig.Folder,
		Actions: generatorConfig.Actions,
		Generator: types.May(generator.New(generator.Config{
			OutputFiler: func(definition string) string {
				return fmt.Sprintf("payloads_%s.go", inflector.Inflect(skel.InflectionLowerSingular))
			},
		})),
		Spec: spec,
		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   "Payloads for " + generatorConfig.Domain,
				Sections:  text.NewGoSections("Payloads"),
				Inflector: inflector,
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
