package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"strings"
)

const (
	AttributeKeyEndpoint = "Endpoint"
	AttributeKeyInputs   = "Inputs"
	AttributeKeyOutputs  = "Outputs"
)

type EndpointProducer struct {
	Transformers []EndpointTransformer
	Endpoints    []Endpoint
}

type EndpointTransformer func(Endpoint) Endpoint

func AddEndpointTag(tag string) EndpointTransformer {
	return func(endpoint Endpoint) Endpoint {
		endpoint.Tags = []string{tag}
		return endpoint
	}
}

type ErrorConverter interface {
	Convert(error) StatusCodeError
}

type ErrorConverterFunc func(error) StatusCodeError

func (f ErrorConverterFunc) Convert(e error) StatusCodeError {
	return f(e)
}

func AddEndpointErrorConverter(converter ErrorConverter) EndpointTransformer {
	return func(endpoint Endpoint) Endpoint {
		endpoint.ErrorConverter = converter
		return endpoint
	}
}

func AddEndpointPathPrefix(prefix string) EndpointTransformer {
	return func(endpoint Endpoint) Endpoint {
		endpoint.Path = path.Join(prefix, endpoint.Path)
		return endpoint
	}
}

func AddEndpointRequestParameter(parameter EndpointRequestParameter) EndpointTransformer {
	return func(endpoint Endpoint) Endpoint {
		endpoint = endpoint.WithRequestParameter(parameter)
		return endpoint
	}
}

func NewEndpoint(method string, p ...string) Endpoint {
	return Endpoint{
		Method: method,
		Path:   path.Join(p...),
	}
}

func applyRouteParamToEndpoint(e Endpoint, r *RouteParam) Endpoint {
	paramData := r.Parameter

	in := "body"
	isParameter, isBody, isForm := false, false, false

	switch paramData.Kind {
	case restful.PathParameterKind:
		in = "path"
		isParameter = true
	case restful.QueryParameterKind:
		in = "query"
		isParameter = true
	case restful.HeaderParameterKind:
		in = "header"
		isParameter = true
	case restful.FormParameterKind:
		in = "form"
		isForm = true
	case restful.BodyParameterKind:
		in = "body"
		isBody = true
	}

	if isParameter {
		p := NewEndpointRequestParameter(paramData.Name, in).
			WithDescription(paramData.Description).
			WithRequired(paramData.Required)

		var schema *openapi3.Schema
		switch paramData.DataType {
		case "string":
			schema = StringSchema()
		case "integer":
			schema = IntegerSchema()
		case "number":
			schema = NumberSchema()
		case "boolean":
			schema = BooleanSchema()
		case "array":
			schema = ArraySchema(openapi3.SchemaOrRef{})
		case "object":
			schema = ObjectSchema()
		}

		if paramData.DataFormat != "" {
			schema = schema.WithFormat(paramData.DataFormat)
		}

		if paramData.DefaultValue != "" {
			// TODO: V2 DefaultValue conversion
			// schema = schema.WithDefault(...)
		}

		p = p.WithSchema(NewSchemaOrRef(schema))

		// CollectionFormat conversion
		if r.Options["csv"] == "true" {
			p = p.WithStyle("form").WithExplode(false)
		} else if r.Options["multi"] == "true" {
			p = p.WithStyle("form").WithExplode(true)
		}

		e.Request = e.Request.WithParameter(p)
	} else if isForm {
		// TODO: Verify
		e.Request.Body = e.Request.Body.WithFormField(EndpointRequestBodyFormField{
			Name:     paramData.Name,
			Required: paramData.Required,
			Schema: NewSchemaOrRefPtr(NewSchemaPtr(
				openapi3.SchemaType(paramData.DataType))),
		})

	} else if isBody {
		var examplePtr *interface{}
		example := types.Instantiate(r.Field.Type)
		schemaOrRef, _ := Reflect(example)

		e.Request = e.Request.WithBody(EndpointRequestBody{
			Description: paramData.Description,
			Required:    r.Field.Type.Kind() == reflect.Ptr,
			Schema:      schemaOrRef,
			Mime:        MIME_JSON,
			Example:     examplePtr,
			// TODO
			//Encoding: nil,
		})
	}

	return e
}

