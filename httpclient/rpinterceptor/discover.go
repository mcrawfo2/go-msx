package rpinterceptor

import (
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	return func(req *http.Request) (response *http.Response, e error) {
		url := req.URL
		serviceName := url.Host
		tag := "resourceprovider:" + serviceName

		if !strings.Contains(serviceName, ":") && !strings.Contains(serviceName, ".") {
			if instances, err := discovery.DiscoverAll(req.Context(), true, tag); err != nil {
				return nil, errors.Wrap(err, "Failed to discover resource provider " + serviceName)
			} else if len(instances) == 0 {
				return nil, errors.New(fmt.Sprintf("No healthy instances of %s found", serviceName))
			} else {
				serviceInstance := instances.SelectRandom()
				url.Host = serviceInstance.Address()
				serviceContextPath := serviceInstance.ContextPath()
				if serviceContextPath != "" {
					url.Path = serviceContextPath + url.Path
					url.RawPath = serviceContextPath + url.RawPath
				}
			}
		}

		return fn(req)
	}
}
