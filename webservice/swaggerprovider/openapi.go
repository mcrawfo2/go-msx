package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/fatih/structtag"
	"github.com/pkg/errors"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
	"reflect"
	"strings"
)

type OpenApiProvider struct {
	ctx        context.Context
	cfg        *DocumentationConfig
	spec       *openapi3.Spec
	appInfo    *AppInfo
	customizer OpenApiCustomizer
	requests   map[reflect.Type]interface{}
	responses  map[reflect.Type]interface{}
}

func (p *OpenApiProvider) initReflector(reflector *openapi3.Reflector) {
	// Custom mapping for types.UUID
	customMappingTypesUuid := jsonschema.Schema{}
	customMappingTypesUuid.AddType(jsonschema.String)
	customMappingTypesUuid.WithFormat("uuid")
	reflector.AddTypeMapping(types.UUID{}, customMappingTypesUuid)

	// Custom mapping for types.UUID
	customMappingTypesTime := jsonschema.Schema{}
	customMappingTypesTime.AddType(jsonschema.String)
	customMappingTypesTime.WithFormat("date-time")
	reflector.AddTypeMapping(types.Time{}, customMappingTypesUuid)
}

func (p *OpenApiProvider) getReflectorRequest(reflector *openapi3.Reflector, request interface{}) (interface{}, error) {
	if p.requests == nil {
		p.requests = make(map[reflect.Type]interface{})
	}

	requestType := reflect.TypeOf(request)

	if reflectorRequest, ok := p.requests[requestType]; ok {
		return reflectorRequest, nil
	}

	reflectorStruct, err := p.createReflectorStruct(reflector, request, "req")
	if err != nil {
		return nil, err
	}

	var reflectorRequest = reflect.Zero(reflectorStruct).Interface()
	p.requests[requestType] = reflectorRequest
	return reflectorRequest, nil
}

func (p *OpenApiProvider) getEmbeddedResponse(payload interface{}) interface{} {
	type embedded struct {
		Body interface{} `resp:"body"`
	}

	responseType := types.NewParameterizedStruct(
		reflect.TypeOf(embedded{}),
		"Body",
		payload)

	return reflect.Zero(responseType).Interface()
}

func (p *OpenApiProvider) getReflectorResponse(reflector *openapi3.Reflector, response interface{}) (interface{}, error) {
	if p.responses == nil {
		p.responses = make(map[reflect.Type]interface{})
	}

	responseType := reflect.TypeOf(response)

	if reflectorResponse, ok := p.responses[responseType]; ok {
		return reflectorResponse, nil
	}

	reflectorStruct, err := p.createReflectorStruct(reflector, response, "resp")
	if err != nil {
		return nil, err
	}

	var reflectorResponse = reflect.Zero(reflectorStruct).Interface()
	p.responses[responseType] = reflectorResponse
	return reflectorResponse, nil
}

func (p *OpenApiProvider) createReflectorStruct(reflector *openapi3.Reflector, source interface{}, tagName string) (reflect.Type, error) {
	sourceType := reflect.TypeOf(source)

	// Create a new struct for the supplied request
	var reflectorFields []reflect.StructField
	for i := 0; i < sourceType.NumField(); i++ {
		requestField := sourceType.Field(i)
		parameterTag := webservice.NewParameterTag(requestField, tagName)
		tags, err := structtag.Parse(string(requestField.Tag))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse struct tag")
		}

		switch parameterTag.Source {
		case "path", "query", "header", "form":
			// just transform req:"x=name" to x:"name"
			tags.Delete(tagName)

			err = tags.Set(&structtag.Tag{
				Key:     parameterTag.Source,
				Name:    parameterTag.Name,
			})
			if err != nil {
				return nil, errors.Wrap(err, "Failed to apply struct tag")
			}

			requestField.Tag = reflect.StructTag(tags.String())
			reflectorFields = append(reflectorFields, requestField)

		case "body":
			// body -> embed struct members or entity
			requestFieldType := requestField.Type
			for requestFieldType.Kind() == reflect.Ptr {
				requestFieldType = requestFieldType.Elem()
			}

			if envelopeTag, err := tags.Get("envelope"); err == nil && envelopeTag.Value() == "true" {
				requestFieldType = types.NewParameterizedStruct(
					reflect.TypeOf(integration.MsxEnvelope{}),
					"Payload",
					types.Instantiate(requestFieldType))
				requestField.Type = requestFieldType
			}

			if requestFieldType.Kind() == reflect.Struct {
				instance := types.Instantiate(requestField.Type)
				nullOp := openapi3.Operation{}
				_ = reflector.SetJSONResponse(&nullOp, instance, -1)

			}

			requestField.Anonymous = true
			requestField.Tag = ``
			reflectorFields = append(reflectorFields, requestField)
		}
	}

	return reflect.StructOf(reflectorFields), nil
}

