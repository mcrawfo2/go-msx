package stream

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewPublisherService(t *testing.T) {
	mockPublisherService := new(MockPublisherService)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want PublisherService
	}{
		{
			name: "Fresh",
			args: args{
				ctx: context.Background(),
			},
			want: ProductionPublisherService,
		},
		{
			name: "FromContext",
			args: args{
				ctx: func() context.Context {
					return ContextWithPublisherService(context.Background(), mockPublisherService)
				}(),
			},
			want: mockPublisherService,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPublisherService(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
					t.Errorf("NewPublisherService() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_publisherFunc_Publish(t *testing.T) {
	type contextKey int
	const (
		contextKeyArgs contextKey = iota
		contextKeyTestingT
	)

	type args struct {
		ctx      context.Context
		topic    string
		payload  []byte
		metadata map[string]string
	}
	tests := []struct {
		name    string
		p       publisherFunc
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			p: func(ctx context.Context, topic string, payload []byte, metadata map[string]string) (err error) {
				t := ctx.Value(contextKeyTestingT).(*testing.T)
				a := ctx.Value(contextKeyArgs).(args)

				assert.Equal(t, a.topic, topic)
				assert.Equal(t, a.payload, payload)
				assert.Equal(t, a.metadata, metadata)

				return nil
			},
			args:    args{
				ctx:      context.Background(),
				topic:    "mock",
				payload:  []byte("{}"),
				metadata: map[string]string{},
			},
			wantErr: false,
		},
		{
			name:    "Failure",
			p: func(ctx context.Context, topic string, payload []byte, metadata map[string]string) (err error) {
				return errors.New("error")
			},
			args:    args{
				ctx:      context.Background(),
				topic:    "mock",
				payload:  []byte("{}"),
				metadata: map[string]string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(tt.args.ctx, contextKeyArgs, tt.args)
			ctx = context.WithValue(ctx, contextKeyTestingT, t)

			if err := tt.p.Publish(ctx, tt.args.topic, tt.args.payload, tt.args.metadata); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
