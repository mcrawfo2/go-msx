// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testEmit(generator Generator) string {
	emitter := NewEmitter(skel.FileFormatOther)
	generator.Generate(emitter)
	return emitter.String()
}

func TestNewSnippet(t *testing.T) {
	type args struct {
		section    string
		name       string
		content    string
		transforms Transformers
	}
	tests := []struct {
		name        string
		args        args
		wantSnippet Snippet
	}{
		{
			name: "SingleLine",
			args: args{
				section:    "section",
				name:       "name",
				content:    " single line ",
				transforms: nil,
			},
			wantSnippet: Snippet{
				Name:  "name",
				Lines: []string{" single line "},
			},
		},
		{
			name: "MultiLine",
			args: args{
				section: "section",
				name:    "name",
				content: "oneline\n\n\ttwoline\n",
			},
			wantSnippet: Snippet{
				Name: "name",
				Lines: []string{
					"oneline",
					"",
					"\ttwoline",
					"",
				},
			},
		},
		{
			name: "Empty",
			args: args{
				section: "section",
				name:    "name",
				content: "  \n  \t  \n  \t  \n  ",
			},
			wantSnippet: Snippet{
				Name:  "name",
				Lines: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t,
				tt.wantSnippet,
				NewSnippet(tt.args.name, tt.args.content, tt.args.transforms),
				"NewSnippet(%v, %v, %v)",
				tt.args.name, tt.args.content, tt.args.transforms)
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

func TestFile_FindSection(t *testing.T) {
	type testCase[I NamedGenerator] struct {
		name    string
		f       *File[I]
		section string
		want    *Section[I]
	}
	tests := []testCase[Snippet]{
		{
			name:    "FlatNew",
			f:       &File[Snippet]{},
			section: "First",
			want: &Section[Snippet]{
				Name: "First",
			},
		},
		{
			name: "FlatExists",
			f: &File[Snippet]{
				Sections: Sections[Snippet]{{
					Name: "Exists",
					Snippets: []Snippet{{
						Name: "ExistsSnippet",
					}},
					Sections: nil,
				}},
			},
			section: "Exists",
			want: &Section[Snippet]{
				Name: "Exists",
				Snippets: []Snippet{{
					Name: "ExistsSnippet",
				}},
			},
		},
		{
			name: "DeepNew",
			f: &File[Snippet]{
				Sections: Sections[Snippet]{{
					Name: "First",
				}},
			},
			section: "First/Second",
			want: &Section[Snippet]{
				Name: "Second",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.f.FindSection(tt.section), "FindSection(%v)", tt.section)
		})
	}
}

func TestFile_Render(t *testing.T) {
	type testCase[I NamedGenerator] struct {
		name string
		f    File[I]
		want string
	}
	tests := []testCase[Snippet]{
		{
			name: "General",
			f: File[Snippet]{
				Comment: "General",
				Sections: Sections[Snippet]{{
					Name: "One",
					Snippets: []Snippet{{
						Name:  "First",
						Lines: []string{"one", "two", "three"},
					}},
				}},
			},
			want: "// General\n\n// One\n\none\ntwo\nthree\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.f.Render(), "Render()")
		})
	}
}

func TestTextFile_Render(t *testing.T) {
	type fields struct {
		File *File[Snippet]
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Simple",
			fields: fields{
				File: &File[Snippet]{
					Format:  skel.FileFormatSql,
					Comment: "heading comment",
					Sections: Sections[Snippet]{
						{
							Name: "section",
							Snippets: []Snippet{
								{
									Name: "name",
									Lines: []string{
										"one",
										"two",
										"three",
									},
								},
							},
						},
					},
				},
			},
			want: "-- heading comment\n\n-- section\n\none\ntwo\nthree\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TextFile{
				File: tt.fields.File,
			}
			assert.Equalf(t, tt.want, f.Render(), "Render()")
		})
	}
}