func NewEndpointFromRoute(route restful.Route) Endpoint {
	result := Endpoint{
		Method:      route.Method,
		Path:        route.Path,
		OperationID: route.Operation,
		Tags:        TagsFromRoute(route),
		Request: EndpointRequest{
			Body: EndpointRequestBody{
				Mime: "",
			},
		},
	}

	ctx := ContextWithRoute(context.Background(), &route)

	// Request
	legacyInputPortStruct := ParamsFromRoute(route)
	if legacyInputPortStruct != nil {
		inputPort, err := getRouteParams(ctx, legacyInputPortStruct)
		if err != nil {
			logger.WithError(err).Error("Failed to parse RouteParams")
		}

		for _, field := range inputPort.Fields {
			result = applyRouteParamToEndpoint(result, field)
		}
	}

	// Response
	if payload, ok := route.Metadata[MetadataEnvelope]; ok {
		result.Response = result.Response.WithEnvelope(true)
		result.Response = result.Response.WithSuccessPayload(payload)
	} else if payload, ok = route.Metadata[MetadataSuccessResponse]; ok {
		result.Response = result.Response.WithEnvelope(false)
		result.Response = result.Response.WithSuccessPayload(payload)
	} else if payload = route.ReadSample; payload != nil {
		result.Response = result.Response.WithEnvelope(false)
		result.Response = result.Response.WithSuccessPayload(payload)
	}

	for code := range route.ResponseErrors {
		if code < 399 {
			result.Response.Codes.Success = append(result.Response.Codes.Success, code)
		} else {
			result.Response.Codes.Error = append(result.Response.Codes.Error, code)
		}
	}

	return result
}

type Endpoint struct {
	Method         string
	Path           string
	OperationID    string
	Description    string
	Summary        string
	Tags           []string
	Deprecated     bool
	Permissions    []string
	Request        EndpointRequest
	Response       EndpointResponse
	Do             []RouteBuilderFunc
	Func           restful.RouteFunction
	ErrorConverter ErrorConverter
}

func (e Endpoint) WithMethod(method string) Endpoint {
	e.Method = method
	return e
}

func (e Endpoint) WithPath(parts ...string) Endpoint {
	e.Path = path.Join(parts...)
	return e
}

func (e Endpoint) WithOperationId(name string) Endpoint {
	e.OperationID = name
	return e
}

func (e Endpoint) WithDescription(description string) Endpoint {
	e.Description = description
	return e
}

func (e Endpoint) WithSummary(summary string) Endpoint {
	e.Summary = summary
	return e
}

func (e Endpoint) WithResponse(response EndpointResponse) Endpoint {
	e.Response = response
	return e
}

func (e Endpoint) WithRequest(request EndpointRequest) Endpoint {
	e.Request = request
	return e
}

func (e Endpoint) WithRequestParameter(parameter EndpointRequestParameter) Endpoint {
	e.Request = e.Request.WithParameter(parameter)
	return e
}

func (e Endpoint) WithResponseCodes(codes EndpointResponseCodes) Endpoint {
	e.Response = e.Response.WithResponseCodes(codes)
	return e
}

func (e Endpoint) WithResponseSuccessHeader(name string, header EndpointResponseHeader) Endpoint {
	e.Response = e.Response.WithSuccessHeader(name, header)
	return e
}

func (e Endpoint) WithResponseErrorHeader(name string, header EndpointResponseHeader) Endpoint {
	e.Response = e.Response.WithErrorHeader(name, header)
	return e
}

func (e Endpoint) WithResponseHeader(name string, header EndpointResponseHeader) Endpoint {
	e.Response = e.Response.WithHeader(name, header)
	return e
}

func (e Endpoint) WithDo(fn RouteBuilderFunc) Endpoint {
	e.Do = append(e.Do, fn)
	return e
}

func (e Endpoint) WithPermissionAnyOf(perms ...string) Endpoint {
	e.Permissions = perms
	e.WithDo(Permissions(perms...))
	return e
}

func (e Endpoint) WithHandler(fn interface{}) Endpoint {
	var inputPortStruct interface{}
	var outputPortStruct interface{}
	e.Func, inputPortStruct, outputPortStruct = InjectingController(fn)

	if inputPortStruct != nil && e.Request.PortStruct == nil {
		e.Request = e.Request.WithPortStruct(inputPortStruct)
	}

	if outputPortStruct != nil && e.Response.PortStruct == nil {
		e.Response = e.Response.WithPortStruct(outputPortStruct)
		if len(e.Response.Codes.Success) == 0 {
			defaultCodes := DefaultResponseCodes(e.Method)
			e.Response = e.Response.WithResponseCodes(defaultCodes)
		}
	}

	return e
}

