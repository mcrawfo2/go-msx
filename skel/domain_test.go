package skel

import (
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
				"Title Singular":     "Employee",
				"Title Plural":       "Employees",
				"UpperCamelSingular": "Employee",
				"UpperCamelPlural":   "Employees",
				"lowerCamelSingular": "employee",
				"lowerCamelPlural":   "employees",
				"lowersingular":      "employee",
				"lowerplural":        "employees",
			},
		},
		{
			name: "Multiple Words",
			args: args{"employee handbook"},
			want: map[string]string{
				"Title Singular":     "Employee Handbook",
				"Title Plural":       "Employee Handbooks",
				"UpperCamelSingular": "EmployeeHandbook",
				"UpperCamelPlural":   "EmployeeHandbooks",
				"lowerCamelSingular": "employeeHandbook",
				"lowerCamelPlural":   "employeeHandbooks",
				"lowersingular":      "employeehandbook",
				"lowerplural":        "employeehandbooks",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inflect(tt.args.title); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inflect() = %v, want %v", got, tt.want)
			}
		})
	}
}
