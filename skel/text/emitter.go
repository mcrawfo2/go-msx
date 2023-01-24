// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
)

type Generator interface {
	// Generate outputs this component to the specified Emitter
	Generate(out Emitter)
}

type Namer interface {
	// GetName returns the name of this component
	GetName() string
}

type Emptier interface {
	Empty() bool
}

type NamedGenerator interface {
	Namer
	Generator
	Emptier
}

// Generators are a sequence of Generator instances
type Generators []Generator

// Generate outputs the sequence of components, separated by newline
func (g Generators) Generate(out Emitter) {
	for _, generator := range g {
		generator.Generate(out)
		out.Newline()
	}
}

// Emitter is a general interface for gently-formatted text output
type Emitter interface {
	String() string
	Bytes() []byte
	Indent(n int)
	Comment(s string)
	Print(format string, args ...interface{})
	Println(format string, args ...interface{})
	Newline()
	Raw() *codegen.Emitter
}

// textEmitter implements Emitter for non-go languages
type textEmitter struct {
	format skel.FileFormat
	*codegen.Emitter
}

// Raw returns the underlying codegen.Emitter
func (e *textEmitter) Raw() *codegen.Emitter {
	return e.Emitter
}

// Comment outputs a single line comment using the
func (e *textEmitter) Comment(s string) {
	startMarker, endMarker := skel.CommentMarkers(e.format)
	e.Print("%s %s %s", startMarker, s, endMarker)
}

func NewEmitter(format skel.FileFormat) Emitter {
	return &textEmitter{
		Emitter: codegen.NewEmitter(MaxLineLength),
		format:  format,
	}
}

type goEmitter struct {
	*codegen.Emitter
}

func (e *goEmitter) Raw() *codegen.Emitter {
	return e.Emitter
}

func NewGoEmitter() Emitter {
	return &goEmitter{
		Emitter: codegen.NewEmitter(MaxLineLength),
	}
}
