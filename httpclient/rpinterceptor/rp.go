// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rpinterceptor

import (
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var logger = log.NewLogger("httpclient.rpinterceptor")

func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	if !discovery.IsDiscoveryProviderRegistered() {
		logger.Info("Discovery provider not registered.  Skipping rp interceptor.")
		return fn
	}

	return func(req *http.Request) (response *http.Response, e error) {
		url := req.URL
		serviceName := url.Host
		tag := "resourceprovider:" + serviceName

		if !strings.Contains(serviceName, ":") && !strings.Contains(serviceName, ".") {
			if instances, err := discovery.DiscoverAll(req.Context(), true, tag); err != nil {
				return nil, errors.Wrap(err, "Failed to discover resource provider "+serviceName)
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

func ApplyInterceptor() httpclient.ClientConfigurationFunc {
	return func(c *http.Client) {
		httpclient.ApplyInterceptor(c, NewInterceptor)
	}
}
