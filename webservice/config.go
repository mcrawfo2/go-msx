package webservice

type WebServerConfig struct {
	Enabled     bool   `config:"default=false"`
	Host        string `config:"default=0.0.0.0"`
	Port        int    `config:"default=8080"`
	Tls         bool   `config:"default=false"`
	CertFile    string `config:"default=server.crt"`
	KeyFile     string `config:"default=server.key"`
	Cors        bool   `config:"default=true"`
	ContextPath string `config:"default=/app"`
	Swagger     SwaggerConfig
	JWT         UserContextFilterConfig
}

type SwaggerConfig struct {
	Enabled        bool   `config:"default=false"`
	Path           string `config:"default=thirdparty/swagger-ui/dist"`
	WebServicesUrl string
	ApiPath        string
	SwaggerPath    string
}
