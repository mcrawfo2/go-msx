// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mcrawfo2/jennifer/jen"
	"github.com/pkg/errors"
)

func init() {
	AddTarget("generate-domain-openapi", "Create domains from OpenAPI 3.0 manifest", GenerateDomainOpenApi)
}

func GenerateDomainOpenApi(args []string) error {
	if len(args) == 0 {
		return errors.New("No OpenAPI spec provided.")
	}

	bytes, err := os.ReadFile(args[0])
	if err != nil {
		return errors.Wrap(err, "Failed to read OpenAPI spec")
	}

	var loader = openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	loader.LoadSwaggerFromURIFunc = loadSwaggerFromUri

	swagger, err := loader.LoadSwaggerFromData(bytes)
	if err != nil {
		return errors.Wrap(err, "Failed to parse OpenAPI spec")
	}

	var spec = NewSpec(swagger)

	controllers, err := spec.Controllers()
	if err != nil {
		return errors.Wrap(err, "Failed to identify controllers from spec")
	}

	logger.Infof("Parsed controllers %+v", controllers)

	schemas, _ := spec.Schemas()
	for _, schema := range schemas {
		var err error
		if schema.Enum() != nil {
			err = generateEnums(schema)
		} else if schema.IsTypeObject() {
			err = generateSchema(schema)
		} else {
			err = generateTypeSchema(schema)
		}
		if err != nil {
			return errors.Wrapf(err, "Failed to generate schema %q", schema.TypeName())
		}
	}

	for _, controller := range controllers {
		err := generateController(controller)
		if err != nil {
			return errors.Wrap(err, "Failed to generate controller")
		}
	}

	return nil
}

func loadSwaggerFromUri(loader *openapi3.SwaggerLoader, url *url.URL) (*openapi3.Swagger, error) {
	resourcePath := strings.TrimPrefix(url.Path, "/domains/Cisco-Systems46/")
	documentParts := strings.SplitN(resourcePath, "/", 2)

	documentFilePath := fmt.Sprintf("openapi/%s.v%s.yaml", documentParts[0], documentParts[1])

	bytes := []byte(staticFiles[documentFilePath].data)
	return loader.LoadSwaggerFromData(bytes)
}

func generateTypeWithImport(f *jen.File, ns string, s *jen.Statement, schemaType Schema) error {
	imports, err := generateType(s, ns, schemaType)
	if err != nil {
		return err
	}

	f.ImportNames(imports)
	return nil
}

func generateType(s *jen.Statement, ns string, schema Schema) (map[string]string, error) {
	if schema.IsArray() {
		// slice
		s = s.Index()
		schema = schema.ItemType()
		return generateType(s, ns, schema)
	}

	if schema.IsDict() {
		// map
		s = s.Map(jen.String())
		schema = schema.ItemType()
		return generateType(s, ns, schema)
	}

	if schema.IsAny() {
		s = s.Interface()
		return nil, nil
	}

	if schema.IsBuiltIn() {
		// No imports or qualifier
		s = s.Id(schema.TypeName())
		return nil, nil
	}

	sns := schema.Namespace(skeletonConfig.AppPackageUrl())
	var imports map[string]string
	if sns == ns {
		s = s.Id(schema.TypeName())
	} else {
		s = s.Qual(sns, schema.TypeName())
		imports = schema.Imports(skeletonConfig.AppPackageUrl())
	}

	return imports, nil
}

func generateValidators(f *jen.File, schema Schema) ([]jen.Code, error) {
	requiredValidators, err := generateRequiredValidators(f, schema)
	if err != nil {
		return nil, err
	}

	var validators []jen.Code
	if schema.Required() {
		f.ImportName(pkgValidation, "validation")
		validators = append(validators, jen.Qual(pkgValidation, "Required"))
		validators = append(validators, requiredValidators...)
	} else {
		if len(requiredValidators) == 0 {
			return []jen.Code{}, nil
		}

		f.ImportName(pkgValidate, "validate")
		ifNotNilValidator := jen.Qual(pkgValidate, "IfNotNil").Call(requiredValidators...)
		validators = append(validators, ifNotNilValidator)
	}

	return validators, nil
}

func generateRequiredValidators(f *jen.File, schema Schema) ([]jen.Code, error) {
	if schema.IsObject() || schema.IsUuid() {
		f.ImportName(pkgValidate, "validate")

		return []jen.Code{
			jen.Qual(pkgValidate, "Self"),
		}, nil
	}

	if schema.IsArray() || schema.IsDict() {
		f.ImportName(pkgValidation, "validation")
		f.ImportName(pkgValidate, "validate")

		itemValidators, err := generateRequiredValidators(f, schema.ItemType())
		if err != nil {
			return nil, err
		}

		return []jen.Code{
			jen.Qual(pkgValidation, "Each").Call(itemValidators...),
		}, nil
	}

	var validators []jen.Code

	if min := schema.Min(); min != nil {
		f.ImportName(pkgValidation, "validation")
		validators = append(validators, jen.Qual(pkgValidation, "Min").Call(jen.Lit(*min)))
	}

	if max := schema.Max(); max != nil {
		f.ImportName(pkgValidation, "validation")
		validators = append(validators, jen.Qual(pkgValidation, "Max").Call(jen.Lit(*max)))
	}

	if factor := schema.MultipleOf(); factor != nil {
		f.ImportName(pkgValidation, "validation")
		validators = append(validators, jen.Qual(pkgValidation, "MultipleOf").Call(jen.Lit(*factor)))
	}

	if min, max := schema.Length(); min != 0 || max != 0 {
		f.ImportName(pkgValidation, "validation")
		validators = append(validators, jen.Qual(pkgValidation, "Length").Call(jen.Lit(min), jen.Lit(max)))
	}

	if min, max := schema.ArrayLength(); min != 0 || max != 0 {
		f.ImportName(pkgValidation, "validation")
		validators = append(validators, jen.Qual(pkgValidation, "Length").Call(jen.Lit(min), jen.Lit(max)))
	}

	if pattern := schema.Pattern(); pattern != "" {
		f.ImportName(pkgValidation, "validation")
		f.ImportName(pkgRegexp, "regexp")

		validators = append(validators, jen.Qual(pkgValidation, "Match").Call(
			jen.Qual(pkgRegexp, "MustCompile").Params(jen.Lit(pattern))))
	}

	if enum := schema.Enum(); enum != nil {
		f.ImportName(pkgValidation, "validation")
		validators = append(validators, jen.Qual(pkgValidation, "In").Call(anyLiterals(enum)...))
	}

	return validators, nil
}

func stringLiterals(values []string) []jen.Code {
	var literals []jen.Code
	for _, value := range values {
		literals = append(literals, jen.Lit(value))
	}
	return literals
}

func anyLiterals(values []interface{}) []jen.Code {
	var literals []jen.Code
	for _, value := range values {
		if value != nil {
			literals = append(literals, jen.Lit(value))
		} else {
			literals = append(literals, jen.Null())
		}
	}
	return literals
}

func writeFile(targetFileName string, f *jen.File) (err error) {
	err = os.MkdirAll(path.Dir(targetFileName), 0755)
	if err != nil {
		return errors.Wrap(err, "Failed to create directory")
	}

	writer, err := os.Create(targetFileName)
	if err != nil {
		return errors.Wrap(err, "Failed to create file")
	}

	err = f.Render(writer)
	if err != nil {
		return errors.Wrap(err, "Failed to write file")
	}

	return nil
}
