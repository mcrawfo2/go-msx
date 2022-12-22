// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/refl"
	"net/http"
	"reflect"
	"strconv"
)

type EndpointResponseDocumentor struct {
	Endpoint  *restops.Endpoint
	Skip      bool
	Responses *openapi3.Responses
	Mutator   ops.DocumentElementMutator[openapi3.Responses]
}

func (d *EndpointResponseDocumentor) WithSkip(skip bool) *EndpointResponseDocumentor {
	d.Skip = skip
	return d
}

func (d *EndpointResponseDocumentor) WithResponses(op *openapi3.Responses) *EndpointResponseDocumentor {
	d.Responses = op
	return d
}

func (d *EndpointResponseDocumentor) WithMutator(mutator ops.DocumentElementMutator[openapi3.Responses]) *EndpointResponseDocumentor {
	d.Mutator = mutator
	return d
}

func (d *EndpointResponseDocumentor) WithEndpoint(endpoint *restops.Endpoint) *EndpointResponseDocumentor {
	d.Endpoint = endpoint
	return d
}

func (d *EndpointResponseDocumentor) DocType() string {
	return DocType
}

func (d *EndpointResponseDocumentor) Document(r *restops.EndpointResponse) error {
	if d.Skip {
		return nil
	}

	if d.Responses == nil {
		d.Responses = new(openapi3.Responses)
	}

	for _, code := range r.Codes.Success {
		successResponseOrRef, err := d.documentResponseContent(r.Success, r.Envelope, code)
		if err != nil {
			return err
		}

		d.Responses.WithMapOfResponseOrRefValuesItem(
			strconv.Itoa(code),
			*successResponseOrRef)
	}

	for _, code := range r.Codes.Error {
		errorResponseOrRef, err := d.documentResponseContent(r.Error, r.Envelope, code)
		if err != nil {
			return err
		}

		d.Responses.WithMapOfResponseOrRefValuesItem(
			strconv.Itoa(code),
			*errorResponseOrRef)
	}

	if d.Mutator != nil {
		d.Mutator(d.Responses)
	}

	return nil
}

func (d *EndpointResponseDocumentor) Result() *openapi3.Responses {
	return d.Responses
}

func (d *EndpointResponseDocumentor) documentHeader(h restops.EndpointResponseHeader) (*openapi3.HeaderOrRef, error) {
	header := openapi3.Header{
		Description:     h.Description,
		Required:        h.Required,
		Deprecated:      h.Deprecated,
		AllowEmptyValue: h.AllowEmptyValue,
		Explode:         h.Explode,
		AllowReserved:   h.AllowReserved,
		//Content:         h.Content,
	}

	if h.PortField != nil {
		schemaOrRef := openApiSchemaFromPortField(h.PortField)
		header.WithSchema(*schemaOrRef)
	} else if h.Payload.IsPresent() {
		schemaOrRef, err := Reflect(h.Payload.Value())
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to generate header schema for %T", h.Payload.Value())
		}
		header.WithSchema(*schemaOrRef)
	} else {
		header.WithSchema(NewSchemaOrRef(AnySchema()))
	}

	if h.Example.IsPresent() {
		header.WithExample(h.Example.Value())
	}

	var headerOrRef = openapi3.HeaderOrRef{Header: &header}
	if h.Reference != nil {
		Spec().ComponentsEns().HeadersEns().WithMapOfHeaderOrRefValuesItem(
			*h.Reference,
			openapi3.HeaderOrRef{
				Header: &header,
			})
		headerOrRef = NewHeaderRef(*h.Reference)
	}

	return &headerOrRef, nil
}

func (d *EndpointResponseDocumentor) payload(c restops.EndpointResponseContent) (*openapi3.SchemaOrRef, string, interface{}, error) {
	var payload interface{}

	if c.Payload.IsPresent() {
		// Fill in the example
		payload = c.Payload.Value()
	}

	if payload == nil {
		return nil, "", nil, nil
	}

	// Generate Payload schema
	payloadSchemaOrRef, err := Reflect(payload)
	if err != nil {
		return nil, "", nil, errors.Wrapf(err, "Failed to generate payload schema for %T", payload)
	}

	var payloadName string
	if payloadSchemaOrRef.SchemaReference == nil {
		// If it is not a type, derive one for the spec eg api.SomeRequest.List
		payloadName = schema.Namer().TypeName(reflect.TypeOf(payload))
		payloadSchema := payloadSchemaOrRef.Schema
		payloadSchema.Title = types.NewStringPtr(payloadName)
		Spec().ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(
			payloadName,
			NewSchemaOrRef(payloadSchema))
		payloadSchemaOrRef = NewSchemaRefPtr(payloadName)
	}

	return payloadSchemaOrRef, payloadName, payload, nil
}

