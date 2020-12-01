package notification

import (
	"context"
	"errors"
	"time"

	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/topics"
)

type Message struct {
	Context   Context                `json:"context"`
	EventName string                 `json:"event"`
	Payload   map[string]interface{} `json:"payload"`
}

type Context struct {
	User            Identifier   `json:"user"`
	Provider        Identifier   `json:"provider"`
	Tenant          Identifier   `json:"tenant"`
	RecipientEmails []string     `json:"recipientEmails"`
	EmailLocale     string       `json:"emailLocale"`
	Recipients      []Identifier `json:"recipients"`
	ServiceType     ServiceType  `json:"serviceType"`
	Initiator       Identifier   `json:"initiator"`
}

type Identifier struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ServiceType struct {
	LogicName              string `json:"logicName"`
	DisplayNameResourceKey string `json:"diplayNameResourceKey"`
}

//PayloadMessage : its the payload struct used by the notification service to unmarshal the Kafka.Message into
//if the payload field you need to add does not exist in the following payload struct please go ahead and add it.
//PayloadMessage refers to : Message.Payload for notification-go templates.
type PayloadMessage struct {
	User struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		UserID    string `json:"userId"`
		Name      string `json:"name"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Locale    string `json:"locale"`
	} `json:"user"`

	Alert struct {
		Time        string `json:"time"`
		Severity    string `json:"severity"`
		Type        string `json:"type"`
		Subtype     string `json:"subtype"`
		Service     string `json:"service"`
		Description string `json:"description"`
	}

	Tenant struct {
		ID           string      `json:"id"`
		Name         string      `json:"name"`
		DisplayName  string      `json:"displayname"`
		ProviderName string      `json:"providername"`
		Description  string      `json:"description"`
		Url          string      `json:"url"`
		GroupName    string      `json:"groupname"`
		Suspended    interface{} `json:"suspended"`
	} `json:"tenant"`

	Provider struct {
		ID               string `json:"id"`
		Email            string `json:"email"`
		Name             string `json:"name"`
		NotificationType string `json:"notificationType"`
		Locale           string `json:"locale"`
		DisplayName      string `json:"displayname"`
	} `json:"provider"`

	RestDetails   interface{} `json:"restDetails"`
	StatusMessage string      `json:"statusMessage"`
	Users         []struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		UserID    string `json:"userId"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Roles     []struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"roles"`
		Tenants []struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"tenants"`
		Status           string `json:"status"`
		IsDeleted        bool   `json:"isDeleted"`
		ProviderNameDesc string `json:"providerNameDesc"`
		ExpireDays       int    `json:"expireDays"`
		GraceLogins      int    `json:"graceLogins"`
		URL              string `json:"url"`
		TokenTimeoutMins string `json:"token_timeout_mins"`
	} `json:"users"`

	Service struct {
		ServiceID string `json:"serviceId"`
	} `json:"service,omitempty"`

	Devices []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		SN    string `json:"sn"`
		Model string `json:"model"`
		Type  string `json:"type"`
	} `json:"devices"`

	Address struct {
		DisplayName string `json:"displayname"`
	} `json:"address"`

	DisplayName string `json:"displayname"`
}

type MessageProducer interface {
	Message(context.Context) (Message, error)
}

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
			User: Identifier{
				Id:   userContextDetails.UserId,
				Name: *userContextDetails.Username,
			},
			Provider: Identifier{
				Id:   userContextDetails.ProviderId,
				Name: *userContextDetails.ProviderName,
			},
			Tenant: Identifier{
				Id:   userContextDetails.TenantId,
				Name: *userContextDetails.TenantName,
			},
		},
		Payload: make(map[string]interface{}),
	}, nil
}
