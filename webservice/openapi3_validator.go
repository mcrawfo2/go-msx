package webservice

import (
	"bytes"
	"crypto/md5"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/swaggest/jsonschema-go"
	"io/ioutil"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"sync"
)

var parameterValidationCompiler = jsv.NewCompiler()
var parameterValidationSchema = make(map[string]*jsv.Schema, 64)
var parameterValidationSchemaMtx sync.Mutex

type OpenApiRequestValidator struct {
	Decoder RequestDecoder
}

func (v OpenApiRequestValidator) getParameterSchema(parameter EndpointRequestParameter) (result ValidationSchema, err error) {
	parameterValidationSchemaMtx.Lock()
	defer parameterValidationSchemaMtx.Unlock()

	if jsvSchema, ok := parameterValidationSchema[parameter.UniqueId]; ok {
		return NewValidationSchema(jsvSchema), nil
	}

	jsob := parameter.Schema.ToJSONSchema(Reflector.Spec)

	// TODO: trim non-validation fields

	schemaBytes, err := jsob.JSONSchemaBytes()
	if err != nil {
		return
	}

	sum := md5.Sum(schemaBytes)
	hash := hex.EncodeToString(sum[:])
	schemaUrl := "mem:///" + hash + ".json"

	if err = parameterValidationCompiler.AddResource(schemaUrl, bytes.NewReader(schemaBytes)); err != nil {
		return
	}

	jsvSchema, err := parameterValidationCompiler.Compile(schemaUrl)
	if err != nil {
		return
	}

	parameterValidationSchema[parameter.UniqueId] = jsvSchema

	return NewValidationSchema(jsvSchema), nil
}

func (v OpenApiRequestValidator) ValidateRequest(endpoint Endpoint) (err error) {
	errs := &ValidationFailure{
		Path:     "request",
		Children: make(map[string]*ValidationFailure),
	}

	for _, param := range endpoint.Request.Parameters {
		if param.PortField == nil {
			continue
		}

		var validationSchema ValidationSchema
		validationSchema, err = v.getParameterSchema(param)
		if err != nil {
			return err
		}

		var value interface{}

		if param.In == "body" {
			var body json.RawMessage
			body, err = v.Decoder.DecodeBodyToJson(!*param.Required)
			if err != nil {
				return err
			}

			err = json.Unmarshal(body, &value)
			if err != nil {
				return err
			}

		} else {

			switch param.PortField.Shape {
			case FieldShapePrimitive:
				value, err = v.GetPrimitive(param, validationSchema)
				if err != nil {
					return err
				}

			case FieldShapeArray:
				value, err = v.GetArray(param, validationSchema)
				if err != nil {
					return err
				}

			case FieldShapeObject:
				value, err = v.GetObject(param, validationSchema)
				if err != nil {
					return err
				}

			case FieldShapeFile:
				value, err = v.GetFileContentAsString(param)
				if err != nil {
					return err
				}
			}
		}

		err = validationSchema.Validate(value)
		if err != nil {
			switch typedErr := err.(type) {
			case *jsv.ValidationError:
				errs.Children[param.Name] = NewValidationFailure(typedErr.InstanceLocation).Apply(typedErr)
			default:
				return err
			}

		}
	}

	if len(errs.Children) > 0 {
		return errs
	}
	return nil
}

func (v OpenApiRequestValidator) TypesHasType(types []string, simpleType jsonschema.SimpleType) bool {
	for _, t := range types {
		if string(simpleType) == t {
			return true
		}
	}
	return false
}

func (v OpenApiRequestValidator) GetPrimitiveElement(value string, types []string) (interface{}, error) {
	switch {
	case v.TypesHasType(types, jsonschema.String):
		return value, nil

	case v.TypesHasType(types, jsonschema.Number):
		return strconv.ParseFloat(value, 64)

	case v.TypesHasType(types, jsonschema.Integer):
		return strconv.ParseInt(value, 10, 64)

	case v.TypesHasType(types, jsonschema.Boolean):
		return strconv.ParseBool(value)

	case v.TypesHasType(types, jsonschema.Array):
		return nil, errors.New("Cannot convert string to non-primitive type 'array'")

	case v.TypesHasType(types, jsonschema.Object):
		return nil, errors.New("Cannot convert string to non-primitive type 'object'")
	}

	return nil, errors.Errorf("Cannot determine target type of schema %+v", types)

}

