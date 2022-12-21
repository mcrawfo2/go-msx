// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"github.com/pkg/errors"
	"github.com/swaggest/openapi-go/openapi3"
	"path"
	"strings"
)

type EndpointsDocumentor struct {
	customizer SpecificationCustomizer
}

func (d EndpointsDocumentor) Document(endpoints []*restops.Endpoint, contextPath string, servicePath string) error {
	// Document each endpoint
	for _, endpoint := range endpoints {
		endpointPath := d.cleanPath(endpoint.Path, contextPath, servicePath)

		if operation, err := d.DocumentEndpoint(endpoint); err != nil {
			logger.WithError(err).Errorf("Failed to document operation %q %q in openapi3 spec", endpoint.Method, endpoint.Path)
			continue
		} else if err = SetSpecOperation(endpoint.Method, endpointPath, operation); err != nil {
			logger.WithError(err).Errorf("Failed to save documentation for %q %q in openapi3 spec", endpoint.Method, endpointPath)
			continue
		}
	}

	// Customize the specification
	return d.customizer.PostBuildSpec(Spec())
}

func (d EndpointsDocumentor) DocumentEndpoint(e *restops.Endpoint) (result *openapi3.Operation, err error) {
	doc := ops.DocumentorWithType[restops.Endpoint](e, DocType).
		OrElse(new(EndpointDocumentor))

	if err = doc.Document(e); err != nil {
		return
	}

	if resulter, ok := doc.(ops.DocumentResult[openapi3.Operation]); !ok {
		err = errors.Errorf("Unable to retrieve operation documentation for %q %q", e.Method, e.Path)
		return
	} else {
		result = resulter.Result()
	}

	return
}

func (d EndpointsDocumentor) cleanPath(endpointPath string, contextPath string, servicePath string) string {
	endpointPath = strings.TrimPrefix(endpointPath, contextPath)
	endpointPath = strings.TrimPrefix(endpointPath, servicePath)
	endpointPath = path.Join(servicePath, endpointPath)
	return endpointPath
}

func NewEndpointsDocumentor(appInfo *schema.AppInfo, serverUrl string, version string) EndpointsDocumentor {
	return EndpointsDocumentor{
		customizer: SpecificationCustomizer{
			appInfo:   *appInfo,
			serverUrl: serverUrl,
			version:   version,
		},
	}
}
