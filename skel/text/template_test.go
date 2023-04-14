// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestIdentifierLineVisitor_VisitLine(t *testing.T) {
	tpl := NewTemplate("", FileFormatGo, TemplateLanguageSkel, nil)
	markers := tpl.directiveMarkers()

	tests := []struct {
		name        string
		visitor     IdentifierLineVisitor
		line        string
		wantLine    string
		wantOutput  bool
		wantErr     bool
		wantVisitor IdentifierLineVisitor
	}{
		{
			name: "IdentifierDefinition",
			visitor: IdentifierLineVisitor{
				Directive: markers,
			},
			line:       "//#id hello world",
			wantLine:   "//#id hello world",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: IdentifierLineVisitor{
				Definitions: map[string]string{
					"hello": "world",
				},
				Directive: markers,
				marker:    "//#id ",
				identifiers: map[string]*regexp.Regexp{
					"hello": regexp.MustCompile(`\bhello\b`),
				},
			},
		},
		{
			name: "IdentifierRedefinition",
			visitor: IdentifierLineVisitor{
				Definitions: map[string]string{
					"hello": "world",
				},
				Directive: markers,
			},
			line:       "//#id hello world2",
			wantLine:   "//#id hello world2",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: IdentifierLineVisitor{
				Definitions: map[string]string{
					"hello": "world2",
				},
				Directive: markers,
				marker:    "//#id ",
				identifiers: map[string]*regexp.Regexp{
					"hello": regexp.MustCompile(`\bhello\b`),
				},
			},
		},
		{
			name: "IdentifierIndented",
			visitor: IdentifierLineVisitor{
				Directive: markers,
			},
			line:       "  //#id hello world  ",
			wantLine:   "  //#id hello world  ",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: IdentifierLineVisitor{
				Definitions: map[string]string{
					"hello": "world",
				},
				Directive: markers,
				marker:    "//#id ",
				identifiers: map[string]*regexp.Regexp{
					"hello": regexp.MustCompile(`\bhello\b`),
				},
			},
		},
		{
			name: "TextLine",
			visitor: IdentifierLineVisitor{
				Directive: markers,
			},
			line:       "  func hello() {}  ",
			wantLine:   "  func hello() {}  ",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: IdentifierLineVisitor{
				Directive: markers,
				marker:    "//#id ",
			},
		},
		{
			name: "TextLineWithIdentifier",
			visitor: IdentifierLineVisitor{
				Definitions: map[string]string{
					"hello": "world",
				},
				Directive: markers,
				identifiers: map[string]*regexp.Regexp{
					"hello": regexp.MustCompile(`\bhello\b`),
				},
			},
			line:       "func hello() {}",
			wantLine:   "func world() {}",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: IdentifierLineVisitor{
				Definitions: map[string]string{
					"hello": "world",
				},
				Directive: markers,
				marker:    "//#id ",
				identifiers: map[string]*regexp.Regexp{
					"hello": regexp.MustCompile(`\bhello\b`),
				},
			},
		},
		{
			name: "TextLineWithIdentifierVariable",
			visitor: IdentifierLineVisitor{
				Definitions: map[string]string{
					"hello": "${async.upmsgtype}",
				},
				Directive: markers,
				identifiers: map[string]*regexp.Regexp{
					"hello": regexp.MustCompile(`\bhello\b`),
				},
			},
			line:       "func hello() {}",
			wantLine:   "func ${async.upmsgtype}() {}",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: IdentifierLineVisitor{
				Definitions: map[string]string{
					"hello": "${async.upmsgtype}",
				},
				Directive: markers,
				marker:    "//#id ",
				identifiers: map[string]*regexp.Regexp{
					"hello": regexp.MustCompile(`\bhello\b`),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &tt.visitor
			gotLine, gotOutput, err := v.VisitLine(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantLine, gotLine)
			assert.Equal(t, tt.wantOutput, gotOutput)
			assert.Equal(t, tt.wantVisitor, *v)
		})
	}
}

