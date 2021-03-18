package notification

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/securitytest"
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		want    Message
		wantErr bool
	}{
		{
			name: "Success",
			ctx: func() context.Context {
				injector := securitytest.TenantAssignmentInjector()
				ctx := injector(context.Background())
				return ctx
			}(),
			want: Message{
				Context: Context{
					User: Identifier{
						Id:   "67f9b089-532e-4b54-9a06-8e4eade2114e",
						Name: "tester",
					},
					Provider: Identifier{
						Id:   "30b62544-860e-42fb-93ba-bc7e771dff61",
						Name: "cisco",
					},
					Tenant: Identifier{
						Id:   "960272b3-e800-43e6-86ce-7d51672bd80d",
						Name: "test-tenant",
					},
				},
				Payload: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:    "NoUserContextDetails",
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name: "NoProviderId",
			ctx: func() context.Context {
				return securitytest.TokenDetailsProviderCustomizer(func(provider *securitytest.MockTokenDetailsProvider) {
					provider.ProviderId = nil
				})(context.Background())
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMessage(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
