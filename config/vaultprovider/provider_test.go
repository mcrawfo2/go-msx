package vaultprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thejerf/abtime"
	"reflect"
	"testing"
	"time"
)

func TestNewProviderConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ProviderConfig
		wantErr bool
	}{
		{
			name: "Default",
			args: args{configtest.NewInMemoryConfig(nil)},
			want: &ProviderConfig{
				Enabled:          false,
				Backend:          "secret",
				ProfileSeparator: "/",
				DefaultContext:   "defaultapplication",
				Delay:            20 * time.Second,
			},
		},
		{
			name: "Custom",
			args: args{configtest.NewInMemoryConfig(map[string]string{
				"spring.cloud.vault.generic.enabled":           "true",
				"spring.cloud.vault.generic.backend":           "secret-v2",
				"spring.cloud.vault.generic.profile-separator": "_",
				"spring.cloud.vault.generic.default-context":   "thirdpartyservices",
				"spring.cloud.vault.generic.delay":             "30s",
			})},
			want: &ProviderConfig{
				Enabled:          true,
				Backend:          "secret-v2",
				ProfileSeparator: "_",
				DefaultContext:   "thirdpartyservices",
				Delay:            30 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProviderConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProviderConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProviderConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewProvidersFromConfig(t *testing.T) {
	clock := abtime.NewManual()
	ctx := types.ContextWithClock(context.Background(), clock)
	mockConnection := new(vault.MockConnection)

	type args struct {
		ctx context.Context
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    []config.Provider
		wantErr bool
	}{
		{
			name: "ProviderConfigError",
			args: args{
				ctx: ctx,
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.generic.enabled": "falsy",
				}),
			},
			wantErr: true,
		},
		{
			name: "ProviderDisabled",
			args: args{
				ctx: ctx,
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.generic.enabled": "false",
				}),
			},
		},
		{
			name: "AppNameConfigError",
			args: args{
				ctx: ctx,
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.generic.enabled": "true",
				}),
			},
			wantErr: true,
		},
		{
			name: "ContextConnection",
			args: args{
				ctx: func() context.Context {
					return vault.ContextWithConnection(ctx, mockConnection)
				}(),
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.generic.enabled": "true",
					"spring.application.name":            "service",
				}),
			},
			want: []config.Provider{
				&Provider{
					name: "ContextConnection",
					cfg: &ProviderConfig{
						Enabled:          true,
						Backend:          "secret",
						ProfileSeparator: "/",
						DefaultContext:   "defaultapplication",
						Delay:            20 * time.Second,
					},
					contextPath: "defaultapplication",
					connection:  mockConnection,
				},
				&Provider{
					name: "ContextConnection",
					cfg: &ProviderConfig{
						Enabled:          true,
						Backend:          "secret",
						ProfileSeparator: "/",
						DefaultContext:   "defaultapplication",
						Delay:            20 * time.Second,
					},
					contextPath: "service",
					connection:  mockConnection,
				},
			},
		},
		{
			name: "NewConnection",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(ctx, nil),
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.generic.enabled": "true",
					"spring.application.name":            "service",
				}),
			},
			want: []config.Provider{
				&Provider{
					name: "NewConnection",
					cfg: &ProviderConfig{
						Enabled:          true,
						Backend:          "secret",
						ProfileSeparator: "/",
						DefaultContext:   "defaultapplication",
						Delay:            20 * time.Second,
					},
					contextPath: "defaultapplication",
				},
				&Provider{
					name: "NewConnection",
					cfg: &ProviderConfig{
						Enabled:          true,
						Backend:          "secret",
						ProfileSeparator: "/",
						DefaultContext:   "defaultapplication",
						Delay:            20 * time.Second,
					},
					contextPath: "service",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewProvidersFromConfig(tt.name, tt.args.ctx, tt.args.cfg)
			if tt.wantErr {
				if gotErr != nil {
					t.Log(gotErr.Error())
				}
				assert.Error(t, gotErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, gotErr)
				assert.Len(t, got, len(tt.want))
				for i := range tt.want {
					wantProvider, gotProvider := tt.want[i].(*Provider), got[i].(*Provider)
					wantProvider.loaded = gotProvider.loaded
					wantProvider.notify = gotProvider.notify
					wantProvider.clock = gotProvider.clock
					if wantProvider.connection == nil {
						assert.NotNil(t, gotProvider.connection)
						wantProvider.connection = gotProvider.connection
					}
				}
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNewProvider(t *testing.T) {
	cfg := &ProviderConfig{
		Enabled:          true,
		Backend:          "secret",
		ProfileSeparator: "/",
		DefaultContext:   "defaultapplication",
		Delay:            500,
	}

	provider := NewProvider("TestNewProvider", cfg, "context-path", nil, nil)
	assert.Equal(t, "secret/context-path", provider.ContextPath())
	assert.Equal(t, "TestNewProvider: [secret/context-path]", provider.Description())
	assert.NotNil(t, provider.Notify())

}

func TestProvider_Load(t *testing.T) {
	providerConfig, _ := NewProviderConfig(configtest.NewInMemoryConfig(nil))

	settings := map[string]string{
		"key-1": "value-1",
	}

	mockConnection := new(vault.MockConnection)
	mockConnection.
		On("ListSecrets", mock.AnythingOfType("*context.valueCtx"), "secret/defaultapplication").
		Return(settings, nil).
		Once()

	mockClock := abtime.NewManual()

	provider := NewProvider("loadSettings", providerConfig, "defaultapplication", mockConnection, mockClock)

	actualEntries, err := provider.Load(context.Background())
	assert.NoError(t, err)

	expectedEntries := config.ProviderEntries{
		{
			NormalizedName: "key1",
			Name:           "key-1",
			Value:          "value-1",
			Source:         provider,
		},
	}
	assert.Equal(t, expectedEntries, actualEntries)
}

func TestProvider_Run(t *testing.T) {
	providerConfig, _ := NewProviderConfig(configtest.NewInMemoryConfig(nil))

	settings := map[string]string{
		"key-1": "value-1",
	}

	ctx, cancelCtx := context.WithCancel(context.Background())

	mockClock := abtime.NewManual()

	mockConnection := new(vault.MockConnection)
	mockConnection.
		On("ListSecrets", mock.AnythingOfType("*context.cancelCtx"), "secret/defaultapplication").
		Return(settings, nil)

	mockConnection.
		On("ListSecrets", mock.AnythingOfType("*context.emptyCtx"), "secret/defaultapplication").
		Return(nil, errors.New("context cancelled"))

	provider := NewProvider("loadSettings", providerConfig, "defaultapplication", mockConnection, mockClock)

	go provider.Run(ctx)

	<-provider.Notify()

	actualEntries, err := provider.Load(context.Background())
	assert.NoError(t, err)

	expectedEntries := config.ProviderEntries{
		{
			NormalizedName: "key1",
			Name:           "key-1",
			Value:          "value-1",
			Source:         provider,
		},
	}
	assert.Equal(t, expectedEntries, actualEntries)

	cancelCtx()
}