func TestConditionalLineVisitor_VisitLine(t *testing.T) {
	tpl := NewTemplate("", FileFormatGo, TemplateLanguageSkel, nil)
	markers := tpl.directiveMarkers()

	tests := []struct {
		name        string
		visitor     ConditionalLineVisitor
		line        string
		wantLine    string
		wantOutput  bool
		wantErr     bool
		wantVisitor ConditionalLineVisitor
	}{
		{
			name: "Outside",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
			},
			line:       "  func hello() {}  ",
			wantLine:   "  func hello() {}  ",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     0,
			},
		},
		{
			name: "OutsideNot",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     false,
				Directive: markers,
			},
			line:       "  func hello() {}  ",
			wantLine:   "  func hello() {}  ",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     false,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     0,
			},
		},
		{
			name: "If",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
			},
			line:       "  //#if COND  ",
			wantLine:   "  //#if COND  ",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     1,
			},
		},
		{
			name: "InsideIf",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				state:     1,
			},
			line:       "  func hello() {}  ",
			wantLine:   "  func hello() {}  ",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     1,
			},
		},
		{
			name: "InsideIfNot",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     false,
				Directive: markers,
				state:     1,
			},
			line:       "  func hello() {}  ",
			wantLine:   "  func hello() {}  ",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     false,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     1,
			},
		},
		{
			name: "InsideIfNegated",
			visitor: ConditionalLineVisitor{
				Name:      "!COND",
				Value:     true,
				Directive: markers,
				state:     1,
			},
			line:       "  func hello() {}  ",
			wantLine:   "  func hello() {}  ",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "!COND",
				Value:     true,
				Directive: markers,
				markers:   []string{"//#if !COND", "//#else !COND", "//#endif !COND"},
				state:     1,
			},
		},
		{
			name: "Else",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				state:     1,
			},
			line:       "  //#else COND  ",
			wantLine:   "  //#else COND  ",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     2,
			},
		},
		{
			name: "InsideElse",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				state:     2,
			},
			line:       "  func hello() {}  ",
			wantLine:   "  func hello() {}  ",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     2,
			},
		},
		{
			name: "InsideElseNot",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     false,
				Directive: markers,
				state:     2,
			},
			line:       "  func hello() {}  ",
			wantLine:   "  func hello() {}  ",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     false,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     2,
			},
		},
		{
			name: "EndIf",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				state:     1,
			},
			line:       "  //#endif COND  ",
			wantLine:   "  //#endif COND  ",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     0,
			},
		},
		{
			name: "EndIfElse",
			visitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				state:     2,
			},
			line:       "  //#endif COND  ",
			wantLine:   "  //#endif COND  ",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: ConditionalLineVisitor{
				Name:      "COND",
				Value:     true,
				Directive: markers,
				markers:   []string{"//#if COND", "//#else COND", "//#endif COND"},
				state:     0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &tt.visitor
			gotLine, gotOutput, err := v.VisitLine(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantLine, gotLine)
			assert.Equal(t, tt.wantOutput, gotOutput)
			assert.Equal(t, tt.wantVisitor, *v)
		})
	}
}

func TestVariableLineVisitor_VisitLine(t *testing.T) {
	tpl := NewTemplate("", FileFormatGo, TemplateLanguageSkel, nil)
	markers := tpl.directiveMarkers()

	tests := []struct {
		name        string
		visitor     VariableLineVisitor
		line        string
		wantLine    string
		wantOutput  bool
		wantErr     bool
		wantVisitor VariableLineVisitor
	}{
		{
			name: "VariableSubstitution",
			visitor: VariableLineVisitor{
				Variables: map[string]string{
					"hello": "world",
				},
				Directive: markers,
			},
			line:       "//#var hello",
			wantLine:   "world",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: VariableLineVisitor{
				Variables: map[string]string{
					"hello": "world",
				},
				Directive: markers,
				marker:    "//#var ",
			},
		},
		{
			name: "VariableNotExists",
			visitor: VariableLineVisitor{
				Variables: map[string]string{
					"hello": "world",
				},
				Directive: markers,
			},
			line:       "//#var hello2",
			wantLine:   "//#var hello2",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: VariableLineVisitor{
				Variables: map[string]string{
					"hello": "world",
				},
				Directive: markers,
				marker:    "//#var ",
			},
		},
		{
			name: "VariableNotExists",
			visitor: VariableLineVisitor{
				Variables: map[string]string{
					"hello": "world",
				},
				Directive: markers,
			},
			line:       "  //#var hello suffix=,",
			wantLine:   "world,",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: VariableLineVisitor{
				Variables: map[string]string{
					"hello": "world",
				},
				Directive: markers,
				marker:    "//#var ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &tt.visitor
			gotLine, gotOutput, err := v.VisitLine(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantLine, gotLine)
			assert.Equal(t, tt.wantOutput, gotOutput)
			assert.Equal(t, tt.wantVisitor, *v)
		})
	}
}

func TestJoinLineVisitor_VisitLine(t *testing.T) {
	tpl := NewTemplate("", FileFormatGo, TemplateLanguageSkel, nil)
	markers := tpl.directiveMarkers()

	tests := []struct {
		name        string
		visitor     JoinLineVisitor
		line        string
		wantLine    string
		wantOutput  bool
		wantErr     bool
		wantVisitor JoinLineVisitor
	}{
		{
			name: "Outside",
			visitor: JoinLineVisitor{
				Directive: markers,
			},
			line:       "outside",
			wantLine:   "outside",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: JoinLineVisitor{
				Directive: markers,
				markers:   []string{"//#join", "//#endjoin"},
			},
		},
		{
			name: "Join",
			visitor: JoinLineVisitor{
				Directive: markers,
			},
			line:       "//#join",
			wantLine:   "//#join",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: JoinLineVisitor{
				Directive: markers,
				markers:   []string{"//#join", "//#endjoin"},
				state:     1,
			},
		},
		{
			name: "Inside",
			visitor: JoinLineVisitor{
				Directive: markers,
				markers:   []string{"//#join", "//#endjoin"},
				state:     1,
			},
			line:       "inside",
			wantLine:   "inside",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: JoinLineVisitor{
				Directive: markers,
				markers:   []string{"//#join", "//#endjoin"},
				lines:     []string{"inside"},
				state:     1,
			},
		},
		{
			name: "EndJoin",
			visitor: JoinLineVisitor{
				Directive: markers,
				markers:   []string{"//#join", "//#endjoin"},
				lines:     []string{"inside"},
				state:     1,
			},
			line:       "//#endjoin",
			wantLine:   "inside",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: JoinLineVisitor{
				Directive: markers,
				markers:   []string{"//#join", "//#endjoin"},
				lines:     []string{"inside"},
				state:     0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &tt.visitor
			gotLine, gotOutput, err := v.VisitLine(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantLine, gotLine)
			assert.Equal(t, tt.wantOutput, gotOutput)
			assert.Equal(t, tt.wantVisitor, *v)
		})
	}
}

