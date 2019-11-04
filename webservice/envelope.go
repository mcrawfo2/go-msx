package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"github.com/emicklei/go-restful"
	"strings"
)

func WriteErrorEnvelope(req *restful.Request, resp *restful.Response, status int, err error) error {
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

	return resp.WriteHeaderAndJson(status, envelope, "application/json")
}

func WriteJsonEnvelope(req *restful.Request, resp *restful.Response, status int, body interface{}) error {
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

	return resp.WriteHeaderAndJson(status, envelope, MIME_JSON)
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