func (d *EndpointResponseDocumentor) paging(c restops.EndpointResponseContent) (*openapi3.SchemaOrRef, string, interface{}, error) {
	var pagingStruct = refl.DeepIndirect(reflect.TypeOf(c.Paging.ValueInterface()))

	payloadIdx, _, err := schema.FindParameterizedStructField(pagingStruct)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "Failed to locate payload field in paging struct")
	}

	var pagingInstanceValue = reflect.New(pagingStruct).Elem()
	var pagingInstance = pagingInstanceValue.Addr().Interface()

	payloadSchemaOrRef, payloadName, payloadContent, err := d.payload(c)
	if err != nil {
		return nil, "", nil, err
	} else if payloadContent != nil {
		payloadField := pagingInstanceValue.FieldByIndex(payloadIdx)
		payloadField.Set(reflect.ValueOf(payloadContent))
	}

	pagingSchemaOrRef, err := Reflect(pagingInstance)
	if err != nil {
		return nil, "", nil, errors.Wrapf(err, "Failed to reflect schema for paging struct %T", pagingInstance)
	}

	// Create customized envelope schema by merging
	var mergeSchema = openapi3.Schema{}
	if payloadSchemaOrRef != nil {
		// Calculate the JSON property name for the schema
		pagingSchemaField := pagingStruct.FieldByIndex(payloadIdx)
		pagingSchemaFieldName := schema.
			GetJsonFieldName(pagingSchemaField).
			OrElse(strcase.ToLowerCamel(pagingSchemaField.Name))

		mergeSchema.AllOf = []openapi3.SchemaOrRef{
			*pagingSchemaOrRef,
			NewSchemaOrRef(ObjectSchema().
				WithPropertiesItem(
					pagingSchemaFieldName,
					*payloadSchemaOrRef)),
		}
	} else {
		mergeSchema.AllOf = []openapi3.SchemaOrRef{
			*pagingSchemaOrRef,
		}
	}

	var mergeSchemaName string
	if payloadName == "" {
		mergeSchemaName = schema.Namer().ParameterizedTypeName(
			reflect.TypeOf(pagingInstance),
			reflect.TypeOf(payloadContent))
	} else {
		mergeSchemaName = schema.Namer().ParameterizedTypeNameWithWrappedName(
			reflect.TypeOf(pagingInstance),
			payloadName)
	}
	mergeSchema.Title = types.NewStringPtr(mergeSchemaName)

	// Store the customized paging schema
	Spec().ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(
		mergeSchemaName,
		NewSchemaOrRef(&mergeSchema))

	// Create a reference to the customized paging schema
	pagingSchemaOrRef = NewSchemaRefPtr(mergeSchemaName)

	return pagingSchemaOrRef, mergeSchemaName, pagingInstance, nil
}

func (d *EndpointResponseDocumentor) envelope(c restops.EndpointResponseContent, code int) (*openapi3.SchemaOrRef, interface{}, error) {
	var envelopeInstance = &integration.MsxEnvelope{
		Command:    d.Endpoint.OperationID,
		HttpStatus: integration.GetSpringStatusNameForCode(code),
		Message:    "Successfully executed " + d.Endpoint.OperationID,
		Params:     restops.EndpointRequestDescriber{Endpoint: *d.Endpoint}.Parameters(),
		Success:    code < 400,
	}
	var envelopeInstanceValue = reflect.ValueOf(envelopeInstance).Elem()
	var envelopeStruct = envelopeInstanceValue.Type()

	payloadIdx, _, err := schema.FindParameterizedStructField(envelopeStruct)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to locate payload field in envelope struct")
	}

	var payloadSchemaOrRef *openapi3.SchemaOrRef
	var payloadName string
	var payloadContent interface{}

	if c.Paging.IsPresent() {
		payloadSchemaOrRef, payloadName, payloadContent, err = d.paging(c)
	} else {
		payloadSchemaOrRef, payloadName, payloadContent, err = d.payload(c)
	}

	if err != nil {
		return nil, nil, err
	}

	envelopeSchemaOrRef, err := Reflect(envelopeInstance)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to generate envelope schema for %T", envelopeInstance)
	}

	if payloadContent != nil {
		payloadField := envelopeInstanceValue.FieldByIndex(payloadIdx)
		payloadField.Set(reflect.ValueOf(payloadContent))
	}

	// Create customized envelope schema by merging
	var mergeSchema = openapi3.Schema{}
	if payloadSchemaOrRef != nil {
		mergeSchema.AllOf = []openapi3.SchemaOrRef{
			*envelopeSchemaOrRef,
			NewSchemaOrRef(ObjectSchema().
				WithPropertiesItem(
					"responseObject",
					*payloadSchemaOrRef)),
		}
	} else {
		mergeSchema.AllOf = []openapi3.SchemaOrRef{
			*envelopeSchemaOrRef,
		}
	}

	var mergeSchemaName string
	if payloadName == "" {
		mergeSchemaName = schema.Namer().ParameterizedTypeName(
			reflect.TypeOf(envelopeInstance),
			reflect.TypeOf(payloadContent))
	} else {
		mergeSchemaName = schema.Namer().ParameterizedTypeNameWithWrappedName(
			reflect.TypeOf(envelopeInstance),
			payloadName)
	}
	mergeSchema.Title = types.NewStringPtr(mergeSchemaName)

	// Store the customized envelope schema
	Spec().ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(
		mergeSchemaName,
		NewSchemaOrRef(&mergeSchema))

	// Create a reference to the customized envelope schema
	envelopeSchemaOrRef = NewSchemaRefPtr(mergeSchemaName)

	if code >= 400 {
		envelopeInstance.Message = "Failed to execute " + d.Endpoint.OperationID
		envelopeInstance.Errors = []string{
			"Service returned " + http.StatusText(code),
		}
		envelopeInstance.Throwable = &integration.Throwable{
			Message: "Service returned " + http.StatusText(code),
		}
	}

	return envelopeSchemaOrRef, &envelopeInstance, nil
}

