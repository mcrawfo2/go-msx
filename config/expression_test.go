// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"reflect"
	"testing"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		wantExpr expression
		wantErr  bool
	}{
		{
			name:     "Literal",
			value:    "a literal value",
			wantExpr: literalExpression("a literal value"),
			wantErr:  false,
		},
		{
			name:  "Concatenated",
			value: "a ${literal} value",
			wantExpr: concatenateExpression{
				Parts: []expression{
					literalExpression("a "),
					variableExpression{
						Name: "literal",
					},
					literalExpression(" value"),
				},
			},
			wantErr: false,
		},
		{
			name:  "Variable",
			value: "${variable}",
			wantExpr: variableExpression{
				Name: "variable",
			},
			wantErr: false,
		},
		{
			name:  "VariableWithDefault",
			value: "${variable:default}",
			wantExpr: variableExpression{
				Name:    "variable",
				Default: literalExpression("default"),
			},
			wantErr: false,
		},
		{
			name:  "VariableWithRequiredKey",
			value: "${variable:?}",
			wantExpr: variableExpression{
				Name:        "variable",
				RequiredKey: true,
			},
			wantErr: false,
		},
		{
			name:  "VariableWithRequiredValue",
			value: "${variable:!}",
			wantExpr: variableExpression{
				Name:          "variable",
				RequiredKey:   true,
				RequiredValue: true,
			},
			wantErr: false,
		},
		{
			name:  "NestedVariable",
			value: "${variable:a ${nested:b ${crested} c} default}",
			wantExpr: variableExpression{
				Name: "variable",
				Default: concatenateExpression{
					Parts: []expression{
						literalExpression("a "),
						variableExpression{
							Name: "nested",
							Default: concatenateExpression{
								Parts: []expression{
									literalExpression("b "),
									variableExpression{
										Name: "crested",
									},
									literalExpression(" c"),
								},
							},
						},
						literalExpression(" default"),
					},
				},
			},
		},
		{
			name:    "UnmatchedBraces",
			value:   "${value",
			wantErr: true,
		},
		{
			name:    "UnmatchedBracesInDefault",
			value:   "${value:",
			wantErr: true,
		},
		{
			name:    "DeeplyUnmatchedBraces",
			value:   "${a:${b:${c:d}",
			wantErr: true,
		},
		{
			name:    "DynamicVariableName",
			value:   "${a${b}c:d}",
			wantErr: true,
		},
		{
			name:  "Log",
			value: "%clr(%d{yyyy-MM-dd'T'HH:mm:ss.SSS,UTC}){faint} %clr(%5p) %clr(${PID:- }){magenta} %clr(---){faint} %clr([%15.15t]){faint} %clr(%-40.40logger{39}){cyan} %clr(:){faint} [%mdc] %msg%n%ex{full}",
			wantExpr: concatenateExpression{
				Parts: []expression{
					literalExpression("%clr(%d{yyyy-MM-dd'T'HH:mm:ss.SSS,UTC}){faint} %clr(%5p) %clr("),
					variableExpression{
						Name:    "PID",
						Default: literalExpression("- "),
					},
					literalExpression("){magenta} %clr(---){faint} %clr([%15.15t]){faint} %clr(%-40.40logger{39}){cyan} %clr(:){faint} [%mdc] %msg%n%ex{full}"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExpr, err := parseExpression(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotExpr, tt.wantExpr) {
				t.Errorf("parseExpression() gotExpr = %v, want %v", gotExpr, tt.wantExpr)
			}
		})
	}
}
