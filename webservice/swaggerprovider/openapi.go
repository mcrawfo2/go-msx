package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/swaggest/openapi-go/openapi3"
	yaml2 "gopkg.in/yaml.v2"
)

type OpenApiProvider struct {
	ctx        context.Context
	cfg        *DocumentationConfig
	spec       *openapi3.Spec
	appInfo    *AppInfo
	customizer OpenApiCustomizer
	reflector  *openapi3.Reflector
}

func (p *OpenApiProvider) BuildSpec(container *restful.Container, contextPath string) (err error) {
	p.spec.Openapi = p.cfg.Version
	return nil
}

func (p *OpenApiProvider) PostBuildSpec(container *restful.Container) error {
	return types.ErrorList{
		p.customizer.PostBuildInfo(p.spec, p.appInfo),
		p.customizer.PostBuildServers(p.spec, p.cfg.Server),
		p.customizer.PostBuildTags(p.spec),
	}.Filter()
}

func (p *OpenApiProvider) GetSecurity(req *restful.Request) (body interface{}, err error) {
	return struct{}{}, nil
}

func (p *OpenApiProvider) GetSwaggerResources(req *restful.Request) (body interface{}, err error) {
	return []struct {
		Name           string `json:"name"`
		Location       string `json:"location"`
		Url            string `json:"url"`
		SwaggerVersion string `json:"swaggerVersion"`
	}{
		{
			Name:           "platform",
			Location:       p.cfg.SwaggerPath + p.cfg.ApiPath,
			Url:            p.cfg.SwaggerPath + p.cfg.ApiPath,
			SwaggerVersion: p.cfg.Version,
		},
	}, nil
}

func (p *OpenApiProvider) GetUi(req *restful.Request) (body interface{}, err error) {
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

func (p *OpenApiProvider) GetSsoSecurity(req *restful.Request) (body interface{}, err error) {
	sso := p.cfg.Security.Sso
	return struct {
		AuthorizeUrl string `json:"authorizeUrl"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		TokenUrl     string `json:"tokenUrl"`
	}{
		AuthorizeUrl: sso.BaseUrl + sso.AuthorizePath,
		ClientId:     sso.ClientId,
		ClientSecret: sso.ClientSecret,
		TokenUrl:     sso.BaseUrl + sso.TokenPath,
	}, nil
}

func (p *OpenApiProvider) Spec(_ *restful.Request) (body interface{}, err error) {
	return p.spec, nil
}

func (p *OpenApiProvider) YamlSpec(_ *restful.Request, response *restful.Response) {
	specYamlBytes, err := p.RenderYamlSpec()
	if err != nil {
		_ = response.WriteError(500, errors.Wrap(err, "Failed to serialize spec to YAML"))
		return
	}

	response.AddHeader("Content-Type", webservice.MIME_YAML_CHARSET)
	response.WriteHeader(200)
	_, _ = response.Write(specYamlBytes)
}

func (p *OpenApiProvider) RenderYamlSpec() ([]byte, error) {
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

func (p *OpenApiProvider) Actuate(container *restful.Container, webService *restful.WebService) error {
	contextPath := webService.RootPath()
	servicePath := contextPath + p.cfg.SwaggerPath
	//apiPath := servicePath + p.cfg.ApiPath

	webService.Path(servicePath)

	// Build spec
	if err := p.BuildSpec(container, contextPath); err != nil {
		return err
	}
	if err := p.PostBuildSpec(container); err != nil {
		return err
	}

	webService.Route(webService.GET("/configuration/security").
		Operation("swagger.configuration.security").
		To(webservice.RawController(p.GetSecurity)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	webService.Route(webService.GET("/configuration/ui").
		Operation("swagger.configuration.ui").
		To(webservice.RawController(p.GetUi)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	webService.Route(webService.GET("").
		Operation("swagger.configuration.swagger-resources").
		To(webservice.RawController(p.GetSwaggerResources)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	webService.Route(webService.GET("/configuration/security/sso").
		Operation("swagger.configuration.security.sso").
		To(webservice.RawController(p.GetSsoSecurity)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	webService.Route(webService.GET(p.cfg.ApiPath).
		To(webservice.RawController(p.Spec)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	webService.Route(webService.GET(p.cfg.ApiYamlPath).
		To(p.YamlSpec).
		Produces(webservice.MIME_YAML).
		Do(webservice.Returns(200, 401)))

	if p.cfg.Ui.Enabled {
		webServer := webservice.WebServerFromContext(p.ctx)
		webServer.RegisterAlias(p.cfg.Ui.Endpoint, p.cfg.Ui.View)

		logger.Infof("Serving OpenApi %s on http://%s:%d%s%s",
			p.cfg.Version,
			p.cfg.Server.Host,
			p.cfg.Server.Port,
			contextPath,
			p.cfg.Ui.Endpoint)

	}

	return nil
}

type OpenApiCustomizer struct {
	cfg *DocumentationConfig
}

func (c OpenApiCustomizer) PostBuildInfo(spec *openapi3.Spec, appInfo *AppInfo) error {
	spec.Info = openapi3.Info{
		Title: "MSX API Documentation for " + appInfo.Name,
		Description: types.NewStringPtr("<h3>This is the REST API documentation for " + appInfo.Name + "</h3>\n \n" +
			appInfo.Description + "\n" +
			"+ API Authorization \n" +
			"    + Authorization header is <b>required</b>. \n" +
			"    + It should be in Bearer authentication scheme </br>(e.g <b> Authorization: BEARER &lt;access token&gt; </b>)\n"),
		TermsOfService: types.NewStringPtr("http://www.cisco.com"),
		Contact: &openapi3.Contact{
			Name:  types.NewStringPtr("Cisco Systems Inc."),
			URL:   types.NewStringPtr("http://www.cisco.com"),
			Email: types.NewStringPtr("somecontact@cisco.com"),
		},
		License: &openapi3.License{
			Name: "Apache License Version 2.0",
			URL:  types.NewStringPtr("http://www.apache.org/licenses/LICENSE-2.0.html"),
		},
		Version: appInfo.Version,
	}
	return nil
}

func (c OpenApiCustomizer) PostBuildServers(spec *openapi3.Spec, cfg DocumentationServerConfig) error {
	spec.Servers = []openapi3.Server{
		{
			URL: fmt.Sprintf("http://%s:%d%s",
				cfg.Host,
				cfg.Port,
				cfg.ContextPath),
		},
	}
	return nil
}

func (c OpenApiCustomizer) PostBuildTags(spec *openapi3.Spec) error {
	return nil
}
