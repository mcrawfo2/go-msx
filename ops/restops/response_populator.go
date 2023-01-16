// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"net/http"
	"reflect"
)

type OutputsPopulator struct {
	Endpoint *Endpoint
	Outputs  *interface{}
	Error    error

	Observer  ResponseObserver
	Encoder   ResponseEncoder
	Describer RequestDescriber
}

func (p *OutputsPopulator) notifyObserver(code int) {
	if p.Error != nil {
		p.Observer.Error(code, p.Error)
	} else {
		p.Observer.Success(code)
	}
}

func (p *OutputsPopulator) extractor(pf *ops.PortField) ops.PortFieldExtractor {
	return ops.NewPortFieldExtractor(pf, *p.Outputs)
}

func (p *OutputsPopulator) PopulateOutputs() (err error) {
	// Calculate code
	code, err := p.EvaluateResponseCode()
	if err != nil {
		return err
	}

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

	// Evaluate content type
	if code != 204 {
		var mediaType string
		mediaType, err = p.EvaluateMediaType(code)
		if err != nil {
			return errors.Wrap(err, "Failed to generate media type")
		}
		err = p.Encoder.EncodeMime(mediaType)
		if err != nil {
			return errors.Wrap(err, "Failed to set media type")
		}
	}

	// Evaluate body
	var body interface{}
	if p.Error != nil {
		body, err = p.EvaluateErrorBody(code)
	} else {
		body, err = p.EvaluateSuccessBody(code)
	}

	if err != nil {
		return errors.Wrap(err, "Failed to generate response body")
	}

	err = p.Encoder.EncodeBody(body)
	if err != nil {
		return errors.Wrap(err, "Failed to set response body")
	}

	return nil
}

type causer interface {
	Cause() error
}

func (p *OutputsPopulator) EvaluateResponseCode() (code int, err error) {
	code = p.Endpoint.Response.Codes.DefaultCode()

	if p.Error != nil {
		code = http.StatusBadRequest
		if codeErr, ok := p.Error.(webservice.StatusCodeProvider); ok {
			code = codeErr.StatusCode()
			if causeErr, ok := codeErr.(causer); ok {
				// Unwrap the error
				p.Error = causeErr.Cause()
			}
		}
		return
	}

	if !p.Endpoint.Response.HasBody() && !p.Endpoint.Response.Envelope {
		code = http.StatusNoContent
	}

	// Check the code output port
	port := p.Endpoint.Response.Port
	var codePortField *ops.PortField
	if port != nil {
		codePortField = port.Fields.First(PortFieldIsCode)
	}
	if codePortField != nil {
		var specified types.Optional[string]
		var value int

		// See if the field has a value
		specified, err = p.extractor(codePortField).ExtractPrimitive()
		if err != nil {
			return
		} else if specified.IsPresent() {
			value, err = cast.ToIntE(specified.Value())
			if err != nil {
				return
			}

			if value != 0 {
				code = value
				return
			}
		}
	}

	return
}

func (p *OutputsPopulator) EvaluateSuccessBody(code int) (interface{}, error) {
	port := p.Endpoint.Response.Port

	var bodyPortField *ops.PortField
	var body interface{}
	if port != nil {
		bodyPortField = port.Fields.First(PortFieldIsSuccessBody)
	}
	if bodyPortField != nil {
		fv, err := p.extractor(bodyPortField).ExtractValue()
		if err != nil {
			return nil, err
		}
		body = fv.Interface()
	} else if p.Endpoint.Response.Success.Payload.IsPresent() {
		body = p.Endpoint.Response.Success.Payload.Value()
	}

	var pagingPortField *ops.PortField
	if port != nil {
		pagingPortField = p.Endpoint.Response.Port.Fields.First(PortFieldIsPaging)
	}
	if pagingPortField != nil {
		fv, err := p.extractor(pagingPortField).ExtractValue()
		if err != nil {
			return nil, err
		}

		pfv := reflect.New(pagingPortField.Type.Type)
		pfv.Elem().Set(fv)

		contentIndices, _, err := schema.FindParameterizedStructField(pagingPortField.Type.Type)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to identify paging contents field")
		}

		if body != nil {
			target := pfv.Elem().FieldByIndex(contentIndices)
			target.Set(reflect.ValueOf(body))
		}

		body = pfv.Elem().Interface()
	}

	if p.Endpoint.Response.Envelope {
		// Automatically generate the envelope
		bodyType := reflect.TypeOf(body)
		bodyValue := reflect.ValueOf(body)
		for bodyValue.Kind() == reflect.Ptr {
			bodyValue = bodyValue.Elem()
			bodyType = bodyValue.Type()
		}

		var envelope integration.MsxEnvelope
		if bodyType == reflect.TypeOf(envelope) {
			envelope = bodyValue.Interface().(integration.MsxEnvelope)
		} else {
			envelope = integration.MsxEnvelope{
				Success: true,
				Payload: body,
				Command: p.Endpoint.OperationID,
				Params:  p.Describer.Parameters(),
				Message: p.Endpoint.OperationID + " succeeded",
			}
		}

		if envelope.HttpStatus == "" {
			envelope.HttpStatus = integration.GetSpringStatusNameForCode(code)
		}

		body = envelope
	}

	return body, nil
}