func (d *EndpointResponseDocumentor) documentHeaders(c restops.EndpointResponseContent, code int, response *openapi3.Response) error {
	for name, header := range c.Headers {
		headerOrRef, err := d.documentHeader(header)
		if err != nil {
			return errors.Wrap(err, name)
		}

		response.WithHeadersItem(name, *headerOrRef)
	}

	return nil
}

func (d *EndpointResponseDocumentor) documentEnvelopeResponse(c restops.EndpointResponseContent, code int) (*openapi3.ResponseOrRef, error) {
	var result = new(openapi3.Response)

	result.Description = http.StatusText(code)

	if err := d.documentHeaders(c, code, result); err != nil {
		return nil, err
	}

	schemaOrRef, example, err := d.envelope(c, code)
	if err != nil {
		return nil, err
	}

	result.WithContentItem(webservice.MIME_JSON, // Envelope is always JSON
		openapi3.MediaType{
			Schema:  schemaOrRef,
			Example: &example,
			//Encoding: c.Encoding,
		})

	return &openapi3.ResponseOrRef{
		Response: result,
	}, nil
}

var errorRawInstance webservice.ErrorRaw
var errorRawType = reflect.TypeOf(&errorRawInstance).Elem()
var errorApplierInstance webservice.ErrorApplier
var errorApplierType = reflect.TypeOf(&errorApplierInstance).Elem()

func (d *EndpointResponseDocumentor) documentRawResponse(c restops.EndpointResponseContent, code int) (*openapi3.ResponseOrRef, error) {
	var result = new(openapi3.Response)

	result.Description = http.StatusText(code)

	if err := d.documentHeaders(c, code, result); err != nil {
		return nil, err
	}

	mime := c.Mime
	payload := c.Payload
	if !payload.IsPresent() {
		if code >= 400 {
			payload = types.OptionalOf[interface{}](new(webservice.ErrorV8))
		}
		if mime == "" {
			mime = webservice.MIME_JSON
		}
	}

	if c.Paging.IsPresent() {
		schemaOrRef, _, example, err := d.paging(c)
		if err != nil {
			return nil, err
		} else {
			var examplePtr *interface{}
			if example != nil {
				examplePtr = &example
			}

			result.WithContentItem(mime,
				openapi3.MediaType{
					Schema:  schemaOrRef,
					Example: examplePtr,
					//Encoding: c.Encoding,
				})
		}
	} else if payload.IsPresent() && payload.Value() != nil {
		schemaOrRef, err := Reflect(payload.Value())
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to generate payload schema for %T", payload.Value())
		}

		example := c.Example
		payloadType := reflect.PtrTo(reflect.TypeOf(payload.Value()))
		if code >= 400 && example.IsPresent() && example.Value() != nil {
			switch {
			case payloadType.Implements(errorRawType):
				errorRaw := reflect.New(payloadType.Elem()).Interface()
				errorRaw.(webservice.ErrorRaw).SetError(code, errors.New(result.Description), d.Endpoint.Path)
				example = types.OptionalOf(errorRaw)

			case payloadType.Implements(errorApplierType):
				errorRaw := reflect.New(payloadType.Elem()).Interface()
				errorRaw.(webservice.ErrorApplier).ApplyError(errors.New(result.Description))
				example = types.OptionalOf(errorRaw)
			}
		}

		exampleValue := example.ValuePtrInterface()
		var exampleValuePtr *interface{}
		if example.IsPresent() {
			exampleValuePtr = &exampleValue
		}

		result.WithContentItem(mime,
			openapi3.MediaType{
				Schema:  schemaOrRef,
				Example: exampleValuePtr,
				//Encoding: c.Encoding,
			})
	}

	return &openapi3.ResponseOrRef{
		Response: result,
	}, nil
}

func (d *EndpointResponseDocumentor) documentResponseContent(c restops.EndpointResponseContent, envelope bool, code int) (*openapi3.ResponseOrRef, error) {
	if envelope {
		return d.documentEnvelopeResponse(c, code)
	} else {
		return d.documentRawResponse(c, code)
	}
}
