// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/restfulcontext"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"strings"
)

//go:generate mockery --name=RestfulWebServer --case=snake --inpackage --testonly --with-expecter
type RestfulWebServer interface {
	NewService(root string) (*restful.WebService, error)
}

type EndpointsRegisterer interface {
	RegisterEndpoints(producer EndpointsProducer) (err error)
}

type RestfulEndpointRegisterer struct {
	w RestfulWebServer
}

// RegisterEndpoints registers the contents of an EndpointProducer
func (s *RestfulEndpointRegisterer) RegisterEndpoints(producer EndpointsProducer) (err error) {
	svc, err := s.w.NewService(PathApiRoot)
	if err != nil {
		return err
	}

	var t EndpointTransformers
	if tp, ok := producer.(EndpointTransformersProducer); ok {
		t = tp.EndpointTransformers()
	}

	endpoints, err := producer.Endpoints()
	if err != nil {
		return err
	}

	t.Transform(endpoints)

	return endpoints.Each(func(endpoint *Endpoint) error {
		svc.Route(RouteBuilderFromEndpoint(svc, endpoint))
		RegisterEndpoint(endpoint)
		return nil
	})
}

func NewRestfulEndpointRegisterer(w RestfulWebServer) EndpointsRegisterer {
	return &RestfulEndpointRegisterer{
		w: w,
	}
}

// ContextEndpointRegisterer returns an EndpointsRegisterer for the context's WebServer
func ContextEndpointRegisterer(ctx context.Context) EndpointsRegisterer {
	ws := webservice.WebServerFromContext(ctx)
	return NewRestfulEndpointRegisterer(ws)
}

// InjectRequestEndpointFilter supplies the endpoint
func InjectRequestEndpointFilter(e *Endpoint) restful.FilterFunction {
	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
		request = RequestWithEndpoint(request, e)
		chain.ProcessFilter(request, response)
	}
}

// InjectEndpointContextFilter executes the injectors against the request context
func InjectEndpointContextFilter(injectors types.ContextInjectors) restful.FilterFunction {
	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
		ctx := request.Request.Context()
		ctx = injectors.Inject(ctx)
		request.Request = request.Request.WithContext(ctx)
		chain.ProcessFilter(request, response)
	}
}

// InjectEndpointRequestDecoder supplies the request data source
func InjectEndpointRequestDecoder(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	dataSource := NewRestfulRequestDataSource(request)
	decoder := NewRequestDecoder(dataSource)
	request = RequestWithEndpointRequestDecoder(request, decoder)
	chain.ProcessFilter(request, response)
}

// InjectEndpointResponseEncoder supplies the response data sink
func InjectEndpointResponseEncoder(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	dataSink := NewRestfulResponseDataSink(response)
	encoder := NewResponseEncoder(dataSink)
	request = RequestWithEndpointResponseEncoder(request, encoder)
	chain.ProcessFilter(request, response)
}

// EndpointMiddlewaresFilter executes a chain of http middleware
func EndpointMiddlewaresFilter(m Middlewares) restful.FilterFunction {
	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
		handler := m.Compose(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			response.ResponseWriter = w
			request.Request = req
			chain.ProcessFilter(request, response)
		}))

		handler.ServeHTTP(response, request.Request)
	}
}

type RequestPopulator interface {
	PopulateInputs(endpoint Endpoint) (interface{}, error)
}

// EndpointRequestFilter validates and injects the inputs into the request
func EndpointRequestFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	var err error

	defer func() {
		if err != nil {
			webservice.RequestWithError(request, err)
		}
	}()

	// Retrieve the endpoint
	e := EndpointFromRequest(request)

	// Retrieve the decoder
	decoder := EndpointRequestDecoderFromRequest(request)

	// Validate request according to the endpoint port schemas
	validator := NewRequestValidator(e.Request.Port, decoder)
	if err = validator.ValidateRequest(); err != nil {
		err = webservice.NewBadRequestError(err)
		return
	}

	// Populate inputs
	var inputs interface{}
	populator := ops.NewInputsPopulator(e.Request.Port, decoder)
	if inputs, err = populator.PopulateInputs(); err != nil {
		err = webservice.NewBadRequestError(err)
		return
	}

	// Custom validation for args
	if e.Request.Validator != nil && inputs != nil {
		if err = e.Request.Validator(inputs); err != nil {
			err = webservice.NewBadRequestError(err)
			return
		}
	}

	// Inject inputs
	request = RequestWithInputs(request, inputs)

	// Continue request processing
	chain.ProcessFilter(request, response)
}

