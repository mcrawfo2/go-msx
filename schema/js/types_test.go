// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
	"reflect"
	"testing"
)

func TestStringFormatDuration_JSONSchema(t *testing.T) {
	want := jsonschema.Schema{
		ID:    types.NewStringPtr("Duration"),
		Title: types.NewStringPtr("Duration"),
		Examples: []interface{}{
			"30s", "10m", "1h5m", "15d",
		},
		Pattern: types.NewStringPtr("^(\\d+(\\.\\d+)?h)?(\\d+(\\.\\d+)?m)?(\\d+(\\.\\d+)?s)?(\\d+(\\.\\d+)?ms)?(\\d+(\\.\\d+)?us)?(\\d+ns)?$"),
		Type:    NewType(jsonschema.String),
		Format:  types.NewStringPtr("duration"),
	}
	got, err := StringFormatDuration{}.JSONSchema()
	assert.NoError(t, err)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestStringFormatDuration_JSONSchemaDefName(t *testing.T) {
	want := "Duration"
	got := StringFormatDuration{}.JSONSchemaDefName()
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestStringFormatTime_JSONSchema(t *testing.T) {
	want := jsonschema.Schema{
		ID:    types.NewStringPtr("Time"),
		Title: types.NewStringPtr("Time"),
		Examples: []interface{}{
			"1995-12-17T03:24:56.778899Z",
		},
		Pattern: types.NewStringPtr(`^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\.[0-9]+)?(([Zz])|([\+|\-]([01][0-9]|2[0-3]):[0-5][0-9]))$`),
		Type:    NewType(jsonschema.String),
		Format:  types.NewStringPtr("date-time"),
	}
	got, err := StringFormatTime{}.JSONSchema()
	assert.NoError(t, err)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestStringFormatTime_JSONSchemaDefName(t *testing.T) {
	want := "Time"
	got := StringFormatTime{}.JSONSchemaDefName()
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestStringFormatUuid_JSONSchema(t *testing.T) {
	want := jsonschema.Schema{
		ID:    types.NewStringPtr("UUID"),
		Title: types.NewStringPtr("UUID"),
		Examples: []interface{}{
			"123e4567-e89b-12d3-a456-426614174000",
		},
		Pattern: types.NewStringPtr(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`),
		Type:    NewType(jsonschema.String),
		Format:  types.NewStringPtr("uuid"),
	}
	got, err := StringFormatUuid{}.JSONSchema()
	assert.NoError(t, err)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestStringFormatUuid_JSONSchemaDefName(t *testing.T) {
	want := "UUID"
	got := StringFormatUuid{}.JSONSchemaDefName()
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}