type pojoer interface {
	ToPojo() types.Pojo
}

func (p *OutputsPopulator) EvaluateErrorBody(code int) (interface{}, error) {
	var payload interface{}

	if p.Endpoint.Response.Envelope {
		payload = new(integration.MsxEnvelope)
	} else if p.Endpoint.Response.Error.Payload.IsPresent() {
		payload = p.Endpoint.Response.Error.Payload.Value()
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
		} else {
			envelope.Errors = []string{p.Error.Error()}
		}

		var pojo pojoer
		if errors.As(p.Error, &pojo) {
			envelope.Debug = pojo.ToPojo()
		}

		return envelope, nil

	case webservice.ErrorApplier:
		envelope := reflect.New(reflect.TypeOf(payload).Elem()).Interface().(webservice.ErrorApplier)
		envelope.ApplyError(p.Error)
		return envelope, nil

	case webservice.ErrorRaw:
		envelope := reflect.New(reflect.TypeOf(payload).Elem()).Interface().(webservice.ErrorRaw)
		envelope.SetError(code, p.Error, p.Describer.Path())
		return envelope, nil

	default:
		return nil, errors.Wrapf(p.Error, "Response serialization failed - invalid error payload type %T", payload)
	}
}

func (p *OutputsPopulator) PopulateHeaders() (err error) {
	if p.Endpoint.Response.Port == nil {
		return nil
	}

	// TODO: Code-Specific Headers
	var headerPortFields ops.PortFields
	if p.Error != nil {
		headerPortFields = p.Endpoint.Response.Port.Fields.All(PortFieldIsErrorHeader)
	} else {
		headerPortFields = p.Endpoint.Response.Port.Fields.All(PortFieldIsSuccessHeader)
	}

	// Set headers from outputs
	for _, headerPortField := range headerPortFields {
		extractor := p.extractor(headerPortField)
		style := headerPortField.Options["style"]
		explode, _ := headerPortField.BoolOption("explode")

		switch headerPortField.Type.Shape {
		case ops.FieldShapePrimitive:
			var value types.Optional[string]
			value, err = extractor.ExtractPrimitive()
			if err != nil {
				return err
			}
			err = p.Encoder.EncodeHeaderPrimitive(headerPortField.Peer, value, style, explode)

		case ops.FieldShapeArray:
			var value []string
			value, err = extractor.ExtractArray()
			if err != nil {
				return err
			}
			err = p.Encoder.EncodeHeaderArray(headerPortField.Peer, value, style, explode)

		case ops.FieldShapeObject:
			var value types.Pojo
			value, err = extractor.ExtractObject()
			if err != nil {
				return err
			}
			err = p.Encoder.EncodeHeaderObject(headerPortField.Peer, value, style, explode)

		default:
			err = errors.Errorf("Unable to encode shape %s as header", headerPortField.Type.Shape)
			return
		}

		if err != nil {
			err = errors.Wrapf(err, "Failed to encode header %q", headerPortField.Name)
			return
		}
	}

	return nil
}

func (p *OutputsPopulator) EvaluateMediaType(code int) (mediaType string, err error) {
	if p.Endpoint.Response.Envelope {
		return MediaTypeJson, nil
	} else if code <= 399 {
		return p.Endpoint.Response.Success.Mime, nil
	} else {
		return p.Endpoint.Response.Error.Mime, nil
	}
}
