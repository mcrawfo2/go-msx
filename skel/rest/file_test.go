package rest

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

func testEmit(decl codegen.Decl) string {
	emitter := codegen.NewEmitter(maxLineLength)
	decl.Generate(emitter)
	return emitter.String()
}

func TestDecls_Generate(t *testing.T) {
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
			got := testEmit(tt.d)

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

func TestSnippet_GetName(t *testing.T) {
	snippet := Snippet{Name: "snippet"}
	assert.Equal(t, "snippet", snippet.GetName())
}

func TestSnippet_Generate(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			name:  "OneLine",
			lines: []string{"abc"},
			want:  "abc\n",
		},
		{
			name:  "TwoLines",
			lines: []string{"abc", "def"},
			want:  "abc\ndef\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Snippet{Lines: tt.lines}
			got := testEmit(s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewDeclSnippet(t *testing.T) {
	snippet := NewDeclSnippet(
		"section",
		"name",
		TestDecl("  content  "),
		[]codegen.Import{
			importContext,
			importTypes,
		},
		Transformers{
			strings.TrimSpace,
		})

	assert.Equal(t, "section", snippet.Section)
	assert.Equal(t, "name", snippet.Name)
	assert.Equal(t, []string{"content"}, snippet.Lines)
	assert.Equal(t, []codegen.Import{importContext, importTypes}, snippet.Imports)

	got := testEmit(snippet)
	assert.Equal(t, "content\n", got)
}

func TestNewTextSnippet(t *testing.T) {
	snippet := NewTextSnippet(
		"section",
		"name",
		"  content  ",
		[]codegen.Import{
			importContext,
			importTypes,
		},
		Transformers{
			strings.TrimSpace,
		})

	assert.Equal(t, "section", snippet.Section)
	assert.Equal(t, "name", snippet.Name)
	assert.Equal(t, []string{"content"}, snippet.Lines)
	assert.Equal(t, []codegen.Import{importContext, importTypes}, snippet.Imports)

	got := testEmit(snippet)
	assert.Equal(t, "content\n", got)
}

func TestNewStatementSnippet(t *testing.T) {
	snippet, err := NewStatementSnippet(
		"section",
		"name",
		jen.Var().Id("content").Id("string"),
		Transformers{
			strings.TrimSpace,
		})
	assert.NoError(t, err)

	assert.Equal(t, "section", snippet.Section)
	assert.Equal(t, "name", snippet.Name)
	assert.Equal(t, []string{"var content string"}, snippet.Lines)
	assert.Len(t, snippet.Imports, 0)

	got := testEmit(snippet)
	assert.Equal(t, "var content string\n", got)
}

func TestConstants_Generate(t *testing.T) {
	tests := []struct {
		name string
		c    Constants
		want string
	}{
		{
			name: "OneConstant",
			c: Constants{
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
			c: Constants{
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
			got := testEmit(tt.c)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComment_Generate(t *testing.T) {
	tests := []struct {
		name string
		c    Comment
		want string
	}{
		{
			name: "Standard",
			c:    Comment("a single line comment"),
			want: "// a single line comment\n",
		},
		{
			name: "Generate",
			c:    Comment("go:generate this"),
			want: "//go:generate this\n",
		},
		{
			name: "Conditional",
			c:    Comment("#if REPOSITORY_COCKROACH"),
			want: "//#if REPOSITORY_COCKROACH\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testEmit(tt.c)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSection_AddSnippet(t *testing.T) {
	section := Section{Name: "section"}
	section.AddSnippet(Snippet{Section: "section", Name: "snippet"})
	assert.Len(t, section.Snippets, 1)
}

func TestSection_Generate(t *testing.T) {
	tests := []struct {
		name string
		c    *Section
		want string
	}{
		{
			name: "Flat",
			c: &Section{
				Name: "section",
				Snippets: []Snippet{
					{Name: "first", Lines: []string{"one line"}},
					{Name: "second", Lines: []string{"two line"}},
				},
			},
			want: "// section\n\none line\n\ntwo line\n\n",
		},
		{
			name: "Nested",
			c: &Section{
				Name:     "section",
				Snippets: []Snippet{{Lines: []string{"one line"}}},
				Sections: Sections{{
					Name:     "subsection",
					Snippets: []Snippet{{Lines: []string{"two line"}}},
				}},
			},
			want: "// section\n\none line\n\n// subsection\n\ntwo line\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testEmit(tt.c)
			assert.Equal(t, tt.want, got)
		})
	}

}

func TestSection_Empty(t *testing.T) {
	tests := []struct {
		name string
		c    *Section
		want bool
	}{
		{
			name: "FlatEmpty",
			c:    &Section{},
			want: true,
		},
		{
			name: "FlatNotEmpty",
			c: &Section{
				Name: "section",
				Snippets: []Snippet{
					{Name: "first", Lines: []string{"one line"}},
				},
			},
			want: false,
		},
		{
			name: "NestedEmpty",
			c: &Section{
				Name: "section",
				Sections: Sections{{
					Name: "subsection",
				}},
			},
			want: true,
		},
		{
			name: "NestedNotEmpty",
			c: &Section{
				Name: "section",
				Sections: Sections{{
					Name: "subsection",
					Snippets: []Snippet{
						{Name: "second", Lines: []string{"two line"}},
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