func (e Endpoint) WithController(fn ControllerFunction) Endpoint {
	e.Func = Controller(fn)
	return e
}

func (e Endpoint) WithRouteFunc(fn restful.RouteFunction) Endpoint {
	e.Func = fn
	return e
}

func (e Endpoint) RouteBuilder(svc *restful.WebService) *restful.RouteBuilder {
	rb := svc.
		Method(e.Method).
		Path(strings.TrimPrefix(e.Path, svc.RootPath())).
		Operation(e.OperationID).
		To(e.Func).
		Do(InjectEndpoint(e)).          // Add the endpoint to the route
		Filter(InjectEndpointFilter(e)) // Add the endpoint to the request

	// Non-Body Request Parameters
	for _, param := range e.Request.Parameters {
		switch param.In {
		case string(openapi3.ParameterInPath):
			rb.Param(restful.PathParameter(param.Name, ""))
		case string(openapi3.ParameterInQuery):
			rb.Param(restful.QueryParameter(param.Name, ""))
		case string(openapi3.ParameterInHeader):
			rb.Param(restful.HeaderParameter(param.Name, ""))
		}
	}

	// Request Body
	if e.Request.HasBody() {
		rb.Consumes(e.Request.Consumes()...)
	}

	// Success Response
	if e.Response.HasBody() || e.Response.Envelope {
		rb.Produces(e.Response.Produces()...)
	} else {
		rb.Produces(MIME_JSON)
	}
	defaultCode := e.Response.Codes.DefaultCode()
	var defaultPayload interface{} = struct{}{}

	if e.Response.Success.Payload != nil {
		defaultPayload = *e.Response.Success.Payload
	}

	rb.Do(DefaultReturns(defaultCode))
	rb.DefaultReturns(
		http.StatusText(defaultCode),
		defaultPayload)

	rb.Writes(defaultPayload)

	// Error Response
	if e.Response.Error.Payload != nil {
		errorPayload := *e.Response.Error.Payload
		errorPayloadType := reflect.TypeOf(errorPayload)
		if errorPayloadType.Kind() != reflect.Ptr {
			errorPayload = reflect.New(errorPayloadType).Interface()
		}
		rb.Do(ErrorPayload(errorPayload))
	}

	for _, doer := range e.Do {
		rb.Do(doer)
	}

	rb.Filter(EndpointResponseFilter)
	rb.Filter(EndpointRequestFilter)

	return rb
}

type EndpointRequestBodyFormField struct {
	Name     string
	Required bool
	Schema   *openapi3.SchemaOrRef
}

func (f EndpointRequestBodyFormField) Parameter() EndpointRequestParameter {
	return EndpointRequestParameter{
		Name:     f.Name,
		Required: types.NewBoolPtr(f.Required),
		Schema:   f.Schema,
		Style:    types.NewStringPtr(""),
		In:       "form",
		Explode:  types.NewBoolPtr(true),
	}
}

type EndpointRequestBody struct {
	Description string
	Required    bool
	Schema      *openapi3.SchemaOrRef
	Mime        string
	Example     *interface{}
	Encoding    map[string]openapi3.Encoding
	Reference   *string
}

func (b EndpointRequestBody) WithFormField(field EndpointRequestBodyFormField) EndpointRequestBody {
	b.Mime = MIME_MULTIPART_FORM
	b.Required = true

	if b.Schema == nil {
		b.Schema = &openapi3.SchemaOrRef{
			Schema: NewSchemaPtr(openapi3.SchemaTypeObject),
		}
	}

	b.Schema.Schema.WithPropertiesItem(field.Name, *field.Schema)

	if field.Required {
		req := types.StringStack(b.Schema.Schema.Required)
		if !req.Contains(field.Name) {
			req = append(req, field.Name)
		}
		b.Schema.Schema.Required = req
	}

	return b
}

func NewSchemaAllOf(first openapi3.SchemaOrRef, rest ...openapi3.SchemaOrRef) *openapi3.Schema {
	if first.Schema != nil {
		if first.Schema.Type == nil && first.Schema.AllOf != nil {
			all := append([]openapi3.SchemaOrRef{}, first.Schema.AllOf...)
			all = append(all, rest...)
			return first.Schema.WithAllOf(all...)
		}
	}

	all := append([]openapi3.SchemaOrRef{}, first)
	all = append(all, rest...)
	return new(openapi3.Schema).WithAllOf(all...)
}

