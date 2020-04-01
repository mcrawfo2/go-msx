package webservice

import (
	"context"
	"github.com/emicklei/go-restful"
	"net/http"
	"strings"
)

type ControllerFunction func(req *restful.Request) (body interface{}, err error)
type ContextFunction func(ctx context.Context) (body interface{}, err error)

// Force only error into an envelope
func RawController(fn ControllerFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		body, err := fn(req)
		RawResponse(req, resp, body, err)
	}
}

func RawContextController(fn ContextFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()

		body, err := fn(ctx)
		RawResponse(req, resp, body, err)
	}
}

func Controller(fn ControllerFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		body, err := fn(req)
		EnvelopeResponse(req, resp, body, err)
	}
}

func ContextController(fn ContextFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()

		body, err := fn(ctx)
		EnvelopeResponse(req, resp, body, err)
	}
}

func HttpHandlerController(fn http.HandlerFunc) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		fn(resp, req.Request)
	}
}

func StaticFileAlias(alias StaticAlias, handler http.HandlerFunc) restful.RouteFunction {
	return HttpHandlerController(func(writer http.ResponseWriter, request *http.Request) {
		request.URL.Path = alias.Alias(request.URL.Path)
		handler(writer, request)
	})
}

func EnsureSlash(destination http.HandlerFunc) restful.RouteFunction {
	return HttpHandlerController(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.HasSuffix(request.URL.Path, "/") {
			http.RedirectHandler(request.URL.Path+"/", 302).ServeHTTP(writer, request)
			return
		}

		destination(writer, request)
	})
}
