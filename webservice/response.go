// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"errors"
	"net/http"
	"reflect"

	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
)

// EnvelopeResponse forces the response body or error into an envelope
func EnvelopeResponse(req *restful.Request, resp *restful.Response, body interface{}, err error) {
	if err != nil {
		status := http.StatusBadRequest
		if statusErr, ok := err.(StatusCodeProvider); ok {
			status = statusErr.StatusCode()
		}
		WriteError(req, resp, status, err)
		return
	}

	status := http.StatusOK
	if req.Attribute(AttributeDefaultReturnCode) != nil {
		status = req.Attribute(AttributeDefaultReturnCode).(int)
	}
	if body != nil {
		if statusProvider, ok := body.(StatusCodeProvider); ok {
			status = statusProvider.StatusCode()
		}
	}

	WriteSuccessEnvelope(req, resp, status, body)
}

func RawResponse(req *restful.Request, resp *restful.Response, body interface{}, err error) {
	if err != nil {
		status := http.StatusBadRequest
		if statusErr, ok := err.(StatusCodeProvider); ok {
			status = statusErr.StatusCode()
		}
		WriteError(req, resp, status, err)
		return
	}

	status := http.StatusOK
	if req.Attribute(AttributeDefaultReturnCode) != nil {
		status = req.Attribute(AttributeDefaultReturnCode).(int)
	}
	if statusProvider, ok := body.(StatusCodeProvider); ok {
		status = statusProvider.StatusCode()
	}

	err = resp.WriteHeaderAndJson(status, body, MIME_JSON_CHARSET)
	if err != nil {
		logger.WithError(err).Error("Failed to write body")
	}
}

func WriteError(req *restful.Request, resp *restful.Response, status int, err error) {
	span := trace.SpanFromContext(req.Request.Context())
	if span != nil {
		span.LogFields(trace.Error(resp.Error()))
	}
	RequestWithError(req, err)

	logger.
		WithContext(req.Request.Context()).
		WithField("status", status).
		WithError(err).
		Error("Request failed")

	payload := ErrorPayloadFromRequest(req)
	if payload == nil {
		WriteErrorEnvelope(req, resp, status, err)
		return
	}

	switch payload.(type) {
	case *integration.MsxEnvelope:
		WriteErrorEnvelope(req, resp, status, err)

	case ErrorApplier:
		WriteErrorApplier(req, resp, status, err, payload.(ErrorApplier))

	case ErrorRaw:
		WriteErrorRaw(req, resp, status, err, payload.(ErrorRaw))

	default:
		logger.
			WithContext(req.Request.Context()).
			WithError(err).
			Errorf("Response serialization failed - invalid error payload type %q", reflect.TypeOf(payload).String())
		WriteErrorRaw(req, resp, status, err, new(integration.ErrorDTO))
	}
}

func WriteErrorRaw(req *restful.Request, resp *restful.Response, status int, err error, payload ErrorRaw) {
	envelope := reflect.New(reflect.TypeOf(payload).Elem()).Interface().(ErrorRaw)
	envelope.SetError(status, err, req.Request.URL.Path)

	err2 := resp.WriteHeaderAndJson(status, envelope, MIME_JSON_CHARSET)
	if err2 != nil {
		logger.WithContext(req.Request.Context()).WithError(err2).Error("Failed to write error payload")
	}
}

func WriteErrorApplier(req *restful.Request, resp *restful.Response, status int, err error, payload ErrorApplier) {
	envelope := reflect.New(reflect.TypeOf(payload).Elem()).Interface().(ErrorApplier)
	envelope.ApplyError(err)

	err2 := resp.WriteHeaderAndJson(status, envelope, MIME_JSON_CHARSET)
	if err2 != nil {
		logger.WithContext(req.Request.Context()).WithError(err2).Error("Failed to write error payload")
	}
}

func WriteErrorEnvelope(req *restful.Request, resp *restful.Response, status int, err error) {
	envelope := integration.MsxEnvelope{
		Success:    false,
		Message:    err.Error(),
		Command:    RouteOperationFromContext(req.Request.Context()),
		Params:     parameters(req),
		HttpStatus: integration.GetSpringStatusNameForCode(status),
		Throwable:  integration.NewThrowable(err),
	}

	var errorList types.ErrorList
	if errors.As(err, &errorList) {
		envelope.Errors = errorList.Strings()
	}

	if err := resp.WriteHeaderAndJson(status, envelope, MIME_JSON_CHARSET); err != nil {
		logger.WithContext(req.Request.Context()).WithError(err).Error("Failed to write error envelope")
	}
}

func WriteSuccessEnvelope(req *restful.Request, resp *restful.Response, status int, body interface{}) {
	if body == nil {
		body = struct{}{}
	}

	var envelope integration.MsxEnvelope
	if bodyEnvelope, ok := body.(integration.MsxEnvelope); ok {
		envelope = bodyEnvelope
	} else if bodyPointerEnvelope, ok := body.(*integration.MsxEnvelope); ok {
		envelope = *bodyPointerEnvelope
	} else {
		envelope = integration.MsxEnvelope{
			Success: true,
			Payload: body,
			Command: RouteOperationFromContext(req.Request.Context()),
			Params:  parameters(req),
		}
	}

	if envelope.HttpStatus == "" {
		envelope.HttpStatus = integration.GetSpringStatusNameForCode(status)
	}

	if err := resp.WriteHeaderAndJson(status, envelope, MIME_JSON_CHARSET); err != nil {
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
