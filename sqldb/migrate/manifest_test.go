package migrate

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewManifestConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ManifestConfig
		wantErr bool
	}{
		{
			name: "Default",
			args: args{
				cfg: configtest.NewStaticConfig(nil),
			},
			want: &ManifestConfig{
				PostUpgrade: "",
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"migrate.post-upgrade": "some-value",
				}),
			},
			want: &ManifestConfig{
				PostUpgrade: "some-value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManifestConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManifestConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManifestConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
