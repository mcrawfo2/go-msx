package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"github.com/emicklei/go-restful"
	"net/http"
	"strings"
)

// Force response body or error into an envelope
func EnvelopeResponse(req *restful.Request, resp *restful.Response, body interface{}, err error) {
	if err != nil {
		status := http.StatusBadRequest
		if statusErr, ok := err.(StatusCodeProvider); ok {
			status = statusErr.StatusCode()
		}
		WriteErrorEnvelope(req, resp, status, err)
		return
	}

	status := http.StatusOK
	if statusProvider, ok := body.(StatusCodeProvider); ok {
		status = statusProvider.StatusCode()
	}

	WriteSuccessEnvelope(req, resp, status, body)
}

func RawResponse(req *restful.Request, resp *restful.Response, body interface{}, err error) {
	if err != nil {
		status := http.StatusBadRequest
		if statusErr, ok := err.(StatusCodeProvider); ok {
			status = statusErr.StatusCode()
		}
		WriteErrorEnvelope(req, resp, status, err)
		return
	}

	status := http.StatusOK
	if statusProvider, ok := body.(StatusCodeProvider); ok {
		status = statusProvider.StatusCode()
	}

	err = resp.WriteHeaderAndJson(status, body, MIME_JSON)
	if err != nil {
		logger.WithError(err).Error("Failed to write body")
	}
}

func WriteErrorEnvelope(req *restful.Request, resp *restful.Response, status int, err error) {
	message := func(err error) string {
		errMessage := err.Error()
		lines := strings.Split(errMessage, "\n")
		parts := strings.Split(lines[0], ": ")
		return parts[0]
	}(err)

	envelope := integration.MsxEnvelope{
		Success:    false,
		Command:    RouteOperationFromContext(req.Request.Context()),
		Params:     parameters(req),
		HttpStatus: integration.GetSpringStatusNameForCode(status),
		Message:    message,
		Throwable:  err.Error(), // TODO: stack trace
	}

	logger.WithError(err).Error("Request failed")

	err2 := resp.WriteHeaderAndJson(status, envelope, MIME_JSON)
	if err2 != nil {
		logger.WithError(err2).Error("Failed to write error envelope")
	}
}

func WriteSuccessEnvelope(req *restful.Request, resp *restful.Response, status int, body interface{}) {
	if body == nil {
		body = struct{}{}
	}

	var envelope integration.MsxEnvelope
	if bodyEnvelope, ok := body.(integration.MsxEnvelope); ok {
		envelope = bodyEnvelope
		if envelope.HttpStatus == "" {
			envelope.HttpStatus = integration.GetSpringStatusNameForCode(status)
		}
	} else if bodyPointerEnvelope, ok := body.(*integration.MsxEnvelope); ok {
		envelope = *bodyPointerEnvelope
	} else {
		envelope = integration.MsxEnvelope{
			Success:    true,
			Payload:    body,
			Command:    RouteOperationFromContext(req.Request.Context()),
			Params:     parameters(req),
			HttpStatus: integration.GetSpringStatusNameForCode(status),
		}

	}

	if err := resp.WriteHeaderAndJson(status, envelope, MIME_JSON); err != nil {
		logger.WithError(err).Error("Failed to write body")
	}
}

func parameters(req *restful.Request) (params map[string]interface{}) {
	params = make(map[string]interface{})
	for k, v := range req.PathParameters() {
		params[k] = v
	}
	if req.Request.Form != nil {
		for k, v := range req.Request.Form {
			params[k] = v
		}
	}
	if req.Request.PostForm != nil {
		for k, v := range req.Request.PostForm {
			params[k] = v
		}
	}
	return
}
