package webservice

import (
	"github.com/emicklei/go-restful"
	"strings"
)

type ResponseEnvelope struct {
	Command    string                 `json:"command"`
	Debug      map[string]interface{} `json:"debug,omitEmpty"`
	Errors     []string               `json:"errors,omitEmpty"`
	HttpStatus string                 `json:"httpStatus"`
	Message    string                 `json:"message"`
	Params     map[string]interface{} `json:"params"`
	Payload    interface{}            `json:"responseObject"`
	Success    bool                   `json:"success"`
	Throwable  interface{}            `json:"throwable"`
}

func (r ResponseEnvelope) StatusCode() int {
	return getSpringStatusCodeForName(r.HttpStatus)
}

func WriteErrorEnvelope(req *restful.Request, resp *restful.Response, status int, err error) error {
	message := func(err error) string {
		errMessage := err.Error()
		lines := strings.Split(errMessage, "\n")
		parts := strings.Split(lines[0], ": ")
		return parts[0]
	}(err)

	envelope := ResponseEnvelope{
		Success:    false,
		Command:    RouteOperationFromContext(req.Request.Context()),
		Params:     Parameters(req),
		HttpStatus: getSpringStatusNameForCode(status),
		Message:    message,
		Throwable:  err.Error(), // TODO: stack trace
	}

	return resp.WriteHeaderAndJson(status, envelope, "application/json")
}

func WriteJsonEnvelope(req *restful.Request, resp *restful.Response, status int, body interface{}) error {
	if body == nil {
		body = struct{}{}
	}

	var envelope ResponseEnvelope
	if bodyEnvelope, ok := body.(ResponseEnvelope); ok {
		envelope = bodyEnvelope
		if envelope.HttpStatus == "" {
			envelope.HttpStatus = getSpringStatusNameForCode(status)
		}
	} else {
		envelope = ResponseEnvelope{
			Success:    true,
			Payload:    body,
			Command:    RouteOperationFromContext(req.Request.Context()),
			Params:     Parameters(req),
			HttpStatus: getSpringStatusNameForCode(status),
		}

	}

	return resp.WriteHeaderAndJson(status, envelope, MIME_JSON)
}

func Parameters(req *restful.Request) (params map[string]interface{}) {
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
