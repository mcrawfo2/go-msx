package httpclient

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/mocks"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestConfigurerFromContext(t *testing.T) {
	customizer := new(ClientConfigurer)

	tests := []struct {
		name    string
		ctx     context.Context
		want    Configurer
		wantLog logtest.Check
	}{
		{
			name: "NotExists",
			ctx:  context.Background(),
			want: nil,
		},
		{
			name: "Exists",
			ctx:  context.WithValue(context.Background(), contextKeyHttpClientConfigurer, customizer),
			want: customizer,
		},
		{
			name: "Invalid",
			ctx:  context.WithValue(context.Background(), contextKeyHttpClientConfigurer, "configurer"),
			want: nil,
			wantLog: logtest.Check{
				Validators: []logtest.EntryPredicate{
					logtest.HasLevel(logrus.WarnLevel),
					logtest.HasMessage(`Context http client configurer is the wrong type`),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording := logtest.RecordLogging()

			if got := ConfigurerFromContext(tt.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigurerFromContext() = %v, want %v", got, tt.want)
			}

			errors := tt.wantLog.Check(recording)
			assert.Len(t, errors, 0)
		})
	}
}

func TestContextWithConfigurer(t *testing.T) {
	customizer := new(ClientConfigurer)
	ctx := ContextWithConfigurer(context.Background(), customizer)
	assert.NotNil(t, ctx)
	assert.Equal(t, customizer, ctx.Value(contextKeyHttpClientConfigurer))
}

func TestFactoryFromContext(t *testing.T) {
	factory := new(mocks.HttpClientFactory)

	tests := []struct {
		name    string
		ctx     context.Context
		want    Factory
		wantLog logtest.Check
	}{
		{
			name: "NotExists",
			ctx:  context.Background(),
			want: nil,
		},
		{
			name: "Exists",
			ctx:  context.WithValue(context.Background(), contextKeyHttpClientFactory, factory),
			want: factory,
		},
		{
			name: "Invalid",
			ctx:  context.WithValue(context.Background(), contextKeyHttpClientFactory, "factory"),
			wantLog: logtest.Check{
				Validators: []logtest.EntryPredicate{
					logtest.HasLevel(logrus.WarnLevel),
					logtest.HasMessage(`Context http client factory value is the wrong type`),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FactoryFromContext(tt.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigurerFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContextWithFactory(t *testing.T) {
	factory := new(mocks.HttpClientFactory)
	ctx := ContextWithFactory(context.Background(), factory)
	assert.NotNil(t, ctx)
	assert.Equal(t, factory, ctx.Value(contextKeyHttpClientFactory))
}

func TestOperationNameFromContext(t *testing.T) {
	const operationName = "my-operation-name"

	tests := []struct {
		name    string
		ctx     context.Context
		want    string
		wantLog logtest.Check
	}{
		{
			name: "NotExists",
			ctx:  context.Background(),
			want: "",
		},
		{
			name: "Exists",
			ctx:  context.WithValue(context.Background(), contextKeyOperationName, operationName),
			want: operationName,
		},
		{
			name: "Invalid",
			ctx:  context.WithValue(context.Background(), contextKeyOperationName, 311),
			want: "",
			wantLog: logtest.Check{
				Validators: []logtest.EntryPredicate{
					logtest.HasLevel(logrus.WarnLevel),
					logtest.HasMessage(`Context http client operation name is the wrong type`),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OperationNameFromContext(tt.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigurerFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContextWithOperationName(t *testing.T) {
	const operationName = "my-operation-name"
	ctx := ContextWithOperationName(context.Background(), operationName)
	assert.NotNil(t, ctx)
	assert.Equal(t, operationName, ctx.Value(contextKeyOperationName))
}
