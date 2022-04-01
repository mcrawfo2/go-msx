// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package notification

import (
	"context"
	"errors"

	"cto-github.cisco.com/NFV-BU/go-msx/security"
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
		Context: Context{
			User: Identifier{
				Id:   userContextDetails.UserId.String(),
				Name: *userContextDetails.Username,
			},
			Provider: Identifier{
				Id:   userContextDetails.ProviderId.String(),
				Name: *userContextDetails.ProviderName,
			},
			Tenant: Identifier{
				Id:   userContextDetails.TenantId.String(),
				Name: *userContextDetails.TenantName,
			},
		},
		Payload: make(map[string]interface{}),
	}, nil
}
