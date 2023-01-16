// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"github.com/pkg/errors"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/swaggest/jsonschema-go"
	"strconv"
)

type MessageValidator struct {
	port           *ops.Port
	messageDecoder MessageDecoder
}

func NewMessageValidator(port *ops.Port, decoder MessageDecoder) MessageValidator {
	return MessageValidator{
		port:           port,
		messageDecoder: decoder,
	}
}

func (v MessageValidator) ValidateMessage() (err error) {
	errs := &ops.ValidationFailure{
		Path:     "message",
		Children: make(map[string]*ops.ValidationFailure),
	}

	for _, field := range v.port.Fields {
		// Skip validation if disabled for this field
		if do, ok := field.BoolOption(ops.PortFieldTagValidate); ok && !do {
			continue
		}

		var validationSchema js.ValidationSchema
		validationSchema, err = GetPortFieldValidationSchema(field)
		if err != nil {
			return err
		}

		var value interface{}
		switch field.Type.Shape {
		case ops.FieldShapePrimitive:
			value, err = v.GetPrimitive(field, validationSchema)
			if err != nil {
				return err
			}

		case ops.FieldShapeContent:
			value, err = v.GetPayloadAsParsedJson(field)
			if err != nil {
				return err
			}

		}

		err = validationSchema.Validate(value)
		if err != nil {
			switch typedErr := err.(type) {
			case *jsv.ValidationError:
				errs.Children[field.Name] = ops.NewValidationFailure(typedErr.InstanceLocation).Apply(typedErr)
			default:
				return err
			}

		}
	}

	if len(errs.Children) > 0 {
		return errors.Wrap(errs, "Validation Failure")
	}
	return nil
}

func (v MessageValidator) TypesHasType(types []string, simpleType jsonschema.SimpleType) bool {
	for _, t := range types {
		if string(simpleType) == t {
			return true
		}
	}
	return false
}

func (v MessageValidator) GetPrimitiveElement(value string, types []string) (interface{}, error) {
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

func (v MessageValidator) GetPrimitive(field *ops.PortField, schema js.ValidationSchema) (interface{}, error) {
	optionalValue, err := v.messageDecoder.DecodePrimitive(field)
	if err != nil {
		return nil, err
	}

	if !optionalValue.IsPresent() {
		return nil, nil
	}

	return v.GetPrimitiveElement(optionalValue.Value(), schema.Types())
}

func (v MessageValidator) GetPayloadAsParsedJson(field *ops.PortField) (interface{}, error) {
	content, err := v.messageDecoder.DecodeContent(field)
	if err != nil {
		return nil, err
	}

	contentType, err := content.BaseMediaType()
	if err != nil {
		return nil, err
	}

	if contentType != httpclient.MimeTypeApplicationJson {
		// Need to round-trip via the field DTO
		return nil, errors.Errorf("Unsupported content format for JSON Schema validation: %s", contentType)
	}

	var parsed interface{}
	if err = content.ReadEntity(&parsed); err != nil {
		return nil, err
	}

	return parsed, nil
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