type EndpointRequestParameter struct {
	Name            string                        `json:"name"` // Required.
	In              string                        `json:"in"`   // Required.
	Description     *string                       `json:"description,omitempty"`
	Required        *bool                         `json:"required,omitempty"`
	Deprecated      *bool                         `json:"deprecated,omitempty"`
	AllowEmptyValue *bool                         `json:"allowEmptyValue,omitempty"`
	Style           *string                       `json:"style,omitempty"`
	Explode         *bool                         `json:"explode,omitempty"`
	AllowReserved   *bool                         `json:"allowReserved,omitempty"`
	Schema          *openapi3.SchemaOrRef         `json:"schema,omitempty"`
	Content         map[string]openapi3.MediaType `json:"content,omitempty"`
	Example         *interface{}                  `json:"example,omitempty"`
	Reference       *string                       `json:"reference,omitempty"`
	PortField       *EndpointPortField            `json:"-"`
	UniqueId        string                        `json:"-"`
}

func (p EndpointRequestParameter) Merge(o EndpointRequestParameter) EndpointRequestParameter {
	if p.Description == nil && o.Description != nil {
		p.Description = o.Description
	}
	if p.Required == nil && o.Required != nil {
		p.Required = o.Required
	}
	if p.Deprecated == nil && o.Deprecated != nil {
		p.Description = o.Description
	}
	if p.AllowEmptyValue == nil && o.AllowEmptyValue != nil {
		p.AllowEmptyValue = o.AllowEmptyValue
	}
	if p.Style == nil && o.Style != nil {
		p.Style = o.Style
	}
	if p.Explode == nil && o.Explode != nil {
		p.Explode = o.Explode
	}
	if p.AllowReserved == nil && o.AllowReserved != nil {
		p.AllowReserved = o.AllowReserved
	}
	if p.Content == nil && o.Content != nil {
		p.Content = o.Content
	}
	if p.Schema == nil && o.Schema != nil {
		p.Schema = o.Schema
	} else if p.Schema != nil && o.Schema != nil {
		// Merge the two schemas
		allOfSchema := NewSchemaAllOf(*p.Schema, *o.Schema)
		allOf := NewSchemaOrRef(allOfSchema)
		p.Schema = &allOf
	}
	if p.Example == nil && o.Example != nil {
		p.Example = o.Example
	}
	if p.Reference == nil && o.Reference != nil {
		p.Reference = o.Reference
	}
	if p.PortField == nil && o.PortField != nil {
		p.PortField = o.PortField
	}
	return p
}

func (p EndpointRequestParameter) WithDescription(description string) EndpointRequestParameter {
	p.Description = &description
	return p
}

func (p EndpointRequestParameter) WithRequired(required bool) EndpointRequestParameter {
	p.Required = &required
	return p
}

func (p EndpointRequestParameter) WithDeprecated(deprecated bool) EndpointRequestParameter {
	p.Deprecated = &deprecated
	return p
}

func (p EndpointRequestParameter) WithStyle(style string) EndpointRequestParameter {
	p.Style = &style
	return p
}

func (p EndpointRequestParameter) WithAllowEmptyValue(allowEmptyValue bool) EndpointRequestParameter {
	p.AllowEmptyValue = &allowEmptyValue
	return p
}

func (p EndpointRequestParameter) WithExplode(explode bool) EndpointRequestParameter {
	p.Explode = &explode
	return p
}

func (p EndpointRequestParameter) WithAllowReserved(allowReserved bool) EndpointRequestParameter {
	p.AllowReserved = &allowReserved
	return p
}

func (p EndpointRequestParameter) WithSchema(schemaOrRef openapi3.SchemaOrRef) EndpointRequestParameter {
	p.Schema = &schemaOrRef
	return p
}

func (p EndpointRequestParameter) WithContentItem(key string, value openapi3.MediaType) EndpointRequestParameter {
	p.Content[key] = value
	return p
}

func (p EndpointRequestParameter) WithExample(example interface{}) EndpointRequestParameter {
	p.Example = &example
	return p
}

func (p EndpointRequestParameter) WithReference(reference string) EndpointRequestParameter {
	p.Reference = &reference
	return p
}

