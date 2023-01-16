package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/davecgh/go-spew/spew"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestValidationFailure_Error(t *testing.T) {
	message := (&ValidationFailure{}).Error()
	assert.Equal(t, "Validation failure", message)
}

func TestValidationFailure_ToPojo(t *testing.T) {
	tests := []struct {
		name    string
		failure ValidationFailure
		want    types.Pojo
	}{
		{
			name:    "Empty",
			failure: ValidationFailure{},
			want:    nil,
		},
		{
			name: "Failures",
			failure: ValidationFailure{
				Path: "/name",
				Failures: []string{
					"minimum length of 1",
				},
			},
			want: types.Pojo{
				".failures": []string{"minimum length of 1"},
			},
		},
		{
			name: "Children",
			failure: ValidationFailure{
				Path:     "/",
				Failures: nil,
				Children: map[string]*ValidationFailure{
					"FirstChild": {
						Path:     "/FirstChild",
						Failures: []string{"first child failure"},
					},
					"SecondChild": {
						Path:     "/SecondChild",
						Failures: []string{"second child failure"},
					},
				},
			},
			want: types.Pojo{
				"FirstChild": types.Pojo{
					".failures": []string{"first child failure"},
				},
				"SecondChild": types.Pojo{
					".failures": []string{"second child failure"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.failure.ToPojo()
			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.Diff(tt.want, got))
		})
	}

}

func TestValidationFailure_Apply(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want *ValidationFailure
	}{
		{
			name: "ChildFailure",
			err: &jsv.ValidationError{
				InstanceLocation: "/Root/FirstChild",
				Message:          "First child failure",
			},
			want: &ValidationFailure{
				Path:     "/Root",
				Failures: nil,
				Children: map[string]*ValidationFailure{
					"FirstChild": {
						Path: "/Root/FirstChild",
						Failures: []string{
							"First child failure",
						},
						Children: make(map[string]*ValidationFailure),
					},
				},
			},
		},
		{
			name: "NestedErrors",
			err: types.ErrorMap{
				"firstField": types.ErrorMap{
					"1": types.ErrorMap{
						"secondField": validation.Errors{
							"2": types.ErrorList{
								errors.New("Missing required value"),
							},
						},
					},
				},
			},
			want: &ValidationFailure{
				Path: "/Root",
				Children: map[string]*ValidationFailure{
					"firstField": {
						Path: "/Root/firstField",
						Children: map[string]*ValidationFailure{
							"1": {
								Path: "/Root/firstField/1",
								Children: map[string]*ValidationFailure{
									"secondField": {
										Path: "/Root/firstField/1/secondField",
										Children: map[string]*ValidationFailure{
											"2": {
												Path: "/Root/firstField/1/secondField/2",
												Failures: []string{
													"Missing required value",
												},
												Children: map[string]*ValidationFailure{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	var spewConfig = spew.ConfigState{
		Indent:                  " ",
		DisablePointerAddresses: true,
		DisableMethods:          true,
		DisableCapacities:       true,
		SortKeys:                true,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failure := NewValidationFailure("Root")
			got := failure.Apply(tt.err)

			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.DiffWithConfig(tt.want, got, spewConfig))
		})
	}
}
