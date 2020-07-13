package lru

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"reflect"
	"testing"
	"time"
)

func testCacheConfig(ttl, expireLimit, expireFrequency, root string) *config.Config {
	provider := config.NewStatic("static", map[string]string{
		root + ".ttl":              ttl,
		root + ".expire-limit":     expireLimit,
		root + ".expire-frequency": expireFrequency,
	})

	cfg := config.NewConfig(provider)
	_ = cfg.Load(context.Background())
	return cfg
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
				cfg:  testCacheConfig("1m", "10", "30s", "cache"),
				root: "cache",
			},
			want: &CacheConfig{
				Ttl:             1 * time.Minute,
				ExpireLimit:     10,
				ExpireFrequency: 30 * time.Second,
			},
		},
		{
			name: "Defaults",
			args: args{
				cfg:  testCacheConfig("1m", "10", "30s", "cache"),
				root: "elsewhere",
			},
			want: &CacheConfig{
				Ttl:             5 * time.Minute,
				ExpireLimit:     100,
				ExpireFrequency: 30 * time.Second,
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
