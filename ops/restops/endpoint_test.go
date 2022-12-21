// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type TestDocumentor[I any] struct{}

func (d TestDocumentor[I]) DocType() string {
	return "Test"
}

func (d TestDocumentor[I]) Document(_ *I) error {
	return nil
}

func TestEndpoints_Each(t *testing.T) {
	tests := []struct {
		name    string
		e       Endpoints
		fn      EndpointFunc
		verify  EndpointsPredicate
		wantErr bool
	}{
		{
			name: "All",
			e: Endpoints{
				{
					Method: http.MethodGet,
				},
				{
					Method: http.MethodPost,
				},
			},
			fn: func(e *Endpoint) error {
				e.Description = strings.ToLower(e.Method)
				return nil
			},
			verify: func(e Endpoints) bool {
				for _, endpoint := range e {
					if endpoint.Description != strings.ToLower(endpoint.Method) {
						return false
					}
				}
				return true
			},
			wantErr: false,
		},
		{
			name: "Error",
			e: Endpoints{
				{
					Method: http.MethodGet,
				},
				{
					Method: http.MethodPost,
				},
			},
			fn: func(e *Endpoint) error {
				return errors.New("some error")
			},
			wantErr: true,
		},
		{
			name: "NoError",
			e:    Endpoints{},
			fn: func(e *Endpoint) error {
				return errors.New("some error")
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.e.Each(tt.fn)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.verify != nil {
				assert.True(t, tt.verify(tt.e))
			}
		})
	}
}

func TestNewEndpoint(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	assert.NotNil(t, e)
	assert.Equal(t, http.MethodGet, e.Method)
	assert.Equal(t, "a/b/c", e.Path)
}

func TestEndpoint_WithDocumentor(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithDocumentor(TestDocumentor[Endpoint]{})
	assert.NotNil(t, f)
}

func TestEndpoint_WithMethod(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithMethod(http.MethodPost)
	assert.Equal(t, http.MethodPost, f.Method)
}

func TestEndpoint_WithPath(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithPath("d", "e", "f")
	assert.Equal(t, "d/e/f", f.Path)
}

func TestEndpoint_WithOperationId(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithOperationId("test-operation-id")
	assert.Equal(t, "test-operation-id", f.OperationID)
}

func TestEndpoint_WithDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{
			name:        "Simple",
			description: "simple-description",
			want:        "simple-description",
		},
		{
			name: "Indented",
			description: `
				# My indented description
				* Another line
`,
			want: "\n" +
				"# My indented description\n" +
				"* Another line\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEndpoint(http.MethodGet, "a", "b", "c")
			f := e.WithDescription(tt.description)
			assert.True(t,
				tt.want == f.Description,
				testhelpers.Diff(tt.want, f.Description))
		})
	}
}

func TestEndpoint_WithSummary(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithSummary("test-summary")
	assert.Equal(t, "test-summary", f.Summary)
}

func TestEndpoint_WithTags(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithTags("test-tag-1", "test-tag-2")
	assert.Equal(t, []string{"test-tag-1", "test-tag-2"}, f.Tags)
}

func TestEndpoint_WithoutTags(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithTags("test-tag-1", "test-tag-2")
	g := f.WithoutTags()
	assert.Equal(t, []string(nil), g.Tags)
}

func TestEndpoint_WithDeprecated(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithDeprecated(true)
	assert.Equal(t, true, f.Deprecated)
}

func TestEndpoint_WithResponse(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithResponse(EndpointResponse{
		Envelope: true,
	})

	assert.Equal(t,
		EndpointResponse{
			Envelope: true,
		},
		f.Response)
}

func TestEndpoint_WithRequest(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithRequest(EndpointRequest{
		Description: "request-description",
	})

	assert.Equal(t,
		EndpointRequest{
			Description: "request-description",
		},
		f.Request)
}

