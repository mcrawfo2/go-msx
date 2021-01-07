package tokensource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewKubernetesConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *KubernetesConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.token-source.kubernetes.role": "PlatformMicroservice",
				}),
			},
			want: &KubernetesConfig{
				Path:    "/auth/kubernetes/login",
				JWTPath: "/run/secrets/kubernetes.io/serviceaccount/token",
				Role:    "PlatformMicroservice",
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.token-source.kubernetes.path":     "/subpath",
					"spring.cloud.vault.token-source.kubernetes.jwt-path": "/another/path/to/token",
					"spring.cloud.vault.token-source.kubernetes.role":     "PlatformMicroservice",
				}),
			},
			want: &KubernetesConfig{
				Path:    "/subpath",
				JWTPath: "/another/path/to/token",
				Role:    "PlatformMicroservice",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKubernetesConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKubernetesConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKubernetesConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
