package skel

import (
	"github.com/gedex/inflector"
	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
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

	inflectionTitleSingular          = InflectionTitleSingular
	inflectionScreamingSnakeSingular = InflectionScreamingSnakeSingular
	inflectionLowerSingular          = InflectionLowerSingular
	inflectionLowerPlural            = InflectionLowerPlural
)

type Inflector map[string]string

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
	caser := cases.Title(TitlingLanguage)
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
