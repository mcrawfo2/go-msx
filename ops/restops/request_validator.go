// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/swaggest/jsonschema-go"
	"io"
	"mime/multipart"
	"strconv"
)

type RequestValidator struct {
	port    *ops.Port
	decoder ops.InputDecoder
}

func NewRequestValidator(port *ops.Port, decoder ops.InputDecoder) RequestValidator {
	return RequestValidator{
		port:    port,
		decoder: decoder,
	}
}

func (v RequestValidator) ValidateRequest() (err error) {
	if v.port == nil {
		return nil
	}

	errs := &ops.ValidationFailure{
		Path:     "request",
		Children: make(map[string]*ops.ValidationFailure),
	}

	for _, field := range v.port.Fields {
		// Skip validation if disabled for this field
		if do, ok := field.BoolOption(ops.PortFieldTagValidate); ok && !do {
			continue
		}

		validationErr := v.ValidateField(field)
		if validationErr != nil {
			switch typedErr := validationErr.(type) {
			case *jsv.ValidationError:
				errs.Children[field.Name] = ops.NewValidationFailure(typedErr.InstanceLocation).Apply(typedErr)
			default:
				return validationErr
			}

		}
	}

	if len(errs.Children) > 0 {
		return errs
	}
	return nil
}

func (v RequestValidator) ValidateField(field *ops.PortField) (err error) {
	var validationSchema js.ValidationSchema
	validationSchema, err = GetPortFieldValidationSchema(field)
	if err != nil {
		return err
	}

	var value interface{}
	value, err = v.GetFieldValue(field, validationSchema)
	if err != nil {
		return
	}

	if value == nil && field.Optional {
		return
	}

	return validationSchema.Validate(value)
}

func (v RequestValidator) GetFileContents(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	fileContents, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(fileContents), nil
}

func (v RequestValidator) GetFieldValue(field *ops.PortField, validationSchema js.ValidationSchema) (value interface{}, err error) {
	switch field.Type.Shape {
	case ops.FieldShapePrimitive:
		value, err = v.GetPrimitive(field, validationSchema)
		if err != nil {
			return
		}

	case ops.FieldShapeArray:
		value, err = v.GetArray(field, validationSchema)
		if err != nil {
			return
		} else if value.([]interface{}) == nil {
			value = nil
		}

	case ops.FieldShapeObject:
		value, err = v.GetObject(field, validationSchema)
		if err != nil {
			return
		} else if value.(types.Pojo) == nil {
			value = nil
		}

	case ops.FieldShapeFile:
		var fileHeader *multipart.FileHeader
		fileHeader, err = v.GetFile(field, validationSchema)
		if err != nil {
			return
		} else if fileHeader == nil {
			value = nil
		} else {
			// Convert the file contents to a string
			value, err = v.GetFileContents(fileHeader)
			if err != nil {
				return nil, err
			}
		}

	case ops.FieldShapeFileArray:
		var fileHeaders []*multipart.FileHeader
		fileHeaders, err = v.GetFileArray(field, validationSchema)
		if err != nil {
			return
		} else if len(fileHeaders) == 0 {
			value = nil
		} else {
			// Convert the file contents to an array of strings
			var results []string
			for _, fileHeader := range fileHeaders {
				var fileContents string
				fileContents, err = v.GetFileContents(fileHeader)
				if err != nil {
					return nil, err
				}
				results = append(results, fileContents)
			}
			value = results
		}

	case ops.FieldShapeContent:
		value, err = v.GetPayloadAsParsedJson(field)
		if err != nil {
			return
		}

	case ops.FieldShapeAny:
		value, err = v.GetAny(field)
		if err != nil {
			return
		} else if value.(any) == nil {
			value = nil
		}
	}

	return value, nil
}

func (v RequestValidator) TypesHasType(types []string, simpleType jsonschema.SimpleType) bool {
	for _, t := range types {
		if string(simpleType) == t {
			return true
		}
	}
	return false
}

