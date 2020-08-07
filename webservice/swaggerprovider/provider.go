package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

var (
	ErrDisabled = errors.New("Swagger disabled")
	logger      = log.NewLogger("msx.webservice.swaggerprovider")
)

type SwaggerProvider struct {
	ctx  context.Context
	cfg  *DocumentationConfig
	spec *spec.Swagger
	info spec.Info
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
		ClientId:     sso.ClientId,
		ClientSecret: sso.ClientSecret,
		TokenUrl:     sso.BaseUrl + sso.TokenPath,
	}, nil
}

func (p SwaggerProvider) GetSpec(req *restful.Request) (body interface{}, err error) {
	return p.spec, nil
}

func (p SwaggerProvider) PostBuildSpec(container *restful.Container, svc *restful.WebService, contextPath string) func(spec *spec.Swagger) {
	return func(swagger *spec.Swagger) {
		// Inject info blob
		swagger.Info = &p.info

		// Register tags definitions from all of the routes
		var existingTags = types.StringStack{}
		for _, svc := range container.RegisteredWebServices() {
			for _, route := range svc.Routes() {
				if routeTagDefinitionInterface, ok := route.Metadata[webservice.MetadataTagDefinition]; ok {
					routeTagDefinition := routeTagDefinitionInterface.(spec.TagProps)
					if !existingTags.Contains(routeTagDefinition.Name) {
						existingTags = append(existingTags, routeTagDefinition.Name)
						swagger.Tags = append(swagger.Tags, spec.Tag{TagProps: routeTagDefinition})
					}
				}
			}
		}

		// Factor out contextPath into basePath
		if contextPath != "/" {
			newPaths := make(map[string]spec.PathItem)
			for path, pathItem := range swagger.Paths.Paths {
				if strings.HasPrefix(path, contextPath) {
					path = strings.TrimPrefix(path, contextPath)
				}
				newPaths[path] = pathItem
			}
			swagger.Paths.Paths = newPaths
			swagger.BasePath = contextPath
		}

		// Sort tags
		sort.Slice(swagger.Tags, func(i, j int) bool {
			iTagName := swagger.Tags[i].Name
			jTagName := swagger.Tags[j].Name
			return strings.Compare(iTagName, jTagName) < 0
		})
	}
}

func (p SwaggerProvider) Actuate(container *restful.Container, swaggerService *restful.WebService) error {
	contextPath := swaggerService.RootPath()
	swaggerService.Path(swaggerService.RootPath() + p.cfg.SwaggerPath)

	p.spec = restfulspec.BuildSwagger(restfulspec.Config{
		WebServices:                   container.RegisteredWebServices(),
		APIPath:                       swaggerService.RootPath() + p.cfg.SwaggerPath + p.cfg.ApiPath,
		ModelTypeNameHandler:          webservice.ResponseTypeName,
		PostBuildSwaggerObjectHandler: p.PostBuildSpec(container, swaggerService, contextPath),
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

	appInfo, err := AppInfoFromConfig(config.MustFromContext(ctx))
	if err != nil {
		return err
	}

	server.AddDocumentationProvider(&SwaggerProvider{
		ctx: ctx,
		cfg: cfg,
		info: spec.Info{
			InfoProps: spec.InfoProps{
				Title: "MSX API Documentation for " + appInfo.Name,
				Description: "<h3>This is the REST API documentation for " + appInfo.Name + "</h3>\n \n" +
					appInfo.Description + "\n" +
					"+ API Authorization \n" +
					"    + Authorization header is <b>required</b>. \n" +
					"    + It should be in Bearer authentication scheme </br>(e.g <b> Authorization: BEARER &lt;access token&gt; </b>)\n",
				TermsOfService: "http://www.cisco.com",
				Contact: &spec.ContactInfo{
					Name:  "Cisco Systems Inc.",
					URL:   "http://www.cisco.com",
					Email: "somecontact@cisco.com",
				},
				License: &spec.License{
					Name: "Apache License Version 2.0",
					URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
				},
				Version: appInfo.Version,
			},
		},
	})
	return nil
}
