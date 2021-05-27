package consulprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
	"time"
)

func TestNewRegistrationProviderConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *RegistrationProviderConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"info.app.name": "TestNewRegistrationProviderConfigFromConfig/Defaults",
				}),
			},
			want: &RegistrationProviderConfig{
				Name:                "TestNewRegistrationProviderConfigFromConfig/Defaults",
				Register:            true,
				Scheme:              "http",
				RegisterHealthCheck: true,
				HealthCheckPath:     "/admin/health",
				HealthCheckInterval: 10 * time.Second,
				HealthCheckTimeout:  10 * time.Second,
				InstanceId:          "local",
				InstanceName:        "TestNewRegistrationProviderConfigFromConfig/Defaults",
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.consul.discovery.enabled":               "true",
					"spring.cloud.consul.discovery.name":                  "custom",
					"spring.cloud.consul.discovery.register":              "false",
					"spring.cloud.consul.discovery.ip-address":            "10.10.10.10",
					"spring.cloud.consul.discovery.interface":             "en7",
					"spring.cloud.consul.discovery.port":                  "9999",
					"spring.cloud.consul.discovery.scheme":                "https",
					"spring.cloud.consul.discovery.register-health-check": "false",
					"spring.cloud.consul.discovery.health-check-path":     "/admin/alive",
					"spring.cloud.consul.discovery.health-check-interval": "30s",
					"spring.cloud.consul.discovery.health-check-timeout":  "30s",
					"spring.cloud.consul.discovery.tags":                  "tag1,tag2=bravo",
					"spring.cloud.consul.discovery.instance-id":           "uuid",
					"spring.cloud.consul.discovery.instance-name":         "custom",
					"spring.cloud.consul.discovery.hidden-api-listing":    "true",
				}),
			},
			want: &RegistrationProviderConfig{
				Enabled:             true,
				Name:                "custom",
				Register:            false,
				IpAddress:           "10.10.10.10",
				Interface:           "en7",
				Port:                9999,
				Scheme:              "https",
				RegisterHealthCheck: false,
				HealthCheckPath:     "/admin/alive",
				HealthCheckInterval: 30 * time.Second,
				HealthCheckTimeout:  30 * time.Second,
				Tags:                "tag1,tag2=bravo",
				InstanceId:          "uuid",
				InstanceName:        "custom",
				HiddenApiListing: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRegistrationProviderConfigFromConfig(tt.args.cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewRegistrationProviderConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRegistrationProviderConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
