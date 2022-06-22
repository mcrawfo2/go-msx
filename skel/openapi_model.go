// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"

	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/gedex/inflector"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

const extKeyMsxPermissions = "x-msx-permissions"

type Parameter struct {
	SchemaType Schema
	parameter  *openapi3.Parameter
}

func (p Parameter) VarName() string {
	return strcase.ToLowerCamel(fmt.Sprintf("param %s %s", p.parameter.In, p.parameter.Name))
}

func (p Parameter) Name() string {
	return strcase.ToCamel(p.parameter.Name)
}

func (p Parameter) JsonName() string {
	return p.parameter.Name
}

func (p Parameter) In() string {
	return p.parameter.In
}

func (p Parameter) Description() string {
	return p.parameter.Description
}

func NewParameter(parameter *openapi3.Parameter) (Parameter, error) {
	schemaType, err := NewSchemaType(parameter.Schema, parameter.Required)
	if err != nil {
		return Parameter{}, err
	}

	return Parameter{
		SchemaType: schemaType,
		parameter:  parameter,
	}, nil
}

func NewParameters(parameters openapi3.Parameters) ([]Parameter, error) {
	var results []Parameter
	for _, parameterRef := range parameters {
		result, err := NewParameter(parameterRef.Value)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

type Property struct {
	Schema Schema
	name   string
}

func (p Property) StructFieldName() string {
	return strcase.ToCamel(p.name)
}

func (p Property) JsonName() string {
	return strcase.ToLowerCamel(p.name)
}

func NewProperty(name string, schema *openapi3.SchemaRef, required bool) (Property, error) {
	schemaType, err := NewSchemaType(schema, required)
	if err != nil {
		return Property{}, err
	}

	return Property{
		Schema: schemaType,
		name:   name,
	}, nil
}

type Schema struct {
	schemaRef   *openapi3.SchemaRef
	builtin     bool
	array       bool
	dict        bool
	object      bool
	qual        string
	name        string
	pkg         string
	externalPkg string
	value       *Schema
	required    bool
}

var _any = Schema{
	schemaRef: &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "any",
		},
	},
}

func (s Schema) IsBuiltIn() bool {
	return s.builtin
}

func (s Schema) IsArray() bool {
	return s.array
}

func (s Schema) IsDict() bool {
	return s.dict
}

func (s Schema) IsObject() bool {
	return s.object
}

func (s Schema) IsAny() bool {
	return s.schemaRef.Value.Type == "any"
}

func (s Schema) IsReference() bool {
	return s.array || s.dict
}

func (s Schema) IsUuid() bool {
	return s.externalPkg == pkgTypes && s.name == "UUID"
}

func (s Schema) ItemType() Schema {
	if s.value == nil {
		return _any
	}
	return *s.value
}

func (s Schema) TypeName() string {
	return s.name
}

func (s Schema) TypeQualifier() string {
	if !s.builtin {
		return ""
	}

	return s.qual
}

func (s Schema) Namespace(appNamespace string) string {
	if s.builtin {
		return ""
	}

	if s.externalPkg != "" {
		return s.externalPkg
	}

	return path.Join(appNamespace, "pkg", s.pkg)
}

func (s Schema) Imports(appNamespace string) map[string]string {
	pkg := s.Namespace(appNamespace)
	return map[string]string{
		pkg: s.pkg,
	}
}

func (s Schema) Required() bool {
	return s.required
}

func (s Schema) Properties() ([]Property, error) {
	var (
		propMap    = make(map[string]Property)
		props, err = addNestedProperties(nil, s.schemaRef)

		names []string
	)

	if err != nil {
		return nil, err
	}

	for _, p := range props {
		if _, ok := propMap[p.name]; ok {
			logger.Printf("Overlapping property name on %s: %s", s.name, p.name)
			continue
		}

		names = append(names, p.name)
		propMap[p.name] = p
	}

	sort.Strings(names)

	results := make([]Property, 0, len(propMap))
	for _, name := range names {
		results = append(results, propMap[name])
	}

	return results, nil
}

func addNestedProperties(properties []Property, subschema *openapi3.SchemaRef) ([]Property, error) {
	properties, err := addSchemaProperties(properties, subschema)
	if err != nil {
		return nil, err
	}

	for _, subschema := range subschema.Value.AllOf {
		properties, err = addNestedProperties(properties, subschema)
		if err != nil {
			return nil, err
		}
	}

	return properties, nil
}

