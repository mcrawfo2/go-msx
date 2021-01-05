package populate

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewServiceConfigPopulatorConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ServiceConfigPopulatorConfig
		wantErr bool
	}{
		{
			name: "StructDefaults",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"populate.root": "/platform-common",
				}),
			},
			want: &ServiceConfigPopulatorConfig{
				Enabled: false,
				Root:    "/platform-common/serviceconfig",
			},
		},
		{
			name: "CustomOptions",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"populate.root":                  "/platform-common",
					"populate.serviceconfig.enabled": "true",
				}),
			},
			want: &ServiceConfigPopulatorConfig{
				Enabled: true,
				Root:    "/platform-common/serviceconfig",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServiceConfigPopulatorConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServiceConfigPopulatorConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServiceConfigPopulatorConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
