// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package httpclient

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
)

type Factory interface {
	NewHttpClient() *http.Client
}

type ContextFactory interface {
	Factory
	NewHttpClientWithConfigurer(context.Context, Configurer) *http.Client
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type DoFunc func(req *http.Request) (*http.Response, error)

func (d DoFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return d(request)
}

type RequestInterceptor func(fn DoFunc) DoFunc

func New(ctx context.Context, configurer Configurer) (*http.Client, error) {
	factory := FactoryFromContext(ctx)
	if factory == nil {
		return nil, errors.New("Failed to retrieve http client factory from context")
	}

	var httpClient *http.Client
	if contextFactory, ok := factory.(ContextFactory); ok {
		httpClient = contextFactory.NewHttpClientWithConfigurer(ctx, configurer)
	} else {
		httpClient = factory.NewHttpClient()
	}

	return httpClient, nil
}