func (v OpenApiRequestValidator) GetPrimitive(param EndpointRequestParameter, schema ValidationSchema) (interface{}, error) {
	optionalValue, _ := v.Decoder.DecodePrimitive(
		param.In,
		param.Name,
		*param.Style,
		*param.Explode)

	if !optionalValue.IsPresent() {
		return nil, nil
	}

	return v.GetPrimitiveElement(optionalValue.String(), schema.Types())
}

func (v OpenApiRequestValidator) GetArray(param EndpointRequestParameter, schema ValidationSchema) ([]interface{}, error) {
	values, _ := v.Decoder.DecodeArray(
		param.In,
		param.Name,
		*param.Style,
		*param.Explode)

	if len(values) == 0 {
		return nil, nil
	}

	var results = make([]interface{}, len(values))
	for i, str := range values {
		items := schema.Items()
		if value, err := v.GetPrimitiveElement(str, items.Types()); err != nil {
			return nil, err
		} else {
			results[i] = value
		}
	}

	return results, nil
}

var anyType = NewValidationSchema(&jsv.Schema{
	Types: []string{
		string(jsonschema.Null),
		string(jsonschema.Array),
		string(jsonschema.Object),
		string(jsonschema.String),
		string(jsonschema.Integer),
		string(jsonschema.Number),
		string(jsonschema.Boolean),
	},
})

func (v OpenApiRequestValidator) GetObject(param EndpointRequestParameter, schema ValidationSchema) (types.Pojo, error) {
	values, _ := v.Decoder.DecodeObject(
		param.In,
		param.Name,
		*param.Style,
		*param.Explode)

	if len(values) == 0 {
		return nil, nil
	}

	var results = make(types.Pojo, len(values))
	for key, str := range values {
		if !schema.IsPropertyAllowed(key) {
			continue
		}

		ps := schema.Property(key)

		if value, err := v.GetPrimitiveElement(str.(string), ps.Types()); err != nil {
			return nil, err
		} else {
			results[key] = value
		}
	}

	return results, nil
}

func (v OpenApiRequestValidator) GetFileHeader(param EndpointRequestParameter) (*multipart.FileHeader, error) {
	value, err := v.Decoder.DecodeFormFile(param.Name, *param.Required)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (v OpenApiRequestValidator) GetFile(param EndpointRequestParameter) (multipart.File, error) {
	fileHeader, err := v.GetFileHeader(param)
	if err != nil {
		return nil, err
	}

	var file multipart.File
	file, err = fileHeader.Open()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (v OpenApiRequestValidator) GetFileContentAsString(param EndpointRequestParameter) (string, error) {
	file, err := v.GetFile(param)

	var valueBytes []byte
	valueBytes, err = ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(valueBytes), nil
}

func NewValidationFailure(p string) *ValidationFailure {
	if len(p) > 0 && p[0] != '/' {
		p = "/" + p
	}
	return &ValidationFailure{
		Path:     p,
		Children: make(map[string]*ValidationFailure),
	}
}

var ErrValidationFailed = errors.New("Validation failure")

type ValidationErrors types.Pojo

func (v ValidationErrors) Error() string {
	return ErrValidationFailed.Error()
}

func (v ValidationErrors) ToPojo() types.Pojo {
	return types.Pojo(v)
}

func (v ValidationErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}

type ValidationFailure struct {
	Path     string
	Failures []string
	Children map[string]*ValidationFailure
}

func (e *ValidationFailure) Error() string {
	return ErrValidationFailed.Error()
}

func (e *ValidationFailure) ToPojo() types.Pojo {
	if len(e.Failures) == 0 && len(e.Children) == 0 {
		return nil
	}

	var result = make(types.Pojo)
	if len(e.Failures) > 0 {
		result[".failures"] = e.Failures
	}
	for k, v := range e.Children {
		result[k] = v.ToPojo()
	}
	return result
}

func (e *ValidationFailure) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ToPojo())
}

