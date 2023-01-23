package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema/openapi"
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/payloads"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	_ "embed"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
	"io/ioutil"
	"path"
)

//go:embed default-domain-openapi-v8.yaml
var defaultDomainV8 []byte

//go:embed default-domain-openapi-v2.yaml
var defaultDomainV2 []byte

type Spec struct {
	Raw        *openapi3.Spec
	Operations Operations
	Payloads   Payloads
}

func (g Spec) ResolveParameter(p openapi3.ParameterOrRef) openapi3.ParameterOrRef {
	if p.Parameter != nil {
		return p
	}
	if p.ParameterReference == nil {
		return p
	}

	refName := openapi.ParameterRefName(p.ParameterReference)
	if target, ok := g.Raw.ComponentsEns().ParametersEns().MapOfParameterOrRefValues[refName]; ok {
		return g.ResolveParameter(target)
	}

	return p
}

func (g Spec) ResolveRequestBody(b openapi3.RequestBodyOrRef) openapi3.RequestBodyOrRef {
	if b.RequestBody != nil {
		return b
	}
	if b.RequestBodyReference == nil {
		return b
	}

	refName := openapi.RequestBodyRefName(b.RequestBodyReference)
	if target, ok := g.Raw.ComponentsEns().RequestBodiesEns().MapOfRequestBodyOrRefValues[refName]; ok {
		return g.ResolveRequestBody(target)
	}

	return b
}

func (g Spec) ResolveSchema(s openapi3.SchemaOrRef) openapi3.SchemaOrRef {
	if s.Schema != nil {
		return s
	}
	if s.SchemaReference == nil {
		return s
	}

	refName := openapi.SchemaRefName(s.SchemaReference)
	if target, ok := g.Raw.ComponentsEns().SchemasEns().MapOfSchemaOrRefValues[refName]; ok {
		return g.ResolveSchema(target)
	}

	return s
}

func (g Spec) ResolveResponse(r openapi3.ResponseOrRef) openapi3.ResponseOrRef {
	if r.Response != nil {
		return r
	}
	if r.ResponseReference == nil {
		return r
	}

	refName := openapi.ResponseRefName(r.ResponseReference)
	if target, ok := g.Raw.ComponentsEns().ResponsesEns().MapOfResponseOrRefValues[refName]; ok {
		return g.ResolveResponse(target)
	}

	return r
}

func (g Spec) ResolveHeader(r openapi3.HeaderOrRef) openapi3.HeaderOrRef {
	if r.Header != nil {
		return r
	}
	if r.HeaderReference == nil {
		return r
	}

	refName := openapi.HeaderRefName(r.HeaderReference)
	if target, ok := g.Raw.ComponentsEns().HeadersEns().MapOfHeaderOrRefValues[refName]; ok {
		return g.ResolveHeader(target)
	}

	return r
}

func (g Spec) CollectJsonSchema(schemaName string, schema jsonschema.Schema) (result jsonschema.Schema, err error) {
	return payloads.CollectSchema(schemaName, schema, g.ResolveJsonSchema)
}

func (g Spec) ResolveJsonSchema(refName string) jsonschema.Schema {
	ref := openapi.NewSchemaRef(refName)
	refSchema := *g.GetJsonSchema(&ref)
	return *refSchema.Definitions[refName].TypeObject
}

func (g Spec) GetJsonSchema(schema *openapi3.SchemaOrRef) *jsonschema.Schema {
	schemaOrBool := schema.ToJSONSchema(g.Raw)

	// Move components into definitions
	components := payloads.ComponentsForSchema(schemaOrBool.TypeObject)
	components.Schemas().Each(func(key string, value jsonschema.SchemaOrBool) {
		schemaOrBool.TypeObject.WithDefinitionsItem(key, value)
	})
	delete(schemaOrBool.TypeObject.ExtraProperties, payloads.ExtraPropertiesComponents)

	// Convert the reference to retain indirect type handling
	if schema.SchemaReference != nil {
		// Turn it back into a schema ref
		refName := openapi.SchemaRefName(schema.SchemaReference)
		schemaOrBool = jsonschema.SchemaOrBool{
			TypeObject: &jsonschema.Schema{
				Definitions: map[string]jsonschema.SchemaOrBool{
					refName: schemaOrBool,
				},
				Ref: types.NewStringPtr(payloads.PrefixDefinitions + refName),
			},
		}
	}

	return schemaOrBool.TypeObject
}

type Operations []Operation

func (o Operations) ForAction(action string) Operations {
	var results Operations
	for _, operation := range o {
		if operation.Action == action {
			results = append(results, operation)
		}
	}
	return results
}

