// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package main

import (
	"bytes"
	"fmt"
	"testing"
)

const (
	ls1 = `this is a file
this is a line in the file
we should be able to compare 1 and 2 and find them the same`
	ls2 = `this is a file
this is a line in the file
we should be able to compare 1 and 2 and find them the same`
	ls3 = `this is a file
this is a line in the file
compare 1 and 3: it is different`
	ls4 = `// Code generated by mockery v2.21.1. DO NOT EDIT.

package complianceevent

import (
	context "context"

	api "cto-github.cisco.com/NFV-BU/someservice/internal/stream/complianceevent/api"

	mock "github.com/stretchr/testify/mock"
)`
	ls5 = `// Code generated by mockery v2.21.1. DO NOT EDIT.  if all is well, this will not be cmpd

package complianceevent

import (
	context "context"

	api "cto-github.cisco.com/NFV-BU/someservice/internal/stream/complianceevent/api"

	mock "github.com/stretchr/testify/mock"
)`
	ls6 = `only line`
	ls7 = `first line
only line
third line`
)

type ltest struct {
	name   string
	a1, a2 string
	rsame  bool
	reg    string
	err    error
}

var ltests = []ltest{
	{name: "same", a1: ls1, a2: ls2, rsame: true},
	{name: "not same", a1: ls1, a2: ls3, rsame: false},
	{name: "same with reg", a1: ls1, a2: ls3, rsame: true, reg: "compare"},
	{name: "same with reg/neg", a1: ls1, a2: ls3, rsame: false, reg: ""},
	{name: "still same 1/3", a1: ls1, a2: ls2, rsame: true, reg: "should not match anything"},
	{name: "where diff failed", a1: ls4, a2: ls5, rsame: true, reg: "DO NOT EDIT"},
	{name: "weak 2 item", a1: ls6, a2: ls7, rsame: true, reg: "(first)|(third)"},
	{name: "weak 2 item/neg", a1: ls6, a2: ls7, rsame: false},
	{name: "weak 2 item/neg2", a1: ls6, a2: ls7, rsame: false, reg: "first"},
}

func TestSameWithLines(t *testing.T) {

	for _, tst := range ltests {

		l1 := []byte(tst.a1)
		l2 := []byte(tst.a2)

		s1, d, err1 := same(bytes.NewReader(l1), bytes.NewReader(l2), RegexString(tst.reg))
		if err1 != tst.err {
			t.Errorf("%s: error should be: %s, was %s", tst.name, tst.err, err1)
		}
		if s1 != tst.rsame || err1 != nil {
			fmt.Printf(">---%s\n---\n%s\n<---", l1, l2)
			t.Errorf("%s: should be same: %t, was: %t, diff: %s)", tst.name, tst.rsame, s1, d)
		}

	}

}
