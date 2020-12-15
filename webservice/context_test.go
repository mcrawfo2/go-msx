package webservice

import (
	"context"
	"github.com/emicklei/go-restful"
	"reflect"
	"testing"
)

func TestAuthenticationProviderFromContext(t *testing.T) {
	provider := new(MockAuthenticationProvider)

	type args struct {
		ctx      context.Context
	}
	tests := []struct {
		name string
		args args
		want AuthenticationProvider
	}{
		{
			name: "Exists",
			args: args{
				ctx: ContextWithSecurityProvider(context.Background(), provider),
			},
			want: provider,
		},
		{
			name: "NotExists",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthenticationProviderFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthenticationProviderFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainerFromContext(t *testing.T) {
	container := new(restful.Container)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *restful.Container
	}{
		{
			name: "Exists",
			args: args{
				ctx: ContextWithContainer(context.Background(), container),
			},
			want: container,
		},
		{
			name: "NotExists",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainerFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContainerFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteFromContext(t *testing.T) {
	route := new(restful.Route)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *restful.Route
	}{
		{
			name: "Exists",
			args: args{
				ctx: ContextWithRoute(context.Background(), route),
			},
			want: route,
		},
		{
			name: "NotExists",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RouteFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouteFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteOperationFromContext(t *testing.T) {
	routeOperation := "route-operation"

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Exists",
			args: args{
				ctx: ContextWithRouteOperation(context.Background(), routeOperation),
			},
			want: routeOperation,
		},
		{
			name: "NotExists",
			args: args{
				ctx: context.Background(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RouteOperationFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouteOperationFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterFromContext(t *testing.T) {
	route := new(restful.CurlyRouter)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want restful.RouteSelector
	}{
		{
			name: "Exists",
			args: args{
				ctx: ContextWithRouter(context.Background(), route),
			},
			want: route,
		},
		{
			name: "NotExists",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RouterFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouterFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceFromContext(t *testing.T) {
	svc := new(restful.WebService)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *restful.WebService
	}{
		{
			name: "Exists",
			args: args{
				ctx: ContextWithService(context.Background(), svc),
			},
			want: svc,
		},
		{
			name: "NotExists",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ServiceFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServerFromContext(t *testing.T) {
	server := new(WebServer)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *WebServer
	}{
		{
			name: "Exists",
			args: args{
				ctx: ContextWithWebServerValue(context.Background(), server),
			},
			want: server,
		},
		{
			name: "NotExists",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WebServerFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WebServerFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
