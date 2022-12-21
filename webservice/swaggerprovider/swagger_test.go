// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSwaggerProvider_GetSecurity(t *testing.T) {
	p := &SwaggerProvider{}
	result, err := p.GetSecurity(nil)
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestSwaggerProvider_GetSwaggerResources(t *testing.T) {
	cfg, err := DocumentationConfigFromConfig(configtest.NewInMemoryConfig(map[string]string{
		"swagger.version": "3",
	}))
	assert.NoError(t, err)

	p := &SwaggerProvider{
		cfg: cfg,
	}
	result, err := p.GetSwaggerResources(restful.NewRequest(&http.Request{
		URL: &url.URL{
			Path: "/swagger-resources",
		},
	}))
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestSwaggerProvider_GetUi(t *testing.T) {
	p := &SwaggerProvider{}
	result, err := p.GetUi(nil)
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestSwaggerProvider_GetSsoSecurity(t *testing.T) {
	cfg, err := DocumentationConfigFromConfig(configtest.NewInMemoryConfig(map[string]string{
		"swagger.version": "2",
	}))
	assert.NoError(t, err)

	p := &SwaggerProvider{
		cfg: cfg,
	}
	result, err := p.GetSsoSecurity(nil)
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestSwaggerProvider_Spec(t *testing.T) {
	p := &SwaggerProvider{
		spec: &spec.Swagger{
			SwaggerProps: spec.SwaggerProps{
				Swagger: "2.0",
			},
		},
	}
	result, err := p.Spec(nil)
	assert.Equal(t, p.spec, result)
	assert.NoError(t, err)
}

func TestSwaggerProvider_YamlSpec(t *testing.T) {
	p := &SwaggerProvider{
		spec: &spec.Swagger{
			SwaggerProps: spec.SwaggerProps{
				Swagger: "2.0",
			},
		},
	}
	r := httptest.NewRecorder()
	resp := restful.NewResponse(r)

	p.YamlSpec(nil, resp)

	result := r.Result()
	assert.NotNil(t, result)

	resultBody, err := ioutil.ReadAll(result.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, resultBody)
}

func TestSwaggerProvider_Actuate(t *testing.T) {
	cfg, err := DocumentationConfigFromConfig(configtest.NewInMemoryConfig(map[string]string{
		"swagger.version":    "3",
		"swagger.ui.enabled": "false",
	}))
	assert.NoError(t, err)

	p := NewSwaggerProvider(
		context.Background(),
		cfg,
		&schema.AppInfo{
			Name:        "app-name",
			DisplayName: "display-name",
			Description: "description",
			Version:     "version",
		})

	c := new(restful.Container)
	c.ServeMux = http.NewServeMux()

	ws := new(restful.WebService)
	ws.Route(ws.GET("/api/v1/entity").
		To(func(request *restful.Request, response *restful.Response) {}))

	c.Add(ws)

	err = p.Actuate(c, ws)
	assert.NoError(t, err)
}
