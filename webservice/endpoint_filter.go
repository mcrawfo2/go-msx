package webservice

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

func InjectEndpointFilter(e Endpoint) restful.FilterFunction {
	return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
		request = RequestWithEndpoint(request, e)
		chain.ProcessFilter(request, response)
	}
}

type RequestValidator interface {
	ValidateRequest(endpoint Endpoint) error
}

type RequestPopulator interface {
	PopulateInputs(endpoint Endpoint) (interface{}, error)
}

func EndpointRequestFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	var err error

	e := EndpointFromRequest(request)

	// Configure the request data source
	dataSource := &RestfulRequestDataSource{Request: request}
	decoder := &OpenApiRequestDecoder{DataSource: dataSource}

	// Validate request according to the operation schemas
	validator := OpenApiRequestValidator{Decoder: decoder}
	if err = validator.ValidateRequest(e); err != nil {
		request = RequestWithError(request, NewBadRequestError(err))
		return
	}

	// Populate inputs
	var inputs interface{}
	populator := OpenApiRequestPopulator{Decoder: decoder}
	if inputs, err = populator.PopulateInputs(e); err != nil {
		request = RequestWithError(request, NewBadRequestError(err))
		return
	}

	// Custom validation for args
	if e.Request.Validator != nil && inputs != nil {
		if err = e.Request.Validator(inputs); err != nil {
			request = RequestWithError(request, NewBadRequestError(err))
			return
		}
	}

	// TODO: Defaults

	// Inject inputs
	request = RequestWithInputs(request, inputs)

	// Continue request processing
	chain.ProcessFilter(request, response)
}

func EndpointResponseFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	var err error

	// Perform request processing
	chain.ProcessFilter(request, response)

	// Retrieve the endpoint definition
	e := EndpointFromRequest(request)

	// Configure the response data sink
	dataSink := &RestfulResponseDataSink{Response: response}
	encoder := OpenApiResponseEncoder{Sink: dataSink}

	// Convert the error to HTTP if required
	responseError := ErrorFromRequest(request)
	if responseError != nil && e.ErrorConverter != nil {
		responseError = e.ErrorConverter.Convert(responseError)
	}

	// Retrieve the output port struct
	outputs := OutputsFromRequest(request)
	if outputs == nil && responseError == nil {
		// Already handled by the controller
		//return
	}

	// Create a new populator
	observer := CompositeResponseObserver{
		LogResponseObserver{Context: request.Request.Context()},
		TracingResponseObserver{Context: request.Request.Context()},
	}

	// Create a new request describer (for the envelope fields)
	describer := RestfulRequestDescriber{
		Request: request,
	}

	populator := OpenApiResponsePopulator{
		Endpoint: e,
		Outputs:  outputs,
		Error:    responseError,

		Observer:  observer,
		Encoder:   encoder,
		Describer: describer,
	}

	// Populate response from outputs
	if err = populator.PopulateOutputs(); err != nil {
		WriteError(request, response, http.StatusInternalServerError, err)
		return
	}
}

func RequestWithEndpoint(request *restful.Request, e Endpoint) *restful.Request {
	request.SetAttribute(AttributeKeyEndpoint, e)
	return request
}

func EndpointFromRequest(request *restful.Request) Endpoint {
	return request.Attribute(AttributeKeyEndpoint).(Endpoint)
}

func InjectEndpoint(e Endpoint) RouteBuilderFunc {
	return func(builder *restful.RouteBuilder) {
		builder.Metadata(AttributeKeyEndpoint, e)
	}
}

func EndpointFromRoute(route restful.Route) Endpoint {
	if val, ok := route.Metadata[AttributeKeyEndpoint]; ok {
		return val.(Endpoint)
	}
	return Endpoint{}
}

func RequestWithInputs(request *restful.Request, inputs interface{}) *restful.Request {
	request.SetAttribute(AttributeKeyInputs, inputs)
	return request
}

func InputsFromRequest(request *restful.Request) interface{} {
	return request.Attribute(AttributeKeyInputs)
}

func RequestWithOutputs(request *restful.Request, outputs interface{}) *restful.Request {
	request.SetAttribute(AttributeKeyOutputs, outputs)
	return request
}

func OutputsFromRequest(request *restful.Request) interface{} {
	return request.Attribute(AttributeKeyOutputs)
}

func RequestWithError(request *restful.Request, err error) *restful.Request {
	request.SetAttribute(AttributeError, err)
	return request
}

func ErrorFromRequest(request *restful.Request) error {
	val, _ := request.Attribute(AttributeError).(error)
	return val
}

func RequestWithErrorPayload(request *restful.Request, payload interface{}) {
	request.SetAttribute(AttributeErrorPayload, payload)
}

func ErrorPayloadFromRequest(request *restful.Request) interface{} {
	return request.Attribute(AttributeErrorPayload)
}
