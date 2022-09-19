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

	schema, exists := LookupSchema("AsyncapiTestStructA")
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
