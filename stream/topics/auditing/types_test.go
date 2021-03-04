package auditing

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/topics"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/securitytest"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestMessage_AddDetail(t *testing.T) {
	injector := securitytest.TenantAssignmentInjector()
	ctx := injector(context.Background())

	m, err := NewMessage(ctx)
	assert.NoError(t, err)
	m.AddDetail("key", "value")
	assert.Equal(t, "value", m.Details["key"])
}

func TestMessage_AddDetails(t *testing.T) {
	injector := securitytest.TenantAssignmentInjector()
	ctx := injector(context.Background())

	m, err := NewMessage(ctx)
	assert.NoError(t, err)
	m.AddDetails(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	assert.Equal(t, "value1", m.Details["key1"])
	assert.Equal(t, "value2", m.Details["key2"])
}

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
				ctx = log.ExtendContext(ctx, log.LogContext{
					log.FieldSpanId: "span-id",
					log.FieldTraceId: "trace-id",
					log.FieldParentId: "parent-id",
				})
				return ctx
			}(),
			want: Message{
				Time:  topics.Time(time.Now()),
				Type:  "GP",
				Trace: TraceAuditContext{
					TraceId:  "trace-id",
					SpanId:   "span-id",
					ParentId: "parent-id",
				},
				Security: SecurityAuditContext{
					ClientId:         "client-id",
					UserId:           "67f9b089-532e-4b54-9a06-8e4eade2114e",
					Username:         "tester",
					TenantId:         "960272b3-e800-43e6-86ce-7d51672bd80d",
					TenantName:       "test-tenant",
					ProviderId:       "30b62544-860e-42fb-93ba-bc7e771dff61",
					OriginalUsername: "tester",
				},
				Details: Details{},
			},
			wantErr: false,
		},
		{
			name: "NoUserContextDetails",
			ctx: context.Background(),
			wantErr: true,
		},
		{
			name: "NoLogContext",
			ctx: func() context.Context {
				injector := securitytest.TenantAssignmentInjector()
				ctx := injector(context.Background())
				return ctx
			}(),
			want: Message{
				Time:  topics.Time(time.Now()),
				Type:  "GP",
				Trace: TraceAuditContext{},
				Security: SecurityAuditContext{
					ClientId:         "client-id",
					UserId:           "67f9b089-532e-4b54-9a06-8e4eade2114e",
					Username:         "tester",
					TenantId:         "960272b3-e800-43e6-86ce-7d51672bd80d",
					TenantName:       "test-tenant",
					ProviderId:       "30b62544-860e-42fb-93ba-bc7e771dff61",
					OriginalUsername: "tester",
				},
				Details: Details{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMessage(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.Time = tt.want.Time
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
