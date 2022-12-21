// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

type TestRequestDescriber struct {
	ReturnParameters map[string]interface{}
	ReturnPath       string
}

func (t TestRequestDescriber) Parameters() map[string]interface{} {
	return t.ReturnParameters
}

func (t TestRequestDescriber) Path() string {
	return t.ReturnPath
}

func TestEndpointRequestDescriber_Parameters(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
		want     map[string]interface{}
	}{
		{
			name: "Examples",
			endpoint: *NewEndpoint("GET", "a", "b", "c").
				WithRequestParameter(NewEndpointRequestParameter("p1", FieldGroupHttpPath).
					WithPayload("bob")).
				WithRequestParameter(NewEndpointRequestParameter("p2", FieldGroupHttpPath).
					WithType("string")).
				WithRequestParameter(NewEndpointRequestParameter("p3", FieldGroupHttpPath).
					WithType("integer")).
				WithRequestParameter(NewEndpointRequestParameter("p4", FieldGroupHttpPath).
					WithType("boolean")).
				WithRequestParameter(NewEndpointRequestParameter("p5", FieldGroupHttpPath).
					WithType("number")).
				WithRequestParameter(NewEndpointRequestParameter("p6", FieldGroupHttpPath)),
			want: map[string]interface{}{
				"p1": "bob",
				"p2": "example",
				"p3": 42,
				"p4": true,
				"p5": 3.14,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EndpointRequestDescriber{
				Endpoint: tt.endpoint,
			}
			assert.Equal(t, tt.want, e.Parameters())
		})
	}
}

func TestEndpointRequestDescriber_Path(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
		want     string
	}{
		{
			name:     "Path",
			endpoint: *NewEndpoint(http.MethodGet, "a", "b", "c"),
			want:     "a/b/c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EndpointRequestDescriber{
				Endpoint: tt.endpoint,
			}
			assert.Equalf(t, tt.want, e.Path(), "Path()")
		})
	}
}

func TestRestfulRequestDescriber_Parameters(t *testing.T) {
	tests := []struct {
		name    string
		request *restful.Request
		want    map[string]interface{}
	}{
		{
			name:    "Nothing",
			request: restful.NewRequest(&http.Request{}),
			want:    map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RestfulRequestDescriber{
				Request: tt.request,
			}
			assert.Equal(t, tt.want, r.Parameters())
		})
	}
}

func TestRestfulRequestDescriber_Path(t *testing.T) {
	tests := []struct {
		name    string
		request *restful.Request
		want    string
	}{
		{
			name: "Path",
			request: restful.NewRequest(&http.Request{
				URL: &url.URL{
					Path: "/a/b/c",
				},
			}),
			want: "/a/b/c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RestfulRequestDescriber{
				Request: tt.request,
			}
			assert.Equal(t, tt.want, r.Path())
		})
	}
}
