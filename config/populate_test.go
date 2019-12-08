package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProperties_Populate_CaseInsensitive(t *testing.T) {
	p := PartialConfig{
		local: map[string]string{
			"arg":       "something",
			"prop.test": "42",
		},
	}

	type TestStruct struct {
		Arg  string
		Prop struct {
			Test int `config:"Test,default=40"`
		} `config:"Prop"`
	}

	s := &TestStruct{}

	err := p.Populate(s)
	assert.Nil(t, err)

	assert.Equal(t, "something", s.Arg)
	assert.Equal(t, int(42), s.Prop.Test)
}

func TestProperties_Populate_Slice(t *testing.T) {
	p := PartialConfig{
		local: map[string]string{
			"whitelist[0]": "value0",
			"whitelist[1]": "value1",
		},
	}

	type TestStruct struct {
		Whitelist []string
	}

	s := &TestStruct{}

	err := p.Populate(s)
	assert.Nil(t, err)

	assert.Len(t, s.Whitelist, 2)
	assert.Contains(t, s.Whitelist, "value0")
	assert.Contains(t, s.Whitelist, "value1")
}
