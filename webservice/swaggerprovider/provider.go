package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-swagger12"
	"github.com/pkg/errors"
)

const (
	clientId = "nfv-service"
	clientSecret = "nfv-service-secret"
)

var (
	ErrDisabled = errors.New("Swagger disabled")
)

type SwaggerProvider struct {
	cfg *DocumentationConfig
}

func (p SwaggerProvider) GetSecurity(req *restful.Request) (body interface{}, err error) {
	return struct {
		ApiKeyVehicle  string `json:"apiKeyVehicle"`
		ScopeSeparator string `json:"scopeSeparator"`
		ApiKeyName     string `json:"apiKeyName"`
	}{
		ApiKeyVehicle:  "header",
		ScopeSeparator: ",",
		ApiKeyName:     "api_key",
	}, nil
}

func (p SwaggerProvider) GetSwaggerResources(req *restful.Request) (body interface{}, err error) {
	return []struct {
		Name           string `json:"name"`
		Location       string `json:"location"`
		SwaggerVersion string `json:"swaggerVersion"`
	}{
		{
			Name:           "alertservice",
			Location:       "/apidocs.json",
			SwaggerVersion: "2.0",
		},
	}, nil
}

func (p SwaggerProvider) GetUi(req *restful.Request) (body interface{}, err error) {
	return struct {
		ValidatorUrl           *string  `json:"validatorUrl"`
		DocExpansion           string   `json:"docExpansion"`
		ApisSorter             string   `json:"apisSorter"`
		DefaultModelRendering  string   `json:"defaultModelRendering"`
		SupportedSubmitMethods []string `json:"supportedSubmitMethods"`
		JsonEditor             bool     `json:"jsonEditor"`
		ShowRequestHeaders     bool     `json:"showRequestHeaders"`
	}{
		ValidatorUrl:          nil,
		DocExpansion:          "none",
		ApisSorter:            "alpha",
		DefaultModelRendering: "schema",
		SupportedSubmitMethods: []string{
			"get", "post", "put", "delete", "patch",
		},
		JsonEditor:         false,
		ShowRequestHeaders: true,
	}, nil
}

func (p SwaggerProvider) GetUserSecurity(req *restful.Request) (body interface{}, err error) {
	return struct {
		Enabled                  bool   `json:"enabled"`
		AuthenticationUrl        string `json:"authenticationUrl"`
		AuthServerBaseUrl        string `json:"authServerBaseUrl"`
		AuthServerLoginEndpoint  string `json:"authServerLoginEndpoint"`
		AuthServerLogoutEndpoint string `json:"authServerLogoutEndpoint"`
		TokenHeader              string `json:"tokenHeader"`
		TokenPrefix              string `json:"tokenPrefix"`
	}{
		Enabled:                  true,
		AuthenticationUrl:        "",
		AuthServerBaseUrl:        "http://usermanagementservice/idm",
		AuthServerLoginEndpoint:  "/api/v1/accesstoken",
		AuthServerLogoutEndpoint: "/api/v1/users/logout",
		TokenHeader:              "Authorization",
		TokenPrefix:              "Bearer ",
	}, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseOauth2 struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

type LoginResponse struct {
	Oauth2 LoginResponseOauth2 `json:"oauth2"`
}

func (p SwaggerProvider) UserLogin(req *restful.Request) (body interface{}, err error) {
	var dto LoginRequest
	err = req.ReadEntity(&dto)
	if err != nil {
		return
	}

	usermanagementIntegration, err := usermanagement.NewIntegration(req.Request.Context())
	if err != nil {
		return nil, err
	}

	msxResponse, err := usermanagementIntegration.Login(clientId, clientSecret, dto.Username, dto.Password)
	if err != nil {
		return nil, err
	}

	payload := msxResponse.Payload.(integration.Pojo)


	var response = LoginResponse{
		Oauth2: LoginResponseOauth2{
			AccessToken: payload["Token"].(string),
			TokenType:   "bearer",
			Scope:       "read,write",
			ExpiresIn:   int(payload["Expires"].(float64)),
		},
	}

	return response, nil
}

func (p SwaggerProvider) Actuate(container *restful.Container, swaggerService *restful.WebService) error {
	swaggerConfig := swagger.Config{
		WebServices:     container.RegisteredWebServices(),
		WebServicesUrl:  p.cfg.WebServicesUrl,
		ApiPath:         p.cfg.ApiPath,
		SwaggerPath:     p.cfg.SwaggerPath,
		SwaggerFilePath: p.cfg.Path,
	}

	swagger.RegisterSwaggerService(swaggerConfig, container)

	swaggerService.Consumes(restful.MIME_JSON)
	swaggerService.Produces(restful.MIME_JSON)
	swaggerService.Path(swaggerService.RootPath() + "/swagger")

	swaggerService.Route(swaggerService.GET("/configuration/security").
		To(webservice.RawController(p.GetSecurity)).
		Do(webservice.Returns(200, 401)))

	swaggerService.Route(swaggerService.GET("/configuration/ui").
		To(webservice.RawController(p.GetUi)).
		Do(webservice.Returns(200, 401)))

	swaggerService.Route(swaggerService.GET("/configuration/swagger-resources").
		To(webservice.RawController(p.GetSwaggerResources)).
		Do(webservice.Returns(200, 401)))

	swaggerService.Route(swaggerService.GET("/configuration/user-security").
		To(webservice.RawController(p.GetUserSecurity)).
		Do(webservice.Returns(200, 401)))

	swaggerService.Route(swaggerService.POST("/user/login").
		To(webservice.RawController(p.UserLogin)).
		Reads(LoginRequest{}).
		Do(webservice.StandardReturns))

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

	server.SetDocumentationProvider(&SwaggerProvider{cfg:cfg})
	return nil
}
