package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const ExtraPropertyNamePermissions = "x-msx-permissions"

func NewEndpointOpenApiDocumentor(e Endpoint) EndpointOpenApiDocumentor {
	return EndpointOpenApiDocumentor{
		Endpoint: e,
	}
}

type EndpointOpenApiDocumentor struct {
	Endpoint
}

func (e EndpointOpenApiDocumentor) RequestParameter(p EndpointRequestParameter) openapi3.ParameterOrRef {
	parameter := openapi3.Parameter{
		Name:            p.Name,
		In:              openapi3.ParameterIn(p.In),
		Description:     p.Description,
		Required:        p.Required,
		Deprecated:      p.Deprecated,
		AllowEmptyValue: p.AllowEmptyValue,
		Style:           p.Style,
		Explode:         p.Explode,
		AllowReserved:   p.AllowReserved,
		Schema:          p.Schema,
		Content:         p.Content,
		Example:         p.Example,
	}

	parameterOrRef := openapi3.ParameterOrRef{
		Parameter: &parameter,
	}

	if p.Reference != nil && (p.PortField == nil || p.PortField.BoolOption("ref")) {
		params := Reflector.SpecEns().ComponentsEns().ParametersEns()
		if _, ok := params.MapOfParameterOrRefValues[*p.Reference]; !ok {
			params.WithMapOfParameterOrRefValuesItem(*p.Reference, parameterOrRef)
		}

		parameterOrRef = openapi3.ParameterOrRef{
			ParameterReference: &openapi3.ParameterReference{
				Ref: "#/components/parameters/" + *p.Reference,
			},
		}
	}

	return parameterOrRef
}

func (e EndpointOpenApiDocumentor) RequestParameters() []openapi3.ParameterOrRef {
	r := e.Endpoint.Request
	var result []openapi3.ParameterOrRef
	for _, p := range r.Parameters {
		parameterOrRef := e.RequestParameter(p)
		result = append(result, parameterOrRef)
	}
	return result
}

func (e EndpointOpenApiDocumentor) EnvelopeResponse(c EndpointResponseContent, code int) openapi3.Response {
	var result openapi3.Response

	result.Description = http.StatusText(code)

	for name, header := range c.Headers {
		headerOrRef := e.Header(header)
		result.WithHeadersItem(name, headerOrRef)
	}

	var exampleEnvelope = integration.MsxEnvelope{
		Command:    e.OperationID,
		HttpStatus: integration.GetSpringStatusNameForCode(code),
		Message:    "Successfully executed " + e.OperationID,
		Params:     EndpointRequestDescriber{Endpoint: e.Endpoint}.Parameters(),
		Success:    code < 400,
	}
	var example interface{} = &exampleEnvelope

	schemaOrRef, _ := GetSchemaOrRefEns(exampleEnvelope)

	if c.Payload != nil {
		// Fill in the example
		payload := *c.Payload
		exampleEnvelope.Payload = payload

		// Generate Payload schema
		payloadSchemaOrRef, err := Reflect(payload)
		if err != nil {
			logger.WithError(err).Errorf("Failed to generate payload schema for %T", payload)
			return result
		}

		if payloadSchemaOrRef.SchemaReference != nil {
			payloadName := SchemaRefName(payloadSchemaOrRef.SchemaReference)
			payloadTitle := payloadName
			if payloadDirectSchemaOrRef, ok := LookupSchema(payloadName); ok {
				payloadTitle = types.
					NewOptionalString(payloadDirectSchemaOrRef.Title).
					OrElse(payloadName)
			}

			// Create customized envelope schema by merging
			var mergeSchema = openapi3.Schema{}
			mergeSchema.AllOf = []openapi3.SchemaOrRef{
				schemaOrRef,
				NewSchemaOrRef(ObjectSchema().
					WithPropertiesItem(
						"responseObject",
						*payloadSchemaOrRef)),
			}
			mergeSchema.Title = types.NewStringPtr("Envelope«" + payloadTitle + "»")
			mergeSchemaName := "Envelope" + payloadName

			// Store the customized envelope schema
			Reflector.SpecEns().ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(
				mergeSchemaName,
				NewSchemaOrRef(&mergeSchema))

			// Create a reference to the customized envelope schema
			schemaOrRef = NewSchemaRef(mergeSchemaName)
		}
	}

	if code >= 400 {
		exampleEnvelope.Message = "Failed to execute " + e.OperationID
		exampleEnvelope.Errors = []string{
			"Service returned " + http.StatusText(code),
		}
		exampleEnvelope.Throwable = &integration.Throwable{
			Message: "Service returned " + http.StatusText(code),
		}
	}

	result.WithContentItem(MIME_JSON, // Envelope is always JSON
		openapi3.MediaType{
			Schema:   &schemaOrRef,
			Example:  &example,
			Encoding: c.Encoding,
		})

	return result
}

