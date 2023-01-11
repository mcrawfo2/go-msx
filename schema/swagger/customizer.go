// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swagger

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
	"reflect"
	"sort"
	"strings"
)

type SpecificationCustomizer struct {
	container   *restful.Container
	service     *restful.WebService
	appInfo     schema.AppInfo
	contextPath string
}

func (c SpecificationCustomizer) PostBuildSpec(swagger *spec.Swagger) {
	c.CustomizeInfo(swagger)
	c.CustomizeTags(swagger)
	c.CustomizeBasePath(swagger)
	c.CustomizeTypeDefinitions(swagger)
	c.SortTags(swagger)
}

func (c SpecificationCustomizer) CustomizeInfo(swagger *spec.Swagger) {
	swagger.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title: "MSX API Documentation for " + c.appInfo.Name,
			Description: "<h3>This is the REST API documentation for " + c.appInfo.Name + "</h3>\n \n" +
				c.appInfo.Description + "\n" +
				"+ API Authorization \n" +
				"    + Authorization header is <b>required</b>. \n" +
				"    + It should be in Bearer authentication scheme </br>(e.g <b> Authorization: BEARER &lt;access token&gt; </b>)\n",
			TermsOfService: "https://www.cisco.com",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "Cisco Systems Inc.",
					URL:   "https://www.cisco.com/c/en/us/support/cloud-systems-management/managed-services-accelerator/series.html",
					Email: "tac@cisco.com",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "https://opensource.org/licenses/MIT",
				},
			},
			Version: c.appInfo.Version,
		},
	}
}

func (c SpecificationCustomizer) CustomizeTags(swagger *spec.Swagger) {
	// Register tags definitions from all routes
	var existingTags = types.StringStack{}
	for _, svc := range c.container.RegisteredWebServices() {
		for _, route := range svc.Routes() {
			if routeTagDefinition, ok := webservice.TagDefinitionFromRoute(route); !ok {
				continue
			} else if !existingTags.Contains(routeTagDefinition.Name) {
				existingTags = append(existingTags, routeTagDefinition.Name)
				swagger.Tags = append(swagger.Tags, spec.Tag{TagProps: routeTagDefinition})
			}
		}
	}
}

func (c SpecificationCustomizer) CustomizeBasePath(swagger *spec.Swagger) {
	// Factor out contextPath into basePath
	if c.contextPath != "/" {
		newPaths := make(map[string]spec.PathItem)
		for path, pathItem := range swagger.Paths.Paths {
			if strings.HasPrefix(path, c.contextPath) {
				path = strings.TrimPrefix(path, c.contextPath)
			}
			newPaths[path] = pathItem
		}
		swagger.Paths.Paths = newPaths
		swagger.BasePath = c.contextPath
	}
}

func (c SpecificationCustomizer) SortTags(swagger *spec.Swagger) {
	sort.Slice(swagger.Tags, func(i, j int) bool {
		iTagName := swagger.Tags[i].Name
		jTagName := swagger.Tags[j].Name
		return strings.Compare(iTagName, jTagName) < 0
	})
}

func (c SpecificationCustomizer) CustomizeTypeDefinitions(swagger *spec.Swagger) {
	var schemaSources = []SchemaSource{
		new(types.Time),
		new(types.UUID),
		new(types.Empty),
	}

	for _, schemaSource := range schemaSources {
		schemaType := reflect.TypeOf(schemaSource).Elem()
		typeName := schema.NewSpringTypeNamer().TypeName(schemaType)
		schemaJson := schemaSource.SwaggerSchemaJson()

		var schemaDef *spec.Schema
		if err := json.Unmarshal([]byte(schemaJson), &schemaDef); err != nil {
			logger.WithError(err).Errorf("Failed to parse Swagger Schema for %q", typeName)
			continue
		}

		swagger.Definitions[typeName] = *schemaDef
	}
}

type SchemaSource interface {
	SwaggerSchemaJson() string
}
