// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestEnvelopeResponse(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "ErrorStatusCodeProvider",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					err := NewStatusError(errors.New("some error"), 444)
					EnvelopeResponse(request, response, nil, err)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(444)),
		},
		{
			name: "NoDefaultReturnCode",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					EnvelopeResponse(request, response, nil, nil)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)),
		},
		{
			name: "DefaultReturnCode",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteBuilderDo(DefaultReturns(202)).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					EnvelopeResponse(request, response, nil, nil)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(202)),
		},
		{
			name: "BodyStatusCodeProvider",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					body := NewStatusCodeProvider(nil, 204)
					EnvelopeResponse(request, response, body, nil)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(204)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestRawResponse(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "ErrorStatusCodeProvider",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					err := NewStatusError(errors.New("some error"), 444)
					RawResponse(request, response, nil, err)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(444)),
		},
		{
			name: "NoDefaultReturnCode",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					RawResponse(request, response, nil, nil)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)),
		},
		{
			name: "DefaultReturnCode",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteBuilderDo(DefaultReturns(202)).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					RawResponse(request, response, nil, nil)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(202)),
		},
		{
			name: "BodyStatusCodeProvider",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					body := NewStatusCodeProvider(nil, 204)
					RawResponse(request, response, body, nil)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(204)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoErrorPayload",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					WriteError(request, response, 409, errors.New("some error"))
				}).
				WithLogCheck(logtest.Check{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.ErrorLevel),
						logtest.HasFieldValue("status", 409),
						logtest.HasMessage("Request failed"),
					},
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(409)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`success`, false)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`message`, "some error")),
		},
		{
			name: "EnvelopeErrorPayload",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteBuilderDo(ErrorPayload(new(integration.MsxEnvelope))).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					WriteError(request, response, 400, errors.New("some error"))
				}).
				WithLogCheck(logtest.Check{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.ErrorLevel),
						logtest.HasFieldValue("status", 400),
						logtest.HasMessage("Request failed"),
					},
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(400)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`success`, false)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`message`, "some error")),
		},
		{
			name: "RawErrorPayload",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteBuilderDo(ErrorPayload(new(integration.ErrorDTO))).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					WriteError(request, response, 401, errors.New("some error"))
				}).
				WithLogCheck(logtest.Check{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.ErrorLevel),
						logtest.HasFieldValue("status", 401),
						logtest.HasMessage("Request failed"),
					},
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(401)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`code`, "401")).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`message`, "some error")),
		},
		{
			name: "UnknownErrorPayload",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteBuilderDo(ErrorPayload(new(struct{}))).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					WriteError(request, response, 403, errors.New("some error"))
				}).
				WithLogCheck(logtest.Check{
					Filters: []logtest.EntryPredicate{
						logtest.HasMessage("Request failed"),
					},
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.ErrorLevel),
						logtest.HasFieldValue("status", 403),
						logtest.HasMessage("Request failed"),
					},
				}).
				WithLogCheck(logtest.Check{
					Filters: []logtest.EntryPredicate{
						logtest.HasMessage(`Response serialization failed - invalid error payload type "*struct {}"`),
					},
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.ErrorLevel),
						logtest.HasMessage(`Response serialization failed - invalid error payload type "*struct {}"`),
					},
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(403)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`code`, "403")).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`message`, "some error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestWriteErrorEnvelope(t *testing.T) {
	t.Skipped()
}

func TestWriteErrorRaw(t *testing.T) {
	t.Skipped()
}

func TestWriteSuccessEnvelope(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoPayload",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					WriteSuccessEnvelope(request, response, 200, nil)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(200)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`success`, true)),
		},
		{
			name: "StructPayload",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					WriteSuccessEnvelope(request, response, 201, struct{}{})
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(201)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`success`, true)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`responseObject`, map[string]interface{}{})),
		},
		{
			name: "EnvelopePayload",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					envelope := &integration.MsxEnvelope{
						Command: "command",
						Message: "message",
						Payload: struct{}{},
						Success: true,
					}
					WriteSuccessEnvelope(request, response, 202, envelope)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(202)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`success`, true)).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`message`, "message")).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`responseObject`, map[string]interface{}{})).
				WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(`command`, "command")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func Test_parameters(t *testing.T) {
	t.Skipped()
}