func (e EndpointOpenApiDocumentor) Header(h EndpointResponseHeader) openapi3.HeaderOrRef {
	header := openapi3.Header{
		Description:     h.Description,
		Required:        h.Required,
		Deprecated:      h.Deprecated,
		AllowEmptyValue: h.AllowEmptyValue,
		Explode:         h.Explode,
		AllowReserved:   h.AllowReserved,
		Schema:          h.Schema,
		Content:         h.Content,
		Example:         h.Example,
	}

	headerOrRef := openapi3.HeaderOrRef{Header: &header}

	if h.Reference != nil && (h.PortField == nil || h.PortField.BoolOption("ref")) {
		params := Reflector.SpecEns().ComponentsEns().HeadersEns()
		if _, ok := params.MapOfHeaderOrRefValues[*h.Reference]; !ok {
			params.WithMapOfHeaderOrRefValuesItem(*h.Reference, headerOrRef)
		}

		headerOrRef = openapi3.HeaderOrRef{
			HeaderReference: &openapi3.HeaderReference{
				Ref: "#/components/headers/" + *h.Reference,
			},
		}
	}

	return headerOrRef
}

var errorRawInstance ErrorRaw
var errorRawType = reflect.TypeOf(&errorRawInstance).Elem()
var errorApplierInstance ErrorApplier
var errorApplierType = reflect.TypeOf(&errorApplierInstance).Elem()

func (e EndpointOpenApiDocumentor) RawResponse(c EndpointResponseContent, code int) openapi3.Response {
	var result openapi3.Response

	result.Description = http.StatusText(code)

	for name, header := range c.Headers {
		headerOrRef := e.Header(header)
		result.WithHeadersItem(name, headerOrRef)
	}

	mime := c.Mime
	payloadPtr := c.Payload
	if payloadPtr == nil {
		var payload interface{} = nil
		if code >= 400 {
			payload = new(ErrorV8)
		}
		payloadPtr = &payload
		if mime == "" {
			mime = MIME_JSON
		}
	}

	if payloadPtr != nil && *payloadPtr != nil {
		schemaOrRef, err := Reflect(*payloadPtr)
		if err != nil {
			logger.WithError(err).Errorf("Failed to generate schema for %T", *c.Payload)
			return result
		}

		example := c.Example
		payloadType := reflect.PtrTo(reflect.TypeOf(*payloadPtr))
		if code >= 400 && example == nil {
			switch {
			case payloadType.Implements(errorRawType):
				errorRaw := reflect.New(payloadType.Elem()).Interface()
				errorRaw.(ErrorRaw).SetError(code, errors.New(result.Description), e.Path)
				example = &errorRaw

			case payloadType.Implements(errorApplierType):
				errorRaw := reflect.New(payloadType.Elem()).Interface()
				errorRaw.(ErrorApplier).ApplyError(errors.New(result.Description))
				example = &errorRaw

			}
		}

		result.WithContentItem(mime,
			openapi3.MediaType{
				Schema:   schemaOrRef,
				Example:  example,
				Encoding: c.Encoding,
			})
	}

	return result
}

func (e EndpointOpenApiDocumentor) Response(c EndpointResponseContent, envelope bool, code int) openapi3.Response {
	if envelope {
		return e.EnvelopeResponse(c, code)
	} else {
		return e.RawResponse(c, code)
	}
}

func (e EndpointOpenApiDocumentor) Responses() openapi3.Responses {
	result := &openapi3.Responses{}
	r := e.Endpoint.Response

	for _, code := range r.Codes.Success {
		successResponse := e.Response(r.Success, r.Envelope, code)
		result = result.WithMapOfResponseOrRefValuesItem(
			strconv.Itoa(code),
			openapi3.ResponseOrRef{Response: &successResponse})
	}

	for _, code := range r.Codes.Error {
		errorResponse := e.Response(r.Error, r.Envelope, code)
		result = result.WithMapOfResponseOrRefValuesItem(
			strconv.Itoa(code),
			openapi3.ResponseOrRef{Response: &errorResponse})
	}

	return *result
}

