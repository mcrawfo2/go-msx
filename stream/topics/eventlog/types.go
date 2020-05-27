package eventlog

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/topics"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"time"
)

type Type string

const (
	TypeSubscription  Type = "SUBSCRIPTION"
	TypeSite          Type = "SITE"
	TypeDevice        Type = "DEVICE"
	TypeUser          Type = "USER"
	TypeTenant        Type = "TENANT"
	TypeSystem        Type = "SYSTEM"
	TypeScheduledTask Type = "SCHEDULE_TASK"
)

type Severity string

const (
	SeverityInformational Severity = "Informational"
	SeverityWarning       Severity = "Warning"
	SeverityCritical      Severity = "Critical"
)

type Message struct {
	Timestamp   topics.Time       `json:"timestamp"`
	Version     *string           `json:"version"`
	UserId      types.UUID        `json:"userId"`
	ProviderId  types.UUID        `json:"providerId"`
	TenantId    types.UUID        `json:"tenantId"`
	ObjectType  *string           `json:"objectType"`
	ObjectId    string            `json:"objectId"`
	EventType   Type              `json:"eventType"`
	Details     map[string]string `json:"details"`
	Severity    Severity          `json:"severity"`
	Description string            `json:"description"`
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
	Message(context.Context) Message
}

type Details map[string]string

func NewMessage(ctx context.Context) Message {
	var userId = types.EmptyUUID() // TODO
	var tenantId = types.EmptyUUID()
	var providerId = types.EmptyUUID() // TODO

	// TODO: UserContextDetails
	userContext := security.UserContextFromContext(ctx)
	if userContext.TenantId != nil {
		tenantId = userContext.TenantId
	}

	return Message{
		Timestamp:  topics.Time(time.Now().UTC()),
		UserId:     userId,
		ProviderId: providerId,
		TenantId:   tenantId,
		Details:    make(Details),
		Severity:   SeverityInformational,
	}
}
