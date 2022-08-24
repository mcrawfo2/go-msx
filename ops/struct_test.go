// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"fmt"
	"reflect"
	"testing"
)

type testStructFieldVisitor struct {
	Stack []int
	Found map[string]string
}

func newTestStructFieldVisitor() *testStructFieldVisitor {
	return &testStructFieldVisitor{
		Stack: []int{0},
		Found: make(map[string]string),
	}
}

func (t *testStructFieldVisitor) addFound(name string) {
	key := t.indices()
	t.Found[key] = name
}

func (t *testStructFieldVisitor) indices() string {
	return fmt.Sprintf("%v", t.Stack)
}

func (t *testStructFieldVisitor) increment() {
	lastIndex := t.Stack[len(t.Stack)-1]
	t.Stack[len(t.Stack)-1] = lastIndex + 1
}

func (t *testStructFieldVisitor) VisitField(f reflect.StructField) error {
	t.addFound(f.Name)
	t.increment()
	return nil
}

func (t *testStructFieldVisitor) EnterAnonymousStructField(f reflect.StructField) {
	// push
	t.Stack = append(t.Stack, 0)
}

func (t *testStructFieldVisitor) ExitAnonymousStructField(f reflect.StructField) {
	// pop
	t.Stack = t.Stack[:len(t.Stack)-1]
	t.increment()
}

type testFlatStruct struct {
	SimpleField string `req:"body"`
}

type testNestedStruct struct {
	testFlatStruct
	AnotherField testFlatStruct `req:"body"`
}

type testRecursiveStruct struct {
	*testRecursiveStruct
}

type testEmptyStruct struct{}

type testParallelStruct struct {
	A testEmptyStruct `req:"body"`
	B testEmptyStruct `req:"body"`
}

func TestWalkStruct(t *testing.T) {
	tests := []struct {
		name      string
		st        reflect.Type
		wantFound map[string]string
		wantErr   bool
	}{
		{
			name: "Simple",
			st:   reflect.TypeOf(testFlatStruct{}),
			wantFound: map[string]string{
				"[0]": "SimpleField",
			},
		},
		{
			name: "Nested",
			st:   reflect.TypeOf(testNestedStruct{}),
			wantFound: map[string]string{
				"[0 0]": "SimpleField",
				"[1]":   "AnotherField",
			},
		},
		{
			name:    "Recursive",
			st:      reflect.TypeOf(testRecursiveStruct{}),
			wantErr: true,
		},
		{
			name: "Parallel",
			st:   reflect.TypeOf(testParallelStruct{}),
			wantFound: map[string]string{
				"[0]": "A",
				"[1]": "B",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := newTestStructFieldVisitor()
			if err := WalkStruct(tt.st, visitor); (err != nil) != tt.wantErr {
				t.Errorf("WalkStruct() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && !reflect.DeepEqual(visitor.Found, tt.wantFound) {
				t.Errorf("Incorrect result found.\n%s", testhelpers.Diff(tt.wantFound, visitor.Found))
			}
		})
	}
}
