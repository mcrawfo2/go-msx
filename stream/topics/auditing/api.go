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
	SubTypeActivity      = "ACTIVITY"
)

const (
	DetailsObjectType          = "objectType"
	DetailsObjectId            = "objectId"
	DetailsUserId              = "userId"
	DetailsTenantId            = "tenantId"
	DetailsProviderId          = "providerId"
	DetailsSeverity            = "severity"
	DetailsObjectName          = "objectName"
	DetailsTargetId            = "targetId"
	DetailsTargetType          = "targetType"
	DetailsTargetName          = "targetName"
	DetailsDescriptionResource = "descriptionResource"
	DetailsServiceType         = "serviceType"
	DetailsAction              = "action"
)

const (
	SeverityInformational = "Informational"
	SeverityWarning       = "Warning"
	SeverityCritical      = "Critical"

	DisplaySeverityPoor     = "POOR"
	DisplaySeverityCritical = "CRITICAL"
	DisplaySeverityFair     = "FAIR"
	DisplaySeverityGood     = "GOOD"
	DisplaySeverityUnknown  = "UNKNOWN"
)

const (
	ActionCreate      = "CREATE"
	ActionUpdate      = "UPDATE"
	ActionDelete      = "DELETE"
	ActionForceDelete = "FORCEDELETE"
)

const (
	ObjectTypeTenant = "TENANT"
)
