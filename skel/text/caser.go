package text

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var TitlingLanguage = language.English

func NewTitleCaser() cases.Caser {
	return cases.Title(TitlingLanguage)
}
