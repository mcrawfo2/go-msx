// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.
//
// txtargen
// see help text for description
// this program loads everything into memory, so it's not suitable for very large files or trees

package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"flag"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"golang.org/x/tools/txtar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const help1 = `
%s is designed to facilitate creating testscript tests for programs that produce
trees of files as their output. 

Usage of %s:
		%s [flags] prefix archive 
where:
  - prefix is a string that has been prepended to all the files in the archive that are to be considered
    "golden"
  - archive is a txtar archive (see below) that contains the golden files (prefixed), 
    input files for the program (unprefixed), and a testscript script that runs the program under test

note: a txtar consists of comments (ordinary file text) and a list of files with their contents
The txtar archive format is described here: https://godoc.org/golang.org/x/tools/txtar

Flags:
`
const help2 = `
The program outputs a testscript script that tests any prefixed files against their unprefixed
  namesakes (for 'same' checks). This allows both the golden files, and input files for the
  program under test, to be included in the same archive and bundled into a single txtar file 
  along with the script needed to run the program and generated test program.

(for testscript syntax see: https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript)

Each file is matched against the GenGlobs, which are prefix:doublestar glob pairs, pairs separated by spaces 

  exist:<**glob>    - generate an existance-only check for these
  notexist:<**glob> - generate a non-existance check 
  same:<**glob>     - generate a sameness check (contents must have been included)
  notsame:<**glob>  - generate a non-sameness check

The first match wins, evaluated in the order provided, any number of them may be given 

Doublestar globs are described here: https://github.com/bmatcuk/doublestar

Output will be sent to <stdout> in the order: archive comments, testscript script, archive files

