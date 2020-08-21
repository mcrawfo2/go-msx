package config

import (
	"context"
	"fmt"
	"testing"
)

func TestPrecedence(t *testing.T) {
	low := map[string]string{
		"without.override": "false",
		"with.override":    "false",
	}

	high := map[string]string{
		"with.override": "true",
	}

	c := NewConfig(
		[]Provider{
			NewStatic("low", low),
			NewStatic("high", high),
		}...,
	)

	if err := c.Load(context.Background()); err != nil {
		t.Error(err)
	}

	without, err := c.Bool("without.override")
	if err != nil {
		t.Error(err)
	}

	if without == true {
		t.Errorf("Setting 'without.override' was true, expected false")
	}

	with, err := c.Bool("with.override")
	if err != nil {
		t.Error(err)
	}

	if with == false {
		t.Errorf("Setting 'with.override' was 'false', expected 'true'")
	}
}

func TestTypeLookups(t *testing.T) {
	settings := map[string]string{
		"string": "some_string",
		"bool":   "true",
		"int":    "1",
		"float":  "1.5",
	}

	c := NewConfig([]Provider{NewStatic("lookups", settings)}...)

	if err := c.Load(context.Background()); err != nil {
		t.Error(err)
	}

	s, err := c.String("string")
	if err != nil {
		t.Error(err)
	}

	if s != "some_string" {
		t.Errorf("String setting was '%s', expected 'some_string'", s)
	}

	b, err := c.Bool("bool")
	if err != nil {
		t.Error(err)
	}

	if b != true {
		t.Errorf("Bool setting was 'false', expected 'true'")
	}

	i, err := c.Int("int")
	if err != nil {
		t.Error(err)
	}

	if i != 1 {
		t.Errorf("Int setting was '%d', expected '1'", i)
	}

	f, err := c.Float("float")
	if err != nil {
		t.Error(err)
	}

	if f != 1.5 {
		t.Errorf("Float setting was '%f', expected '1.5'", f)
	}
}

func TestTypeOrLookups(t *testing.T) {
	c := NewConfig()

	if err := c.Load(context.Background()); err != nil {
		t.Error(err)
	}

	s, err := c.StringOr("string", "some_string")
	if err != nil {
		t.Error(err)
	}

	if s != "some_string" {
		t.Errorf("String setting was '%s', expected 'some_string'", s)
	}

	b, err := c.BoolOr("bool", true)
	if err != nil {
		t.Error(err)
	}

	if b != true {
		t.Errorf("Bool setting was 'false', expected 'true'")
	}

	i, err := c.IntOr("int", 1)
	if err != nil {
		t.Error(err)
	}

	if i != 1 {
		t.Errorf("Int setting was '%d', expected '1'", i)
	}

	f, err := c.FloatOr("float", 1.5)
	if err != nil {
		t.Error(err)
	}

	if f != 1.5 {
		t.Errorf("Float setting was '%f', expected '1.5'", f)
	}
}

func TestValidate(t *testing.T) {
	c := NewConfig()
	c.Validate = func(map[string]string) error {
		return fmt.Errorf("some error")
	}

	if err := c.Load(context.Background()); err == nil {
		t.Errorf("Error was nil")
	}
}

type resolvevaluetest struct {
	resolved map[string]string
	settings map[string]string
	input    string
	expected string
	desc     string
}

var resolveValueTests = []resolvevaluetest{
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "xx", "xx", "plain text"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "", "", "empty value"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "${test.port:9213}", "9210", "value from settings"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "${server.port:9213}", "9211", "value from resolved"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"},
		"%clr(%d{yyyy-MM-dd'T'HH:mm:ss.SSS,UTC}){faint}%clr(%5p)%clr( ${server.port:-} ){magenta}%clr(---){faint}%clr([%15.15t]){faint}%clr(%-40.40logger{39}){cyan}%clr(:){faint[%mdc]%msg%n%ex{full}",
		"%clr(%d{yyyy-MM-dd'T'HH:mm:ss.SSS,UTC}){faint}%clr(%5p)%clr( 9211 ){magenta}%clr(---){faint}%clr([%15.15t]){faint}%clr(%-40.40logger{39}){cyan}%clr(:){faint[%mdc]%msg%n%ex{full}",
		"mixed in string"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "${TEST.pOrT:9213}", "9210", "value from settings case insensitive"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "${TE-ST.pO-r-T:9213}", "9210", "value from settings ignore '-'"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "${nothing.port:9212}", "9212", "value from default field"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "${nothing.port}", "", "not defined variable"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "${nothing.port:}", "", "empty default value"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"}, "${management.server.port:${local.server.port:${server.port:0}}}", "9211", "nested reference"},
	{ map[string]string{"server.port":"9211"},map[string]string{"test.port":"9210"},
		"${management.server.port:${local.server.port:${server.port:0}}}-${management.server.port:${local.server.port:9214}}",
		"9211-9214", "multiple nested references"},
	{ map[string]string{"test.port":"9211"},map[string]string{"test.port":"9210"}, "${management.server.port:${local.server.port:${server.port:9215}}}", "9215", "nested reference and use default value"},

}

func testResolveValueOutput(t *testing.T, cfg *Config, tc *resolvevaluetest) {

	actual := cfg.resolveValue(tc.resolved,tc.settings, tc.input)
	if actual != tc.expected {
		t.Errorf("Actual:'%s', expected: '%s', test: '%s'", actual, tc.expected, tc.desc)
	}
}

func TestResolveValue(t *testing.T) {
	cfg := NewConfig()
	for _, tc := range resolveValueTests {
		testResolveValueOutput(t, cfg,&tc)
	}

	tc := resolvevaluetest{ map[string]string{"server.port":"${test.port}"},
		map[string]string{"test.port":"9210"},
		"${management.server.port:${local.server.port:${server.port:9215}}}",
		"9210",
		"update resolved map"}

	testResolveValueOutput(t, cfg,&tc)
	if v,ok := tc.resolved["server.port"];!ok || v!="9210" {
		t.Errorf("Resolved map not updated! test: '%s'", tc.desc)
	}



	logger.Info("Negative test cases:")


	tc = resolvevaluetest{ map[string]string{"local.server.port":"${management.server.port}"},
		map[string]string{"test.port":"9210"},
		"${management.server.port:${local.server.port:${server.port:9216}}}",
		"",
		"circular variable reference"}
	testResolveValueOutput(t, cfg,&tc)
}
