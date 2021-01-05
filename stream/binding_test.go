package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewBindingConfigurationFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
		key string
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
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.application.name":     "testservice",
					"spring.application.instance": "XYZABC",
				}),
				key: "SOME_TOPIC",
			},
			want: &BindingConfiguration{
				Destination: "SOME_TOPIC",
				Group:       "SOME_TOPIC-TEST_GP",
				ContentType: "application/json",
				Binder:      "kafka",
				BindingId:   "XYZABC",
				LogMessages: true,
				Retry: retry.RetryConfig{
					Attempts: 3,
					Delay:    500,
					BackOff:  0.0,
					Linear:   true,
				},
				Consumer: ConsumerConfiguration{
					AutoStartup:            true,
					Concurrency:            1,
					Partitioned:            false,
					HeaderMode:             "none",
					MaxAttempts:            3,
					BackOffInitialInterval: 1000,
					BackOffMaxInterval:     10000,
					BackOffMultiplier:      2.0,
					DefaultRetryable:       true,
					InstanceIndex:          -1,
					InstanceCount:          -1,
				},
			},
		},
		{
			name: "ConsumerDefaults",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.cloud.stream.bindings.bravo.destination":                "SOME_TOPIC",
					"spring.cloud.stream.bindings.bravo.group":                      "SOME_TOPIC-TEST2_GP",
					"spring.cloud.stream.bindings.bravo.binder":                     "nats",
					"spring.cloud.stream.bindings.bravo.binding-id":                 "XYZABC",
					"spring.cloud.stream.bindings.bravo.log-messages":               "true",
					"spring.cloud.stream.bindings.bravo.retry.attempts":             "10",
					"spring.cloud.stream.bindings.bravo.retry.delay":                "100",
					"spring.cloud.stream.bindings.bravo.retry.backoff":              "1.0",
					"spring.cloud.stream.bindings.bravo.retry.linear":               "false",
					"spring.cloud.stream.default.consumer.auto-startup":             "false",
					"spring.cloud.stream.default.consumer.concurrency":              "5",
					"spring.cloud.stream.default.consumer.partitioned":              "true",
					"spring.cloud.stream.default.consumer.max-attempts":             "5",
					"spring.cloud.stream.default.consumer.backoff-initial-interval": "100",
					"spring.cloud.stream.default.consumer.backoff-max-interval":     "1000",
					"spring.cloud.stream.default.consumer.default-retryable":        "true",
					"spring.cloud.stream.default.consumer.instance-index":           "1",
					"spring.cloud.stream.default.consumer.instance-count":           "2",
				}),
				key: "bravo",
			},
			want: &BindingConfiguration{
				Destination: "SOME_TOPIC",
				Group:       "bravo-SOME_TOPIC-TEST2_GP",
				ContentType: "application/json",
				Binder:      "nats",
				BindingId:   "XYZABC",
				LogMessages: true,
				Retry: retry.RetryConfig{
					Attempts: 10,
					Delay:    100,
					BackOff:  1.0,
					Linear:   false,
				},
				Consumer: ConsumerConfiguration{
					AutoStartup:            false,
					Concurrency:            5,
					Partitioned:            true,
					HeaderMode:             "none",
					MaxAttempts:            5,
					BackOffInitialInterval: 100,
					BackOffMaxInterval:     1000,
					BackOffMultiplier:      2.0,
					DefaultRetryable:       true,
					InstanceIndex:          1,
					InstanceCount:          2,
				},
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.cloud.stream.bindings.bravo.destination":                       "SOME_TOPIC",
					"spring.cloud.stream.bindings.bravo.group":                             "SOME_TOPIC-TEST2_GP",
					"spring.cloud.stream.bindings.bravo.binder":                            "nats",
					"spring.cloud.stream.bindings.bravo.binding-id":                        "XYZABC",
					"spring.cloud.stream.bindings.bravo.log-messages":                      "true",
					"spring.cloud.stream.bindings.bravo.retry.attempts":                    "10",
					"spring.cloud.stream.bindings.bravo.retry.delay":                       "100",
					"spring.cloud.stream.bindings.bravo.retry.backoff":                     "1.0",
					"spring.cloud.stream.bindings.bravo.retry.linear":                      "false",
					"spring.cloud.stream.bindings.bravo.consumer.auto-startup":             "false",
					"spring.cloud.stream.bindings.bravo.consumer.concurrency":              "5",
					"spring.cloud.stream.bindings.bravo.consumer.partitioned":              "true",
					"spring.cloud.stream.bindings.bravo.consumer.max-attempts":             "5",
					"spring.cloud.stream.bindings.bravo.consumer.backoff-initial-interval": "100",
					"spring.cloud.stream.bindings.bravo.consumer.backoff-max-interval":     "1000",
					"spring.cloud.stream.bindings.bravo.consumer.default-retryable":        "true",
					"spring.cloud.stream.bindings.bravo.consumer.instance-index":           "1",
					"spring.cloud.stream.bindings.bravo.consumer.instance-count":           "2",
				}),
				key: "bravo",
			},
			want: &BindingConfiguration{
				Destination: "SOME_TOPIC",
				Group:       "bravo-SOME_TOPIC-TEST2_GP",
				ContentType: "application/json",
				Binder:      "nats",
				BindingId:   "XYZABC",
				LogMessages: true,
				Retry: retry.RetryConfig{
					Attempts: 10,
					Delay:    100,
					BackOff:  1.0,
					Linear:   false,
				},
				Consumer: ConsumerConfiguration{
					AutoStartup:            false,
					Concurrency:            5,
					Partitioned:            true,
					HeaderMode:             "none",
					MaxAttempts:            5,
					BackOffInitialInterval: 100,
					BackOffMaxInterval:     1000,
					BackOffMultiplier:      2.0,
					DefaultRetryable:       true,
					InstanceIndex:          1,
					InstanceCount:          2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBindingConfigurationFromConfig(tt.args.cfg, tt.args.key)
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
