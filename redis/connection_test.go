package redis

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewConnectionConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ConnectionConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{}),
			},
			want: &ConnectionConfig{
				Enable:   false,
				Host:     "localhost",
				Port:     6379,
				Password: "",
				DB:       0,
				Sentinel: SentinelConfig{
					Enable: false,
					Master: "mymaster",
					Nodes:  []string{"localhost:26379"},
				},
				MaxRetries:  2,
				IdleTimeout: 1,
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.redis.enable":          "true",
					"spring.redis.host":            "remote-vm",
					"spring.redis.port":            "9999",
					"spring.redis.password":        "password",
					"spring.redis.sentinel.enable": "true",
					"spring.redis.sentinel.master": "mymaster",
					"spring.redis.sentinel.nodes":  "remote-vm1,remote-vm2",
					"spring.redis.db":              "2",
					"spring.redis.max-retries":     "3",
					"spring.redis.idle-timeout":    "4",
				}),
			},
			want: &ConnectionConfig{
				Enable:   true,
				Host:     "remote-vm",
				Port:     9999,
				Password: "password",
				DB:       2,
				Sentinel: SentinelConfig{
					Enable: true,
					Master: "mymaster",
					Nodes:  []string{"remote-vm1", "remote-vm2"},
				},
				MaxRetries:  3,
				IdleTimeout: 4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnectionConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnectionConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConnectionConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
