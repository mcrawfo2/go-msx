// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"github.com/emicklei/go-restful"
)

type RequestDescriber interface {
	Parameters() map[string]interface{}
	Path() string
}

type RestfulRequestDescriber struct {
	Request *restful.Request
}

func (r RestfulRequestDescriber) Path() string {
	return r.Request.Request.URL.Path
}

func (r RestfulRequestDescriber) Parameters() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range r.Request.PathParameters() {
		result[k] = v
	}
	return result
}

type EndpointRequestDescriber struct {
	Endpoint Endpoint
}

func (e EndpointRequestDescriber) Parameters() map[string]interface{} {
	result := make(map[string]interface{})
	for _, parameter := range e.Endpoint.Request.Parameters {
		if parameter.In == FieldGroupHttpPath {
			if parameter.Example.IsPresent() {
				result[parameter.Name] = parameter.Example.Value()
			} else if parameter.Payload.IsPresent() {
				result[parameter.Name] = parameter.Payload.Value()
			} else if parameter.Type != nil {
				v := e.exampleValue(*parameter.Type, parameter.Format)
				if v != nil {
					result[parameter.Name] = v
				}
			}
		}
	}
	return result
}

func (e EndpointRequestDescriber) exampleValue(dataType string, format *string) interface{} {
	switch dataType {
	case "integer":
		return 42
	case "number":
		return 3.14
	case "string":
		return "example"
	case "boolean":
		return true
	default:
		return nil
	}
}

func (e EndpointRequestDescriber) Path() string {
	return e.Endpoint.Path
}
