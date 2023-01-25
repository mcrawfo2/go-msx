// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"github.com/mcrawfo2/jennifer/jen"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type TestDecl string

func (t TestDecl) Generate(out *codegen.Emitter) {
	out.Println("%s", string(t))
}

func testGoEmit(generator Generator) string {
	emitter := NewGoEmitter()
	generator.Generate(emitter)
	return emitter.String()
}

func TestDecls_Generate(t *testing.T) {
	tests := []struct {
		name string
		d    Generator
		want string
	}{
		{
			name: "SingleDecl",
			d: Decls{
				TestDecl("SingleDecl"),
			},
			want: "SingleDecl\n\n",
		},
		{
			name: "TwoDecls",
			d: Decls{
				TestDecl("OneDecl"),
				TestDecl("TwoDecl"),
			},
			want: "OneDecl\n\nTwoDecl\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testGoEmit(tt.d)

			assert.True(t,
				tt.want == got,
				testhelpers.Diff(tt.want, got))
		})
	}
}

func TestDecls_Render(t *testing.T) {
	tests := []struct {
		name string
		d    Decls
		want string
	}{
		{
			name: "SingleDecl",
			d: Decls{
				TestDecl("SingleDecl"),
			},
			want: "\tSingleDecl\n\n",
		},
		{
			name: "TwoDecls",
			d: Decls{
				TestDecl("OneDecl"),
				TestDecl("TwoDecl"),
			},
			want: "\tOneDecl\n\n\tTwoDecl\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.Render(1)
			assert.True(t,
				tt.want == got,
				testhelpers.Diff(tt.want, got))
		})
	}
}

func testStringTransformer(pfx string, sfx string) func(string) string {
	return func(a string) string {
		return pfx + a + sfx
	}
}

func TestTransformers_Transform(t *testing.T) {
	tests := []struct {
		name   string
		t      Transformers
		target string
		want   string
	}{
		{
			name: "Transform",
			t: Transformers{
				testStringTransformer("a", "x"),
			},
			target: "m",
			want:   "amx",
		},
		{
			name: "Ordering",
			t: Transformers{
				testStringTransformer("a", "x"),
				testStringTransformer("b", "y"),
			},
			target: "m",
			want:   "bamxy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.t.Transform(tt.target))
		})
	}
}

func TestNewDeclSnippet(t *testing.T) {
	snippet := NewGoGeneratorSnippet(
		"name",
		Decls{TestDecl("  content  ")},
		[]codegen.Import{
			ImportContext,
			ImportTypes,
		},
		Transformers{
			strings.TrimSpace,
		})

	assert.Equal(t, "name", snippet.Name)
	assert.Equal(t, []string{"content"}, snippet.Lines)
	assert.Equal(t, []codegen.Import{ImportContext, ImportTypes}, snippet.Imports)

	got := testGoEmit(snippet)
	assert.Equal(t, "content\n", got)
}

func TestNewTextSnippet(t *testing.T) {
	snippet := NewGoTextSnippet(
		"name",
		"  content  ",
		[]codegen.Import{
			ImportContext,
			ImportTypes,
		},
		Transformers{
			strings.TrimSpace,
		})

	assert.Equal(t, "name", snippet.Name)
	assert.Equal(t, []string{"content"}, snippet.Lines)
	assert.Equal(t, []codegen.Import{ImportContext, ImportTypes}, snippet.Imports)

	got := testGoEmit(snippet)
	assert.Equal(t, "content\n", got)
}

func TestNewStatementSnippet(t *testing.T) {
	snippet, err := NewGoStatementSnippet(
		"name",
		jen.Var().
			Id("content").
			Id("string").
			Op("=").
			Qual(PkgUuid, "NewUUID").
			Call().
			Dot("String").
			Call(),
		Transformers{
			strings.TrimSpace,
		})
	assert.NoError(t, err)

	assert.Equal(t, "name", snippet.Name)
	assert.Equal(t, []string{"var content string = uuid.NewUUID().String()"}, snippet.Lines)
	assert.Len(t, snippet.Imports, 1)
	assert.Equal(t, PkgUuid, snippet.Imports[0].QualifiedName)

	got := testGoEmit(snippet)
	assert.Equal(t, "var content string = uuid.NewUUID().String()\n", got)
}

