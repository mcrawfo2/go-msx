// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swagger

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"reflect"
)

var swagger *spec.Swagger

func Spec() *spec.Swagger {
	return swagger
}

type WebServices struct {
	// api listing is constructed from this list of restful WebServices.
	WebServices []*restful.WebService
	// APIPath is the path where the JSON api is available , e.g. /apidocs.json
	APIPath string
}

type WebServicesDocumentor struct {
	customizer SpecificationCustomizer
	namer      schema.TypeNamer
	swagger    *spec.Swagger
}

func (w *WebServicesDocumentor) DocType() string {
	return DocType
}

func (w *WebServicesDocumentor) Document(s *WebServices) error {
	w.swagger = restfulspec.BuildSwagger(restfulspec.Config{
		WebServices:                   s.WebServices,
		APIPath:                       s.APIPath,
		ModelTypeNameHandler:          w.ModelTypeName,
		PostBuildSwaggerObjectHandler: w.customizer.PostBuildSpec,
	})

	// Set globally
	swagger = w.swagger

	return nil
}

func (w *WebServicesDocumentor) DocumentResult() *spec.Swagger {
	return w.swagger
}

func (w *WebServicesDocumentor) ModelTypeName(t reflect.Type) (string, bool) {
	typeName := w.namer.TypeName(t)
	return typeName, typeName != ""
}

func NewWebServicesDocumentor(container *restful.Container, swaggerService *restful.WebService, contextPath string, appInfo *schema.AppInfo) *WebServicesDocumentor {
	return &WebServicesDocumentor{
		customizer: SpecificationCustomizer{
			container:   container,
			service:     swaggerService,
			appInfo:     *appInfo,
			contextPath: contextPath,
		},
		namer:   schema.NewSpringTypeNamer(),
		swagger: new(spec.Swagger),
	}
}
