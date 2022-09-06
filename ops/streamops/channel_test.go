package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestChannelDocumentor struct{}

func (m TestChannelDocumentor) DocType() string {
	return "mock"
}

func (m TestChannelDocumentor) Document(i *Channel) error {
	return nil
}

func TestChannel_Binder(t *testing.T) {
	sut := Channel{
		binding: &stream.BindingConfiguration{
			Binder: "channel-binding-binder",
		},
	}
	assert.Equal(t, sut.binding.Binder, sut.Binder())
}

func TestChannel_DefaultContentEncoding(t *testing.T) {
	sut := Channel{
		binding: &stream.BindingConfiguration{
			ContentEncoding: "channel-binding-contentEncoding",
		},
	}
	assert.Equal(t, sut.binding.ContentEncoding, sut.DefaultContentEncoding())
}

func TestChannel_DefaultContentType(t *testing.T) {
	sut := Channel{
		binding: &stream.BindingConfiguration{
			ContentType: "channel-binding-contentType",
		},
	}
	assert.Equal(t, sut.binding.ContentType, sut.DefaultContentType())
}

func TestChannel_Destination(t *testing.T) {
	sut := Channel{
		binding: &stream.BindingConfiguration{
			Destination: "channel-binding-destination",
		},
	}
	assert.Equal(t, sut.binding.Destination, sut.Destination())
}

func TestChannel_Documentor(t *testing.T) {
	documentors := ops.Documentors[Channel]{
		TestChannelDocumentor{},
	}

	type fields struct {
		name        string
		binding     *stream.BindingConfiguration
		documentors ops.Documentors[Channel]
	}
	type args struct {
		pred ops.DocumentorPredicate[Channel]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ops.Documentor[Channel]
	}{
		{
			name: "Match",
			fields: fields{
				name:        "channel",
				documentors: documentors,
			},
			args: args{
				pred: ops.WithDocType[Channel]("mock"),
			},
			want: documentors[0],
		},
		{
			name: "NoMatch",
			fields: fields{
				name: "channel",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Channel{
				name:        tt.fields.name,
				binding:     tt.fields.binding,
				documentors: tt.fields.documentors,
			}
			if got := c.Documentor(tt.args.pred); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Documentor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_Name(t *testing.T) {
	sut := Channel{
		name: "channel-name",
	}
	assert.Equal(t, sut.name, sut.Name())
}

func TestChannel_WithDocumentor(t *testing.T) {
	doc := TestChannelDocumentor{}
	type fields struct {
		name        string
		binding     *stream.BindingConfiguration
		documentors ops.Documentors[Channel]
	}
	type args struct {
		d ops.Documentor[Channel]
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Channel
	}{
		{
			name: "Success",
			fields: fields{
				name: "channel",
			},
			args: args{
				d: doc,
			},
			want: &Channel{
				name: "channel",
				documentors: []ops.Documentor[Channel]{
					doc,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Channel{
				name:        tt.fields.name,
				binding:     tt.fields.binding,
				documentors: tt.fields.documentors,
			}
			if got := c.WithDocumentor(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddDocumentor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewChannel(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *Channel
		wantErr bool
	}{
		{
			name: "Redis",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
					"spring.cloud.stream.bindings.channel.binder": "redis",
					"spring.application.name":                     "TestNewChannel",
				}),
				name: "channel",
			},
			want: &Channel{
				name: "channel",
				binding: &stream.BindingConfiguration{
					Destination: "channel",
					Binder:      "redis",
					Group:       "channel-TESTNEWCHANNEL_GP",
					ContentType: "application/json",
					LogMessages: true,
					Retry: retry.RetryConfig{
						Attempts: 3,
						Delay:    500,
						BackOff:  0,
						Linear:   true,
					},
					Consumer: stream.ConsumerConfiguration{
						AutoStartup:            true,
						Concurrency:            1,
						HeaderMode:             "none",
						MaxAttempts:            3,
						BackOffInitialInterval: 1000,
						BackOffMaxInterval:     10000,
						BackOffMultiplier:      2,
						DefaultRetryable:       true,
						InstanceIndex:          -1,
						InstanceCount:          -1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChannel(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChannel() error : %v", testhelpers.Dump(err))
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannel() result difference: %v", testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestRegisterChannel(t *testing.T) {
	c := &Channel{name: "channel"}
	RegisterChannel(c)
	assert.Equal(t, registeredChannels[c.name], c)
}

func TestRegisteredChannel(t *testing.T) {
	registeredChannels = make(map[string]*Channel)
	c := &Channel{name: "channel"}
	RegisterChannel(c)
	d := RegisteredChannel("channel")
	assert.Equal(t, d, c)
}

func TestRegisteredChannels(t *testing.T) {
	registeredChannels = make(map[string]*Channel)
	c := &Channel{name: "channel"}
	RegisterChannel(c)
	m := RegisteredChannels()
	assert.Equal(t, m, map[string]*Channel{
		"channel": c,
	})
}
