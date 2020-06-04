package auditing

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
)

const TopicName = "AUDITING_GENERIC_TOPIC"

func Publish(ctx context.Context, message Message) error {
	return stream.PublishObject(ctx, TopicName, message, nil)
}

func PublishFromProducer(ctx context.Context, producer MessageProducer) error {
	message, err := producer.Message(ctx)
	if err != nil {
		return err
	}

	return Publish(ctx, message)
}

const (
	SubTypeSubscription  = "SUBSCRIPTION"
	SubTypeSite          = "SITE"
	SubTypeDevice        = "DEVICE"
	SubTypeUser          = "USER"
	SubTypeTenant        = "TENANT"
	SubTypeSystem        = "SYSTEM"
	SubTypeScheduledTask = "SCHEDULE_TASK"
)

const (
	DetailsObjectType = "objectType"
	DetailsObjectId   = "objectId"
	DetailsUserId     = "userId"
	DetailsTenantId   = "tenantId"
	DetailsProviderId = "providerId"
	DetailsSeverity   = "severity"
)

const (
	SeverityInformational = "Informational"
	SeverityWarning       = "Warning"
	SeverityCritical      = "Critical"
)
