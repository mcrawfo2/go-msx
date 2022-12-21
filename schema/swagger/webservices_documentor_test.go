// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swagger

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestStruct struct{}

func TestWebServicesDocumentor_ModelTypeName(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "BuiltIn",
			args: args{
				t: reflect.TypeOf("string"),
			},
			want:  "string",
			want1: true,
		},
		{
			name: "Struct",
			args: args{
				t: reflect.TypeOf(TestStruct{}),
			},
			want:  "swagger.TestStruct",
			want1: true,
		},
		{
			name: "ListOfStruct",
			args: args{
				t: reflect.TypeOf([]TestStruct{}),
			},
			want:  "List«swagger.TestStruct»",
			want1: true,
		},
		{
			name: "MapOfStruct",
			args: args{
				t: reflect.TypeOf(map[string]TestStruct{}),
			},
			want:  "Map«string,swagger.TestStruct»",
			want1: true,
		},
		{
			name: "EnvelopedListOfStruct",
			args: args{
				t: schema.NewParameterizedStruct(
					reflect.TypeOf(integration.MsxEnvelope{}),
					[]TestStruct{},
				),
			},
			want:  "Envelope«List«swagger.TestStruct»»",
			want1: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			documentor := WebServicesDocumentor{
				namer: schema.NewSpringTypeNamer(),
			}
			got, got1 := documentor.ModelTypeName(tt.args.t)
			if got != tt.want {
				t.Errorf("ModelTypeName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ModelTypeName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestWebServicesDocumentor_DocType(t *testing.T) {
	wsd := new(WebServicesDocumentor)
	assert.Equal(t, DocType, wsd.DocType())
}

func TestWebServicesDocumentor_Document(t *testing.T) {
	container := restful.NewContainer()

	ws := new(restful.WebService)
	ws.Path("/api")

	rb := ws.GET("/v1/entities").
		Operation("v1.listAllEntities").
		Doc("List all entities").
		Param(restful.QueryParameter("max", "Maximum number of entities to return")).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Writes(types.Poja{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{"Entity", "List"}).
		Metadata(webservice.MetadataTagDefinition, spec.TagProps{
			Name:        "Entity",
			Description: "Entity Methods",
		}).
		Metadata(webservice.MetadataTagDefinition, spec.TagProps{
			Name:        "List",
			Description: "List Methods",
		}).
		To(func(request *restful.Request, response *restful.Response) {
			response.WriteHeader(404)
			_ = response.WriteEntity(types.Poja{})
		})
	ws.Route(rb)

	rb = ws.POST("/v1/entities").
		Operation("v1.createEntity").
		Doc("Create entity").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Reads(types.Pojo{}).
		Writes(types.Pojo{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{"Entity", "Create"}).
		To(func(request *restful.Request, response *restful.Response) {
			response.WriteHeader(404)
			_ = response.WriteEntity(types.Poja{})
		})
	ws.Route(rb)

	container.Add(ws)

	swaggerService := new(restful.WebService).Path("/swagger")

	appInfo := schema.AppInfo{
		Name:        "app-name",
		DisplayName: "display-name",
		Description: "description",
		Version:     "version",
	}

	wsd := NewWebServicesDocumentor(container, swaggerService, "/context", &appInfo)

	err := wsd.Document(&WebServices{
		WebServices: container.RegisteredWebServices(),
		APIPath:     "/apidocs.json",
	})
	assert.NoError(t, err)

	spec := wsd.DocumentResult()
	assert.NotNil(t, spec)
}
