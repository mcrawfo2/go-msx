// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package httpclient

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_NewRetryable(t *testing.T) {
	cfg := config.NewConfig()
	_ = cfg.Load(context.Background())
	ctx := config.ContextWithConfig(context.Background(), cfg)
	factory, _ := NewProductionHttpClientFactory(ctx)
	ctx = ContextWithFactory(ctx, factory)

	r := retry.NewRetry(ctx, DefaultHTTPClientRetryConfig)

	type args struct {
		ctx         context.Context
		configurer  Configurer
		retryConfig retry.RetryConfig
		backoff     retryablehttp.Backoff
		retryPolicy RetryPolicy
	}
	tests := []struct {
		name       string
		args       args
		wantClient bool
		wantErr    bool
	}{
		{
			name: "NoFactoryInContext",
			args: args{
				ctx: context.Background(),
			},
			wantClient: false,
			wantErr:    true,
		},
		{
			name: "OnlyContext",
			args: args{
				ctx: ctx,
			},
			wantClient: true,
		},
		{
			name: "ContextAndConfigurer",
			args: args{
				ctx: ctx,
				configurer: ClientConfigurer{
					ClientFuncs: []ClientConfigurationFunc{
						ClientTimeout(3 * time.Second),
					},
				},
			},
			wantClient: true,
		},
		{
			name: "ContextConfigurerRetry",
			args: args{
				ctx: ctx,
				configurer: ClientConfigurer{
					ClientFuncs: []ClientConfigurationFunc{
						ClientTimeout(3 * time.Second),
					},
				},
				retryConfig: DefaultHTTPClientRetryConfig,
			},
			wantClient: true,
		},
		{
			name: "ContextConfigurerBackoff",
			args: args{
				ctx: ctx,
				configurer: ClientConfigurer{
					ClientFuncs: []ClientConfigurationFunc{
						ClientTimeout(3 * time.Second),
					},
				},
				backoff: func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
					return time.Duration(r.GetCurrentDelay(attemptNum))
				},
			},
			wantClient: true,
		},
		{
			name: "ContextConfigurerRetryPolicy",
			args: args{
				ctx: ctx,
				configurer: ClientConfigurer{
					ClientFuncs: []ClientConfigurationFunc{
						ClientTimeout(3 * time.Second),
					},
				},
				retryConfig: DefaultHTTPClientRetryConfig,
				retryPolicy: retryablehttp.DefaultRetryPolicy,
			},
			wantClient: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRetryable(tt.args.ctx, tt.args.configurer, tt.args.retryConfig, tt.args.backoff, tt.args.retryPolicy)

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != nil) != tt.wantClient {
				t.Errorf("New() got = %v, wantClient %v", got, tt.wantClient)
			}

		})
	}

}

func Test_Retryable_Test429(t *testing.T) {
	cfg := config.NewConfig()
	_ = cfg.Load(context.Background())
	ctx := config.ContextWithConfig(context.Background(), cfg)
	factory, _ := NewProductionHttpClientFactory(ctx)
	ctx = ContextWithFactory(ctx, factory)

	retryable, _ := NewRetryable(
		ctx,
		ClientConfigurer{
			ClientFuncs: []ClientConfigurationFunc{
				applyFakeMerakiInterceptor,
			},
			TransportFuncs: []TransportConfigurationFunc{},
		},
		retry.RetryConfig{
			Attempts: 3,
		},
		nil,
		nil,
	)

	resp, err := retryable.Get("https://api.meraki.com/api/v1/organizations")
	if err != nil {
		logger.Info(err)
	}
	if resp != nil {
		logger.Info(resp.StatusCode)
	}
}

func applyFakeMerakiInterceptor(c *http.Client) {
	c.Transport = func(fn DoFunc) DoFunc {
		return func(req *http.Request) (resp *http.Response, err error) {
			url := req.URL
			host := url.Host

			/* change url straight up
			if strings.Contains(host, "meraki") {
				req.URL, _ = url.Parse("http://localhost:8080/testGet")
			}
			*/

			resp, err = fn(req)

			if strings.Contains(host, "meraki") {
				// resp.StatusCode = http.StatusBadGateway // 502

				// simulate 429
				resp.StatusCode = http.StatusTooManyRequests
				resp.Header.Set("Retry-After", "1")
			}

			return resp, nil
		}
	}(c.Transport.RoundTrip)
}