func (s Schema) Min() *float64 {
	if s.schemaRef == nil {
		return nil
	}

	return s.schemaRef.Value.Min
}

func (s Schema) Max() *float64 {
	if s.schemaRef == nil {
		return nil
	}

	return s.schemaRef.Value.Max
}

func (s Schema) MultipleOf() *float64 {
	if s.schemaRef == nil {
		return nil
	}

	return s.schemaRef.Value.MultipleOf
}

func (s Schema) Enum() []interface{} {
	if s.schemaRef == nil {
		return nil
	}

	if len(s.schemaRef.Value.Enum) == 0 {
		return nil
	}

	return s.schemaRef.Value.Enum
}

func (s Schema) ArrayLength() (int, int) {
	if s.schemaRef == nil {
		return 0, 0
	}

	if s.schemaRef.Value.MinItems == 0 {
		if s.schemaRef.Value.MaxItems == nil {
			return 0, 0
		} else {
			return 0, int(*s.schemaRef.Value.MaxItems)
		}
	} else if s.schemaRef.Value.MaxItems == nil {
		return int(s.schemaRef.Value.MinItems), 0
	} else {
		return int(s.schemaRef.Value.MinItems), int(*s.schemaRef.Value.MaxItems)
	}
}

func (s Schema) Length() (int, int) {
	if s.schemaRef == nil {
		return 0, 0
	}

	if s.schemaRef.Value.MinLength == 0 {
		if s.schemaRef.Value.MaxLength == nil {
			return 0, 0
		} else {
			return 0, int(*s.schemaRef.Value.MaxLength)
		}
	} else if s.schemaRef.Value.MaxLength == nil {
		return int(s.schemaRef.Value.MinLength), 0
	} else {
		return int(s.schemaRef.Value.MinLength), int(*s.schemaRef.Value.MaxLength)
	}
}

func (s Schema) Pattern() string {
	if s.schemaRef == nil {
		return ""
	}

	if s.schemaRef.Value.Pattern == "" {
		if s.schemaRef.Value.Type == "string" {
			return s.FormatPattern()
		}

		return ""
	}

	return s.schemaRef.Value.Pattern
}

func (s Schema) FormatPattern() string {
	if s.schemaRef == nil {
		return ""
	}

	switch s.schemaRef.Value.Format {
	case "uuid":
		return openapi3.FormatOfStringForUUIDOfRFC4122
	case "email", "date":
		return openapi3.SchemaStringFormats[s.schemaRef.Value.Format].String()
	}

	return ""
}

func addSchemaProperties(properties []Property, schemaRef *openapi3.SchemaRef) ([]Property, error) {
	required := types.StringStack(schemaRef.Value.Required)
	for propertyName, schemaRef := range schemaRef.Value.Properties {
		property, err := NewProperty(propertyName, schemaRef, required.Contains(propertyName))
		if err != nil {
			return nil, err
		}
		properties = append(properties, property)
	}
	return properties, nil
}

func NewComponentType(schemaName string, schemaRef *openapi3.SchemaRef) Schema {
	return Schema{
		schemaRef: schemaRef,
		qual:      "",
		name:      schemaName,
		pkg:       "api",
	}
}

func NewExternalType(schemaRef *openapi3.SchemaRef, required bool) Schema {
	refParts := strings.Split(schemaRef.Ref, "/")
	qualifiedName := refParts[len(refParts)-1]

	parts := strings.SplitN(qualifiedName, ".", 2)
	if len(parts) == 1 {
		parts = append([]string{"api"}, parts...)
	}

	qual := parts[0]
	name := parts[1]

	return Schema{
		schemaRef: schemaRef,
		builtin:   false,
		object:    schemaRef.Value.Type == "object",
		qual:      qual,
		name:      name,
		pkg:       "api",
		required:  required,
	}
}

func NewBuiltinType(schemaRef *openapi3.SchemaRef, name string, required bool) Schema {
	return Schema{
		schemaRef: schemaRef,
		builtin:   true,
		name:      name,
		required:  required,
	}
}

const pkgJson = "json"