func NewEndpointRequestParameter(name, in string) EndpointRequestParameter {
	return EndpointRequestParameter{
		Name:     name,
		In:       in,
		UniqueId: types.MustNewUUID().String(),
	}
}

func PathParameter(name string, description string) EndpointRequestParameter {
	return NewEndpointRequestParameter(name, string(openapi3.ParameterInPath)).
		WithDescription(description).
		WithRequired(true)
}

func QueryParameter(name string, description string) EndpointRequestParameter {
	return NewEndpointRequestParameter(name, string(openapi3.ParameterInQuery)).
		WithDescription(description)
}

func HeaderParameter(name string, description string) EndpointRequestParameter {
	return NewEndpointRequestParameter(name, string(openapi3.ParameterInHeader)).
		WithDescription(description)
}

func CookieParameter(name string, description string) EndpointRequestParameter {
	return NewEndpointRequestParameter(name, string(openapi3.ParameterInCookie)).
		WithDescription(description)
}

type EndpointRequest struct {
	PortStruct  interface{}
	Port        EndpointPort
	Description string
	Parameters  []EndpointRequestParameter
	Validator   PortValidatorFunction
	Body        EndpointRequestBody
}

func (r EndpointRequest) bodyParameter() EndpointRequestParameter {
	for _, parameter := range r.Parameters {
		if parameter.In == "body" {
			return parameter
		}
	}
	return EndpointRequestParameter{}
}

func (r EndpointRequest) parameterByName(name string) EndpointRequestParameter {
	for _, p := range r.Parameters {
		if p.Name == name {
			return p
		}
	}
	return EndpointRequestParameter{}
}

func (r EndpointRequest) PatchParameter(name string, fn func(p EndpointRequestParameter) EndpointRequestParameter) EndpointRequest {
	for i, p := range r.Parameters {
		if p.Name == name {
			r.Parameters[i] = fn(p)
			return r
		}
	}
	return r
}

func (r EndpointRequest) WithParameter(p EndpointRequestParameter) EndpointRequest {
	for i, par := range r.Parameters {
		if par.Name == p.Name {
			r.Parameters[i] = r.Parameters[i].Merge(p)
			return r
		}
	}

	r.Parameters = append(r.Parameters, p)
	return r
}

func (r EndpointRequest) HasFormParameter() bool {
	for _, par := range r.Parameters {
		if par.In == "form" {
			return true
		}
	}

	return false
}

func (r EndpointRequest) WithMime(mime string) EndpointRequest {
	r.Body.Mime = mime
	return r
}

func (r EndpointRequest) WithBodySchema(schema openapi3.Schema) EndpointRequest {
	schemaOrRef := NewSchemaOrRef(&schema)
	r.Body.Schema = &schemaOrRef
	return r
}

func (r EndpointRequest) PatchBodySchema(fn func(openapi3.Schema) openapi3.Schema) EndpointRequest {
	// No schema to patch
	if r.Body.Schema != nil {
		return r
	}

	// Custom schema exists
	if r.Body.Schema.Schema != nil {
		bodySchema := *r.Body.Schema.Schema
		bodySchema = fn(bodySchema)
		r.Body.Schema.Schema = &bodySchema
		return r
	}

	// Clone the referenced schema into the body schema
	if r.Body.Schema.SchemaReference != nil {
		typeName := SchemaRefName(r.Body.Schema.SchemaReference)
		bodySchema, ok := LookupSchema(typeName)
		if ok {
			newBodySchema := fn(*bodySchema)
			r.Body.Schema.Schema = &newBodySchema
			r.Body.Schema.SchemaReference = nil
			return r
		}
	}

	return r
}

func (r EndpointRequest) WithBodyExample(example interface{}) EndpointRequest {
	r.Body.Example = &example
	return r
}

func (r EndpointRequest) WithValidator(fn PortValidatorFunction) EndpointRequest {
	r.Validator = fn
	return r
}

func (r EndpointRequest) WithBodyEncodingItem(key string, encoding openapi3.Encoding) EndpointRequest {
	if r.Body.Encoding == nil {
		r.Body.Encoding = make(map[string]openapi3.Encoding)
	}
	r.Body.Encoding[key] = encoding
	return r
}

func (r EndpointRequest) HasBody() bool {
	return r.Body.Mime != ""
}

func (r EndpointRequest) Consumes() []string {
	return []string{r.Body.Mime}
}

