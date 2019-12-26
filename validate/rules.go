package validate

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"time"
)

var (
	Required      = validation.Required
	NilOrNotEmpty = validation.NilOrNotEmpty
	OptionalUuid  = []Rule{validation.NilOrNotEmpty, IfNotNil(is.UUID)}
	ValidScope    = []Rule{validation.Required, validation.In([]interface{}{
		"controlPlaneId",
		"deviceId",
		"deviceType",
		"deviceSubType",
		"providerId",
		"serviceId",
		"serviceType",
		"shardId",
		"siteId",
		"subscriptionId",
		"templateId",
		"tenantGroupId",
		"tenantId",
	}...)}

	IsDuration    = RuleFunc(CheckDuration)
)

type Rule interface {
	Validate(value interface{}) error
}

type RuleFunc func(value interface{}) error

func (f RuleFunc) Validate(value interface{}) error {
	return f(value)
}

var Self = RuleFunc(func(value interface{}) error {
	if validatable, ok := value.(Validatable); ok {
		return Validate(validatable)
	}
	return nil
})

func IfNotNil(rules ...validation.Rule) RuleFunc {
	return func(value interface{}) error {
		result := types.ErrorList{}
		for _, rule := range rules {
			result = append(result, rule.Validate(value))
		}
		return result.Filter()
	}
}

func CheckDuration(value interface{}) error {
	valueString, ok := value.(string)
	if !ok {
		return errors.New("Duration is not a string")
	}

	_, err := time.ParseDuration(valueString)
	return err
}
