package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/sanitize"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"net/http"
	"reflect"
	"strconv"
)

type OpenApiRequestPopulator struct {
	Decoder RequestDecoder
}

func (p OpenApiRequestPopulator) PopulateInputs(e Endpoint) (interface{}, error) {
	result, err := p.PopulatePortStruct(e.Request)
	if err != nil {
		return nil, err
	}

	// Auto-validation for validatable Port Fields
	err = e.Request.Port.Validate(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p OpenApiRequestPopulator) PopulatePortStruct(r EndpointRequest) (interface{}, error) {
	var err error

	inputs := r.Port.NewStruct()
	iv := reflect.ValueOf(inputs).Elem()
	for _, portField := range r.Port.Fields {

		fv := iv.FieldByIndex(portField.Field.Index)

		var param EndpointRequestParameter
		switch portField.In {
		case "body":
			param = r.bodyParameter()
			err = p.populateBody(param, fv)
			continue

		case "form":
			param = portField.Parameter()
			fieldSchema := r.Body.Schema.Schema.Properties[param.Name]
			param.Schema = &fieldSchema

		default:
			param = r.parameterByName(portField.Name)
		}

		switch portField.Shape {
		case FieldShapePrimitive:
			err = p.populatePrimitive(param, fv)
			if err != nil {
				return inputs, err
			}

		case FieldShapeArray:
			err = p.populateArray(param, fv)
			if err != nil {
				return inputs, err
			}

		case FieldShapeObject:
			err = p.populateObject(param, fv)
			if err != nil {
				return inputs, err
			}

		case FieldShapeFile:
			err = p.populateFile(param, fv)
			if err != nil {
				return inputs, err
			}

		case FieldShapeFileArray:
			err = p.populateFiles(param, fv)
			if err != nil {
				return inputs, err
			}

		}

	}

	return inputs, err
}

func (p OpenApiRequestPopulator) populateBody(param EndpointRequestParameter, fieldValue reflect.Value) (err error) {
	return p.Decoder.DecodeBody(fieldValue.Addr().Interface(), !*param.Required)
}

func (p OpenApiRequestPopulator) populatePrimitive(param EndpointRequestParameter, fieldValue reflect.Value) (err error) {
	optionalValue, err := p.Decoder.DecodePrimitive(
		param.In,
		param.Name,
		*param.Style,
		*param.Explode)
	if err != nil {
		return err
	}

	if !optionalValue.IsPresent() && *param.Required {
		return errors.Errorf("Missing non-optional form field %q", param.Name)
	} else if optionalValue.IsPresent() {
		err = p.populateScalar(param, fieldValue, optionalValue.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p OpenApiRequestPopulator) populateArray(param EndpointRequestParameter, fieldValue reflect.Value) error {
	values, err := p.Decoder.DecodeArray(
		param.In, param.Name,
		*param.Style,
		*param.Explode)
	if err != nil {
		return err
	}

	if len(values) == 0 && !*param.Required {
		return nil
	}

	sliceType := fieldValue.Type()
	isPtr := sliceType.Kind() == reflect.Ptr
	if isPtr {
		sliceType = sliceType.Elem()
	}

	var sliceValue reflect.Value
	if sliceType.Kind() == reflect.Slice {
		sliceValue = reflect.MakeSlice(sliceType, len(values), len(values))
	} else {
		// TODO: test this
		sliceValue = reflect.New(sliceType).Elem()
	}

	for i, queryValue := range values {
		err = p.populateScalar(param, sliceValue.Index(i), queryValue)
		if err != nil {
			return err
		}
	}

	if isPtr {
		x := reflect.New(sliceType)
		x.Elem().Set(sliceValue)
		fieldValue.Set(x)
	} else {
		fieldValue.Set(sliceValue)
	}

	return nil
}

func (p OpenApiRequestPopulator) populateFile(param EndpointRequestParameter, fieldValue reflect.Value) error {
	file, err := p.Decoder.DecodeFormFile(param.Name, !*param.Required)
	if err != nil {
		return err
	}

	fieldValue.Set(reflect.ValueOf(file))

	return nil
}

func (p OpenApiRequestPopulator) populateFiles(param EndpointRequestParameter, fieldValue reflect.Value) error {
	files, err := p.Decoder.DecodeFormMultiFile(param.Name)
	if err != nil {
		return err
	}

	fieldValue.Set(reflect.ValueOf(files))

	return nil
}

func (p OpenApiRequestPopulator) populateObject(param EndpointRequestParameter, fieldValue reflect.Value) (err error) {
	pojo, err := p.Decoder.DecodeObject(
		param.In, param.Name,
		*param.Style,
		*param.Explode)
	if err != nil {
		return err
	}

	if pojo == nil && !*param.Required {
		return nil
	}

	defer func() {
		if v := recover(); v != nil {
			err = errors.Errorf("Cannot marshal object %q into field %q: %s", pojo, param.Name, v)
		}
	}()

	objectType := fieldValue.Type()
	isPtr := objectType.Kind() == reflect.Ptr
	if isPtr {
		objectType = objectType.Elem()
	}

	var objectRef reflect.Value
	var objectValue reflect.Value
	switch objectType.Kind() {
	case reflect.Map:
		objectValue = reflect.MakeMapWithSize(objectType, len(pojo))
		objectRef = reflect.New(objectType)
		objectRef.Elem().Set(objectValue)

		keyType := objectType.Key()
		valueType := objectType.Elem()
		for k, v := range pojo {
			entryKey := reflect.New(keyType).Elem()
			if err = p.populateScalar(param, entryKey, k); err != nil {
				return err
			}

			entryValue := reflect.New(valueType).Elem()
			if err = p.populateScalar(param, entryValue, cast.ToString(v)); err != nil {
				return err
			}

			objectValue.SetMapIndex(entryKey, entryValue)
		}
		break
	case reflect.Struct:
		objectRef = reflect.New(objectType)
		objectValue = objectRef.Elem()

		for i := 0; i < objectType.NumField(); i++ {
			var value string
			sf := objectType.Field(i)

			name := strcase.ToLowerCamel(sf.Name)
			value, err = pojo.StringValue(name)
			if err == nil {
				entryValue := reflect.New(sf.Type).Elem()
				if err = p.populateScalar(param, entryValue, cast.ToString(value)); err != nil {
					return err
				}
				objectValue.FieldByIndex(sf.Index).Set(entryValue)
			} else {
				logger.WithError(err).Errorf("Failed to populate field %q", sf.Name)
			}
		}
		break

	}

	if isPtr {
		fieldValue.Set(objectRef)
	} else {
		fieldValue.Set(objectValue)
	}

	return nil
}

func (p OpenApiRequestPopulator) populateScalar(param EndpointRequestParameter, fieldValue reflect.Value, value string) (err error) {
	errorWrapper := func(err error) error {
		return errors.Wrapf(err, "Cannot marshal string %q into field %q", value, param.Name)
	}

	defer func() {
		if v := recover(); v != nil {
			err = errors.Errorf("Cannot marshal string %q into field %q: %s", value, param.Name, v)
		}
	}()

	// Scalars
	switch fieldValue.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))
		if param.PortField.Options["san"] == "true" {
			if err = sanitize.Input(fieldValue.Addr().Interface(), param.PortField.SanitizeOptions); err != nil {
				return err
			}
		}
		return nil

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return errorWrapper(err)
		}
		fieldValue.Set(reflect.ValueOf(boolValue).Convert(fieldValue.Type()))
		return nil

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errorWrapper(err)
		}
		fieldValue.Set(reflect.ValueOf(floatValue).Convert(fieldValue.Type()))
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return errorWrapper(err)
		}
		fieldValue.Set(reflect.ValueOf(intValue).Convert(fieldValue.Type()))
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return errorWrapper(err)
		}
		fieldValue.Set(reflect.ValueOf(uintValue).Convert(fieldValue.Type()))
		return nil
	}

	// Nil Pointers
	fieldType := fieldValue.Type()
	switch fieldValue.Kind() {
	case reflect.Map, reflect.Ptr, reflect.Slice, reflect.Interface:
		if fieldValue.IsNil() {
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue.Set(reflect.New(fieldType.Elem()))
			} else {
				fieldValue.Set(reflect.New(fieldType).Elem())
			}
		}
	}

	// Pointers to scalars
	if fieldValue.Kind() == reflect.Ptr {
		switch fieldValue.Elem().Kind() {
		case reflect.Slice: // bytes
			convertedValue := reflect.ValueOf(value).Convert(fieldValue.Elem().Type()).Interface().([]byte)
			ptrValue := &convertedValue
			fieldValue.Set(reflect.ValueOf(ptrValue))
			if param.PortField.Options["san"] == "true" {
				if err = sanitize.Input(fieldValue.Interface(), param.PortField.SanitizeOptions); err != nil {
					return err
				}
			}
			return nil

		case reflect.String:
			ptrValue := &value
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			if param.PortField.Options["san"] == "true" {
				if err = sanitize.Input(fieldValue.Interface(), param.PortField.SanitizeOptions); err != nil {
					return err
				}
			}
			return nil

		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return errorWrapper(err)
			}
			ptrValue := &boolValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Float32:
			floatValue, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := float32(floatValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Float64:
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := float64(floatValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int:
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int8:
			intValue, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int8(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int16:
			intValue, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int16(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int32:
			intValue, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int32(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Int64:
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := int64(intValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint:
			uintValue, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint8:
			uintValue, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint8(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint16:
			uintValue, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint16(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint32:
			uintValue, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint32(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil

		case reflect.Uint64:
			uintValue, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return errorWrapper(err)
			}
			targetValue := uint64(uintValue)
			ptrValue := &targetValue
			fieldValue.Set(reflect.ValueOf(ptrValue).Convert(fieldValue.Type()))
			return nil
		}
	}

	// Unmarshaler
	fieldInterface := fieldValue.Interface()
	if fieldUnmarshaler, ok := fieldInterface.(types.TextUnmarshaler); ok {
		err = fieldUnmarshaler.UnmarshalText(value)
		if err != nil {
			err = errors.Wrap(err, param.Name)
			return NewBadRequestError(err)
		} else {
			return nil
		}
	} else {
		fieldPointerInterface := fieldValue.Addr().Interface()
		if fieldUnmarshaler, ok := fieldPointerInterface.(types.TextUnmarshaler); ok {
			err = fieldUnmarshaler.UnmarshalText(value)
			if err != nil {
				err = errors.Wrap(err, param.Name)
				return NewBadRequestError(err)
			} else {
				return nil
			}
		}
	}

	return NewInternalError(errors.Errorf("Cannot marshal string %q into field %q", value, param.Name))
}

type OpenApiResponsePopulator struct {
	Endpoint Endpoint
	Outputs  interface{}
	Error    error

	Observer  ResponseObserver
	Encoder   ResponseEncoder
	Describer RequestDescriber
}

func (p OpenApiResponsePopulator) notifyObserver(code int) {
	if p.Error != nil {
		p.Observer.Error(code, p.Error)
	} else {
		p.Observer.Success(code)
	}
}

func (p OpenApiResponsePopulator) PopulateOutputs() (err error) {
	// Calculate code
	code := p.EvaluateResponseCode()

	// Notify observers on exit
	defer p.notifyObserver(code)

	// Encode code
	if err = p.Encoder.EncodeCode(code); err != nil {
		return errors.Wrap(err, "Failed to set response status code")
	}

	// Encode headers
	err = p.PopulateHeaders()
	if err != nil {
		return errors.Wrap(err, "Failed to set response headers")
	}

	// Evaluate body
	var body interface{}
	if p.Error != nil {
		body, err = p.EvaluateErrorBody(code)
	} else {
		body, err = p.EvaluateSuccessBody(code)
	}

	if err != nil {
		return errors.Wrap(err, "Failed to set response body")
	}

	return p.Encoder.EncodeBody(body)
}

func (p OpenApiResponsePopulator) EvaluateResponseCode() (code int) {
	code = p.Endpoint.Response.Codes.DefaultCode()

	if p.Error != nil {
		code = http.StatusBadRequest
		if codeErr, ok := p.Error.(StatusCodeProvider); ok {
			code = codeErr.StatusCode()
		}
		return
	}

	if !p.Endpoint.Response.HasBody() {
		code = http.StatusNoContent
	}

	var ok bool

	// Check the code output port
	var codePort EndpointPortField
	if codePort, ok = p.Endpoint.Response.Port.Fields.Code(); ok {
		// See if the field has a value
		fv := codePort.FieldValue(p.Outputs)
		fvi := fv.Int()
		if fvi != 0 {
			code = int(fvi)
			return
		}

		// See if the field schema has a default
		defaultOverride := codePort.DefaultValue()
		if defaultOverride != nil {
			code = cast.ToInt(*defaultOverride)
			return
		}
	}

	// Check the response body output port
	var bodyPort EndpointPortField
	if bodyPort, ok = p.Endpoint.Response.Port.Fields.Body(); ok {
		bv := bodyPort.FieldValue(p.Outputs)
		bvi := bv.Interface()
		if bodyCodeProvider, ok := bvi.(StatusCodeProvider); ok {
			code = bodyCodeProvider.StatusCode()
			return
		}
	}

	return code
}

func (p OpenApiResponsePopulator) EvaluateSuccessBody(code int) (interface{}, error) {
	bodyPortField, ok := p.Endpoint.Response.Port.Fields.Body()
	if !ok {
		return nil, nil
	}

	body := bodyPortField.FieldValue(p.Outputs).Interface()

	if p.Endpoint.Response.Envelope {
		// Automatically generate the envelope
		if body == nil {
			body = struct{}{}
		}

		var envelope integration.MsxEnvelope
		if bodyEnvelope, ok := body.(integration.MsxEnvelope); ok {
			envelope = bodyEnvelope
		} else if bodyPointerEnvelope, ok := body.(*integration.MsxEnvelope); ok {
			envelope = *bodyPointerEnvelope
		} else {
			envelope = integration.MsxEnvelope{
				Success: true,
				Payload: body,
				Command: p.Endpoint.OperationID,
				Params:  p.Describer.Parameters(),
			}
		}

		if envelope.HttpStatus == "" {
			envelope.HttpStatus = integration.GetSpringStatusNameForCode(code)
		}

		body = envelope
	}

	return body, nil
}

type Pojoer interface {
	ToPojo() types.Pojo
}

func (p OpenApiResponsePopulator) EvaluateErrorBody(code int) (interface{}, error) {
	var payload interface{}

	if p.Endpoint.Response.Envelope {
		payload = new(integration.MsxEnvelope)
	} else if p.Endpoint.Response.Error.Payload == nil {
		payload = new(ErrorV8)
	} else {
		payload = *p.Endpoint.Response.Error.Payload
		payloadType := reflect.TypeOf(payload)
		if payloadType.Kind() != reflect.Ptr {
			payload = reflect.New(payloadType).Interface()
		}
	}

	switch payload.(type) {
	case *integration.MsxEnvelope:
		envelope := integration.MsxEnvelope{
			Success:    false,
			Message:    p.Error.Error(),
			Command:    p.Endpoint.OperationID,
			Params:     p.Describer.Parameters(),
			HttpStatus: integration.GetSpringStatusNameForCode(code),
			Throwable:  integration.NewThrowable(p.Error),
		}

		var errorList types.ErrorList
		if errors.As(p.Error, &errorList) {
			envelope.Errors = errorList.Strings()
		}

		var pojo Pojoer
		if errors.As(p.Error, &pojo) {
			envelope.Debug = pojo.ToPojo()
		}

		return envelope, nil

	case ErrorApplier:
		envelope := reflect.New(reflect.TypeOf(payload).Elem()).Interface().(ErrorApplier)
		envelope.ApplyError(p.Error)
		return envelope, nil

	case ErrorRaw:
		envelope := reflect.New(reflect.TypeOf(payload).Elem()).Interface().(ErrorRaw)
		envelope.SetError(code, p.Error, p.Describer.Path())
		return envelope, nil

	default:
		return nil, errors.Wrapf(p.Error, "Response serialization failed - invalid error payload type %T", payload)
	}
}

func (p OpenApiResponsePopulator) PopulateHeaders() (err error) {
	var headerPortFields EndpointPortFields
	if p.Error != nil {
		headerPortFields = p.Endpoint.Response.Port.Fields.ErrorHeaders()
	} else {
		headerPortFields = p.Endpoint.Response.Port.Fields.Headers()
	}

	// Set headers from outputs
	for _, headerPortField := range headerPortFields {
		param := headerPortField.Parameter()
		style := types.NewOptionalString(param.Style).OrElse("simple")
		explode := types.NewOptionalBool(param.Explode).OrElse(false)

		switch headerPortField.Shape {
		case FieldShapePrimitive:
			value := headerPortField.PrimitiveValue(p.Outputs)
			err = p.Encoder.EncodeHeaderPrimitive(headerPortField.Name, value, style, explode)

		case FieldShapeArray:
			value := headerPortField.ArrayValue(p.Outputs)
			err = p.Encoder.EncodeHeaderArray(headerPortField.Name, value, style, explode)

		case FieldShapeObject:
			value := headerPortField.ObjectValue(p.Outputs)
			err = p.Encoder.EncodeHeaderObject(headerPortField.Name, value, style, explode)

		default:
			err = errors.Errorf("Unable to encode %T as header", headerPortField.FieldValue(p.Outputs).Interface())
		}

		if err != nil {
			logger.WithError(err).Errorf("Failed to encode header %q", headerPortField.Name)
		}
	}

	return nil
}
