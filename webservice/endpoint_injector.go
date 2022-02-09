package webservice

import (
	"context"
	"github.com/emicklei/go-restful"
	"github.com/swaggest/refl"
	"net/http"
	"reflect"
)

type inputGenerator func(request *restful.Request, response *restful.Response) reflect.Value

type inputInjector struct {
	PortStructType reflect.Type
	Generators     []inputGenerator
}

func (i inputInjector) PortStruct() interface{} {
	if i.PortStructType == nil {
		return nil
	}

	return reflect.New(i.PortStructType).Elem().Interface()
}

func (i inputInjector) Args(request *restful.Request, response *restful.Response) []reflect.Value {
	var args []reflect.Value
	for _, generator := range i.Generators {
		value := generator(request, response)
		args = append(args, value)
	}
	return args
}

func (i inputInjector) generateContext(request *restful.Request, _ *restful.Response) reflect.Value {
	return reflect.ValueOf(request.Request.Context())
}

func (i inputInjector) generateRequestPointer(request *restful.Request, _ *restful.Response) reflect.Value {
	return reflect.ValueOf(request)
}

func (i inputInjector) generateResponsePointer(_ *restful.Request, response *restful.Response) reflect.Value {
	return reflect.ValueOf(response)
}

func (i inputInjector) generateHttpRequestPointer(request *restful.Request, _ *restful.Response) reflect.Value {
	return reflect.ValueOf(request.Request)
}

func (i inputInjector) generateHttpResponseWriter(_ *restful.Request, response *restful.Response) reflect.Value {
	return reflect.ValueOf(response.ResponseWriter)
}

func (i inputInjector) generateInputs(request *restful.Request, _ *restful.Response) reflect.Value {
	return reflect.ValueOf(InputsFromRequest(request)).Elem()
}

var contextInstance context.Context
var contextType = reflect.TypeOf(&contextInstance).Elem()

var restfulRequestPointerInstance *restful.Request
var restfulRequestPointerType = reflect.TypeOf(&restfulRequestPointerInstance).Elem()

var restfulResponsePointerInstance *restful.Response
var restfulResponsePointerType = reflect.TypeOf(&restfulResponsePointerInstance).Elem()

var httpRequestPointerInstance *http.Request
var httpRequestPointerType = reflect.TypeOf(&httpRequestPointerInstance).Elem()

var httpResponseWriterInstance http.ResponseWriter
var httpResponseWriterType = reflect.TypeOf(&httpResponseWriterInstance).Elem()

var errorInstance error
var errorType = reflect.TypeOf(&errorInstance).Elem()

func newInputInjector(fn reflect.Type) inputInjector {
	var result inputInjector

	for i := 0; i < fn.NumIn(); i++ {
		ti := fn.In(i)
		switch ti {
		case contextType:
			result.Generators = append(result.Generators, result.generateContext)

		case restfulRequestPointerType:
			result.Generators = append(result.Generators, result.generateRequestPointer)

		case restfulResponsePointerType:
			result.Generators = append(result.Generators, result.generateResponsePointer)

		case httpRequestPointerType:
			result.Generators = append(result.Generators, result.generateHttpRequestPointer)

		case httpResponseWriterType:
			result.Generators = append(result.Generators, result.generateHttpResponseWriter)

		default:
			result.Generators = append(result.Generators, result.generateInputs)
			result.PortStructType = refl.DeepIndirect(ti)
		}
	}

	return result
}

type outputApplicator func(v reflect.Value, request *restful.Request, response *restful.Response)

type outputExtractor struct {
	PortStructType reflect.Type
	Applicators    []outputApplicator
}

func (e outputExtractor) PortStruct() interface{} {
	if e.PortStructType == nil {
		return nil
	}

	return reflect.New(e.PortStructType).Elem().Interface()
}

func (e outputExtractor) injectError(v reflect.Value, request *restful.Request, _ *restful.Response) {
	err, _ := v.Interface().(error)
	RequestWithError(request, err)
}

func (e outputExtractor) injectOutputs(v reflect.Value, request *restful.Request, _ *restful.Response) {
	RequestWithOutputs(request, v.Interface())
}

func (e outputExtractor) ApplyResults(values []reflect.Value, request *restful.Request, response *restful.Response) {
	for i, applicator := range e.Applicators {
		value := values[i]
		applicator(value, request, response)
	}
}

func newOutputExtractor(fn reflect.Type) outputExtractor {
	var result outputExtractor

	for i := 0; i < fn.NumOut(); i++ {
		ti := fn.Out(i)
		switch ti {
		case errorType:
			// Apply error into the request attributes
			result.Applicators = append(result.Applicators, result.injectError)

		default:
			// Apply output port struct into the request attributes
			result.Applicators = append(result.Applicators, result.injectOutputs)
			result.PortStructType = refl.DeepIndirect(ti)
		}
	}

	return result
}

func InjectingController(fn interface{}) (restful.RouteFunction, interface{}, interface{}) {
	fnValue := reflect.ValueOf(fn)

	injector := newInputInjector(fnValue.Type())
	extractor := newOutputExtractor(fnValue.Type())

	routeFunction := func(request *restful.Request, response *restful.Response) {
		argValues := injector.Args(request, response)
		resultValues := fnValue.Call(argValues)
		extractor.ApplyResults(resultValues, request, response)
	}

	return routeFunction, injector.PortStruct(), extractor.PortStruct()
}
