// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
)

func Test_recoveryFilter(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoPanic",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteFilter(recoveryFilter).
				WithRouteFilter(DummyFilter).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					http.NotFound(response, request.Request)
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(404)).
				WithResponsePredicate(DummyFilterResponseCheck),
		},
		{
			name: "PanicString",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteFilter(recoveryFilter).
				WithRouteFilter(DummyFilter).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					panic("panic")
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(500)).
				WithResponsePredicate(DummyFilterResponseCheck).
				WithLogCheck(logtest.Check{
					Filters: []logtest.EntryPredicate{
						logtest.HasMessage("Recovered from panic"),
					},
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.ErrorLevel),
						logtest.HasMessage("Recovered from panic"),
						logtest.FieldValue("error"),
						logtest.FieldValue(log.FieldStack),
					},
				}),
		},
		{
			name: "PanicError",
			test: new(webservicetest.RouteBuilderTest).
				WithRouteFilter(recoveryFilter).
				WithRouteFilter(DummyFilter).
				WithRouteTarget(func(request *restful.Request, response *restful.Response) {
					panic(errors.New("panic"))
				}).
				WithResponsePredicate(webservicetest.ResponseHasStatus(500)).
				WithResponsePredicate(DummyFilterResponseCheck).
				WithLogCheck(logtest.Check{
					Filters: []logtest.EntryPredicate{
						logtest.HasMessage("Recovered from panic"),
					},
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.ErrorLevel),
						logtest.HasMessage("Recovered from panic"),
						logtest.FieldValue("error"),
						logtest.FieldValue(log.FieldStack),
					},
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}
