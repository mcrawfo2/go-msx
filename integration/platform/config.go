package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/discoveryinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/loginterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/traceinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/security/httprequest"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

func tokenInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	return func(req *http.Request) (*http.Response, error) {
		httprequest.InjectToken(req)
		return fn(req)
	}
}

func newPlatformClientConfigFromContext(ctx context.Context, serviceName integration.ServiceName) *platform.Configuration {
	httpClient := httpclient.FactoryFromContext(ctx).NewHttpClient()
	remoteServiceConfig := integration.NewRemoteServiceConfig(ctx, serviceName)
	var remoteServiceName = remoteServiceConfig.ServiceName

	transport := httpClient.Transport.RoundTrip
	transport = tokenInterceptor(transport)
	transport = traceinterceptor.NewInterceptor(transport)
	transport = discoveryinterceptor.NewInterceptor(transport)
	transport = loginterceptor.NewInterceptor(transport)

	httpClient.Transport = httpclient.DoFunc(transport)

	return &platform.Configuration{
		Host:       remoteServiceName,
		Scheme:     "http",
		HTTPClient: httpClient,
	}
}

func newSecurityClientConfigFromContext(ctx context.Context, serviceName integration.ServiceName) *platform.Configuration {
	httpClient := httpclient.FactoryFromContext(ctx).NewHttpClient()
	remoteServiceConfig := integration.NewRemoteServiceConfig(ctx, serviceName)
	var remoteServiceName = remoteServiceConfig.ServiceName
	transport := httpClient.Transport.RoundTrip
	transport = traceinterceptor.NewInterceptor(transport)
	transport = discoveryinterceptor.NewInterceptor(transport)
	transport = loginterceptor.NewInterceptor(transport)

	httpClient.Transport = httpclient.DoFunc(transport)

	return &platform.Configuration{
		Host:       remoteServiceName,
		Scheme:     "http",
		HTTPClient: httpClient,
	}
}
