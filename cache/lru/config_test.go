// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package lru

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
	"time"
)

func testCacheConfig(ttl, expireLimit, expireFrequency, deAge, metrics, metricsPrefix, root string) *config.Config {
	return configtest.NewInMemoryConfig(map[string]string{
		root + ".ttl":              ttl,
		root + ".expire-limit":     expireLimit,
		root + ".expire-frequency": expireFrequency,
		root + ".de-age-on-access": deAge,
		root + ".metrics":          metrics,
		root + ".metrics-prefix":   metricsPrefix,
	})
}

func TestNewCacheConfig(t *testing.T) {
	type args struct {
		cfg  *config.Config
		root string
	}
	tests := []struct {
		name    string
		args    args
		want    *CacheConfig
		wantErr bool
	}{
		{
			name: "Configured",
			args: args{
				cfg: testCacheConfig("1m", "10", "30s", "true",
					"true", "idm", "cache"),
				root: "cache",
			},
			want: &CacheConfig{
				Ttl:             1 * time.Minute,
				ExpireLimit:     10,
				ExpireFrequency: 30 * time.Second,
				DeAgeOnAccess:   true,
				Metrics:         true,
				MetricsPrefix:   "idm",
			},
		},
		{
			name: "Defaults",
			args: args{
				cfg:  testCacheConfig("1m", "10", "30s", "", "", "", "cache"),
				root: "elsewhere",
			},
			want: &CacheConfig{
				Ttl:             5 * time.Minute,
				ExpireLimit:     100,
				ExpireFrequency: 30 * time.Second,
				DeAgeOnAccess:   false,
				Metrics:         false,
				MetricsPrefix:   "cache",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := NewCacheConfig(tt.args.cfg, tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCacheConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCacheConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
