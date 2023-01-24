// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"github.com/lithammer/dedent"
	"strings"
)

type SqlSnippet struct {
	Snippet
}

type SqlComment string

func (s SqlComment) Generate(out Emitter) {
	out.Println("-- %s", string(s))
}

func NewSqlTextSnippet(section, name, content string, transforms Transformers) (snippet SqlSnippet) {
	content = transforms.Transform(content)

	var lines []string
	if strings.TrimSpace(content) != "" {
		lines = strings.Split(content, "\n")
	}

	return SqlSnippet{
		Snippet: Snippet{
			Section: section,
			Name:    name,
			Lines:   lines,
		},
	}
}

type SqlFile struct {
	*File[SqlSnippet]
}

func (f *SqlFile) AddSnippet(snippet SqlSnippet) {
	logger.Infof("  📃 Adding snippet %q to section %q", snippet.Name, snippet.Section)
	section := f.FindSection(snippet.Section)
	section.AddSnippet(snippet)
}

func (f *SqlFile) AddNewText(path, name, content string) error {
	f.AddSnippet(NewSqlTextSnippet(path, name, content, Transformers{
		f.Inflector.Inflect,
		dedent.Dedent,
		strings.TrimSpace,
	}))
	return nil
}

func (f *SqlFile) Render() string {
	out := NewEmitter(f.Format)

	if f.Comment != "" {
		SqlComment(f.Comment).Generate(out)
		out.Newline()
	}

	f.Sections.Generate(out)

	return out.String()
}
