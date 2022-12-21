// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
	"testing"
)

type TestStructA struct {
	B TestStructB `json:"b"`
}

type TestStructB struct {
	C TestEnum `json:"c"`
}

type TestEnum string

func (t TestEnum) Enum() []interface{} {
	return []interface{}{
		"A", "B", "C",
	}
}

func TestRegistrySpec(t *testing.T) {
	got := RegistrySpec()
	assert.NotNil(t, got)
}

func TestLookupSchema(t *testing.T) {
	_, err := Reflect(TestStructA{})
	assert.NoError(t, err)

	schema, exists := LookupSchema("asyncapi.TestStructA")
	assert.True(t, exists)
	assert.True(t, schema.HasType(jsonschema.Object))
}

func TestReflect(t *testing.T) {
	schema, err := Reflect(TestStructA{})
	assert.NoError(t, err)
	assert.NotNil(t, schema.Ref)

	schema, err = Reflect(TestEnum(""))
	assert.NoError(t, err)
	assert.NotNil(t, schema.Ref)
}