func (p *OpenApiProvider) BuildOperationForServiceRoute(reflector *openapi3.Reflector, webService *restful.WebService, contextPath string, route restful.Route) error {
	tags, _ := route.Metadata[restfulspec.KeyOpenAPITags].([]string)

	op := openapi3.Operation{
		Summary:     stringptr(route.Doc),
		Description: stringptr(route.Notes),
		ID:          stringptr(route.Operation),
		Deprecated:  boolptr(route.Deprecated),
		Tags:        tags,
	}

	if request, ok := route.Metadata[webservice.MetadataRequest]; ok {
		reflectorRequest, err := p.getReflectorRequest(reflector, request)
		if err != nil {
			return err
		}

		err = reflector.SetRequest(&op, reflectorRequest, route.Method)
		if err != nil {
			return err
		}
	}


	if response, ok := route.Metadata[webservice.MetadataSuccessResponse]; ok {
		reflectorResponse, err := p.getReflectorResponse(reflector, response)
		if err != nil {
			return err
		}

		defaultResponseCode, ok := route.Metadata[webservice.MetadataDefaultReturnCode].(int)
		if !ok {
			defaultResponseCode = route.DefaultResponse.Code
		}

		err = reflector.SetJSONResponse(&op, reflectorResponse, defaultResponseCode)
		if err != nil {
			return err
		}
	}

	if errorPayload, ok := route.Metadata[webservice.MetadataErrorPayload]; ok {
		reflectorResponse, err := p.getReflectorResponse(reflector, p.getEmbeddedResponse(errorPayload))
		err = reflector.SetJSONResponse(&op, reflectorResponse, -1)
		if err != nil {
			return err
		}

		responseDefinition := op.Responses.MapOfResponseOrRefValues["-1"]
		delete(op.Responses.MapOfResponseOrRefValues, "-1")
		op.Responses.Default = &responseDefinition
		op.Responses.Default.Response.Description = "Error"
	} else {
		// TODO: Auto-calculate error response payload
	}

	routePath := strings.TrimPrefix(route.Path, contextPath)

	return reflector.Spec.AddOperation(route.Method, routePath, op)
}

func (p *OpenApiProvider) BuildOperationsForServiceRoutes(reflector *openapi3.Reflector, webService *restful.WebService, contextPath string) error {
	for _, route := range webService.Routes() {
		if err := p.BuildOperationForServiceRoute(reflector, webService, contextPath, route); err != nil {
			return err
		}
	}

	return nil
}

func (p *OpenApiProvider) BuildSpec(container *restful.Container, contextPath string) error {
	reflector := new(openapi3.Reflector)
	p.initReflector(reflector)
	reflector.Spec = &openapi3.Spec{
		Openapi: p.cfg.Version,
	}

	for _, each := range container.RegisteredWebServices() {
		if err := p.BuildOperationsForServiceRoutes(reflector, each, contextPath); err != nil {
			return err
		}
	}

	p.spec = reflector.Spec
	return nil
}
func (p *OpenApiProvider) PostBuildSpec(container *restful.Container, contextPath string) error {
	return types.ErrorList{
		p.customizer.PostBuildInfo(p.spec, p.appInfo),
		p.customizer.PostBuildServers(p.spec, contextPath),
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

func (p *OpenApiProvider) GetSpec(_ *restful.Request) (body interface{}, err error) {
	return p.spec, nil
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
	if err := p.PostBuildSpec(container, contextPath); err != nil {
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
		To(webservice.RawController(p.GetSpec)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	if p.cfg.Ui.Enabled {
		webServer := webservice.WebServerFromContext(p.ctx)
		webServer.RegisterAlias(p.cfg.Ui.Endpoint, p.cfg.Ui.View)
	}

	return nil
}

func boolptr(val bool) *bool {
	return &val
}

func stringptr(s string) *string {
	return types.NewOptionalStringFromString(s).NilIfEmpty().Ptr()
}

type OpenApiCustomizer struct{}

func (c OpenApiCustomizer) PostBuildInfo(spec *openapi3.Spec, appInfo *AppInfo) error {
	spec.Info = openapi3.Info{
		Title: "MSX API Documentation for " + appInfo.Name,
		Description: stringptr("<h3>This is the REST API documentation for " + appInfo.Name + "</h3>\n \n" +
			appInfo.Description + "\n" +
			"+ API Authorization \n" +
			"    + Authorization header is <b>required</b>. \n" +
			"    + It should be in Bearer authentication scheme </br>(e.g <b> Authorization: BEARER &lt;access token&gt; </b>)\n"),
		TermsOfService: stringptr("http://www.cisco.com"),
		Contact: &openapi3.Contact{
			Name:  stringptr("Cisco Systems Inc."),
			URL:   stringptr("http://www.cisco.com"),
			Email: stringptr("somecontact@cisco.com"),
		},
		License: &openapi3.License{
			Name: "Apache License Version 2.0",
			URL:  stringptr("http://www.apache.org/licenses/LICENSE-2.0.html"),
		},
		Version: appInfo.Version,
	}
	return nil
}

func (c OpenApiCustomizer) PostBuildServers(spec *openapi3.Spec, contextPath string) error {
	spec.Servers = []openapi3.Server{
		{
			URL: "http://localhost:9482" + contextPath,
		},
	}
	return nil
}

func (c OpenApiCustomizer) PostBuildTags(spec *openapi3.Spec) error {
	return nil
}
