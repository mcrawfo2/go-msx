package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

var (
	ErrDisabled = errors.New("Swagger disabled")
	logger      = log.NewLogger("msx.webservice.swaggerprovider")
)

type SwaggerProvider struct {
	ctx  context.Context
	cfg  *DocumentationConfig
	spec *spec.Swagger
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
			Location:       p.cfg.SwaggerPath + p.cfg.ApiPath,
			Url:            p.cfg.SwaggerPath + p.cfg.ApiPath,
			SwaggerVersion: "2.0",
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
		TokenUrl     string `json:"tokenUrl"`
	}{
		AuthorizeUrl: sso.BaseUrl + sso.AuthorizePath,
		ClientId: sso.ClientId,
		ClientSecret: sso.ClientSecret,
		TokenUrl: sso.BaseUrl + sso.TokenPath,
	}, nil
}

func (p SwaggerProvider) GetSpec(req *restful.Request) (body interface{}, err error) {
	return p.spec, nil
}

func (p SwaggerProvider) Actuate(container *restful.Container, swaggerService *restful.WebService) error {
	swaggerService.Path(swaggerService.RootPath() + p.cfg.SwaggerPath)

	p.spec = restfulspec.BuildSwagger(restfulspec.Config{
		WebServices:    container.RegisteredWebServices(),
		APIPath:        swaggerService.RootPath() + p.cfg.SwaggerPath + p.cfg.ApiPath,
	})

	swaggerService.Route(swaggerService.GET(p.cfg.ApiPath).
		To(webservice.RawController(p.GetSpec)).
		Produces(webservice.MIME_JSON).
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
		webServer.RegisterAlias(p.cfg.Ui.Endpoint, p.cfg.Ui.View)
	}

	return nil
}

func RegisterSwaggerProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server == nil {
		return nil
	}

	cfg, err := DocumentationConfigFromConfig(config.MustFromContext(ctx))
	if err != nil {
		return err
	}

	if !cfg.Enabled {
		return ErrDisabled
	}

	server.SetDocumentationProvider(&SwaggerProvider{
		ctx: ctx,
		cfg: cfg,
	})
	return nil
}
