package consulprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewDiscoveryProviderConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *DiscoveryProviderConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{}),
			},
			want: &DiscoveryProviderConfig{
				DefaultQueryTag: "",
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.consul.discovery.default-query-tag": "version=3.9.0",
				}),
			},
			want: &DiscoveryProviderConfig{
				DefaultQueryTag: "version=3.9.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDiscoveryProviderConfigFromConfig(tt.args.cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewDiscoveryProviderConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDiscoveryProviderConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