// EndpointResponseFilter validates and extracts the outputs into the response
func EndpointResponseFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	var err error

	// Perform request processing
	chain.ProcessFilter(request, response)

	// Retrieve the endpoint definition
	e := EndpointFromRequest(request)

	// Convert the error to HTTP if required
	responseError := webservice.ErrorFromRequest(request)
	if responseError != nil {
		if _, ok := responseError.(StatusCodeError); !ok {
			var errorConverter ErrorConverter = ErrorStatusCoderConverter{
				ErrorStatusCoder: ErrorStatusCoderFunc(DefaultErrorStatusCoder),
			}
			if e.ErrorConverter != nil {
				errorConverter = e.ErrorConverter
			}
			responseError = errorConverter.Convert(responseError)
		}
	}

	// Retrieve the output port struct
	outputs := OutputsFromRequest(request)
	if outputs == nil && responseError == nil {
		// Already handled by the controller
		//return
	}

	// Retrieve the response encoder
	encoder := EndpointResponseEncoderFromRequest(request)

	// Create a new observer
	observer := CompositeResponseObserver{
		LoggingResponseObserver{Context: request.Request.Context()},
		TracingResponseObserver{Context: request.Request.Context()},
	}

	// Create a new request describer (for the envelope fields)
	describer := RestfulRequestDescriber{
		Request: request,
	}

	// Create a new populator
	populator := OutputsPopulator{
		Endpoint: e,
		Outputs:  &outputs,
		Error:    responseError,

		Observer:  observer,
		Encoder:   encoder,
		Describer: describer,
	}

	// Populate response from outputs
	if err = populator.PopulateOutputs(); err != nil {
		webservice.WriteError(request, response, http.StatusInternalServerError, err)
		return
	}
}

// EndpointController calls the endpoint handler using the handler context
func EndpointController(endpoint *Endpoint) (fn restful.RouteFunction) {
	return func(request *restful.Request, response *restful.Response) {
		handlerContext := &EndpointHandlerContext{
			request:    request,
			response:   response,
			inputType:  endpoint.Inputs.Value(),
			outputType: endpoint.Outputs.Value(),
		}

		ctx := request.Request.Context()
		ctx = types.ContextWithHandlerContext(ctx, handlerContext)
		request.Request.WithContext(ctx)

		err := endpoint.Handler.Call(ctx)
		if err != nil {
			webservice.RequestWithError(request, err)
		}
	}
}

// RouteBuilderPermissionsFromEndpoint applies the endpoint permissions to the RouteBuilder
func RouteBuilderPermissionsFromEndpoint(e *Endpoint) restfulcontext.RouteBuilderFunc {
	return func(rb *restful.RouteBuilder) {
		if len(e.Permissions) > 0 {
			rb.Do(webservice.Permissions(e.Permissions...))
		}
	}
}

// RouteBuilderRequestParamsFromEndpoint applies the endpoint request parameters to the RouteBuilder
func RouteBuilderRequestParamsFromEndpoint(e *Endpoint) restfulcontext.RouteBuilderFunc {
	return func(rb *restful.RouteBuilder) {
		// Non-Body Request Parameters
		for _, param := range e.Request.Parameters {
			var rbp *restful.Parameter
			switch param.In {
			case ParameterInPath:
				rbp = restful.PathParameter(param.Name,
					types.NewOptionalString(param.Description).OrEmpty())
			case ParameterInQuery:
				rbp = restful.QueryParameter(param.Name,
					types.NewOptionalString(param.Description).OrEmpty())
			case ParameterInHeader:
				rbp = restful.HeaderParameter(param.Name,
					types.NewOptionalString(param.Description).OrEmpty())
			default:
				continue
			}

			if param.Required != nil {
				rbp.Required(*param.Required)
			}

			rb.Param(rbp)
		}
	}
}

// RouteBuilderRequestBodyFromEndpoint applies the endpoint request body to the RouteBuilder
func RouteBuilderRequestBodyFromEndpoint(e *Endpoint) restfulcontext.RouteBuilderFunc {
	return func(rb *restful.RouteBuilder) {
		// Request Body
		if e.Request.HasBody() {
			rb.Consumes(e.Request.Consumes()...)
		} else {
			rb.Consumes("")
			//switch e.Method {
			//case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			//default:
			rb.AllowedMethodsWithoutContentType([]string{e.Method})
			//}
		}
	}
}

