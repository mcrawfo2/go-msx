package auditing

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
)

const TopicName = "AUDITING_GENERIC_TOPIC"
const TopicNameDeviceV4 = "AUDITING_DEVICE_V4_TOPIC"
const TopicNameService = "AUDITING_SERVICE_TOPIC"
const TopicNameSite = "AUDITING_SITE_TOPIC"

func Publish(ctx context.Context, message Message) error {
	return PublishToTopic(ctx, TopicName, message)
}

func PublishFromProducer(ctx context.Context, producer MessageProducer) error {
	return PublishToTopicFromProducer(ctx, TopicName, producer)
}

func PublishToTopic(ctx context.Context, topicName string, message Message) error {
	return stream.PublishObject(ctx, topicName, message, nil)
}

func PublishToTopicFromProducer(ctx context.Context, topicName string, producer MessageProducer) error {
	message, err := producer.Message(ctx)
	if err != nil {
		return err
	}

	return PublishToTopic(ctx, topicName, message)
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

const (
	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
	ActionForceDelete = "FORCEDELETE"
)
