package rest

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
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

const maxLineLength = 120

type Decls []codegen.Decl

func (d Decls) Generate(out *codegen.Emitter) {
	for _, decl := range d {
		decl.Generate(out)
		out.Newline()
	}
}

func (d Decls) Render(indent int) string {
	out := codegen.NewEmitter(maxLineLength)
	out.Indent(indent)

	d.Generate(out)

	out.Indent(-indent)
	return out.String()
}

type Transformer func(string) string

type Transformers []Transformer

func (t Transformers) Transform(target string) string {
	for _, transformer := range t {
		target = transformer(target)
	}

	return target
}

type Snippet struct {
	Section string
	Name    string
	Lines   []string
	Imports []codegen.Import
}

func (s Snippet) GetName() string {
	return s.Name
}

func (s Snippet) Generate(out *codegen.Emitter) {
	for _, line := range s.Lines {
		out.Println("%s", line)
	}
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

func NewDeclSnippet(section, name string, decl codegen.Decl, imports []codegen.Import, transforms Transformers) (snippet Snippet) {
	emitter := codegen.NewEmitter(maxLineLength)
	decl.Generate(emitter)
	content := emitter.String()
	return NewTextSnippet(section, name, content, imports, transforms)
}

func NewTextSnippet(section, name, content string, imports []codegen.Import, transforms Transformers) (snippet Snippet) {
	content = transforms.Transform(content)
	imports = transformImports(imports, transforms)

	return Snippet{
		Section: section,
		Name:    name,
		Lines:   strings.Split(content, "\n"),
		Imports: imports,
	}
}

func NewStatementSnippet(section, name string, stmt *jen.Statement, transforms Transformers) (snippet Snippet, err error) {
	f := jen.NewFile(name)

	w := new(bytes.Buffer)
	if err = stmt.RenderWithFile(w, f); err != nil {
		return
	}

	inflectedSnippet := w.String()
	inflectedSnippet = transforms.Transform(inflectedSnippet)
	statementLines := strings.Split(inflectedSnippet, "\n")

	w.Reset()
	if err = f.Render(w); err != nil {
		return
	}

	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, name+".go", w.String(), parser.ImportsOnly)
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

	return Snippet{
		Section: section,
		Name:    name,
		Lines:   statementLines,
		Imports: imports,
	}, nil
}

type Constants []*codegen.Constant

func (c Constants) Generate(out *codegen.Emitter) {
	out.Println("const (")
	out.Indent(1)

	for _, constant := range c {
		out.Print("%s ", constant.Name)
		if constant.Type != nil {
			constant.Type.Generate(out)
		}
		out.Print(" = %s", litter.Sdump(constant.Value))
		out.Newline()
	}

	out.Indent(-1)
	out.Println(")")
}

type Comment string

func (c Comment) Generate(out *codegen.Emitter) {
	if strings.HasPrefix(string(c), "go:") || strings.HasPrefix(string(c), "#") {
		out.Println("//%s", c)
	} else {
		out.Println("// %s", c)
	}
}

type Section struct {
	Name     string
	Snippets []Snippet
	Sections Sections
}

func (s *Section) AddSnippet(snippet Snippet) {
	s.Snippets = append(s.Snippets, snippet)
}

func (s *Section) Generate(out *codegen.Emitter) {
	if s.Empty() {
		return
	}

	Comment(s.Name).Generate(out)
	out.Newline()

	for _, decl := range s.Snippets {
		decl.Generate(out)
		out.Newline()
	}

	s.Sections.Generate(out)
}

