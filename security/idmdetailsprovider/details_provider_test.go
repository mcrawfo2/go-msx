// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package idmdetailsprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cache/lru"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
	"time"
)

func TestNewIdmTokenDetailsProviderConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *IdmTokenDetailsProviderConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{}),
			},
			want: &IdmTokenDetailsProviderConfig{
				ActiveCache: lru.CacheConfig{
					Ttl:             300 * time.Second,
					ExpireLimit:     100,
					ExpireFrequency: 30 * time.Second,
					Metrics:         false,
					MetricsPrefix:   "cache",
				},
				DetailsCache: lru.CacheConfig{
					Ttl:             300 * time.Second,
					ExpireLimit:     100,
					ExpireFrequency: 30 * time.Second,
					DeAgeOnAccess:   false,
					Metrics:         false,
					MetricsPrefix:   "cache",
				},
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"security.token.details.active-cache.ttl":               "30s",
					"security.token.details.active-cache.expire-limit":      "10",
					"security.token.details.active-cache.expire-frequency":  "5s",
					"security.token.details.active-cache.metrics":           "true",
					"security.token.details.active-cache.metrics-prefix":    "idm_active",
					"security.token.details.details-cache.ttl":              "40s",
					"security.token.details.details-cache.expire-limit":     "20",
					"security.token.details.details-cache.expire-frequency": "15s",
					"security.token.details.details-cache.metrics":          "true",
					"security.token.details.details-cache.de-age-on-access": "true",
					"security.token.details.details-cache.metrics-prefix":   "idm_details",
				}),
			},
			want: &IdmTokenDetailsProviderConfig{
				ActiveCache: lru.CacheConfig{
					Ttl:             30 * time.Second,
					ExpireLimit:     10,
					ExpireFrequency: 5 * time.Second,
					DeAgeOnAccess:   false,
					Metrics:         true,
					MetricsPrefix:   "idm_active",
				},
				DetailsCache: lru.CacheConfig{
					Ttl:             40 * time.Second,
					ExpireLimit:     20,
					ExpireFrequency: 15 * time.Second,
					DeAgeOnAccess:   true,
					Metrics:         true,
					MetricsPrefix:   "idm_details",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewIdmTokenDetailsProviderConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIdmTokenDetailsProviderConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIdmTokenDetailsProviderConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
