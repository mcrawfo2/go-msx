package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/certificate"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"strconv"
)

type WebServerConfig struct {
	Enabled       bool   `config:"default=false"`
	Host          string `config:"default=${network.outbound.address:0.0.0.0}"`
	Port          int    `config:"default=8080"`
	Tls           certificate.TLSConfig
	Cors          CorsConfig
	ContextPath   string `config:"default=/app"`
	StaticPath    string `config:"default=/www"`
	StaticEnabled bool   `config:"default=true"`
	TraceEnabled  bool   `config:"default=false"`
	DebugEnabled  bool   `config:"default=false"`
}

type CorsConfig struct {
	Enabled              bool     `config:"default=true"`
	CustomAllowedHeaders []string `config:"default=${security.cors.allowedHeaders}"`
	CustomExposedHeaders []string `config:"default=${security.cors.exposedHeaders}"`
}

func (c WebServerConfig) Address() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (c WebServerConfig) GenericAddress() string {
	return "0.0.0.0:" + strconv.Itoa(c.Port)
}

func (c WebServerConfig) Url() string {
	if c.Tls.Enabled {
		return "https://" + c.Address() + c.ContextPath
	}
	return "http://" + c.Address() + c.ContextPath
}

func NewWebServerConfig(cfg *config.Config) (*WebServerConfig, error) {
	var webServerConfig WebServerConfig
	if err := cfg.Populate(&webServerConfig, configRootWebServer); err != nil {
		return nil, err
	}
	return &webServerConfig, nil
}
