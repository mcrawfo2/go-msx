package sqldb

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewSqlConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "EmbeddedDefaults",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"cockroach.host":                     "localhost",
					"cockroach.port":                     "26257",
					"spring.datasource.driver":           "postgres",
					"spring.datasource.name":             "",
					"spring.datasource.username":         "root",
					"spring.datasource.password":         "",
					"spring.datasource.data-source-name": "postgresql://${spring.datasource.username}:${spring.datasource.password}@${cockroach.host}:${cockroach.port}/${spring.datasource.name}?sslmode=disable",
				}),
			},
			want: &Config{
				Driver:         "postgres",
				DataSourceName: "postgresql://root:@localhost:26257/?sslmode=disable",
				Enabled:        false,
			},
			wantErr: false,
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.datasource.driver":           "sqlite3",
					"spring.datasource.enabled":          "true",
					"spring.datasource.name":             "TestNewSqlConfigFromConfig",
					"spring.datasource.data-source-name": "file:${spring.datasource.name}?cache=shared&mode=memory",
				}),
			},
			want: &Config{
				Driver:         "sqlite3",
				DataSourceName: "file:TestNewSqlConfigFromConfig?cache=shared&mode=memory",
				Enabled:        true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSqlConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSqlConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSqlConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