func (e *ValidationFailure) Apply(err error) *ValidationFailure {
	switch typedErr := err.(type) {
	case *jsv.ValidationError:
		return e.applyJsvValidationError(typedErr)

	case types.ErrorList:
		for _, v := range typedErr {
			e.Failures = append(e.Failures, v.Error())
		}

	case types.ErrorMap:
		for k, v := range typedErr {
			child := new(ValidationFailure)
			e.Children[k] = child.Apply(v)
		}

	default:
		e.Failures = append(e.Failures, err.Error())
	}

	return e
}

func (e *ValidationFailure) applyJsvValidationError(err *jsv.ValidationError) *ValidationFailure {
	if err.InstanceLocation != e.Path {
		// Walk to next child
		suffix := strings.TrimPrefix(err.InstanceLocation, e.Path)
		suffixParts := strings.SplitN(suffix, "/", 3)
		child, ok := e.Children[suffixParts[1]]
		if !ok {
			// Create new child
			child = NewValidationFailure(path.Join(e.Path, suffixParts[1]))
			e.Children[suffixParts[1]] = child
		}

		child.Apply(err)
		return e
	}

	if len(err.Causes) == 0 {
		e.Failures = append(e.Failures, err.Message)
		return e
	}

	for _, cause := range err.Causes {
		e.Apply(cause)
	}

	return e
}

type ValidationSchema struct {
	schema *jsv.Schema
}

func NewValidationSchema(schema *jsv.Schema) ValidationSchema {
	for schema.Ref != nil {
		schema = schema.Ref
	}
	return ValidationSchema{schema: schema}
}

func (s ValidationSchema) Types() []string {
	var st = make(types.StringSet)
	st.AddAll(anyType.schema.Types...)

	if s.schema.Types != nil {
		schemaTypes := s.schema.Types
		intersectionTypes := st.Intersect(types.NewStringSet(schemaTypes...))
		st = types.NewStringSet(intersectionTypes...)
	}

	for _, allOfOne := range s.AllOf() {
		allOfOneTypes := allOfOne.Types()
		intersectionTypes := st.Intersect(types.NewStringSet(allOfOneTypes...))
		st = types.NewStringSet(intersectionTypes...)
	}

	return st.Values()
}

func (s ValidationSchema) Property(key string) ValidationSchema {
	// Explicitly declared properties
	if ps, ok := s.schema.Properties[key]; ok {
		return NewValidationSchema(ps)
	}

	if s.schema.AdditionalProperties == nil {
		return NewValidationSchema(new(jsv.Schema))
	}

	if allowed, ok := s.schema.AdditionalProperties.(bool); ok {
		if allowed {
			return anyType
		}
	}

	if ps, ok := s.schema.AdditionalProperties.(*jsv.Schema); ok {
		return NewValidationSchema(ps)
	}

	return ValidationSchema{}
}

func (s ValidationSchema) Items() ValidationSchema {
	itemsSchema := s.schema.Items2020
	return NewValidationSchema(itemsSchema)
}

func (s ValidationSchema) AllOf() []ValidationSchema {
	var results []ValidationSchema
	for _, allOfOne := range s.schema.AllOf {
		results = append(results, NewValidationSchema(allOfOne))
	}
	return results
}

func (s ValidationSchema) IsPropertyAllowed(key string) bool {
	// Explicitly declared properties
	if _, ok := s.schema.Properties[key]; ok {
		return true
	}

	if s.schema.AdditionalProperties == nil {
		return false
	}

	if allowed, ok := s.schema.AdditionalProperties.(bool); ok {
		return allowed
	}

	if _, ok := s.schema.AdditionalProperties.(*jsv.Schema); ok {
		return true
	}

	return false

}

func (s ValidationSchema) Validate(value interface{}) error {
	return s.schema.Validate(value)
}
