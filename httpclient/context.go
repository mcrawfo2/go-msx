package httpclient

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
)

var logger = log.NewLogger("msx.httpclient")

type httpClientContextKey int

const (
	contextKeyHttpClientFactory httpClientContextKey = iota
	contextKeyOperationName
	contextKeyHttpClientConfigurer
)

func ContextWithFactory(ctx context.Context, factory Factory) context.Context {
	return context.WithValue(ctx, contextKeyHttpClientFactory, factory)
}

func FactoryFromContext(ctx context.Context) Factory {
	factoryInterface := ctx.Value(contextKeyHttpClientFactory)
	if factoryInterface == nil {
		return nil
	}
	if factory, ok := factoryInterface.(Factory); !ok {
		logger.Warn("Context http client factory value is the wrong type")
		return nil
	} else {
		return factory
	}
}

func ContextWithOperationName(ctx context.Context, operationName string) context.Context {
	return context.WithValue(ctx, contextKeyOperationName, operationName)
}

func OperationNameFromContext(ctx context.Context) string {
	operationNameInterface := ctx.Value(contextKeyOperationName)
	if operationNameInterface == nil {
		return ""
	}
	if operationName, ok := operationNameInterface.(string); !ok {
		logger.Warn("Context http client operation name is the wrong type")
		return ""
	} else {
		return operationName
	}
}

func ContextWithConfigurer(ctx context.Context, configurer Configurer) context.Context {
	return context.WithValue(ctx, contextKeyHttpClientConfigurer, configurer)
}

func ConfigurerFromContext(ctx context.Context) Configurer {
	iface := ctx.Value(contextKeyHttpClientConfigurer)
	if iface == nil {
		return nil
	}
	if configurer, ok := iface.(Configurer); ok {
		return configurer
	}
	logger.Warn("Context http client configurer is the wrong type")
	return nil
}
