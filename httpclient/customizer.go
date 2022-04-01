// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package httpclient

import (
	"net/http"
	"time"
)

type ClientConfigurationFunc func(c *http.Client)
type TransportConfigurationFunc func(c *http.Transport)

type Configurer interface {
	HttpClient(*http.Client)
	HttpTransport(*http.Transport)
}

type ClientConfigurer struct {
	ClientFuncs    []ClientConfigurationFunc
	TransportFuncs []TransportConfigurationFunc
}

func (c ClientConfigurer) HttpClient(client *http.Client) {
	for _, configure := range c.ClientFuncs {
		configure(client)
	}
}

func (c ClientConfigurer) HttpTransport(transport *http.Transport) {
	for _, configure := range c.TransportFuncs {
		configure(transport)
	}
}

type CompositeConfigurer struct {
	Service  Configurer
	Endpoint Configurer
}

func (c CompositeConfigurer) HttpClient(client *http.Client) {
	if c.Service != nil {
		c.Service.HttpClient(client)
	}
	if c.Endpoint != nil {
		c.Endpoint.HttpClient(client)
	}
}

func (c CompositeConfigurer) HttpTransport(transport *http.Transport) {
	if c.Service != nil {
		c.Service.HttpTransport(transport)
	}
	if c.Endpoint != nil {
		c.Endpoint.HttpTransport(transport)
	}
}

// Common customizations

func TlsInsecure(insecure bool) TransportConfigurationFunc {
	return func(c *http.Transport) {
		c.TLSClientConfig.InsecureSkipVerify = insecure
	}
}

func NoProxy() TransportConfigurationFunc {
	return func(c *http.Transport) {
		c.Proxy = nil
	}
}

func ClientTimeout(timeout time.Duration) ClientConfigurationFunc {
	return func(c *http.Client) {
		c.Timeout = timeout
	}
}

func ClientCookieJar(jar http.CookieJar) ClientConfigurationFunc {
	return func(c *http.Client) {
		c.Jar = jar
	}
}
