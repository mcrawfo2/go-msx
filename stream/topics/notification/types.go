package notification

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/topics"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"errors"
	"time"
)

type Message struct {
	Timestamp topics.Time       `json:"timestamp"`
	Version   *string           `json:"version"`
	Context   Context           `json:"context"`
	EventName string            `json:"event"`
	Payload   map[string]string `json:"payload"`
}

type Context struct {
	User            User                     `json:"user"`
	Provider        Provider                 `json:"provider"`
	Tenant          Tenant                   `json:"tenant"`
	RecipientEmails []string                 `json:"recipientEmails"`
	EmailLocale     string                   `json:"emailLocale"`
	Recipients      []map[string]interface{} `json:"recipients"`
	ServiceType     map[string]interface{}   `json:"serviceType"`
	Initiator       Initiator                `json:"initiator"`
}

type User struct {
	Id   types.UUID `json:"id"`
	Name string     `json:"name"`
}

type Provider struct {
	Id   types.UUID `json:"id"`
	Name string     `json:"name"`
}

type Tenant struct {
	Id   types.UUID `json:"id"`
	Name string     `json:"name"`
}

type Initiator struct {
	Id   types.UUID `json:"id"`
	Name string     `json:"name"`
}

type MessageProducer interface {
	Message(context.Context) (Message, error)
}

type Details map[string]string

func NewMessage(ctx context.Context) (Message, error) {
	userContextDetails, err := security.NewUserContextDetails(ctx)
	if err != nil {
		return Message{}, err
	} else if userContextDetails.ProviderId == nil {
		return Message{}, errors.New("provider id not identified")
	}

	return Message{
		Timestamp: topics.Time(time.Now().UTC()),
		Context: Context{
			User: User{
				Id:   userContextDetails.UserId,
				Name: *userContextDetails.Username,
			},
			Provider: Provider{
				Id:   userContextDetails.ProviderId,
				Name: *userContextDetails.ProviderName,
			},
			Tenant: Tenant{
				Id:   userContextDetails.TenantId,
				Name: *userContextDetails.TenantName,
			},
		},
		Payload: map[string]string{},
	}, nil
}
