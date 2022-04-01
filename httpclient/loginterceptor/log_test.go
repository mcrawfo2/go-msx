// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package loginterceptor

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewInterceptor(t *testing.T) {
	request, _ := http.NewRequest("PATCH", "http://manageddevice/api/v1/devices/CPE-XYZ", http.NoBody)

	tests := []struct {
		name     string
		response *http.Response
		err      error
		wantLog  logtest.Check
	}{
		{
			name: "Success",
			response: &http.Response{
				Status:     "OK",
				StatusCode: 200,
				Request:    request,
			},
			err: nil,
			wantLog: logtest.Check{
				Filters: []logtest.EntryPredicate{
					logtest.HasLevel(logrus.InfoLevel),
				},
				Validators: []logtest.EntryPredicate{
					logtest.HasMessage("200 OK : PATCH http://manageddevice/api/v1/devices/CPE-XYZ"),
				},
			},
		},
		{
			name:     "No Response",
			response: nil,
			err:      errors.New("Client Failure"),
			wantLog: logtest.Check{
				Filters: []logtest.EntryPredicate{
					logtest.HasLevel(logrus.ErrorLevel),
				},
				Validators: []logtest.EntryPredicate{
					logtest.HasMessage("000 : PATCH http://manageddevice/api/v1/devices/CPE-XYZ"),
					logtest.HasError("Client Failure"),
				},
			},
		},
		{
			name: "Error Response",
			response: &http.Response{
				Status:        "NOT FOUND",
				StatusCode:    404,
				Body:          ioutil.NopCloser(bytes.NewBufferString("{}")),
				ContentLength: 2,
				Request:       request,
			},
			err: nil,
			wantLog: logtest.Check{
				Filters: []logtest.EntryPredicate{
					logtest.HasLevel(logrus.ErrorLevel),
					logtest.Index(0, 1),
				},
				Validators: []logtest.EntryPredicate{
					logtest.HasMessage("404 NOT FOUND : PATCH http://manageddevice/api/v1/devices/CPE-XYZ"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := logtest.RecordLogging()

			got := NewInterceptor(func(*http.Request) (*http.Response, error) {
				return tt.response, tt.err
			})

			_, _ = got(request)

			testhelpers.ReportErrors(t, "Log", tt.wantLog.Check(r))
		})
	}
}