func (s *Section) Empty() bool {
	content := len(s.Snippets) > 0
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

type Sections []*Section

func (s Sections) Generate(out *codegen.Emitter) {
	for _, section := range s {
		section.Generate(out)
	}
}

func (s Sections) WithSection(name string) (Sections, *Section) {
	for _, section := range s {
		if section.Name == name {
			return s, section
		}
	}

	section := &Section{Name: name}
	return append(s, section), section
}

func NewSections(sections ...any) Sections {
	var results Sections
	for _, section := range sections {
		switch st := section.(type) {
		case string:
			results = append(results, &Section{
				Name: st,
			})

		case *Section:
			results = append(results, st)
		}
	}
	return results
}

type File struct {
	Comment   string
	Package   string
	Imports   []codegen.Import
	Sections  Sections
	Inflector skel.Inflector
}

func (f *File) AddImport(qualified, alias string) {
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

func (f *File) AddSnippet(snippet Snippet) {
	section := f.FindSection(snippet.Section)
	section.AddSnippet(snippet)

	imports := snippet.Imports
	for _, imp := range imports {
		f.AddImport(imp.QualifiedName, imp.Name)
	}
}

func (f *File) AddNewSnippet(snippet Snippet, err error) error {
	if err != nil {
		return err
	}
	f.AddSnippet(snippet)
	return nil
}

func (f *File) AddNewStatement(path, name string, stmt *jen.Statement) error {
	return f.AddNewSnippet(NewStatementSnippet(path, name, stmt, Transformers{f.Inflector.Inflect}))
}

func (f *File) AddNewText(path, name, text string, imports []codegen.Import) error {
	f.AddSnippet(NewTextSnippet(path, name, text, imports,
		Transformers{
			f.Inflector.Inflect,
			dedent.Dedent,
			strings.TrimSpace,
		}))
	return nil
}

func (f *File) AddNewDecl(path, name string, decl codegen.Decl, imports []codegen.Import) error {
	f.AddSnippet(NewDeclSnippet(path, name, decl, imports, Transformers{f.Inflector.Inflect}))
	return nil
}

func (f *File) AddImports(imports []codegen.Import) {
	for _, imp := range imports {
		f.AddImport(imp.QualifiedName, imp.Name)
	}
}

func (f *File) FindSection(section string) *Section {
	if section == "" {
		return nil
	}

	sectionPath := strings.Split(section, "/")

	var prevSection *Section
	f.Sections, prevSection = f.Sections.WithSection(sectionPath[0])
	sectionPath = sectionPath[1:]

	for len(sectionPath) > 0 {
		var nextSection *Section
		prevSection.Sections, nextSection = prevSection.Sections.WithSection(sectionPath[0])
		prevSection = nextSection
		sectionPath = sectionPath[1:]
	}

	return prevSection
}

func (f *File) Render() string {
	out := codegen.NewEmitter(maxLineLength)

	if f.Comment != "" {
		Comment(f.Comment).Generate(out)
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

const PkgApp = "cto-github.cisco.com/NFV-BU/go-msx/app"
const PkgRestops = "cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
const PkgTypes = "cto-github.cisco.com/NFV-BU/go-msx/types"
const PkgContext = "context"
const PkgOpenapi = "cto-github.cisco.com/NFV-BU/go-msx/schema/openapi"
const PkgLog = "cto-github.cisco.com/NFV-BU/go-msx/log"
const PkgRestopsV2 = PkgRestops + "/v8"
const PkgRestopsV8 = PkgRestops + "/v8"
const PkgSqldb = "cto-github.cisco.com/NFV-BU/go-msx/sqldb"
const PkgOpenApi3 = "github.com/swaggest/openapi-go/openapi3"
const PkgUuid = "github.com/google/uuid"
const PkgPrepared = "cto-github.cisco.com/NFV-BU/go-msx/sqldb/prepared"
const PkgPaging = "cto-github.cisco.com/NFV-BU/go-msx/paging"
const PkgRepository = "cto-github.cisco.com/NFV-BU/go-msx/repository"

var importApp = codegen.Import{QualifiedName: PkgApp}
var importRestOps = codegen.Import{QualifiedName: PkgRestops}
var importTypes = codegen.Import{QualifiedName: PkgTypes}
var importContext = codegen.Import{QualifiedName: PkgContext}
var importOpenapi = codegen.Import{QualifiedName: PkgOpenapi}
var importLog = codegen.Import{QualifiedName: PkgLog}
var importRestOpsV2 = codegen.Import{QualifiedName: PkgRestopsV2}
var importRestOpsV8 = codegen.Import{QualifiedName: PkgRestopsV8}
var importOpenApi3 = codegen.Import{QualifiedName: PkgOpenApi3}
var importSqldb = codegen.Import{QualifiedName: PkgSqldb}
var importUuid = codegen.Import{QualifiedName: PkgUuid}
var importPrepared = codegen.Import{QualifiedName: PkgPrepared, Name: "db"}
var importPaging = codegen.Import{QualifiedName: PkgPaging}
var importRepository = codegen.Import{QualifiedName: PkgRepository}
