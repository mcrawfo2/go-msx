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
			if got := inflect(tt.args.title); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inflect() = %v, want %v", got, tt.want)
			}
		})
	}
}
