package kafka

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewBindingConfigurationFromConfig(t *testing.T) {
	type args struct {
		cfg                 *config.Config
		key                 string
		streamBindingConfig *stream.BindingConfiguration
	}
	tests := []struct {
		name    string
		args    args
		want    *BindingConfiguration
		wantErr bool
	}{
		{
			name: "Default",
			args: args{
				cfg:                 configtest.NewStaticConfig(map[string]string{}),
				key:                 "alpha",
				streamBindingConfig: nil,
			},
			want: &BindingConfiguration{
				Producer: BindingProducerConfig{
					Sync: true,
				},
				StreamBindingConfig: nil,
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.cloud.stream.kafka.bindings.alpha.producer.sync": "false",
				}),
				key: "alpha",
				streamBindingConfig: &stream.BindingConfiguration{
					Destination: "DESTINATION_TOPIC",
					Group:       "group",
					ContentType: "content-type",
					Binder:      "kafka",
					BindingId:   "kafka-DESTINATION_TOPIC-0",
					LogMessages: true,
				},
			},
			want: &BindingConfiguration{
				Producer: BindingProducerConfig{
					Sync: false,
				},
				StreamBindingConfig: &stream.BindingConfiguration{
					Destination: "DESTINATION_TOPIC",
					Group:       "group",
					ContentType: "content-type",
					Binder:      "kafka",
					BindingId:   "kafka-DESTINATION_TOPIC-0",
					LogMessages: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBindingConfigurationFromConfig(tt.args.cfg, tt.args.key, tt.args.streamBindingConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBindingConfigurationFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBindingConfigurationFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
