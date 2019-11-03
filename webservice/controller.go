package webservice

import (
	"context"
	"github.com/emicklei/go-restful"
	"net/http"
)

type ControllerFunction func(req *restful.Request) (body interface{}, err error)
type ContextFunction func(ctx context.Context) (body interface{}, err error)

func rawResponse(req *restful.Request, resp *restful.Response, body interface{}, err error) {
	if err != nil {
		status := http.StatusBadRequest
		if statusErr, ok := err.(StatusProvider); ok {
			status = statusErr.Status()
		}
		if err = WriteErrorEnvelope(req, resp, status, err); err != nil {
			logger.WithError(err).Error("Failed to write error envelope")
		}
		return
	}

	status := http.StatusOK
	if statusProvider, ok := body.(StatusProvider); ok {
		status = statusProvider.Status()
	}

	err = resp.WriteHeaderAndJson(status, body, MIME_JSON)
	if err != nil {
		logger.WithError(err).Error("Failed to write body")
	}
}

// Force only error into an envelope
func RawController(fn ControllerFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		body, err := fn(req)
		rawResponse(req, resp, body, err)
	}
}

func RawContextController(fn ContextFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()

		body, err := fn(ctx)
		rawResponse(req, resp, body, err)
	}
}

// Force response body or error into an envelope
func envelopeResponse(req *restful.Request, resp *restful.Response, body interface{}, err error) {
	if err != nil {
		status := http.StatusBadRequest
		if statusErr, ok := err.(StatusProvider); ok {
			status = statusErr.Status()
		}
		err = WriteErrorEnvelope(req, resp, status, err)
		if err != nil {
			logger.WithError(err).Error("Failed to write error envelope")
		}
		return
	}

	status := http.StatusOK
	if statusProvider, ok := body.(StatusProvider); ok {
		status = statusProvider.Status()
	}

	if err = WriteJsonEnvelope(req, resp, status, body); err != nil {
		logger.WithError(err).Error("Failed to write body")
	}
}

func Controller(fn ControllerFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		body, err := fn(req)
		envelopeResponse(req, resp, body, err)
	}
}

func ContextController(fn ContextFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()

		body, err := fn(ctx)
		envelopeResponse(req, resp, body, err)
	}
}

func HttpHandlerController(fn http.HandlerFunc) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		fn(resp.ResponseWriter, req.Request)
	}
}