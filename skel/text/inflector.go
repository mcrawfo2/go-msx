// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"github.com/gedex/inflector"
	"github.com/iancoleman/strcase"
	"strings"
)

const (
	InflectionTitleSingular          = "Title Singular"
	InflectionTitlePlural            = "Title Plural"
	InflectionUpperCamelSingular     = "UpperCamelSingular"
	InflectionUpperCamelPlural       = "UpperCamelPlural"
	InflectionLowerCamelSingular     = "lowerCamelSingular"
	InflectionLowerCamelPlural       = "lowerCamelPlural"
	InflectionLowerSnakeSingular     = "lower_snake_singular"
	InflectionScreamingSnakeSingular = "SCREAMING_SNAKE_SINGULAR"
	InflectionScreamingSnakePlural   = "SCREAMING_SNAKE_PLURAL"
	InflectionLowerSingular          = "lowersingular"
	InflectionLowerPlural            = "lowerplural"
)

type Inflector map[string]string

func (i Inflector) Invert() Inflector {
	result := make(map[string]string)
	for k, v := range i {
		result[v] = k
	}
	return result
}

func (i Inflector) Inflect(target string) string {
	for k, v := range i {
		target = strings.ReplaceAll(target, k, v)
	}
	return target
}

func (i Inflector) Inflections() map[string]string {
	return i
}

func NewInflector(title string) Inflector {
	caser := NewTitleCaser()
	titleSingular := caser.String(inflector.Singularize(title))
	titlePlural := caser.String(inflector.Pluralize(titleSingular))
	upperCamelSingular := strcase.ToCamel(titleSingular)
	upperCamelPlural := strcase.ToCamel(titlePlural)
	lowerCamelSingular := strcase.ToLowerCamel(titleSingular)
	lowerCamelPlural := strcase.ToLowerCamel(titlePlural)
	lowerSingular := strings.ToLower(lowerCamelSingular)
	lowerPlural := strings.ToLower(lowerCamelPlural)
	lowerSnakeSingular := strcase.ToSnake(titleSingular)
	screamingSnakeSingular := strcase.ToScreamingSnake(titleSingular)
	screamingSnakePlural := strcase.ToScreamingSnake(titlePlural)

	return map[string]string{
		InflectionTitleSingular:          titleSingular,
		InflectionTitlePlural:            titlePlural,
		InflectionUpperCamelSingular:     upperCamelSingular,
		InflectionUpperCamelPlural:       upperCamelPlural,
		InflectionLowerCamelSingular:     lowerCamelSingular,
		InflectionLowerCamelPlural:       lowerCamelPlural,
		InflectionLowerSingular:          lowerSingular,
		InflectionLowerPlural:            lowerPlural,
		InflectionLowerSnakeSingular:     lowerSnakeSingular,
		InflectionScreamingSnakeSingular: screamingSnakeSingular,
		InflectionScreamingSnakePlural:   screamingSnakePlural,
	}
}