func TestIgnoreLineVisitor_VisitLine(t *testing.T) {
	tpl := NewTemplate("", FileFormatGo, TemplateLanguageSkel, nil)
	markers := tpl.directiveMarkers()

	tests := []struct {
		name        string
		visitor     IgnoreLineVisitor
		line        string
		wantLine    string
		wantOutput  bool
		wantErr     bool
		wantVisitor IgnoreLineVisitor
	}{
		{
			name: "Outside",
			visitor: IgnoreLineVisitor{
				Directive: markers,
			},
			line:       "outside",
			wantLine:   "outside",
			wantOutput: true,
			wantErr:    false,
			wantVisitor: IgnoreLineVisitor{
				Directive: markers,
				markers:   []string{"//#ignore", "//#endignore"},
			},
		},
		{
			name: "Ignore",
			visitor: IgnoreLineVisitor{
				Directive: markers,
			},
			line:       "//#ignore",
			wantLine:   "//#ignore",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: IgnoreLineVisitor{
				Directive: markers,
				markers:   []string{"//#ignore", "//#endignore"},
				state:     1,
			},
		},
		{
			name: "Inside",
			visitor: IgnoreLineVisitor{
				Directive: markers,
				markers:   []string{"//#ignore", "//#endignore"},
				state:     1,
			},
			line:       "inside",
			wantLine:   "inside",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: IgnoreLineVisitor{
				Directive: markers,
				markers:   []string{"//#ignore", "//#endignore"},
				state:     1,
			},
		},
		{
			name: "EndIgnore",
			visitor: IgnoreLineVisitor{
				Directive: markers,
				markers:   []string{"//#ignore", "//#endignore"},
				state:     1,
			},
			line:       "//#endignore",
			wantLine:   "//#endignore",
			wantOutput: false,
			wantErr:    false,
			wantVisitor: IgnoreLineVisitor{
				Directive: markers,
				markers:   []string{"//#ignore", "//#endignore"},
				state:     0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &tt.visitor
			gotLine, gotOutput, err := v.VisitLine(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantLine, gotLine)
			assert.Equal(t, tt.wantOutput, gotOutput)
			assert.Equal(t, tt.wantVisitor, *v)
		})
	}
}

func TestTemplate_Render(t *testing.T) {
	tests := []struct {
		name       string
		template   Template
		options    TemplateOptions
		wantResult string
		wantErr    bool
	}{
		{
			name: "Skel",
			template: Template{
				Name: "Simple",
				Loader: TemplateStringOption(
					`
						//#id contextKey contextKeyNamed
						//#if COND
						const contextKeyComponent = contextKey("${component.name}")
						//#endif COND
					`,
					DefaultTextTransformers...),
				Format:   FileFormatGo,
				Language: TemplateLanguageSkel,
			},
			options: TemplateOptions{
				Conditions: map[string]bool{
					"COND": true,
				},
				Strings: map[string]string{
					"Component": "Repository",
				},
				Variables: map[string]string{
					"component.name": "Repository",
				},
			},
			wantResult: `const contextKeyRepository = contextKeyNamed("Repository")` + "\n",
			wantErr:    false,
		},
		{
			name: "GoText",
			template: Template{
				Name:     "Simple",
				Loader:   TemplateStringOption(`const contextKey{{- .Strings.Component}} = contextKey("{{- .Strings.Component -}}")`),
				Format:   FileFormatGo,
				Language: TemplateLanguageGoText,
			},
			options: TemplateOptions{
				Strings: map[string]string{
					"Component": "Repository",
				},
			},
			wantResult: `const contextKeyRepository = contextKey("Repository")`,
			wantErr:    false,
		},
		{
			name: "GoHtml",
			template: Template{
				Name:     "Simple",
				Loader:   TemplateStringOption(`<a href="{{- .Strings.Component -}}" />`),
				Format:   FileFormatGo,
				Language: TemplateLanguageGoHtml,
			},
			options: TemplateOptions{
				Strings: map[string]string{
					"Component": "Repository",
				},
			},
			wantResult: `<a href="Repository" />`,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := &tt.template
			gotResult, err := template.Render(tt.options)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.wantResult, gotResult)
		})
	}
}
