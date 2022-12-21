// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRegisterEndpoint(t *testing.T) {
	registeredEndpoints = nil
	assert.Empty(t, registeredEndpoints)
	e := NewEndpoint(http.MethodConnect, "a", "b", "c")
	RegisterEndpoint(e)
	assert.NotEmpty(t, registeredEndpoints)
}

func TestRegisteredEndpoints(t *testing.T) {
	registeredEndpoints = Endpoints{new(Endpoint)}
	assert.Equal(t, registeredEndpoints, RegisteredEndpoints())
}

func TestRegisterRouteEndpoint(t *testing.T) {
	registeredEndpoints = nil
	route := restful.Route{
		Method: http.MethodOptions,
		Path:   "a/b/c",
	}
	RegisterRouteEndpoint(route, "")
	assert.NotEmpty(t, registeredEndpoints)
}
