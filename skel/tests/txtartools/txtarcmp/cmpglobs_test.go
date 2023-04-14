// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"gopkg.in/pipe.v2"
	"strings"
	"testing"
)

type Test struct {
	Name     string
	globby   string
	filename string
	wantErr  bool
	ignores  RegexString
	wantKind clienv.Ix
}

var tests = []Test{
	{Name: "empty", globby: "", filename: "",
		wantKind: noMatchIx, wantErr: false},
	{Name: "exists 1", globby: "exists:**", filename: "anything",
		wantKind: existsIx, wantErr: false},
	{Name: "same 1", globby: "same:**", filename: "weasels",
		wantKind: sameIx, wantErr: false},
	{Name: "notexists 1", globby: "notexists:**", filename: "banana",
		wantKind: notexistsIx, wantErr: false},
	{Name: "notsame 1", globby: "notsame:**", filename: "fruitfly",
		wantKind: notsameIx, wantErr: false},
	{Name: "exists-2", globby: "same:*rat notexists:splung notsame:splorg exists:**", filename: "anything",
		wantKind: existsIx, wantErr: false},
	{Name: "notsame-2", globby: "same:*rat notexists:*splung*  exists:plain notsame:splorg", filename: "splorg",
		wantKind: notsameIx, wantErr: false},
	{Name: "ignorelines1", globby: "same:** ignorelines:*.go:something*", filename: "splorg.go",
		wantKind: sameIx, wantErr: false, ignores: "(something*)"},
}

func TestGlobs(t *testing.T) {

	for _, this := range tests {

		globs, gotErr := ParseGlobs(strings.Split(this.globby, " "))

		_, matchedKind, err := matchToGlobs(globs, this.filename)
		if gotErr != nil && !this.wantErr {
			t.Errorf("error matching globs: %v\n", err)

		}

		if matchedKind != this.wantKind {
			t.Errorf("%s kind mismatch matchCmpGlobs(%q, %q)\ngot %s\nwant %s", this.Name, globs,
				this.filename, ix2prefix[matchedKind], ix2prefix[this.wantKind])
		}

		ig := globs.Ignores(this.filename)
		if err != nil {

		}
		if ig != this.ignores {
			t.Errorf("%s ignores mismatch (%q, %q)\ngot %s\nwant %s", this.Name, globs,
				this.filename, ig, this.ignores)
		}

	}

}

const globz = `same:** ignorelines:*.og:.*EDIT.*`

func TestIgnoreLines(t *testing.T) {
	// this is an integration test
	t.SkipNow()

	args := []string{"-debug", "-cmpglobs=\"" + globz + "\"", "mockery.txtar", "."}

	cmp := pipe.Script(
		pipe.Exec("txtarcmp", args...),
	)

	res, err := pipe.CombinedOutput(cmp)
	if err != nil {
		t.Errorf("ERROR doing:\ntxtarcmp %s\nerror is:%s\nresult is:%s\n", strings.Join(args, " "), err, res)
	}

}

func TestIgnoreLines2(t *testing.T) {
	// this is an integration test
	t.SkipNow()

	args := []string{"-debug", "-cmpglobs=\"" + globz + "\"", "mockery2.txtar", "."}

	cmp := pipe.Script(
		pipe.Exec("txtarcmp", args...),
	)

	res, err := pipe.CombinedOutput(cmp)
	if err != nil {
		t.Errorf("ERROR doing:\ntxtarcmp %s\nerror is:%s\nresult is:%s\n", strings.Join(args, " "), err, res)
	}

}

func TestIgnoreLines3(t *testing.T) {
	// this is an integration test
	t.SkipNow()

	args := []string{"-debug", "-cmpglobs=\"" + globz + "\"", "mockery3.txtar", "."}

	cmp := pipe.Script(
		pipe.Exec("txtarcmp", args...),
	)

	res, err := pipe.CombinedOutput(cmp)
	if err == nil {
		t.Errorf("ERROR doing:\ntxtarcmp %s\nerror is:%s\nresult is:%s\n", strings.Join(args, " "), err, res)
	}

}