func (r EndpointRequest) WithBody(body EndpointRequestBody) EndpointRequest {
	r.Body = body
	return r
}

func (r EndpointRequest) WithPortStruct(portStruct interface{}) EndpointRequest {
	result := r

	port, _ := NewEndpointPort(PortTypeRequest, portStruct)
	result.Port = port

	for _, field := range port.Fields {
		if field.IsBody() {
			result = result.WithBody(field.RequestBody())
		} else if field.IsForm() {
			result.Body = result.Body.WithFormField(field.RequestBodyFormField())
		} else {
			result = result.WithParameter(field.Parameter())
		}
	}

	return result
}

func NewEndpointRequest() EndpointRequest {
	return EndpointRequest{}
}

type EndpointResponse struct {
	PortStruct interface{}
	Port       EndpointPort
	Codes      EndpointResponseCodes
	Envelope   bool
	Success    EndpointResponseContent
	Error      EndpointResponseContent
}

func (r EndpointResponse) WithEnvelope(envelope bool) EndpointResponse {
	r.Envelope = envelope
	if envelope {
		// Error is rendered as part of the envelope
		r.Error.Payload = nil
	}
	return r
}

func (r EndpointResponse) WithResponseCodes(codes EndpointResponseCodes) EndpointResponse {
	r.Codes = codes
	return r
}

func (r EndpointResponse) WithSuccessPayload(body interface{}) EndpointResponse {
	if r.Success.Mime == "" {
		r.Success.Mime = MIME_JSON
		r.Error.Mime = MIME_JSON
	}
	r.Success = r.Success.WithPayload(body)
	return r
}

func (r EndpointResponse) WithSuccessHeader(name string, header EndpointResponseHeader) EndpointResponse {
	r.Success = r.Success.WithHeader(name, header)
	return r
}

func (r EndpointResponse) WithErrorPayload(body interface{}) EndpointResponse {
	r.Error = r.Error.WithPayload(body)
	return r
}

func (r EndpointResponse) WithErrorHeader(name string, header EndpointResponseHeader) EndpointResponse {
	r.Error = r.Error.WithHeader(name, header)
	return r
}

func (r EndpointResponse) WithMime(mime string) EndpointResponse {
	r.Success = r.Success.WithMime(mime)
	r.Error = r.Error.WithMime(mime)
	return r
}

func (r EndpointResponse) WithSuccess(success EndpointResponseContent) EndpointResponse {
	r.Success = success
	return r
}

func (r EndpointResponse) WithError(error EndpointResponseContent) EndpointResponse {
	r.Error = error
	return r
}

func (r EndpointResponse) WithHeader(name string, header EndpointResponseHeader) EndpointResponse {
	return r.
		WithSuccessHeader(name, header).
		WithErrorHeader(name, header)
}

func (r EndpointResponse) HasBody() bool {
	return r.Success.Mime != "" || r.Error.Mime != ""
}

func (r EndpointResponse) Produces() []string {
	var mimeTypes = make(types.StringSet)
	mimeTypes.AddAll(r.Success.Mime, r.Error.Mime)
	if _, ok := mimeTypes[""]; ok {
		delete(mimeTypes, "")
	}
	return mimeTypes.Values()
}

func (r EndpointResponse) withEnumCodes(enum string) EndpointResponse {
	vals := strings.Split(enum, ",")
	codes := EndpointResponseCodes{}
	for _, v := range vals {
		iv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			logger.WithError(err).Errorf("Invalid response code %q", v)
			continue
		}
		if iv < 200 || iv > 599 {
			logger.Errorf("Invalid response code %d", iv)
			continue
		}
		if iv <= 399 {
			codes.Success = append(codes.Success, int(iv))
		} else {
			codes.Error = append(codes.Error, int(iv))
		}
	}
	return r.WithResponseCodes(codes)
}

func (r EndpointResponse) withPortField(field EndpointPortField) EndpointResponse {
	result := r
	if field.IsBody() {
		if field.Options["envelope"] == "true" {
			result = result.WithEnvelope(true)
		}

		if field.Options["error"] == "true" {
			result.Error = result.Error.withPortFieldBody(field)
		} else {
			result.Success = result.Success.withPortFieldBody(field)
		}

	} else if field.IsCode() {

		tags := field.Tags()
		if enumTag, ok := tags.Lookup("enum"); ok {
			result = result.withEnumCodes(enumTag)
		}

	} else if field.IsHeader() {
		fieldHeader := field.ResponseHeader()

		errorOnly := field.BoolOption("error")
		if !errorOnly {
			result = result.WithSuccessHeader(field.Name, fieldHeader)
		}

		successOnly := field.BoolOption("success")
		if !successOnly {
			result = result.WithErrorHeader(field.Name, fieldHeader)
		}
	}

	return result
}

