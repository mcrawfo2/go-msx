// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package gochannel

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
				cfg:                 configtest.NewInMemoryConfig(map[string]string{}),
				key:                 "alpha",
				streamBindingConfig: nil,
			},
			want: &BindingConfiguration{
				Producer: BindingProducerConfig{
					OutputChannelBuffer:            16,
					Persistent:                     false,
					BlockPublishUntilSubscriberAck: false,
				},
				StreamBindingConfig: nil,
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.stream.gochannel.bindings.alpha.producer.output-channel-buffer":              "10",
					"spring.cloud.stream.gochannel.bindings.alpha.producer.persistent":                         "true",
					"spring.cloud.stream.gochannel.bindings.alpha.producer.block-publish-until-subscriber-ack": "true",
				}),
				key: "alpha",
				streamBindingConfig: &stream.BindingConfiguration{
					Destination: "DESTINATION_TOPIC",
					Group:       "group",
					ContentType: "content-type",
					Binder:      "gochannel",
					BindingId:   "gochannel-DESTINATION_TOPIC-0",
					LogMessages: true,
				},
			},
			want: &BindingConfiguration{
				Producer: BindingProducerConfig{
					OutputChannelBuffer:            10,
					Persistent:                     true,
					BlockPublishUntilSubscriberAck: true,
				},
				StreamBindingConfig: &stream.BindingConfiguration{
					Destination: "DESTINATION_TOPIC",
					Group:       "group",
					ContentType: "content-type",
					Binder:      "gochannel",
					BindingId:   "gochannel-DESTINATION_TOPIC-0",
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
