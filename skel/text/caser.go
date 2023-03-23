// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var TitlingLanguage = language.English

func NewTitleCaser() cases.Caser {
	return cases.Title(TitlingLanguage)
}
