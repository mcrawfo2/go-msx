// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestIntransientPublisher_Close(t *testing.T) {
	type fields struct {
		publisher Publisher
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				publisher: func() Publisher {
					mockPublisher := new(MockPublisher)
					mockPublisher.
						On("Close").
						Return(nil)
					return mockPublisher
				}(),
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				publisher: func() Publisher {
					mockPublisher := new(MockPublisher)
					mockPublisher.
						On("Close").
						Return(errors.New("error"))
					return mockPublisher
				}(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &IntransientPublisher{
				publisher: tt.fields.publisher,
			}
			if err := n.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntransientPublisher_Publish(t *testing.T) {
	type fields struct {
		publisher Publisher
	}
	type args struct {
		msg *message.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				publisher: func() Publisher {
					mockPublisher := new(MockPublisher)
					mockPublisher.
						On("Publish", mock.AnythingOfType("*message.Message")).
						Return(nil)
					return mockPublisher
				}(),
			},
			args: args{
				msg: message.NewMessage(watermill.NewUUID(), []byte("{}")),
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				publisher: func() Publisher {
					mockPublisher := new(MockPublisher)
					mockPublisher.
						On("Publish", mock.AnythingOfType("*message.Message")).
						Return(errors.New("error"))
					return mockPublisher
				}(),
			},
			args: args{
				msg: message.NewMessage(watermill.NewUUID(), []byte("{}")),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &IntransientPublisher{
				publisher: tt.fields.publisher,
			}
			if err := n.Publish(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewIntransientPublisher(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := NewIntransientPublisher(mockPublisher)
	assert.Equal(t, publisher.(*IntransientPublisher).publisher, mockPublisher)
}

func TestNewTopicPublisher(t *testing.T) {
	mockMessagePublisher := new(MockMessagePublisher)
	cfg := configtest.NewInMemoryConfig(map[string]string{
		"spring.application.name":                       "TestNewTopicPublisher",
		"spring.cloud.stream.bindings.mybinding.binder": "mock",
	})
	bindingConfiguration, err := NewBindingConfigurationFromConfig(cfg, "mybinding")
	assert.NoError(t, err)

	actualPublisher := NewTopicPublisher(mockMessagePublisher, bindingConfiguration)
	assert.NotNil(t, actualPublisher)

	tracePublisher := actualPublisher.(*TracePublisher)
	topicPublisher := tracePublisher.publisher

	assert.Equal(t, mockMessagePublisher, topicPublisher.(*TopicPublisher).publisher)
	assert.Equal(t, bindingConfiguration, topicPublisher.(*TopicPublisher).cfg)
}

func TestPublish(t *testing.T) {
	type args struct {
		ctx      context.Context
		binding  string
		payload  []byte
		metadata map[string]string
	}
	tests := []struct {
		name      string
		args      args
		publisher func(publisher *MockPublisher)
		wantErr   bool
	}{
		{
			name: "Success",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
					"spring.application.name":                       "TestPublish",
					"spring.cloud.stream.bindings.mybinding.binder": "mock",
				}),
				binding: "mybinding",
				payload: []byte("{}"),
				metadata: map[string]string{
					"X-MSX-Tenant-ID": "72b94158-20dd-48d2-b87a-52e2ae0f4910",
				},
			},
			publisher: func(publisher *MockPublisher) {
				publisher.
					On("Publish", mock.AnythingOfType("*message.Message")).
					Return(nil)
				publisher.
					On("Close").
					Return(nil)
			},
		},
		{
			name: "ErrNoConfig",
			args: args{
				ctx:     context.Background(),
				binding: "mybinding",
			},
			wantErr: true,
		},
		{
			name: "ErrNoBinding",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
					"spring.application.name":                       "TestPublish",
					"spring.cloud.stream.bindings.mybinding.binder": "missing",
				}),
				binding: "mybinding",
			},
			wantErr: true,
		},
		{
			name: "ErrPublish",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
					"spring.application.name":                       "TestPublish",
					"spring.cloud.stream.bindings.mybinding.binder": "mock",
				}),
				binding: "mybinding",
			},
			publisher: func(publisher *MockPublisher) {
				publisher.
					On("Publish", mock.AnythingOfType("*message.Message")).
					Return(errors.New("publish error"))
				publisher.
					On("Close").
					Return(nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPublisher, _ := registerMockProvider()
			if tt.publisher != nil {
				tt.publisher(mockPublisher)
			}

			if err := Publish(tt.args.ctx, tt.args.binding, tt.args.payload, tt.args.metadata); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPublishObject(t *testing.T) {
	mockPublisher, _ := registerMockProvider()
	mockPublisher.
		On("Publish", mock.AnythingOfType("*message.Message")).
		Return(nil)
	mockPublisher.
		On("Close").
		Return(nil)

	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name":                       "TestPublish",
		"spring.cloud.stream.bindings.mybinding.binder": "mock",
	})
	binding := "mybinding"
	payload := struct{}{}
	metadata := map[string]string{
		"X-MSX-Tenant-ID": "72b94158-20dd-48d2-b87a-52e2ae0f4910",
	}

	err := PublishObject(ctx, binding, payload, metadata)
	assert.NoError(t, err)
}

func TestTopicPublisher_Close(t *testing.T) {
	expectedError := errors.New("error")

	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.
		On("Close").
		Return(expectedError)

	cfg := configtest.NewInMemoryConfig(map[string]string{
		"spring.application.name":                       "TestTopicPublisher_Close",
		"spring.cloud.stream.bindings.mybinding.binder": "mock",
	})
	bindingConfiguration, err := NewBindingConfigurationFromConfig(cfg, "mybinding")
	assert.NoError(t, err)

	topicPublisher := NewTopicPublisher(mockMessagePublisher, bindingConfiguration)
	assert.NotNil(t, topicPublisher)

	actualError := topicPublisher.Close()
	assert.Error(t, actualError)
}

func TestTopicPublisher_Publish(t *testing.T) {
	expectedError := errors.New("error")

	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.
		On("Publish",
			"mybinding",
			mock.AnythingOfType("*message.Message")).
		Return(expectedError)

	cfg := configtest.NewInMemoryConfig(map[string]string{
		"spring.application.name":                       "TestTopicPublisher_Close",
		"spring.cloud.stream.bindings.mybinding.binder": "mock",
	})
	bindingConfiguration, err := NewBindingConfigurationFromConfig(cfg, "mybinding")
	assert.NoError(t, err)

	topicPublisher := NewTopicPublisher(mockMessagePublisher, bindingConfiguration)
	assert.NotNil(t, topicPublisher)

	actualError := topicPublisher.Publish(message.NewMessage(watermill.NewUUID(), nil))
	assert.Error(t, actualError)
}