func (r EndpointResponse) WithPortStruct(portStruct interface{}) EndpointResponse {
	result := r
	port, _ := NewEndpointPort(PortTypeResponse, portStruct)
	result.Port = port

	for _, field := range port.Fields {
		result = result.withPortField(field)
	}

	return result
}

func NewEndpointResponse() EndpointResponse {
	return EndpointResponse{}
}

func NewEnvelopeEndpointResponse(portStruct interface{}) EndpointResponse {
	r := NewEndpointResponse().WithPortStruct(portStruct).WithEnvelope(true)
	if portStruct == nil {
		r = r.WithResponseCodes(DeleteResponseCodes)
		r = r.WithMime(MIME_JSON)
	}
	return r
}

type EndpointResponseCodes struct {
	Success []int
	Error   []int
}

func (c EndpointResponseCodes) DefaultCode() int {
	if len(c.Success) == 0 {
		return http.StatusOK
	}
	return c.Success[0]
}

var (
	StandardErrorCodes = []int{
		http.StatusBadRequest,   // Some parameter issue
		http.StatusUnauthorized, // No tenant access
		http.StatusForbidden,    // No authentication token
		http.StatusNotFound,     // Target resource does not exist
		http.StatusConflict,     // State issue (eg. already exists)
	}
	ListResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusOK},
		Error:   []int{http.StatusBadRequest, http.StatusForbidden, http.StatusUnauthorized},
	}
	GetResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusOK},
		Error:   []int{http.StatusBadRequest, http.StatusForbidden, http.StatusUnauthorized, http.StatusNotFound},
	}
	CreateResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusCreated},
		Error:   []int{http.StatusBadRequest, http.StatusForbidden, http.StatusUnauthorized, http.StatusConflict},
	}
	UpdateResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusOK},
		Error:   []int{http.StatusBadRequest, http.StatusForbidden, http.StatusUnauthorized, http.StatusNotFound},
	}
	DeleteResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusOK},
		Error:   StandardErrorCodes,
	}
	AcceptResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusAccepted},
		Error:   StandardErrorCodes,
	}
	NoContentResponseCodes = EndpointResponseCodes{
		Success: []int{http.StatusNoContent},
		Error:   StandardErrorCodes,
	}
)

func DefaultResponseCodes(method string) EndpointResponseCodes {
	switch method {
	case http.MethodGet:
		return GetResponseCodes
	case http.MethodPost:
		return CreateResponseCodes
	case http.MethodPut:
		return UpdateResponseCodes
	case http.MethodDelete:
		return DeleteResponseCodes
	default:
		return EndpointResponseCodes{
			Success: []int{http.StatusOK},
			Error:   []int{http.StatusBadRequest, http.StatusForbidden, http.StatusUnauthorized},
		}
	}
}

type EndpointResponseContent struct {
	Mime     string
	Headers  map[string]EndpointResponseHeader
	Payload  *interface{}
	Example  *interface{}
	Encoding map[string]openapi3.Encoding
}

func (c EndpointResponseContent) WithHeader(name string, header EndpointResponseHeader) EndpointResponseContent {
	if c.Headers == nil {
		c.Headers = make(map[string]EndpointResponseHeader)
	} else if h, ok := c.Headers[name]; ok {
		header = h.Merge(header)
	}
	c.Headers[name] = header
	return c
}

func (c EndpointResponseContent) WithEncoding(key string, encoding openapi3.Encoding) EndpointResponseContent {
	if c.Encoding == nil {
		c.Encoding = make(map[string]openapi3.Encoding)
	}
	c.Encoding[key] = encoding
	return c
}

func (c EndpointResponseContent) WithPayload(payload interface{}) EndpointResponseContent {
	c.Payload = &payload
	if c.Mime == "" || c.Mime == "*/*" {
		return c.WithMime(MIME_JSON)
	}

	return c
}

func (c EndpointResponseContent) WithExample(example interface{}) EndpointResponseContent {
	c.Example = &example
	return c
}