`

const (
	genGlobsIx clienv.Ix = iota
	sameIx
	notsameIx
	existsIx
	notexistsIx
	noMatchIx
)

// glob prefixes

var prefix2ix = map[string]clienv.Ix{
	"same":      sameIx,
	"notsame":   notsameIx,
	"exists":    existsIx,
	"notexists": notexistsIx,
}

// Variables for command line flags
var debugOn bool // debug output
var genGlobsRaw string
var genGlobs []string

var flag2env = clienv.Flagvars{
	genGlobsIx: {Name: "genglobs", Env: clienv.EnvCmp,
		Default: "same:**", RawVar: &genGlobsRaw, FinVar: &genGlobs,
		Use: "globs that describe how to generate test lines"},
}

const (
	fromStdin = "-"
)

func main() {

	var err error

	flag.Usage = func() {
		_, err = fmt.Fprintf(flag.CommandLine.Output(), help1, strings.ToUpper(os.Args[0]), os.Args[0], os.Args[0])
		flag.PrintDefaults()
		_, err = fmt.Fprintf(flag.CommandLine.Output(), help2)
	}

	flag2env.Register()
	flag.BoolVar(&debugOn, "debug", false, "debug output")
	flag.Parse()
	flag2env.Fallback()

	debug("txtargen: args: %s\n", flag.Args())
	debug(flag2env.String())
	debug("############# genglobs: %s\n", genGlobs)

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	generate := flag.Arg(0)
	if !strings.HasSuffix(generate, string(filepath.Separator)) {
		generate += string(filepath.Separator)
	}
	debug("generate: %s\n", generate)

	goldens := flag.Arg(1)
	debug("goldens: %s\n", goldens)

	arch, err := readArchive(goldens)
	comment := arch.Comment
	files := arch.Files
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "error reading archive: %v\n", err)
		os.Exit(1)
	}
	debug("read archive: %s, comment: %d, files: %d\n", goldens, len(comment), len(files))

	generated := fmt.Sprintf("\n\n# ========== Generated, unwise to edit ==========\n")

	for _, file := range files {

		goldpath := file.Name
		path := file.Name
		if strings.HasPrefix(path, generate) { // only generate for these
			path = path[len(generate):]
		} else {
			continue // next file please
		}

		kind, err := matchGenGlobs(path, genGlobs)

		if err != nil {
			_, err = fmt.Fprintf(os.Stderr, "error matching globs: %v\n", err)
			os.Exit(1)
		}

		var line string
		if kind == noMatchIx {
			continue
		}

		switch kind {
		case existsIx:
			line = "exists " + path + "\n"
		case notexistsIx:
			line = "! exists " + path + "\n"
		case sameIx:
			line = "cmp " + path + " " + goldpath + "\n"
		case notsameIx:
			line = "! cmp " + path + " " + goldpath + "\n"
		}
		debug("generated: %s", line)
		generated = generated + line

	}

	//	debug("all generated: %s", generated+"<<")

	generated = generated + fmt.Sprintf("# ========== End of generated code ==========\n\n")

	finalcomm := string(comment) + generated
	final := txtar.Archive{Files: files, Comment: []byte(finalcomm)}

	_, err = os.Stdout.Write(txtar.Format(&final))
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "error writing archive: %v\n", err)
		os.Exit(1)
	}

	return
}

func readArchive(filename string) (arch *txtar.Archive, err error) {
	var f *os.File
	arch = &txtar.Archive{}
	if filename == fromStdin {
		f = os.Stdin
		err = nil
	} else {
		f, err = os.Open(filename)
	}
	if err != nil {
		return arch, fmt.Errorf("error reading the archive: %v\n", err)
	}
	contents, err := io.ReadAll(f)
	if err != nil {
		return arch, fmt.Errorf("error reading the archive: %v\n", err)
	}
	arch = txtar.Parse(contents)
	_ = f.Close()
	return arch, err
}

func debug(format string, args ...interface{}) {
	if debugOn {
		fmt.Printf("#debug> "+format, args...)
	}
}

func matchGenGlobs(filename string, genglobs []string) (matchIx clienv.Ix, err error) {
	for _, glob := range genglobs {
		kind, matchTo := globSplit(glob)
		if kind == noMatchIx {
			continue
		}
		matched, err := doublestar.Match(matchTo, filename)
		if err != nil {
			return noMatchIx, err
		}
		if matched {
			return kind, nil
		}
	}
	return noMatchIx, nil
}

func globSplit(gg string) (kind clienv.Ix, match string) {
	before, after, found := strings.Cut(gg, ":")
	if !found {
		return noMatchIx, ""
	}
	kind, found = prefix2ix[before]
	if !found {
		return noMatchIx, ""
	}
	return kind, after
}

//sameIx: {Name: "same", Env: EnvSame, Default: "**", RawVar: &sameRaw, FinVar: &same,
//	Use: "globs for files that must be the same"},
//notsameIx: {Name: "notsame", Env: EnvNotSame, Default: "", RawVar: &notsameRaw, FinVar: &notsame,
//	Use: "globs for files that must not be the same"},
//existsIx: {Name: "exists", Env: EnvExists, Default: "", RawVar: &existsRaw, FinVar: &exists,
//	Use: "globs for files that must exist"},
//notexistsIx: {Name: "notexists", Env: EnvNotExists, Default: "", RawVar: &notexistsRaw, FinVar: &notexists,
//	Use: "globs for files that must not exist"},
//var sameRaw string
//var notsameRaw string
//var existsRaw string
//var notexistsRaw string
//var same []string
//var notsame []string
//var exists []string
//var notexists []string
//EnvSame      = "TXTAR_SAME"
//EnvNotSame   = "TXTAR_NOTSAME"
//EnvExists    = "TXTAR_EXISTS"
//EnvNotExists = "TXTAR_NOTEXISTS"
//func matchGlobs(fileName string, globs []string) (isMatch bool, err error) {
//	isMatch = false
//	for _, m := range globs {
//		j, err := doublestar.Match(m, fileName)
//		if err != nil {
//			return false, err
//		}
//		if j {
//			return true, nil
//		}
//	}
//	return isMatch, nil
//}

//var ix2prefix = map[clienv.Ix]string{
//	sameIx:      "same",
//	notsameIx:   "notsame",
//	existsIx:    "exists",
//	notexistsIx: "notexists",
//}