func NewFrameworkType(schemaRef *openapi3.SchemaRef, pkg, qual, name string, required bool) Schema {
	return Schema{
		schemaRef:   schemaRef,
		object:      schemaRef.Value.Type == "object",
		qual:        qual,
		name:        name,
		externalPkg: pkg,
		required:    required,
	}
}

func NewArrayType(schemaRef *openapi3.SchemaRef, required bool) (Schema, error) {
	value, err := NewSchemaType(schemaRef.Value.Items, true)
	if err != nil {
		return Schema{}, err
	}

	return Schema{
		array:    true,
		value:    &value,
		required: required,
	}, nil
}

func NewDictType(schemaRef *openapi3.SchemaRef, required bool) Schema {
	return Schema{
		schemaRef: schemaRef.Value.Items,
		dict:      true,
		required:  required,
	}
}

func NewSchemaType(schemaRef *openapi3.SchemaRef, required bool) (Schema, error) {
	if schemaRef.Ref != "" {
		return NewExternalType(schemaRef, required), nil
	}

	switch schemaRef.Value.Type {
	case "string":
		switch schemaRef.Value.Format {
		case "uuid":
			return NewFrameworkType(schemaRef, pkgTypes, "types", "UUID", required), nil
		case "date-time":
			return NewFrameworkType(schemaRef, pkgTypes, "types", "Time", required), nil
		default:
			return NewBuiltinType(schemaRef, "string", required), nil
		}

	case "integer":
		switch schemaRef.Value.Format {
		case "int32":
			return NewBuiltinType(schemaRef, "int32", required), nil
		case "int64":
			return NewBuiltinType(schemaRef, "int64", required), nil
		default:
			return NewBuiltinType(schemaRef, "int", required), nil
		}

	case "number":
		switch schemaRef.Value.Format {
		case "float32":
			return NewBuiltinType(schemaRef, "float32", required), nil
		default:
			return NewBuiltinType(schemaRef, "float64", required), nil
		}

	case "boolean":
		return NewBuiltinType(schemaRef, "bool", required), nil

	case "array":
		return NewArrayType(schemaRef, required)

	case "object":
		switch schemaRef.Value.Format {
		case "json":
			return NewFrameworkType(schemaRef, pkgJson, "api", "RawMessage", required), nil
		default:
			return NewDictType(schemaRef, required), nil
		}

	default:
		return Schema{}, errors.Errorf("Unknown property type %q", schemaRef.Value.Type)
	}
}

type Body struct {
	Exists       bool
	Schema       Schema
	ContentTypes []string
}

func NewRequestBody(requestBody *openapi3.RequestBodyRef) (Body, error) {
	if requestBody == nil || requestBody.Value == nil {
		return Body{}, nil
	}

	schemaType, err := NewSchemaType(
		requestBody.Value.Content.Get("application/json").Schema,
		requestBody.Value.Required)
	if err != nil {
		return Body{}, err
	}

	var contentTypes []string
	for contentType := range requestBody.Value.Content {
		contentTypes = append(contentTypes, contentType)
	}

	return Body{
		Exists:       true,
		Schema:       schemaType,
		ContentTypes: contentTypes,
	}, nil
}

func NewResponseBody(responses openapi3.Responses) (Body, error) {
	responseBody := responses.Default()
	if responseBody == nil {
		// Find the first success response
		for code, codeResponse := range responses {
			// Only success responses
			if code[0] != '2' {
				continue
			}

			// Only responses with bodies
			if codeResponse.Value.Content.Get("application/json") == nil {
				continue
			}

			responseBody = codeResponse
			break
		}
	}

	if responseBody == nil || len(responseBody.Value.Content) == 0 {
		return Body{
			Exists: false,
		}, nil
	}

	schemaType, err := NewSchemaType(
		responseBody.Value.Content.Get("application/json").Schema,
		true)
	if err != nil {
		return Body{}, err
	}

	var contentTypes []string
	for contentType := range responseBody.Value.Content {
		contentTypes = append(contentTypes, contentType)
	}

	return Body{
		Exists:       true,
		Schema:       schemaType,
		ContentTypes: contentTypes,
	}, nil
}

type Endpoint struct {
	Path         string
	Method       string
	RequestBody  Body
	ResponseBody Body // Default success response
	Parameters   []Parameter
	operation    *openapi3.Operation
}

func (e Endpoint) OperationId() string {
	return e.operation.OperationID
}

