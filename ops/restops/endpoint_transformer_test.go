// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEndpointTransformers_Transform(t *testing.T) {
	i := 0
	x := func(e *Endpoint) { i = i + 1 }
	e := Endpoints{
		{Method: http.MethodGet},
		{Method: http.MethodPut},
	}
	tfs := EndpointTransformers{x, x, x}
	tfs.Transform(e)
	assert.Equal(t, 6, i)
}

func TestAddEndpointTag(t *testing.T) {
	e := Endpoint{Method: http.MethodGet}
	tf := AddEndpointTag("tag")
	tf(&e)
	assert.Equal(t, []string{"tag"}, e.Tags)
}

func TestAddEndpointPathPrefix(t *testing.T) {
	e := Endpoint{Method: http.MethodGet, Path: "a/b/c"}
	tf := AddEndpointPathPrefix("v1")
	tf(&e)
	assert.Equal(t, "v1/a/b/c", e.Path)
}

func TestAddEndpointRequestParameter(t *testing.T) {
	e := Endpoint{Method: http.MethodGet, Path: "a/b/c"}
	p := EndpointRequestParameter{
		Name: "cookie",
		In:   FieldGroupHttpCookie,
	}
	tf := AddEndpointRequestParameter(p)
	tf(&e)
	assert.Equal(t, p, e.Request.Parameters[0])
}

func TestAddEndpointErrorConverter(t *testing.T) {
	e := Endpoint{Method: http.MethodGet, Path: "a/b/c"}
	tf := AddEndpointErrorConverter(ErrorConverterFunc(func(err error) StatusCodeError {
		return webservice.NewStatusCodeError(err, 400)
	}))
	tf(&e)
	assert.NotNil(t, e.ErrorConverter)
}

func TestAddEndpointErrorCoder(t *testing.T) {
	e := Endpoint{Method: http.MethodGet, Path: "a/b/c"}
	tf := AddEndpointErrorCoder(ErrorStatusCoderFunc(func(err error) int {
		return 400
	}))
	tf(&e)
	assert.NotNil(t, e.ErrorConverter)
}

func TestAddEndpointMiddleware(t *testing.T) {
	e := Endpoint{Method: http.MethodGet, Path: "a/b/c"}
	tf := AddEndpointMiddleware(func(next http.Handler) http.Handler {
		return next
	})
	tf(&e)
	assert.NotNil(t, e.Middleware)
}

func TestAddEndpointContextInjector(t *testing.T) {
	e := Endpoint{Method: http.MethodGet, Path: "a/b/c"}
	tf := AddEndpointContextInjector(func(ctx context.Context) context.Context {
		return ctx
	})
	tf(&e)
	assert.NotNil(t, e.Injectors)
}
