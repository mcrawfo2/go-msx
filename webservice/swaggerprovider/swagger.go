// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/swagger"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"encoding/json"
	"github.com/emicklei/go-restful"
	spec "github.com/go-openapi/spec"
	"github.com/pkg/errors"
	yaml2 "gopkg.in/yaml.v2"
	"path"
)

var (
	ErrDisabled = errors.New("Swagger disabled")
)

type SwaggerProvider struct {
	ctx     context.Context
	cfg     *DocumentationConfig
	appInfo *schema.AppInfo
	spec    *spec.Swagger
}

func (p SwaggerProvider) GetSecurity(req *restful.Request) (body interface{}, err error) {
	return struct{}{}, nil
}

func (p SwaggerProvider) GetSwaggerResources(req *restful.Request) (body interface{}, err error) {
	return []struct {
		Name           string `json:"name"`
		Location       string `json:"location"`
		Url            string `json:"url"`
		SwaggerVersion string `json:"swaggerVersion"`
	}{
		{
			Name:           "platform",
			Location:       path.Join(p.cfg.SwaggerPath, p.cfg.ApiPath),
			Url:            path.Join(p.cfg.SwaggerPath, p.cfg.ApiPath),
			SwaggerVersion: p.cfg.Version,
		},
	}, nil
}

func (p SwaggerProvider) GetUi(req *restful.Request) (body interface{}, err error) {
	return struct {
		ApisSorter               string   `json:"apisSorter"`
		DeepLinking              bool     `json:"deepLinking"`
		DefaultModelExpandDepth  int      `json:"defaultModelExpandDepth"`
		DefaultModelRendering    string   `json:"defaultModelRendering"`
		DefaultModelsExpandDepth int      `json:"defaultModelsExpandDepth"`
		DisplayOperationId       bool     `json:"displayOperationId"`
		DisplayRequestDuration   bool     `json:"displayRequestDuration"`
		DocExpansion             string   `json:"docExpansion"`
		Filter                   bool     `json:"filter"`
		JsonEditor               bool     `json:"jsonEditor"`
		OperationsSorter         string   `json:"operationsSorter"`
		ShowExtensions           bool     `json:"showExtensions"`
		ShowRequestHeaders       bool     `json:"showRequestHeaders"`
		SupportedSubmitMethods   []string `json:"supportedSubmitMethods"`
		TagsSorter               string   `json:"tagsSorter"`
		ValidatorUrl             string   `json:"validatorUrl"`
	}{
		ApisSorter:               "alpha",
		DeepLinking:              true,
		DefaultModelExpandDepth:  1,
		DefaultModelRendering:    "example",
		DefaultModelsExpandDepth: 1,
		DisplayOperationId:       false,
		DisplayRequestDuration:   false,
		DocExpansion:             "none",
		Filter:                   false,
		JsonEditor:               false,
		OperationsSorter:         "alpha",
		ShowExtensions:           false,
		ShowRequestHeaders:       false,
		SupportedSubmitMethods:   []string{"get", "post", "put", "delete", "patch", "head", "options", "trace"},
		TagsSorter:               "alpha",
		ValidatorUrl:             "",
	}, nil
}