func (e EndpointOpenApiDocumentor) ExtensionProperties() map[string]interface{} {
	result := make(types.Pojo)
	if len(e.Permissions) > 0 {
		result[ExtraPropertyNamePermissions] = types.Pojo{
			"anyOf": e.Permissions,
		}
	}
	return result
}

func (e EndpointOpenApiDocumentor) RequestBody() openapi3.RequestBodyOrRef {
	b := e.Endpoint.Request.Body

	result := openapi3.RequestBodyOrRef{
		RequestBody: &openapi3.RequestBody{},
	}

	result.RequestBody = result.RequestBody.
		WithRequired(b.Required).
		WithContentItem(b.Mime, openapi3.MediaType{
			Schema:   b.Schema,
			Example:  b.Example,
			Encoding: b.Encoding,
		})

	if b.Description != "" {
		result.RequestBody = result.RequestBody.WithDescription(b.Description)
	}

	return result
}

func (e EndpointOpenApiDocumentor) DocumentOpenApiOperation(op *openapi3.Operation) error {
	op = op.
		WithID(e.OperationID).
		WithTags(e.Tags...).
		WithParameters(e.RequestParameters()...).
		WithResponses(e.Responses()).
		WithMapOfAnything(e.ExtensionProperties())

	if e.Request.HasBody() {
		op = op.WithRequestBody(e.RequestBody())
	} else {
		op = op.WithRequestBody(
			*(&openapi3.RequestBodyOrRef{}).WithRequestBody(
				*(&openapi3.RequestBody{}).WithContentItem(
					"*/*", openapi3.MediaType{})))
	}

	if e.Deprecated {
		op = op.WithDeprecated(e.Deprecated)
	}

	if e.Description != "" {
		op = op.WithDescription(e.Description)
	}

	if e.Summary != "" {
		op = op.WithSummary(e.Summary)
	}

	return nil
}

// DocumentEndpoint creates an operation
func (e EndpointOpenApiDocumentor) DocumentEndpoint(s *openapi3.Spec) error {
	method := strings.ToLower(e.Method)
	path := e.Path

	pathParametersSubmatches := regexFindPathParameter.FindAllStringSubmatch(path, -1)

	switch method {
	case "get", "put", "post", "delete", "options", "head", "patch", "trace":
		break
	default:
		return fmt.Errorf("unexpected http method: %s", method)
	}

	pathItem := s.Paths.MapOfPathItemValues[e.Path]
	pathParams := map[string]bool{}

	if len(pathParametersSubmatches) > 0 {
		for _, submatch := range pathParametersSubmatches {
			pathParams[submatch[1]] = true

			if submatch[2] != "" { // Remove gorilla.Mux-style regexp in path
				path = strings.Replace(path, submatch[0], "{"+submatch[1]+"}", 1)
			}
		}
	}

	var errs []string

	operation := pathItem.MapOfOperationValues[method]

	setup := []func(*openapi3.Operation) error{
		e.DocumentOpenApiOperation,
	}

	for _, f := range setup {
		err := f(&operation)
		if err != nil {
			return err
		}
	}

	paramIndex := make(map[string]bool, len(operation.Parameters))

	for _, p := range operation.Parameters {
		p = ResolveParameter(p)
		if p.Parameter == nil {
			continue
		}

		if found := paramIndex[p.Parameter.Name+string(p.Parameter.In)]; found {
			errs = append(errs, "duplicate parameter in "+string(p.Parameter.In)+": "+p.Parameter.Name)

			continue
		}

		if found := pathParams[p.Parameter.Name]; !found && p.Parameter.In == openapi3.ParameterInPath {
			errs = append(errs, "missing path parameter placeholder in url: "+p.Parameter.Name)

			continue
		}

		paramIndex[p.Parameter.Name+string(p.Parameter.In)] = true
	}

	for pathParam := range pathParams {
		if !paramIndex[pathParam+string(openapi3.ParameterInPath)] {
			errs = append(errs, "undefined path parameter: "+pathParam)
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	pathItem.WithMapOfOperationValuesItem(method, operation)

	s.Paths.WithMapOfPathItemValuesItem(path, pathItem)

	return nil
}

func DocumentRoute(route restful.Route) {
	endpoint := EndpointFromRoute(route)
	if endpoint.Method == "" {
		endpoint = NewEndpointFromRoute(route)
	}

	// Add an endpoint to the openapi3 spec
	err := NewEndpointOpenApiDocumentor(endpoint).DocumentEndpoint(Reflector.SpecEns())
	if err != nil {
		logger.WithError(err).Errorf("Failed to register route %q %q with openapi spec", endpoint.Method, endpoint.Path)
	}
}
