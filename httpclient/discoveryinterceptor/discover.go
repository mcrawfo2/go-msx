package discoveryinterceptor

import (
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var logger = log.NewLogger("msx.httpclient.discoveryinterceptor")

func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	if !discovery.IsDiscoveryProviderRegistered() {
		logger.Info("Discovery provider not registered.  Skipping discovery interceptor.")
		return fn
	}

	return func(req *http.Request) (response *http.Response, e error) {
		url := req.URL
		serviceName := url.Host

		if !strings.Contains(serviceName, ":") && !strings.Contains(serviceName, ".") {
			if instances, err := discovery.Discover(req.Context(), serviceName, true); err != nil {
				return nil, errors.Wrap(err, "Failed to discover service "+serviceName)
			} else if len(instances) == 0 {
				return nil, errors.New(fmt.Sprintf("No healthy instances of %s found", serviceName))
			} else {
				serviceInstance := instances.SelectRandom()
				url.Host = serviceInstance.Address()
				serviceContextPath := serviceInstance.ContextPath()
				if serviceContextPath != "" {
					// Normalize leading slash
					if !strings.HasPrefix(url.Path, "/") {
						url.Path = "/" + url.Path
						url.RawPath = "/" + url.RawPath
					}

					// Strip duplicate context path
					if strings.HasPrefix(url.Path, serviceContextPath) {
						url.Path = strings.TrimPrefix(url.Path, serviceContextPath)
						url.RawPath = strings.TrimPrefix(url.RawPath, serviceContextPath)
					}

					url.Path = serviceContextPath + url.Path
					url.RawPath = serviceContextPath + url.RawPath
				}
			}
		}

		return fn(req)
	}
}
