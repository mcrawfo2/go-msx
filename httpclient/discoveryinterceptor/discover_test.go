// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package discoveryinterceptor

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func setMockDiscoveryProvider() *discovery.MockDiscoveryProvider {
	provider := new(discovery.MockDiscoveryProvider)
	discovery.RegisterDiscoveryProvider(provider)
	return provider
}

func TestNewInterceptor(t *testing.T) {
	tests := []struct {
		name           string
		provider       func(provider *discovery.MockDiscoveryProvider)
		url            string
		wantHost       string
		wantPathPrefix string
		wantErr        bool
	}{
		{
			name: "Success",
			provider: func(provider *discovery.MockDiscoveryProvider) {
				provider.On("Discover",
					mock.AnythingOfType("*context.emptyCtx"),
					"myservice",
					true).
					Return(discovery.ServiceInstances{
						{
							ID:   "myservice-XYZ",
							Name: "myservice",
							Host: "10.10.10.10",
							Tags: []string{
								"contextPath=/my",
							},
							Port: 8080,
						},
					}, nil)
			},
			url:            "http://myservice",
			wantHost:       "10.10.10.10:8080",
			wantPathPrefix: "/my",
			wantErr:        false,
		},
		{
			name: "DuplicateContextPath",
			provider: func(provider *discovery.MockDiscoveryProvider) {
				provider.On("Discover",
					mock.AnythingOfType("*context.emptyCtx"),
					"myservice",
					true).
					Return(discovery.ServiceInstances{
						{
							ID:   "myservice-XYZ",
							Name: "myservice",
							Host: "10.10.10.10",
							Tags: []string{
								"contextPath=/my",
							},
							Port: 8080,
						},
					}, nil)
			},
			url:            "http://myservice/my/api/v1/test",
			wantHost:       "10.10.10.10:8080",
			wantPathPrefix: "/my/api/v1/test",
			wantErr:        false,
		},
		{
			name:           "NoDiscovery",
			url:            "http://10.10.10.10:8080/api/v1/test",
			wantHost:       "10.10.10.10:8080",
			wantPathPrefix: "",
			wantErr:        false,
		},
		{
			name: "DiscoveryFailed",
			url:  "http://myservice/api/v1/test",
			provider: func(provider *discovery.MockDiscoveryProvider) {
				provider.On("Discover",
					mock.AnythingOfType("*context.emptyCtx"),
					"myservice",
					true).
					Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "NoInstances",
			url:  "http://myservice/api/v1/test",
			provider: func(provider *discovery.MockDiscoveryProvider) {
				provider.On("Discover",
					mock.AnythingOfType("*context.emptyCtx"),
					"myservice",
					true).
					Return(discovery.ServiceInstances{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := setMockDiscoveryProvider()
			if tt.provider != nil {
				tt.provider(provider)
			}

			req := (&http.Request{}).WithContext(context.Background())
			req.URL, _ = url.Parse(tt.url)

			fn := func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, tt.wantHost, req.URL.Host)
				assert.True(t, strings.HasPrefix(req.URL.Path, tt.wantPathPrefix))
				return nil, nil
			}

			decorated := NewInterceptor(fn)

			resp, err := decorated(req)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			assert.Nil(t, resp)

		})
	}
}
