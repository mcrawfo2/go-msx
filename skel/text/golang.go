// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package text

import (
	"bytes"
	"github.com/lithammer/dedent"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"github.com/mcrawfo2/jennifer/jen"
	"github.com/sanity-io/litter"
	"go/parser"
	"go/token"
	"path"
	"strconv"
	"strings"
)

// Decls is a sequence of declarations
type Decls []codegen.Decl

// Generate
func (d Decls) Generate(emitter Emitter) {
	e := emitter.Raw()

	for _, decl := range d {
		decl.Generate(e)
		e.Newline()
	}
}

func (d Decls) Render(indent int) string {
	out := NewGoEmitter()
	out.Indent(indent)

	d.Generate(out)

	out.Indent(-indent)
	return out.String()
}

type GoSnippet struct {
	Snippet
	Imports []codegen.Import
}

func transformImports(imports []codegen.Import, transforms Transformers) []codegen.Import {
	for n, imp := range imports {
		if imp.Name != "" {
			imports[n].Name = transforms.Transform(imp.Name)
		}
		imports[n].QualifiedName = transforms.Transform(imp.QualifiedName)
	}
	return imports
}

func NewGoGeneratorSnippet(name string, gen Generator, imports []codegen.Import, transforms Transformers) (snippet GoSnippet) {
	emitter := NewGoEmitter()
	gen.Generate(emitter)
	content := emitter.String()
	return NewGoTextSnippet(name, content, imports, transforms)
}

func NewGoTextSnippet(name, content string, imports []codegen.Import, transforms Transformers) (snippet GoSnippet) {
	content = transforms.Transform(content)
	imports = transformImports(imports, transforms)

	var lines []string
	if strings.TrimSpace(content) != "" {
		lines = strings.Split(content, "\n")
	}

	return GoSnippet{
		Snippet: Snippet{
			Name:  name,
			Lines: lines,
		},
		Imports: imports,
	}
}

func NewGoStatementSnippet(name string, stmt *jen.Statement, transforms Transformers) (snippet GoSnippet, err error) {
	if stmt == nil {
		return GoSnippet{
			Snippet: Snippet{
				Name: name,
			},
		}, nil
	}

	f := jen.NewFile(name)

	w := new(bytes.Buffer)
	if err = stmt.RenderWithFile(w, f); err != nil {
		return
	}

	inflectedSnippet := w.String()
	inflectedSnippet = transforms.Transform(inflectedSnippet)

	var lines []string
	if strings.TrimSpace(inflectedSnippet) != "" {
		lines = strings.Split(inflectedSnippet, "\n")
	}

	w.Reset()
	if err = f.Render(w); err != nil {
		return
	}

	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, name+".go", w.Bytes(), parser.ImportsOnly)
	if err != nil {
		return
	}

	var imports []codegen.Import
	for _, importSpec := range af.Imports {
		importName := ""
		if importSpec.Name != nil {
			importName = importSpec.Name.Name
		}

		qualifiedName := ""
		qualifiedName, err = strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			return
		}

		imports = append(imports, codegen.Import{
			Name:          importName,
			QualifiedName: qualifiedName,
		})
	}
	imports = transformImports(imports, transforms)

	return GoSnippet{
		Snippet: Snippet{
			Name:  name,
			Lines: lines,
		},
		Imports: imports,
	}, nil
}

type GoConstants []*codegen.Constant

func (c GoConstants) Generate(out Emitter) {
	out.Println("const (")
	out.Indent(1)

	for _, constant := range c {
		out.Print("%s", constant.Name)
		if constant.Type != nil {
			out.Print(" ")
			constant.Type.Generate(out.Raw())
		}
		out.Print(" = %s", litter.Sdump(constant.Value))
		out.Newline()
	}

	out.Indent(-1)
	out.Println(")")
}

type GoComment string

func (c GoComment) Generate(out Emitter) {
	if strings.HasPrefix(string(c), "go:") || strings.HasPrefix(string(c), "#") {
		out.Println("//%s", c)
	} else {
		out.Println("// %s", c)
	}
}

func NewGoSections(sections ...any) Sections[GoSnippet] {
	return NewSections[GoSnippet](sections...)
}

type GoFile struct {
	*File[GoSnippet]
	Package string
	Imports []codegen.Import
}

func (f *GoFile) AddImport(qualified, alias string) {
	_, pkg := path.Split(qualified)
	if pkg == alias {
		alias = ""
	}

	for n, imp := range f.Imports {
		if imp.QualifiedName == qualified {
			if imp.Name == "_" {
				f.Imports[n].Name = alias
			}
			return
		}
	}

	f.Imports = append(f.Imports, codegen.Import{
		Name:          alias,
		QualifiedName: qualified,
	})
}

func (f *GoFile) AddNewStatement(sectionPath, name string, stmt *jen.Statement) error {
	snippet, err := NewGoStatementSnippet(name, stmt, Transformers{f.Inflector.Inflect})
	if err != nil {
		return err
	}

	f.AddSnippet(sectionPath, snippet)
	return nil
}

func (f *GoFile) AddNewText(sectionPath, name, body string, imports []codegen.Import) error {
	f.AddSnippet(sectionPath, NewGoTextSnippet(name, body, imports, Transformers{
		f.Inflector.Inflect,
		dedent.Dedent,
		strings.TrimSpace,
	}))
	return nil
}

func (f *GoFile) AddNewDecl(sectionPath, name string, decl codegen.Decl, imports []codegen.Import) error {
	f.AddSnippet(sectionPath, NewGoGeneratorSnippet(name, Decls{decl}, imports, Transformers{f.Inflector.Inflect}))
	return nil
}

func (f *GoFile) AddNewGenerator(sectionPath, name string, decl Generator, imports []codegen.Import) error {
	f.AddSnippet(sectionPath, NewGoGeneratorSnippet(name, decl, imports, Transformers{f.Inflector.Inflect}))
	return nil
}

func (f *GoFile) AddImports(imports []codegen.Import) {
	for _, imp := range imports {
		f.AddImport(imp.QualifiedName, imp.Name)
	}
}

func (f *GoFile) Render() string {
	out := NewGoEmitter()

	if f.Comment != "" {
		GoComment(f.Comment).Generate(out)
		out.Newline()
	}

	out.Println("package %s", f.Package)
	out.Newline()

	if len(f.Imports) > 1 {
		out.Println("import (")
		out.Indent(1)
		for _, imp := range f.Imports {
			if imp.Name != "" {
				out.Println("%s %q", imp.Name, imp.QualifiedName)
			} else {
				out.Println("%q", imp.QualifiedName)
			}
		}
		out.Indent(-1)
		out.Println(")")
		out.Newline()
	} else if len(f.Imports) == 1 {
		imp := f.Imports[0]
		qualifiedName := strconv.Quote(imp.QualifiedName)
		if imp.Name != "" {
			out.Println("import %s %s", imp.Name, qualifiedName)
		} else {
			out.Println("import %s", qualifiedName)
		}
		out.Newline()
	}

	f.Sections.Generate(out)

	return out.String()
}
