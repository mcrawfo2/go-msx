package populate

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewPermissionPopulatorConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *PermissionPopulatorConfig
		wantErr bool
	}{
		{
			name: "StructDefaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"populate.root": "/platform-common",
				}),
			},
			want: &PermissionPopulatorConfig{
				Enabled: false,
				Root:    "/platform-common/usermanagement",
			},
		},
		{
			name: "CustomOptions",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"populate.root": "/platform-common",
					"populate.usermanagement.permission.enabled": "true",
				}),
			},
			want: &PermissionPopulatorConfig{
				Enabled: true,
				Root:    "/platform-common/usermanagement",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPermissionPopulatorConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPermissionPopulatorConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPermissionPopulatorConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