type Operation struct {
	Method    string
	SubPath   string
	Action    string
	Operation openapi3.Operation
}

type schemaVisitor func(s *openapi3.SchemaOrRef, spec *Spec)

func (o Operation) WalkSchemas(spec *Spec, visitor schemaVisitor) {
	if o.Operation.RequestBody != nil {
		requestBodyOrRef := spec.ResolveRequestBody(*o.Operation.RequestBody)
		for _, mediaType := range requestBodyOrRef.RequestBody.Content {
			if mediaType.Schema == nil {
				continue
			}

			visitor(mediaType.Schema, spec)
		}
	}

	for _, parameterOrRef := range o.Operation.Parameters {
		parameterOrRef = spec.ResolveParameter(parameterOrRef)
		if parameterOrRef.Parameter == nil {
			continue
		}

		parameter := parameterOrRef.Parameter
		if parameter.Schema == nil {
			continue
		}

		visitor(parameter.Schema, spec)
	}

	for _, responseOrRef := range o.Operation.Responses.MapOfResponseOrRefValues {
		responseOrRef = spec.ResolveResponse(responseOrRef)
		if responseOrRef.Response == nil {
			continue
		}

		response := responseOrRef.Response
		for _, headerOrRef := range response.Headers {
			headerOrRef = spec.ResolveHeader(headerOrRef)
			if headerOrRef.Header == nil {
				continue
			}

			header := headerOrRef.Header
			if header.Schema == nil {
				continue
			}

			visitor(header.Schema, spec)
		}

		for _, mediaType := range response.Content {
			if mediaType.Schema == nil {
				continue
			}

			visitor(mediaType.Schema, spec)
		}
	}
}

type Payload struct {
	Actions []string
	Schema  *openapi3.SchemaOrRef
}

type Payloads []Payload

func (p Payloads) ForActions(actions ...string) Payloads {
	var results Payloads
	for _, payload := range p {
		if types.ComparableSlice[string](payload.Actions).ContainsAny(actions...) {
			results = append(results, payload)
		}
	}
	return results
}

func loadSpecification(tagName string, inflector skel.Inflector) (spec Spec, err error) {
	spec.Raw, err = loadSpecificationFile(inflector)
	if err != nil {
		return
	}

	for urlPath, pathItem := range spec.Raw.Paths.MapOfPathItemValues {
		for method, openApiOperation := range pathItem.MapOfOperationValues {
			if !types.ComparableSlice[string](openApiOperation.Tags).Contains(tagName) {
				// skip any endpoints outside our chosen domain
				continue
			}

			op := Operation{
				Method:    method,
				SubPath:   urlPath,
				Operation: openApiOperation,
			}

			actionAny, ok := openApiOperation.MapOfAnything[ExtraPropertiesMsxAction]
			if ok {
				op.Action = actionAny.(string)
			}

			if op.Operation.ID == nil {
				// TODO: auto-generate operation id
			}

			spec.Operations = append(spec.Operations, op)

			op.WalkSchemas(&spec, func(s *openapi3.SchemaOrRef, g *Spec) {
				if s.SchemaReference == nil {
					return
				}

				// add the payload
				payload := Payload{
					Actions: nil,
					Schema:  s,
				}

				if op.Action != "" {
					payload.Actions = []string{op.Action}
				}

				g.Payloads = append(g.Payloads, payload)
			})
		}
	}

	return
}

func loadSpecificationFile(inflector skel.Inflector) (spec *openapi3.Spec, err error) {
	var specBytes []byte
	var format = skel.FileFormatYaml

	if generatorConfig.Spec == "" {
		switch generatorConfig.Style {
		case StyleV8:
			specBytes = defaultDomainV8

		case StyleV2:
			specBytes = defaultDomainV2
		}
	} else {
		switch path.Ext(generatorConfig.Spec) {
		case ".json", ".json5":
			format = skel.FileFormatJson
		case ".yaml":
			format = skel.FileFormatYaml
		}

		specBytes, err = ioutil.ReadFile(generatorConfig.Spec)
		if err != nil {
			return nil, err
		}

		inflector = nil
	}

	if format == skel.FileFormatYaml {
		if specBytes, err = yaml.YAMLToJSON(specBytes); err != nil {
			return
		}
	}

	if inflector != nil {
		template := skel.Template{
			Name:       "spec",
			SourceData: specBytes,
			Format:     skel.FileFormatJson,
		}

		renderOptions := skel.NewRenderOptions()
		renderOptions.Strings = inflector

		renderedSpec, err := template.RenderContents(renderOptions)
		if err != nil {
			return nil, err
		}

		specBytes = []byte(renderedSpec)
	}

	err = json.Unmarshal(specBytes, &spec)
	return spec, err
}
