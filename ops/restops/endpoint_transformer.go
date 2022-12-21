// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"path"
)

type EndpointTransformer func(*Endpoint)

type EndpointTransformersProducer interface {
	EndpointTransformers() EndpointTransformers
}

type EndpointTransformers []EndpointTransformer

func (t EndpointTransformers) Transform(endpoints Endpoints) {
	for _, endpoint := range endpoints {
		for _, transformer := range t {
			transformer(endpoint)
		}
	}
}

// Transformers

func AddEndpointTag(tag string) EndpointTransformer {
	return func(endpoint *Endpoint) {
		endpoint.Tags = []string{tag}
	}
}

func AddEndpointPathPrefix(prefix string) EndpointTransformer {
	return func(endpoint *Endpoint) {
		endpoint.Path = path.Join(prefix, endpoint.Path)
	}
}

func AddEndpointRequestParameter(parameter EndpointRequestParameter) EndpointTransformer {
	return func(endpoint *Endpoint) {
		endpoint.WithRequestParameter(parameter)
	}
}

func AddEndpointErrorConverter(converter ErrorConverter) EndpointTransformer {
	return func(endpoint *Endpoint) {
		endpoint.ErrorConverter = converter
	}
}

func AddEndpointErrorCoder(coder ErrorStatusCoder) EndpointTransformer {
	return func(endpoint *Endpoint) {
		endpoint.ErrorConverter = ErrorStatusCoderConverter{ErrorStatusCoder: coder}
	}
}
