// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ddl

import "strings"

type OptionsQueryPartBuilder struct{}

func (b *OptionsQueryPartBuilder) Options(optionsMaps ...map[string]string) string {
	sb := new(strings.Builder)
	sb.WriteRune('{')

	n := 0
	for _, options := range optionsMaps {
		for k, v := range options {
			if n > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("'")
			sb.WriteString(k)
			sb.WriteString("': '")
			sb.WriteString(v)
			sb.WriteString("'")
			n++
		}
	}

	sb.WriteRune('}')
	return sb.String()
}
