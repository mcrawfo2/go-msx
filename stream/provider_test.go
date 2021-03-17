package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func registerMockProvider() (*MockPublisher, *MockSubscriber) {
	mockPublisher := new(MockPublisher)
	errPublisher := errors.New("publisher error")

	mockSubscriber := new(MockSubscriber)
	errSubscriber := errors.New("subscriber error")

	mockProvider := new(MockProvider)
	mockProvider.
		On("NewPublisher",
			mock.AnythingOfType("*config.Config"),
			"mybinding",
			mock.AnythingOfType("*stream.BindingConfiguration")).
		Return(mockPublisher, nil)
	mockProvider.
		On("NewPublisher",
			mock.AnythingOfType("*config.Config"),
			"errbinding",
			mock.AnythingOfType("*stream.BindingConfiguration")).
		Return(nil, errPublisher)

	mockProvider.
		On("NewSubscriber",
			mock.AnythingOfType("*config.Config"),
			"mybinding",
			mock.AnythingOfType("*stream.BindingConfiguration")).
		Return(mockSubscriber, nil)
	mockProvider.
		On("NewSubscriber",
			mock.AnythingOfType("*config.Config"),
			"errbinding",
			mock.AnythingOfType("*stream.BindingConfiguration")).
		Return(nil, errSubscriber)

	RegisterProvider("mock", mockProvider)

	return mockPublisher, mockSubscriber
}

func TestNewPublisher(t *testing.T) {
	mockPublisher, _ := registerMockProvider()

	type args struct {
		cfg  *config.Config
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Publisher
		wantErr bool
	}{
		{
			name: "Exists",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.stream.bindings.mybinding.binder": "mock",
					"spring.application.name":                       "TestNewPublisher",
				}),
				name: "mybinding",
			},
			want:    mockPublisher,
			wantErr: false,
		},
		{
			name: "NotExists",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.application.name": "TestNewPublisher",
				}),
				name: "anotherbinding",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "PublisherFailed",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.stream.bindings.errbinding.binder": "mock",
					"spring.application.name":                        "TestNewPublisher",
				}),
				name: "errbinding",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPublisher(tt.args.cfg, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPublisher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				assert.Equal(t, tt.want, got)
				return
			}

			if !reflect.DeepEqual(got.(*StatsPublisher).publisher, tt.want) {
				t.Errorf("NewPublisher() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSubscriber(t *testing.T) {
	_, mockSubscriber := registerMockProvider()

	type args struct {
		cfg  *config.Config
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    message.Subscriber
		wantErr bool
	}{
		{
			name: "Exists",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.stream.bindings.mybinding.binder": "mock",
					"spring.application.name":                       "TestNewSubscriber",
				}),
				name: "mybinding",
			},
			want:    mockSubscriber,
			wantErr: false,
		},
		{
			name: "NotExists",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.application.name": "TestNewSubscriber",
				}),
				name: "anotherbinding",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "NotEnabled",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.application.name":                           "TestNewSubscriber",
					"spring.cloud.stream.default.consumer.auto-startup": "false",
				}),
				name: "anotherbinding",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "SubscriberFailed",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.stream.bindings.errbinding.binder": "mock",
					"spring.application.name":                        "TestNewSubscriber",
				}),
				name: "errbinding",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSubscriber(tt.args.cfg, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSubscriber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				assert.Equal(t, tt.want, got)
				return
			}

			if !reflect.DeepEqual(got.(*StatsSubscriber).subscriber, tt.want) {
				t.Errorf("NewSubscriber() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterProvider(t *testing.T) {
	type args struct {
		name     string
		provider Provider
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Mock",
			args: args{
				name:     "mock",
				provider: new(MockProvider),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterProvider(tt.args.name, tt.args.provider)
			assert.Equal(t, tt.args.provider, providers[tt.args.name])
		})
	}
}
