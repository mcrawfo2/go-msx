package targetinterceptor

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"net/http"
	"strings"
	"sync"
)

type RemoteServiceConfig struct {
	ServiceName string `config:"key=service,default="`
}

type RemoteServicesConfig struct {
	RemoteService map[string]RemoteServiceConfig
}

func (s RemoteServicesConfig) MapServiceName(serviceName string) string {
	if rsc, ok := s.RemoteService[serviceName]; ok {
		if rsc.ServiceName != "" {
			return rsc.ServiceName
		}
	}
	return serviceName
}

func NewRemoteServicesConfig(ctx context.Context) (RemoteServicesConfig, error) {
	cfg := config.FromContext(ctx)
	var remoteServices RemoteServicesConfig
	if err := cfg.Populate(&remoteServices, ""); err != nil {
		return RemoteServicesConfig{}, err
	}
	return remoteServices, nil
}

// NewInterceptor returns an HTTP transport interceptor to map
// the default service name for an endpoint to the configured
// service name.
func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	var rsc RemoteServicesConfig
	var once sync.Once

	var initialize = func(req *http.Request) {
		var err error
		ctx := req.Context()
		rsc, err = NewRemoteServicesConfig(ctx)
		if err != nil {
			rsc = RemoteServicesConfig{}
		}
	}

	return func(req *http.Request) (*http.Response, error) {
		once.Do(func() {
			initialize(req)
		})

		url := req.URL
		serviceName := url.Host

		if !strings.Contains(serviceName, ":") && !strings.Contains(serviceName, ".") {
			serviceName = rsc.MapServiceName(serviceName)
			url.Host = serviceName
		}

		return fn(req)
	}
}
