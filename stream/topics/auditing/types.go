package auditing

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/topics"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"time"
)

type Details map[string]string

type Message struct {
	Time        topics.Time          `json:"timestamp"`
	Service     string               `json:"service"`
	Type        string               `json:"type"`
	Subtype     string               `json:"subtype"`
	Trace       TraceAuditContext    `json:"trace"`
	Security    SecurityAuditContext `json:"security"`
	Details     Details              `json:"details"`
	Description string               `json:"description"`
	Keywords    string               `json:"keywords"`
}

type SecurityAuditContext struct {
	ClientId         string `json:"clientId"`
	UserId           string `json:"userId"`
	Username         string `json:"userName"`
	TenantId         string `json:"tenantId"`
	TenantName       string `json:"tenantName"`
	ProviderId       string `json:"providerId"`
	OriginalUsername string `json:"originalUsername"`
}

type TraceAuditContext struct {
	TraceId  string `json:"traceId"`
	SpanId   string `json:"spanId"`
	ParentId string `json:"parentId"`
}

func (m *Message) AddDetail(key, value string) {
	m.Details[key] = value
}

func (m *Message) AddDetails(kv map[string]string) {
	for k, v := range kv {
		m.Details[k] = v
	}
}

type MessageProducer interface {
	Message(context.Context) (Message, error)
}

func NewMessage(ctx context.Context) (Message, error) {
	userContext, err := security.NewUserContextDetails(ctx)
	if err != nil {
		return Message{}, err
	}

	var tenantId = types.EmptyUUID()
	var securityAudit SecurityAuditContext
	var trace TraceAuditContext

	if userContext.TenantId != nil {
		tenantId = userContext.TenantId

		securityAudit.ClientId = types.NewOptionalString(userContext.ClientId).OrEmpty()
		securityAudit.UserId = types.NewOptional(userContext.UserId).OrElse(types.EmptyUUID()).(types.UUID).String()
		securityAudit.Username = types.NewOptionalString(userContext.Username).OrEmpty()
		securityAudit.TenantId = tenantId.String()
		securityAudit.TenantName = types.NewOptionalString(userContext.TenantName).OrEmpty()
		securityAudit.OriginalUsername = types.NewOptionalString(userContext.Username).OrEmpty()
		securityAudit.ProviderId = types.NewOptional(userContext.ProviderId).OrElse(types.EmptyUUID()).(types.UUID).String()

		logContext, exists := log.LogContextFromContext(ctx)

		if exists {
			trace.ParentId = logContext[log.FieldParentId].(string)
			trace.SpanId = logContext[log.FieldSpanId].(string)
			trace.TraceId = logContext[log.FieldTraceId].(string)
		}
	}

	return Message{
		Time:     topics.Time(time.Now().UTC()),
		Type:     "GP",
		Trace:    trace,
		Security: securityAudit,
		Details:  make(Details),
	}, nil
}