func TestEndpoint_WithRequestParameter(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	rp := EndpointRequestParameter{
		Description: types.NewStringPtr("endpoint-request-parameter"),
	}
	f := e.WithRequestParameter(rp)
	assert.Len(t, f.Request.Parameters, 1)
	assert.Equal(t, types.NewStringPtr("endpoint-request-parameter"), f.Request.Parameters[0].Description)
}

func TestEndpoint_WithResponseCodes(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithResponseCodes(GetResponseCodes)
	assert.True(t,
		reflect.DeepEqual(GetResponseCodes, f.Response.Codes),
		testhelpers.Diff(GetResponseCodes, f.Response.Codes))
}

func TestEndpoint_WithResponseSuccessHeader(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	h := EndpointResponseHeader{
		Description: types.NewStringPtr("response-success-header"),
	}
	f := e.WithResponseSuccessHeader("X-Success", h)
	assert.Len(t, f.Response.Success.Headers, 1)
	assert.Len(t, f.Response.Error.Headers, 0)
	assert.Equal(t, h, f.Response.Success.Headers["X-Success"])
}

func TestEndpoint_WithResponseErrorHeader(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	h := EndpointResponseHeader{
		Description: types.NewStringPtr("response-error-header"),
	}
	f := e.WithResponseErrorHeader("X-Error", h)
	assert.Len(t, f.Response.Success.Headers, 0)
	assert.Len(t, f.Response.Error.Headers, 1)
	assert.Equal(t, h, f.Response.Error.Headers["X-Error"])
}

func TestEndpoint_WithResponseHeader(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	h := EndpointResponseHeader{
		Description: types.NewStringPtr("response-header"),
	}
	f := e.WithResponseHeader("X-Header", h)
	assert.Len(t, f.Response.Success.Headers, 1)
	assert.Equal(t, h, f.Response.Success.Headers["X-Header"])
	assert.Len(t, f.Response.Error.Headers, 1)
	assert.Equal(t, h, f.Response.Error.Headers["X-Header"])
}

func TestEndpoint_WithPermissionAnyOf(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithPermissionAnyOf("VIEW_TESTS")
	assert.Len(t, f.Permissions, 1)
	assert.Equal(t, "VIEW_TESTS", f.Permissions[0])
}

func TestEndpoint_WithHandler(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithHandler(func() {})
	assert.True(t, f.Func.IsPresent())
	g := f.WithHandler(nil)
	assert.False(t, g.Func.IsPresent())
}

func TestEndpoint_WithHttpHandler(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithHttpHandler(func(http.ResponseWriter, *http.Request) {})
	assert.True(t, f.Func.IsPresent())
	g := f.WithHttpHandler(nil)
	assert.False(t, g.Func.IsPresent())
}

func TestEndpoint_WithInputs(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithInputs(struct{}{})
	assert.True(t, f.Inputs.IsPresent())
	g := f.WithInputs(nil)
	assert.False(t, g.Inputs.IsPresent())
}

func TestEndpoint_WithOutputs(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c")
	f := e.WithOutputs(struct{}{})
	assert.True(t, f.Outputs.IsPresent())
	g := f.WithOutputs(nil)
	assert.False(t, g.Outputs.IsPresent())
}

func TestEndpoint_Build_ExplicitPorts(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c").
		WithOperationId("testOperation").
		WithInputs(struct{}{}).
		WithOutputs(struct{}{})

	_, err := e.Build()
	assert.Error(t, err)

	e.WithHandler(func() {})
	f, err := e.Build()
	assert.NotNil(t, f)
}

func TestEndpoint_Build_ImplicitPorts(t *testing.T) {
	e := NewEndpoint(http.MethodGet, "a", "b", "c").
		WithOperationId("testOperation").
		WithHandler(func(inp struct{}) (out struct{}, err error) {
			return struct{}{}, nil
		})
	f, err := e.Build()
	assert.NoError(t, err)
	assert.NotNil(t, f)
}