func (c EndpointResponseContent) WithMime(mime string) EndpointResponseContent {
	c.Mime = mime
	return c
}

func (c EndpointResponseContent) withPortFieldBody(field EndpointPortField) EndpointResponseContent {
	payload := types.Instantiate(field.Field.Type)
	content := c.WithPayload(payload)

	mime := MIME_JSON
	if mimeOverride := field.Options["mime"]; mimeOverride != "" {
		mime = mimeOverride
	}
	content = content.WithMime(mime)

	return content
}

type EndpointResponseHeader struct {
	Description     *string                       `json:"description,omitempty"`
	Required        *bool                         `json:"required,omitempty"`
	Deprecated      *bool                         `json:"deprecated,omitempty"`
	AllowEmptyValue *bool                         `json:"allowEmptyValue,omitempty"`
	Explode         *bool                         `json:"explode,omitempty"`
	AllowReserved   *bool                         `json:"allowReserved,omitempty"`
	Schema          *openapi3.SchemaOrRef         `json:"schema,omitempty"`
	Content         map[string]openapi3.MediaType `json:"content,omitempty"`
	Example         *interface{}                  `json:"example,omitempty"`
	Reference       *string                       `json:"reference,omitempty"`
	PortField       *EndpointPortField            `json:"-"`
}

func (p EndpointResponseHeader) Merge(o EndpointResponseHeader) EndpointResponseHeader {
	if p.Description == nil && o.Description != nil {
		p.Description = o.Description
	}
	if p.Required == nil && o.Required != nil {
		p.Required = o.Required
	}
	if p.Deprecated == nil && o.Deprecated != nil {
		p.Description = o.Description
	}
	if p.AllowEmptyValue == nil && o.AllowEmptyValue != nil {
		p.AllowEmptyValue = o.AllowEmptyValue
	}
	if p.Explode == nil && o.Explode != nil {
		p.Explode = o.Explode
	}
	if p.AllowReserved == nil && o.AllowReserved != nil {
		p.AllowReserved = o.AllowReserved
	}
	if p.Content == nil && o.Content != nil {
		p.Content = o.Content
	}
	if p.Schema == nil && o.Schema != nil {
		p.Schema = o.Schema
	} else if p.Schema != nil && o.Schema != nil {
		// Merge the two schemas
		allOfSchema := NewSchemaAllOf(*p.Schema, *o.Schema)
		allOf := NewSchemaOrRef(allOfSchema)
		p.Schema = &allOf
	}
	if p.Example == nil && o.Example != nil {
		p.Example = o.Example
	}
	if p.Reference == nil && o.Reference != nil {
		p.Reference = o.Reference
	}
	if p.PortField == nil && o.PortField != nil {
		p.PortField = o.PortField
	}
	return p
}

func (p EndpointResponseHeader) WithDescription(description string) EndpointResponseHeader {
	p.Description = &description
	return p
}

func (p EndpointResponseHeader) WithRequired(required bool) EndpointResponseHeader {
	p.Required = &required
	return p
}

func (p EndpointResponseHeader) WithDeprecated(deprecated bool) EndpointResponseHeader {
	p.Deprecated = &deprecated
	return p
}

func (p EndpointResponseHeader) WithAllowEmptyValue(allowEmptyValue bool) EndpointResponseHeader {
	p.AllowEmptyValue = &allowEmptyValue
	return p
}

func (p EndpointResponseHeader) WithExplode(explode bool) EndpointResponseHeader {
	p.Explode = &explode
	return p
}

func (p EndpointResponseHeader) WithAllowReserved(allowReserved bool) EndpointResponseHeader {
	p.AllowReserved = &allowReserved
	return p
}

func (p EndpointResponseHeader) WithSchema(schemaOrRef openapi3.SchemaOrRef) EndpointResponseHeader {
	p.Schema = &schemaOrRef
	return p
}

func (p EndpointResponseHeader) WithContentItem(key string, value openapi3.MediaType) EndpointResponseHeader {
	p.Content[key] = value
	return p
}

func (p EndpointResponseHeader) WithExample(example interface{}) EndpointResponseHeader {
	p.Example = &example
	return p
}

func (p EndpointResponseHeader) WithReference(reference string) EndpointResponseHeader {
	p.Reference = &reference
	return p
}

func NewEndpointResponseHeader() EndpointResponseHeader {
	return EndpointResponseHeader{}
}
