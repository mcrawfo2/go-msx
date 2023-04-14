// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.
//
// txtargen
// see help text for description
// this program loads everything into memory, so it's not suitable for very large files or trees

package main

import (
	"bufio"
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"flag"
	"fmt"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/tools/txtar"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const help1 = `
%s is designed to be used in tests for programs that produce
trees of files as their output. It uses the diff command to compare the files and txtar archives to
store entires tree of 'golden' files

Usage of %s:
		%s [flags] archive|- directory 
where:
  - archive is a txtar archive (see below) that contains the "golden" files to which
    files created by the program under test will be compared
  - directory is the root of a directory tree, whose files will be read and compared

note: a txtar consists of comments (ordinary file text) and a list of files with their contents
The txtar archive format is described here: https://godoc.org/golang.org/x/tools/txtar

Flags:
`
const help2 = `
The program outputs the result of comparing the files (and perhaps their contents) in the txtar archive
with their namesakes in the directory tree.

Each file is matched against the cmpglobs, which are prefix:doublestar glob pairs (or trios) separated by spaces 

  exist:<**glob>    - perform an existance-only check for these
  notexist:<**glob> - perform a non-existance check 
  same:<**glob>     - perform a sameness check (contents must have been included in the archive)
  notsame:<**glob>  - perform a non-sameness check

Using ignorelines: if the regexp matches a line in a file that matches the glob, that line will be ignored in 
the comparison. If multple globs match a file the regexps are OR'd together using the | operator. Do not put literal
spaces in regexs, use [:space:] instead. Globs should not contain colons.

  ignorelines:<**glob>:<regex>

The first match wins, evaluated in the order provided, any number of them may be given 
Doublestar globs are described here: https://github.com/bmatcuk/doublestar

Any output will be sent to <stdout>
Reurns 0 if all tests pass, 1 if any fail, or if any produces an error

`

const (
	cmpGlobsIx clienv.Ix = iota
	sameIx
	notsameIx
	existsIx
	notexistsIx
	noMatchIx
	ignoreLinesIx
)

// glob prefixes

var prefix2ix = map[string]clienv.Ix{
	"same":        sameIx,
	"notsame":     notsameIx,
	"exists":      existsIx,
	"notexists":   notexistsIx,
	"ignorelines": ignoreLinesIx,
}

var ix2prefix = map[clienv.Ix]string{
	sameIx:        "same",
	notsameIx:     "notsame",
	existsIx:      "exists",
	notexistsIx:   "notexists",
	ignoreLinesIx: "ignorelines",
}

type RegexString string // a string that is a valid regex
type Glob string

type CmpFilters struct {
	whichFiles []Glob
	kinds      []clienv.Ix
	ignore     map[Glob]RegexString
}

func (g CmpFilters) Ignores(file string) (r RegexString) {
	r = ""
	for glob, regex := range g.ignore {
		if ok, _ := doublestar.Match(string(glob), file); ok {
			return regex
		}
	}
	return ""
}

func ParseGlobs(globs []string) (g CmpFilters, err error) {
	g.ignore = make(map[Glob]RegexString, 0)
	for _, p := range globs {
		piece := strings.Trim(p, "\"' ")
		if len(piece) == 0 {
			continue
		}
		prefix, after, found := strings.Cut(piece, ":")
		if !found {
			return g, fmt.Errorf("invalid glob expression: %s", piece)
		}
		ix, ok := prefix2ix[prefix]
		if !ok {
			return g, fmt.Errorf("invalid glob prefix: %s", prefix)
		}
		switch ix {
		case ignoreLinesIx:
			filesel, regex, found := strings.Cut(after, ":")
			if !found {
				return g, fmt.Errorf("invalid ignorelines glob expression: %s", after)
			}
			if !doublestar.ValidatePattern(filesel) {
				return g, fmt.Errorf("invalid ignorelines glob: %s", filesel)
			}
			_, err = regexp.Compile(regex)
			if err != nil {
				return g, fmt.Errorf("invalid ignorelines regex: %s", regex)
			}
			reg, already := g.ignore[Glob(filesel)]
			if already {
				regex = fmt.Sprintf("(%s)|(%s)", reg, regex)
			} else {
				regex = fmt.Sprintf("(%s)", regex)
			}
			g.ignore[Glob(filesel)] = RegexString(regex)
		default:
			g.whichFiles = append(g.whichFiles, Glob(after))
			g.kinds = append(g.kinds, ix)
		}
	}
	return g, err
}

// Variables for command line flags
var debugOn bool // debug output
var quiet bool   // produce no output, just exit code (0=tests all pass)
var cmpGlobsRaw string
var cmpGlobs []string

var flag2env = clienv.Flagvars{
	cmpGlobsIx: {Name: "cmpglobs", Env: clienv.EnvCmp,
		Default: "same:**", RawVar: &cmpGlobsRaw, FinVar: &cmpGlobs,
		Use: "globs that describe what tests to perform"},
}

const (
	fromStdin = "-"
)

func matchToGlobs(filters CmpFilters, filename string) (matched bool, matchKind clienv.Ix, err error) {
	matched = false
	matchKind = noMatchIx
	for i, glob := range filters.whichFiles {
		matched, err = doublestar.Match(string(glob), filename)
		if err != nil {
			return false, noMatchIx, fmt.Errorf("error matching glob: %v", err)
		}
		if matched {
			matchKind = filters.kinds[i]
			break
		}
	}
	return matched, matchKind, err
}

func main() {

	var err error

	flag.Usage = func() {
		_, err = fmt.Fprintf(flag.CommandLine.Output(), help1, strings.ToUpper(os.Args[0]), os.Args[0], os.Args[0])
		flag.PrintDefaults()
		_, err = fmt.Fprintf(flag.CommandLine.Output(), help2)
	}

	flag2env.Register()
	flag.BoolVar(&debugOn, "debug", false, "debug output")
	flag.BoolVar(&quiet, "quiet", false, "produce no output, just exit code (0=tests all pass)")
	flag.Parse()
	flag2env.Fallback()

	debug("txtarcmp: args: %s\n", flag.Args())
	debug(flag2env.String())
	debug("cmpglobs: %s\n", cmpGlobs)

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	questionDir := flag.Arg(1)
	if !strings.HasSuffix(questionDir, string(filepath.Separator)) {
		questionDir += string(filepath.Separator)
	}
	debug("question dir: %s\n", questionDir)

	goldens := flag.Arg(0)
	debug("golden files in: %s\n", goldens)

	filters, err := ParseGlobs(cmpGlobs)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "error parsing globs: %v\n", err)
		os.Exit(1)
	}
	debug("filters: %v\n", filters)

	arch, err := readArchive(goldens)
	comment := arch.Comment
	files := arch.Files
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "error reading archive: %v\n", err)
		os.Exit(1)
	}
	debug("read archive: %s, comment: %d, files: %d\n", goldens, len(comment), len(files))

	var reasons []string
	var allOK = true

	for _, goldfile := range files {

		goldpath := goldfile.Name
		questionpath := filepath.Join(filepath.Base(questionDir), goldpath)
		debug("comparing: goldpath: %s (archive %s), questionpath: %s\n", goldpath, goldens, questionpath)

		matched, matchedKind, err := matchToGlobs(filters, goldpath)
		if err != nil {
			_, err = fmt.Fprintf(os.Stderr, "error matching globs: %v\n", err)
			os.Exit(1)
		}

		if !matched {
			continue
		}

		var ok bool
		var reason string

		switch matchedKind {
		case existsIx:
			debug("exists: %s, %v\n", questionpath, ok)
			ok = fileExists(questionpath)
			if !ok {
				reason = fmt.Sprintf("should exist: %s", questionpath)
			}
		case notexistsIx:
			debug("notexists: %s, %v\n", questionpath, ok)
			ok = !fileExists(questionpath)
			if !ok {
				reason = fmt.Sprintf("%s exists", questionpath)
			}
		case sameIx:
			debug("same: %s, %s (ignore: %s)\n", questionpath, goldpath, filters.Ignores(questionpath))
			questionData, err := os.ReadFile(questionpath)
			if err != nil {
				reason = fmt.Sprintf("could not read: %s, %s, %v", questionpath, goldpath, err)
				break
			}
			same, diff, err := same(bytes.NewReader(questionData), bytes.NewReader(goldfile.Data), filters.Ignores(questionpath))
			ok = same
			if err != nil {
				reason = fmt.Sprintf("error comparing files: %s, %s, %v", questionpath, goldpath, err)
				break
			}
			if !ok {
				reason = fmt.Sprintf("files wrongly differ: %s, %s diffs:\n%s\n\n",
					questionpath, goldpath, diff)
			}
		case notsameIx:
			debug("notsame: %s, %s (ignore: %s)\n", questionpath, goldpath, filters.Ignores(questionpath))
			questionData, err := os.ReadFile(questionpath)
			if err != nil {
				reason = fmt.Sprintf("could not read: %s, %s, %v", questionpath, goldpath, err)
				break
			}
			same, _, err := same(bytes.NewReader(questionData), bytes.NewReader(goldfile.Data), filters.Ignores(questionpath))
			ok = !same
			if err != nil {
				reason = fmt.Sprintf("error comparing files: %s, %s, %v", questionpath, goldpath, err)
				break
			}
			if !ok {
				reason = fmt.Sprintf("files wrongly same: %s, %s", questionpath, goldpath)
			}
		}

		if !ok {
			if len(reason) > 0 {
				reasons = append(reasons, reason)
			}
			allOK = false
		}

	}

	if !allOK {
		if !quiet {
			fmt.Printf("Comparison failed:\n")
			for _, reason := range reasons {
				fmt.Printf("%s\n", reason)
			}
		}
		os.Exit(1)
	}

	os.Exit(0)
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
		fmt.Printf("#debug 💥> "+format, args...)
	}
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	if err == os.ErrNotExist || err != nil {
		return false
	}
	return true
}

func same(actual, expected io.Reader, ignoreLines RegexString) (same bool, diffText string, err error) {
	var rg *regexp.Regexp = nil
	if len(ignoreLines) > 0 {
		rg = regexp.MustCompile(string(ignoreLines))
	}

	actualLines, err := matchingLines(actual, rg)
	if err != nil {
		return false, "", err
	}

	expectedLines, err := matchingLines(expected, rg)
	if err != nil {
		return false, "", err
	}

	diff := difflib.UnifiedDiff{
		A:        expectedLines,
		B:        actualLines,
		FromFile: "Expected",
		ToFile:   "Actual",
		Context:  2,
	}

	diffText, err = difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return false, "", err
	}

	same = diffText == ""
	return
}

func matchingLines(source io.Reader, rg *regexp.Regexp) (matching []string, err error) {
	scanr := bufio.NewScanner(source)
	scanr.Split(bufio.ScanLines)

	for scanr.Scan() {
		t := scanr.Text()
		if rg == nil || !rg.MatchString(t) {
			matching = append(matching, t+"\n")
		}
	}
	if scanr.Err() != nil {
		return nil, scanr.Err()
	}

	return matching, nil
}