// RouteBuilderResponsesFromEndpoint applies the endpoint responses to the RouteBuilder
func RouteBuilderResponsesFromEndpoint(e *Endpoint) restfulcontext.RouteBuilderFunc {
	return func(rb *restful.RouteBuilder) {
		// Success Response
		if e.Response.HasBody() {
			rb.Produces(e.Response.Produces()...)
		} else {
			rb.Produces(MediaTypeJson)
		}
		defaultCode := e.Response.Codes.DefaultCode()
		var payload interface{}

		if e.Response.Success.Payload.IsPresent() {
			payload = e.Response.Success.Payload.Value()
		}

		if e.Response.Success.Paging.IsPresent() {
			paging := e.Response.Success.Paging.Value()
			example := webservice.NewEnvelopedResponse(
				reflect.TypeOf(paging),
				payload)
			webservice.RouteBuilderWithSuccessResponse(rb, example)
			webservice.RouteBuilderWithEnvelopedPayload(rb, payload)
			payload = example
		} else {
			webservice.RouteBuilderWithSuccessResponse(rb, payload)
		}

		if e.Response.Envelope {
			example := webservice.NewEnvelopedResponse(
				reflect.TypeOf(integration.MsxEnvelope{}),
				payload)
			webservice.RouteBuilderWithSuccessResponse(rb, example)
			webservice.RouteBuilderWithEnvelopedPayload(rb, payload)
			rb.DefaultReturns("Success", example)
			rb.Writes(example)
			rb.Do(webservice.ErrorPayload(new(integration.MsxEnvelope)))
		} else {
			rb.Do(webservice.DefaultReturns(defaultCode))
			rb.DefaultReturns(
				http.StatusText(defaultCode),
				payload)
			rb.Writes(payload)
		}

		// Error Response
		if e.Response.Error.Payload.IsPresent() {
			errorPayload := e.Response.Error.Payload.Value()
			errorPayloadType := reflect.TypeOf(errorPayload)
			if errorPayloadType.Kind() != reflect.Ptr {
				errorPayload = reflect.New(errorPayloadType).Interface()
			}
			rb.Do(webservice.ErrorPayload(errorPayload))
		}
	}
}

// RouteBuilderFromEndpoint creates a RouteBuilder for the endpoint on the supplied WebService
func RouteBuilderFromEndpoint(svc *restful.WebService, e *Endpoint) *restful.RouteBuilder {
	return svc.
		Method(e.Method).
		Path(strings.TrimPrefix(e.Path, svc.RootPath())).
		Operation(e.OperationID).
		To(EndpointController(e)).
		Do(RouteWithEndpoint(e)).
		Filter(InjectRequestEndpointFilter(e)).
		Filter(InjectEndpointContextFilter(e.Injectors)).
		Do(RouteBuilderPermissionsFromEndpoint(e)).
		Do(RouteBuilderRequestParamsFromEndpoint(e)).
		Do(RouteBuilderRequestBodyFromEndpoint(e)).
		Do(RouteBuilderResponsesFromEndpoint(e)).
		Filter(EndpointMiddlewaresFilter(e.Middleware)).
		Filter(InjectEndpointRequestDecoder).
		Filter(InjectEndpointResponseEncoder).
		Filter(EndpointResponseFilter).
		Filter(EndpointRequestFilter)
}

// EndpointHandlerTypesAnalyzer determines the types for arguments and return values for the EndpointHandler
type EndpointHandlerTypesAnalyzer struct {
	endpoint    *Endpoint
	handlerFunc interface{}
}

func (i *EndpointHandlerTypesAnalyzer) getInputsType(argsTypes types.TypeSet) (inputsType reflect.Type) {
	// Find the input port type from the handler arguments
	knownArgTypes := types.NewTypeSet().With(types.DefaultHandlerArgumentValueTypeSet).With(argsTypes)
	handlerFuncType := reflect.ValueOf(i.handlerFunc).Type()
	for n := 0; n < handlerFuncType.NumIn(); n++ {
		argType := handlerFuncType.In(n)
		if _, ok := knownArgTypes[argType]; !ok {
			return argType
		}
	}

	return nil
}

func (i *EndpointHandlerTypesAnalyzer) getOutputsType(resultsTypes types.TypeSet) (outputsType reflect.Type) {
	// Find the input port type from the handler arguments
	knownResultTypes := types.NewTypeSet().With(types.DefaultHandlerResultValueTypeSet).With(resultsTypes)
	handlerFuncType := reflect.ValueOf(i.handlerFunc).Type()
	for n := 0; n < handlerFuncType.NumOut(); n++ {
		resultType := handlerFuncType.Out(n)
		if _, ok := knownResultTypes[resultType]; !ok {
			return resultType
		}
	}

	return nil
}

