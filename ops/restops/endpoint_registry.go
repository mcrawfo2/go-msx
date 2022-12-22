// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"github.com/emicklei/go-restful"
	"sync"
)

var registeredEndpoints Endpoints
var registeredEndpointsMtx sync.Mutex

func RegisterEndpoint(endpoint *Endpoint) {
	registeredEndpointsMtx.Lock()
	defer registeredEndpointsMtx.Unlock()
	registeredEndpoints = append(registeredEndpoints, endpoint)
}

func RegisteredEndpoints() Endpoints {
	return registeredEndpoints
}

func RegisterRouteEndpoint(route restful.Route, basePath string) {
	endpoint := EndpointFromRoute(route)
	if endpoint.Method == "" {
		endpoint = NewEndpointFromRoute(route, basePath)
		RegisterEndpoint(endpoint)
	}
}
