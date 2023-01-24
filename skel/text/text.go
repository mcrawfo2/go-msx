// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"strings"
)

const MaxLineLength = 120

// Snippet stores a named series of text lines attached to a File section.
type Snippet struct {
	Section string
	Name    string
	Lines   []string
}

func (s Snippet) GetName() string {
	return s.Name
}

// Generate outputs the snippet to the supplied Emitter
func (s Snippet) Generate(out Emitter) {
	for _, line := range s.Lines {
		out.Println("%s", line)
	}
}

func (s Snippet) Empty() bool {
	return len(s.Lines) == 0
}

// Section is a portion of a File
type Section[I NamedGenerator] struct {
	Name     string
	Snippets []I
	Sections Sections[I]
}

// AddSnippet adds a snippet to the section
func (s *Section[I]) AddSnippet(snippet I) {
	s.Snippets = append(s.Snippets, snippet)
}

// Generate outputs the section to the supplied Emitter
func (s *Section[I]) Generate(out Emitter) {
	if s.Empty() {
		return
	}

	out.Comment(s.Name)
	out.Newline()

	for _, snippet := range s.Snippets {
		snippet.Generate(out)
		out.Newline()
	}

	s.Sections.Generate(out)
}

// Empty returns true if the section contains no snippets and only empty sections
func (s *Section[I]) Empty() bool {
	content := false
	for _, snip := range s.Snippets {
		content = content || !snip.Empty()
		if content {
			break
		}
	}
	if !content {
		for _, sub := range s.Sections {
			content = content || !sub.Empty()
			if content {
				break
			}
		}
	}
	return !content
}

// Sections is a sequence of Section instances
type Sections[I NamedGenerator] []*Section[I]

// Generate outputs the sections to the supplied Emitter
func (s Sections[I]) Generate(out Emitter) {
	for _, section := range s {
		section.Generate(out)
	}
}

// WithSection adds a section to the Sections sequence
func (s Sections[I]) WithSection(name string) (Sections[I], *Section[I]) {
	for _, section := range s {
		if section.Name == name {
			return s, section
		}
	}

	section := &Section[I]{Name: name}
	return append(s, section), section
}

// NewSections creates a new sequence of Section instances
func NewSections[I NamedGenerator](sections ...any) Sections[I] {
	var results Sections[I]
	for _, section := range sections {
		switch st := section.(type) {
		case string:
			results = append(results, &Section[I]{
				Name: st,
			})

		case *Section[I]:
			results = append(results, st)
		}
	}
	return results
}

type File[I NamedGenerator] struct {
	Comment   string
	Sections  Sections[I]
	Inflector skel.Inflector
	Format    skel.FileFormat
}

func (f *File[I]) FileFormat() skel.FileFormat {
	return f.Format
}

func (f *File[I]) FindSection(section string) *Section[I] {
	if section == "" {
		return nil
	}

	sectionPath := strings.Split(section, "/")

	var prevSection *Section[I]
	f.Sections, prevSection = f.Sections.WithSection(sectionPath[0])
	sectionPath = sectionPath[1:]

	for len(sectionPath) > 0 {
		var nextSection *Section[I]
		prevSection.Sections, nextSection = prevSection.Sections.WithSection(sectionPath[0])
		prevSection = nextSection
		sectionPath = sectionPath[1:]
	}

	return prevSection
}

func (f *File[I]) Render() string {
	out := NewEmitter(f.Format)

	if f.Comment != "" {
		out.Comment(f.Comment)
		out.Newline()
	}

	f.Sections.Generate(out)

	return out.String()
}
