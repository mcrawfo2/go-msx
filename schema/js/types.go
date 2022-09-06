// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js

import (
	"github.com/swaggest/jsonschema-go"
)

const (
	FormatDateTime = "date-time"
	FormatUuid     = "uuid"
	FormatDuration = "duration"
)

type DefNameExposer interface {
	JSONSchemaDefName() string
}

type StringFormatTime struct{}

func (s StringFormatTime) JSONSchemaDefName() string {
	return "Time"
}

func (s StringFormatTime) JSONSchema() (jsonschema.Schema, error) {
	return NewSchemaPtr(jsonschema.String).
		WithID("Time").
		WithTitle("Time").
		WithFormat(FormatDateTime).
		WithPattern(`^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\.[0-9]+)?(([Zz])|([\+|\-]([01][0-9]|2[0-3]):[0-5][0-9]))$`).
		WithExamples("1995-12-17T03:24:56.778899Z").
		JSONSchema()
}

type StringFormatUuid struct{}

func (s StringFormatUuid) JSONSchemaDefName() string {
	return "UUID"
}

func (s StringFormatUuid) JSONSchema() (jsonschema.Schema, error) {
	return NewSchemaPtr(jsonschema.String).
		WithID("UUID").
		WithTitle("UUID").
		WithFormat(FormatUuid).
		WithPattern(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`).
		WithExamples("123e4567-e89b-12d3-a456-426614174000").
		JSONSchema()
}

type StringFormatDuration struct{}

func (s StringFormatDuration) JSONSchemaDefName() string {
	return "Duration"
}

func (s StringFormatDuration) JSONSchema() (jsonschema.Schema, error) {
	return NewSchemaPtr(jsonschema.String).
		WithID("Duration").
		WithTitle("Duration").
		WithFormat(FormatDuration).
		WithPattern(`^(\d+(\.\d+)?h)?(\d+(\.\d+)m)?(\d+(\.\d+)?s)?(\d+(\.\d+)?ms)?(\d+(\.\d+)?us)?(\d+ns)?$`).
		WithExamples("30s", "10m", "1h", "15d").
		JSONSchema()
}