func (i *EndpointHandlerTypesAnalyzer) ArgsTypeSet() types.TypeSet {
	ts := types.NewTypeSet(
		restfulRequestPointerType,
		httpRequestPointerType,
		restfulResponsePointerType,
		httpResponseWriterType,
		endpointRequestDecoderType,
		endpointResponseEncoderType,
	)

	if i.endpoint.Inputs.IsPresent() {
		ts = ts.WithType(i.endpoint.Inputs.Value())
	}

	return ts
}

func (i *EndpointHandlerTypesAnalyzer) ReturnsTypeSet() types.TypeSet {
	ts := types.NewTypeSet()

	if i.endpoint.Outputs.IsPresent() {
		ts = ts.WithType(i.endpoint.Outputs.Value())
	}

	return ts
}

var contextContextInstance context.Context
var contextContextType = reflect.TypeOf(&contextContextInstance).Elem()

var restfulRequestPointerInstance *restful.Request
var restfulRequestPointerType = reflect.TypeOf(&restfulRequestPointerInstance).Elem()

var restfulResponsePointerInstance *restful.Response
var restfulResponsePointerType = reflect.TypeOf(&restfulResponsePointerInstance).Elem()

var httpRequestPointerInstance *http.Request
var httpRequestPointerType = reflect.TypeOf(&httpRequestPointerInstance).Elem()

var httpResponseWriterInstance http.ResponseWriter
var httpResponseWriterType = reflect.TypeOf(&httpResponseWriterInstance).Elem()

var endpointInputDecoderInstance ops.InputDecoder
var endpointRequestDecoderType = reflect.TypeOf(&endpointInputDecoderInstance).Elem()

var endpointResponseEncoderInstance ResponseEncoder
var endpointResponseEncoderType = reflect.TypeOf(&endpointResponseEncoderInstance).Elem()

var errorInstance error
var errorType = reflect.TypeOf(&errorInstance).Elem()

// EndpointHandlerContext holds the data to be injected into the subscriber's handler function
type EndpointHandlerContext struct {
	request    *restful.Request
	response   *restful.Response
	inputType  reflect.Type
	outputType reflect.Type
}

func (m *EndpointHandlerContext) RestfulRequestPointer() *restful.Request {
	return m.request
}

func (m *EndpointHandlerContext) RestfulResponsePointer() *restful.Response {
	return m.response
}

func (m *EndpointHandlerContext) HttpRequestPointer() *http.Request {
	return m.request.Request
}

func (m *EndpointHandlerContext) HttpResponseWriter() http.ResponseWriter {
	return m.response.ResponseWriter
}

func (m *EndpointHandlerContext) EndpointRequestDecoder() ops.InputDecoder {
	return EndpointRequestDecoderFromRequest(m.request)
}

func (m *EndpointHandlerContext) EndpointResponseEncoder() ResponseEncoder {
	return EndpointResponseEncoderFromRequest(m.request)
}

func (m *EndpointHandlerContext) GenerateArgument(ctx context.Context, t types.HandlerValueType) (result reflect.Value, err error) {
	switch t.ValueType {
	case contextContextType:
		result = reflect.ValueOf(ctx)
	case restfulRequestPointerType:
		result = reflect.ValueOf(m.RestfulRequestPointer())
	case httpRequestPointerType:
		result = reflect.ValueOf(m.HttpRequestPointer())
	case restfulResponsePointerType:
		result = reflect.ValueOf(m.RestfulResponsePointer())
	case httpResponseWriterType:
		result = reflect.ValueOf(m.HttpResponseWriter())
	case endpointRequestDecoderType:
		result = reflect.ValueOf(m.EndpointRequestDecoder())
	case endpointResponseEncoderType:
		result = reflect.ValueOf(m.EndpointResponseEncoder())
	case m.inputType:
		inputs := InputsFromRequest(m.request)
		result = reflect.ValueOf(inputs)
	default:
		err = errors.Wrapf(types.ErrUnknownValueType, "%v", t)
	}
	return

}

func (m *EndpointHandlerContext) HandleResult(t types.HandlerValueType, v reflect.Value) (err error) {
	switch t.ValueType {
	case errorType:
		erri := v.Interface()
		if erri != nil {
			err = erri.(error)
		}
	case m.outputType:
		RequestWithOutputs(m.request, v.Interface())
	default:
		err = errors.Wrapf(types.ErrUnknownValueType, "%v", t)
	}
	return
}