func (e Endpoint) Summary() string {
	return e.operation.Summary
}

func (e Endpoint) ReturnCodes() []string {
	var results []string
	for code := range e.operation.Responses {
		if code[0] >= '0' && code[0] <= '5' {
			results = append(results, code)
		}
	}
	return results
}

func (e Endpoint) Permissions() []string {
	permissionsData, ok := e.operation.Extensions[extKeyMsxPermissions]
	if !ok {
		return nil
	}

	var permissions []string
	err := json.Unmarshal(permissionsData.(json.RawMessage), &permissions)
	if err != nil {
		return nil
	}
	return permissions
}

func NewEndpoint(path, method string, operation *openapi3.Operation) (Endpoint, error) {
	parameters, err := NewParameters(operation.Parameters)
	if err != nil {
		return Endpoint{}, err
	}

	requestBody, err := NewRequestBody(operation.RequestBody)
	if err != nil {
		return Endpoint{}, err
	}

	responseBody, err := NewResponseBody(operation.Responses)
	if err != nil {
		return Endpoint{}, err
	}

	return Endpoint{
		Path:         path,
		Method:       method,
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		Parameters:   parameters,
		operation:    operation,
	}, nil
}

type Controller struct {
	Tag       ControllerTag
	RootPath  string
	Endpoints []Endpoint
	Swagger   *openapi3.Swagger
}

func (c Controller) Name() string {
	modelEnglish := inflector.Singularize(c.Tag.Name)
	return strcase.ToLowerCamel(modelEnglish + " Controller")
}

func (c Controller) Package() string {
	modelEnglish := inflector.Singularize(c.Tag.Name)
	return strings.ReplaceAll(strcase.ToSnake(modelEnglish), "_", "")
}

func (c Controller) Namespace(appNamespace string) string {
	return path.Join(appNamespace, "internal", c.Package())
}

func NewController(tag ControllerTag, swagger *openapi3.Swagger) Controller {
	return Controller{
		Tag:       tag,
		RootPath:  "",
		Endpoints: nil,
		Swagger:   swagger,
	}
}

type ControllerTag struct {
	Name        string
	Description string
}

type Spec struct {
	Swagger *openapi3.Swagger
}

func (s Spec) Tags() map[string]ControllerTag {
	var result = make(map[string]ControllerTag)

	// Explicitly defined tags
	for _, tag := range s.Swagger.Tags {
		result[tag.Name] = ControllerTag{
			Name:        tag.Name,
			Description: tag.Description,
		}
	}

	// Implicitly defined tags
	for _, pathDefinition := range s.Swagger.Paths {
		for _, operation := range pathDefinition.Operations() {
			for _, tagName := range operation.Tags {
				if _, ok := result[tagName]; !ok {
					result[tagName] = ControllerTag{
						Name: tagName,
					}
				}
			}
		}
	}

	return result
}

func (s Spec) Controllers() (map[string]Controller, error) {
	var controllers = make(map[string]Controller)

	for _, tag := range s.Tags() {
		controllers[tag.Name] = NewController(tag, s.Swagger)
	}

	for pathKey, pathDefinition := range s.Swagger.Paths {
		for methodName, operation := range pathDefinition.Operations() {
			err := s.addControllerEndpoint(controllers, pathKey, methodName, operation)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to add controller endpoint %q %q", methodName, pathKey)
			}
		}
	}

	return controllers, nil
}

func (s Spec) addControllerEndpoint(controllers map[string]Controller, pathKey, method string, operation *openapi3.Operation) error {
	if operation == nil {
		return nil
	}

	endpoint, err := NewEndpoint(pathKey, method, operation)
	if err != nil {
		return err
	}

	controllerTagKey := operation.Tags[0]
	controllerEntry := controllers[controllerTagKey]
	controllerEntry.Endpoints = append(controllerEntry.Endpoints, endpoint)
	controllers[controllerTagKey] = controllerEntry
	return nil
}

func (s Spec) Schemas() ([]Schema, error) {
	var results []Schema
	for schemaName, schemaRef := range s.Swagger.Components.Schemas {
		results = append(results, NewComponentType(schemaName, schemaRef))
	}
	return results, nil
}

func NewSpec(swagger *openapi3.Swagger) Spec {
	return Spec{Swagger: swagger}
}
