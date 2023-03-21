// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_inflect(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Simple",
			args: args{"employee"},
			want: map[string]string{
				"SCREAMING_SNAKE_PLURAL":   "EMPLOYEES",
				"SCREAMING_SNAKE_SINGULAR": "EMPLOYEE",
				"Title Singular":           "Employee",
				"Title Plural":             "Employees",
				"UpperCamelSingular":       "Employee",
				"UpperCamelPlural":         "Employees",
				"lowerCamelSingular":       "employee",
				"lowerCamelPlural":         "employees",
				"lowersingular":            "employee",
				"lower_snake_singular":     "employee",
				"lowerplural":              "employees",
			},
		},
		{
			name: "Multiple Words",
			args: args{"employee handbook"},
			want: map[string]string{
				"SCREAMING_SNAKE_SINGULAR": "EMPLOYEE_HANDBOOK",
				"SCREAMING_SNAKE_PLURAL":   "EMPLOYEE_HANDBOOKS",
				"Title Singular":           "Employee Handbook",
				"Title Plural":             "Employee Handbooks",
				"UpperCamelSingular":       "EmployeeHandbook",
				"UpperCamelPlural":         "EmployeeHandbooks",
				"lowerCamelSingular":       "employeeHandbook",
				"lowerCamelPlural":         "employeeHandbooks",
				"lowersingular":            "employeehandbook",
				"lower_snake_singular":     "employee_handbook",
				"lowerplural":              "employeehandbooks",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewInflector(tt.args.title).Inflections()
			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.Diff(tt.want, got))
		})
	}
}
