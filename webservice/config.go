package webservice

import "strconv"

type WebServerConfig struct {
	Enabled     bool   `config:"default=false"`
	Host        string `config:"default=0.0.0.0"`
	Port        int    `config:"default=8080"`
	Tls         bool   `config:"default=false"`
	CertFile    string `config:"default=server.crt"`
	KeyFile     string `config:"default=server.key"`
	Cors        bool   `config:"default=true"`
	ContextPath string `config:"default=/app"`
	JWT         UserContextFilterConfig
}

func (c WebServerConfig) Address() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}