func TestConstants_Generate(t *testing.T) {
	tests := []struct {
		name string
		c    GoConstants
		want string
	}{
		{
			name: "OneConstant",
			c: GoConstants{
				&codegen.Constant{
					Type:  codegen.PrimitiveType{Type: "string"},
					Name:  "single",
					Value: "value",
				},
			},
			want: "const (\n\tsingle string = \"value\"\n)\n",
		},
		{
			name: "TwoConstants",
			c: GoConstants{
				&codegen.Constant{
					Type:  codegen.PrimitiveType{Type: "string"},
					Name:  "single",
					Value: "value",
				},
				&codegen.Constant{
					Name:  "double",
					Value: 42,
				},
			},
			want: "const (\n\tsingle string = \"value\"\n\tdouble = 42\n)\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testGoEmit(tt.c)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComment_Generate(t *testing.T) {
	tests := []struct {
		name string
		c    GoComment
		want string
	}{
		{
			name: "Standard",
			c:    GoComment("a single line comment"),
			want: "// a single line comment\n",
		},
		{
			name: "Generate",
			c:    GoComment("go:generate this"),
			want: "//go:generate this\n",
		},
		{
			name: "Conditional",
			c:    GoComment("#if REPOSITORY_COCKROACH"),
			want: "//#if REPOSITORY_COCKROACH\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testGoEmit(tt.c)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSection_AddSnippet(t *testing.T) {
	section := Section[GoSnippet]{Name: "section"}
	section.AddSnippet(GoSnippet{Snippet: Snippet{Name: "snippet"}})
	assert.Len(t, section.Snippets, 1)
}

func TestSection_Generate(t *testing.T) {
	tests := []struct {
		name string
		c    *Section[GoSnippet]
		want string
	}{
		{
			name: "Flat",
			c: &Section[GoSnippet]{
				Name: "section",
				Snippets: []GoSnippet{
					{Snippet: Snippet{Name: "first", Lines: []string{"one line"}}},
					{Snippet: Snippet{Name: "second", Lines: []string{"two line"}}},
				},
			},
			want: "// section\n\none line\n\ntwo line\n\n",
		},
		{
			name: "Nested",
			c: &Section[GoSnippet]{
				Name:     "section",
				Snippets: []GoSnippet{{Snippet: Snippet{Lines: []string{"one line"}}}},
				Sections: Sections[GoSnippet]{{
					Name:     "subsection",
					Snippets: []GoSnippet{{Snippet: Snippet{Lines: []string{"two line"}}}},
				}},
			},
			want: "// section\n\none line\n\n// subsection\n\ntwo line\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testGoEmit(tt.c)
			assert.Equal(t, tt.want, got)
		})
	}

}

func TestSection_Empty(t *testing.T) {
	tests := []struct {
		name string
		c    *Section[GoSnippet]
		want bool
	}{
		{
			name: "FlatEmpty",
			c:    &Section[GoSnippet]{},
			want: true,
		},
		{
			name: "FlatNotEmpty",
			c: &Section[GoSnippet]{
				Name: "section",
				Snippets: []GoSnippet{
					{Snippet: Snippet{Name: "first", Lines: []string{"one line"}}},
				},
			},
			want: false,
		},
		{
			name: "NestedEmpty",
			c: &Section[GoSnippet]{
				Name: "section",
				Sections: Sections[GoSnippet]{{
					Name: "subsection",
				}},
			},
			want: true,
		},
		{
			name: "NestedNotEmpty",
			c: &Section[GoSnippet]{
				Name: "section",
				Sections: Sections[GoSnippet]{{
					Name: "subsection",
					Snippets: []GoSnippet{
						{Snippet: Snippet{Name: "second", Lines: []string{"two line"}}},
					},
				}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.Empty()
			assert.Equal(t, tt.want, got)
		})
	}

}