func (v RequestValidator) GetPrimitiveElement(value string, types []string) (interface{}, error) {
	switch {
	case v.TypesHasType(types, jsonschema.String):
		return value, nil

	case v.TypesHasType(types, jsonschema.Number):
		return strconv.ParseFloat(value, 64)

	case v.TypesHasType(types, jsonschema.Integer):
		i64, err := strconv.ParseInt(value, 10, 64)
		return int(i64), err

	case v.TypesHasType(types, jsonschema.Boolean):
		return strconv.ParseBool(value)

	case v.TypesHasType(types, jsonschema.Array):
		return nil, errors.New("Cannot convert string to non-primitive type 'array'")

	case v.TypesHasType(types, jsonschema.Object):
		return nil, errors.New("Cannot convert string to non-primitive type 'object'")
	}

	return nil, errors.Errorf("Cannot determine target type of schema %+v", types)

}

func (v RequestValidator) GetPrimitive(field *ops.PortField, schema js.ValidationSchema) (interface{}, error) {
	optionalValue, err := v.decoder.DecodePrimitive(field)
	if err != nil {
		return nil, err
	}

	if !optionalValue.IsPresent() {
		return nil, nil
	}

	return v.GetPrimitiveElement(optionalValue.Value(), schema.Types())
}

func (v RequestValidator) GetPayloadAsParsedJson(field *ops.PortField) (interface{}, error) {
	content, err := v.decoder.DecodeContent(field)
	if err != nil {
		return nil, err
	}

	contentType, err := content.BaseMediaType()
	if err != nil {
		return nil, err
	}

	if contentType != MediaTypeJson {
		// Need to round-trip via the field DTO
		return nil, errors.Errorf("Unsupported content format for JSON Schema validation: %s", contentType)
	}

	var parsed interface{}
	if err = content.ReadEntity(&parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}

func (v RequestValidator) GetArray(field *ops.PortField, schema js.ValidationSchema) ([]interface{}, error) {
	values, err := v.decoder.DecodeArray(field)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, nil
	}

	var results = make([]interface{}, len(values))
	var value interface{}
	for i, str := range values {
		items := schema.Items()
		if value, err = v.GetPrimitiveElement(str, items.Types); err != nil {
			return nil, err
		} else {
			results[i] = value
		}
	}

	return results, nil
}

func (v RequestValidator) GetObject(field *ops.PortField, schema js.ValidationSchema) (types.Pojo, error) {
	values, err := v.decoder.DecodeObject(field)
	if err != nil {
		return nil, err
	}

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

func (v RequestValidator) GetAny(field *ops.PortField) (any, error) {
	optionalValue, err := v.decoder.DecodeAny(field)
	if err != nil {
		return nil, err
	}

	if !optionalValue.IsPresent() {
		return nil, nil
	}

	return optionalValue.Value(), nil
}

func (v RequestValidator) GetFile(field *ops.PortField, schema js.ValidationSchema) (*multipart.FileHeader, error) {
	value, err := v.decoder.DecodeFile(field)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (v RequestValidator) GetFileArray(field *ops.PortField, schema js.ValidationSchema) ([]*multipart.FileHeader, error) {
	value, err := v.decoder.DecodeFileArray(field)
	if err != nil {
		return nil, err
	}

	return value, nil
}

var portFieldValidatorFunc ops.PortFieldValidationSchemaFunc

func RegisterPortFieldValidationSchemaFunc(validatorFunc ops.PortFieldValidationSchemaFunc) {
	portFieldValidatorFunc = validatorFunc
}

func GetPortFieldValidationSchema(field *ops.PortField) (js.ValidationSchema, error) {
	if portFieldValidatorFunc == nil {
		return js.ValidationSchema{}, errors.New("No port field validation schema handler registered.")
	}

	return portFieldValidatorFunc(field)
}