func (p SwaggerProvider) GetSsoSecurity(req *restful.Request) (body interface{}, err error) {
	sso := p.cfg.Security.Sso
	return struct {
		AuthorizeUrl string `json:"authorizeUrl"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		Enabled      bool   `json:"enabled"`
		TokenUrl     string `json:"tokenUrl"`
	}{
		AuthorizeUrl: sso.BaseUrl + sso.AuthorizePath,
		ClientId:     sso.ClientId,
		ClientSecret: sso.ClientSecret,
		Enabled:      true,
		TokenUrl:     sso.BaseUrl + sso.TokenPath,
	}, nil
}

func (p SwaggerProvider) Spec(req *restful.Request) (body interface{}, err error) {
	return p.spec, nil
}

func (p SwaggerProvider) YamlSpec(_ *restful.Request, response *restful.Response) {
	specYamlBytes, err := p.RenderYamlSpec()
	if err != nil {
		_ = response.WriteError(500, errors.Wrap(err, "Failed to serialize spec to YAML"))
		return
	}

	response.AddHeader("Content-Type", webservice.MIME_YAML_CHARSET)
	response.WriteHeader(200)
	_, _ = response.Write(specYamlBytes)
}

func (p SwaggerProvider) RenderYamlSpec() ([]byte, error) {
	specJsonBytes, err := json.Marshal(p.spec)
	if err != nil {
		return nil, err
	}

	var specYaml = yaml2.MapSlice{}
	err = yaml2.Unmarshal(specJsonBytes, &specYaml)
	if err != nil {
		return nil, err
	}

	specYamlBytes, err := yaml2.Marshal(specYaml)
	if err != nil {
		return nil, err
	}

	return specYamlBytes, nil
}

func (p *SwaggerProvider) Actuate(container *restful.Container, swaggerService *restful.WebService) error {
	contextPath := swaggerService.RootPath()
	swaggerService.Path(swaggerService.RootPath() + p.cfg.SwaggerPath)

	wsdoc := swagger.NewWebServicesDocumentor(container, swaggerService, contextPath, p.appInfo)

	ws := swagger.WebServices{
		WebServices: container.RegisteredWebServices(),
		APIPath:     swaggerService.RootPath() + p.cfg.SwaggerPath + p.cfg.ApiPath,
	}

	if err := wsdoc.Document(&ws); err != nil {
		return errors.Wrap(err, "Failed to generate Swagger contract specification")
	}

	p.spec = swagger.Spec()

	swaggerService.Route(swaggerService.GET(p.cfg.ApiPath).
		To(webservice.RawController(p.Spec)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	swaggerService.Route(swaggerService.GET(p.cfg.ApiYamlPath).
		To(p.YamlSpec).
		Produces(webservice.MIME_YAML).
		Do(webservice.Returns(200, 401)))

	swaggerService.Route(swaggerService.GET("/configuration/security").
		Operation("swagger.configuration.security").
		To(webservice.RawController(p.GetSecurity)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	swaggerService.Route(swaggerService.GET("/configuration/ui").
		Operation("swagger.configuration.ui").
		To(webservice.RawController(p.GetUi)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	swaggerService.Route(swaggerService.GET("").
		Operation("swagger.configuration.swagger-resources").
		To(webservice.RawController(p.GetSwaggerResources)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	swaggerService.Route(swaggerService.GET("/configuration/security/sso").
		Operation("swagger.configuration.security.sso").
		To(webservice.RawController(p.GetSsoSecurity)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	if p.cfg.Ui.Enabled {
		webServer := webservice.WebServerFromContext(p.ctx)
		webServer.RegisterAlias(p.cfg.Ui.StaticView+"/{subPath:*}", p.cfg.Ui.StaticFiles)
		for _, rootFile := range p.cfg.Ui.RootFiles {
			webServer.RegisterAlias(rootFile, path.Join(p.cfg.Ui.StaticFiles, rootFile))
		}
		webServer.RegisterAlias(p.cfg.Ui.Endpoint, path.Join(p.cfg.Ui.StaticFiles, p.cfg.Ui.View))

		logger.Infof("Serving Swagger %s on %s://%s:%d%s%s",
			p.cfg.Version,
			p.cfg.Server.Scheme(),
			p.cfg.Server.Host,
			p.cfg.Server.Port,
			contextPath,
			p.cfg.Ui.Endpoint)
	}

	return nil
}

func NewSwaggerProvider(ctx context.Context, cfg *DocumentationConfig, appInfo *schema.AppInfo) *SwaggerProvider {
	return &SwaggerProvider{
		ctx:     ctx,
		cfg:     cfg,
		appInfo: appInfo,
	}
}
