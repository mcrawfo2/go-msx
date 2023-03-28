// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// This program creates all the txtar files corresponding to normal tests given in
// the testlist package -- it skips the special build tests

// Requires the FIXT env var to be set to the path of the fixtures directory

package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/fixtures/build/testlist"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"flag"
	"fmt"
	"golang.org/x/tools/txtar"
	"gopkg.in/pipe.v2"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	ignore  = "**/.git/** **/go.sum" // glob of files we will leave out of golden sets and tests
	envFIXT = "FIXT"
)

var nooverwrite bool

func main() {

	fmt.Printf("Fixtures directory: %s, paths below are relative to that\n", testlist.FIXT)

	flag.BoolVar(&nooverwrite, "nooverwrite", false, "do not overwrite existing files")
	flag.Parse()

	licenceTemplate := `# Copyright © %d, Cisco Systems Inc.
# Use of this source code is governed by an MIT-style license that can be
# found in the LICENSE file or at https://opensource.org/licenses/MIT.
`

	license := fmt.Sprintf(licenceTemplate, time.Now().Year())

	flames := "🔥🔥🔥🔥🔥 Error  🔥🔥🔥🔥🔥"

	// random order turns out to be aggravating, and we are not actually running the tests
	// so, lets sort them
	keys := make([]string, 0)
	if len(flag.Args()) > 0 {
		keys = flag.Args()
	} else {
		for k := range testlist.Tests {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	for _, tname := range keys {

		test, ok := testlist.Tests[tname]
		if !ok {
			fmt.Printf("Skipping    ⏭  No such test: %s\n", tname)
			continue
		}

		if test.Disabled { // skip broken ones
			fmt.Printf("Skipping    ⏭  disabled: test %s of command: %s %s\n", tname, test.Command, test.Args)
			continue
		}

		switch test.SpecialBuild {
		case testlist.SpecBuildNone:
			{
				fmt.Printf("Skipping    ⏭  special, none required: test %s of command: %s %s\n", tname, test.Command, test.Args)
				continue
			}
		}

		fname := testlist.Test2Filename("golden", tname, "-test.txtar")
		if nooverwrite && testlist.FileExists(fname) {
			fmt.Printf("Skipping    ⏭  Nooverwrite set & exists: %s for test %s of command: %s %s\n", fname, tname, test.Command, test.Args)
			continue
		}

		tmpdir, err := os.MkdirTemp("", "msx_skel_test_")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := os.RemoveAll(tmpdir)
			if err != nil {
				log.Fatal(err)
			}
		}()

		cmd := strings.TrimSpace(test.Command)
		args := strings.TrimSpace(test.Args)
		outputFName := testlist.Test2Filename("golden", tname, "-test.txtar")
		relname := strings.TrimPrefix(outputFName, testlist.FIXT+"/")
		fmt.Printf("Making test 🚧 skel %s, args:%s in 📁 %s\n", cmd, args, relname)

		var makeIt pipe.Pipe

		switch test.SpecialBuild {
		case testlist.SpecBuildNone:
			{
				fmt.Printf("Skipping    ⏭  special, none required: test %s of command: %s %s\n", tname, test.Command, test.Args)
				continue
			}
		case testlist.SpecBuildScript:
			{
				makeIt = pipe.Script(
					pipe.SetEnvVar(envFIXT, testlist.FIXT),
					pipe.ChDir(tmpdir),
					pipe.Exec(test.SBScript, args),
				)
			}
		case testlist.OrdinaryBuild:
			{
				rootFiles := testlist.Test2Filename("before", "plain-subroot.txtar", "")

				if len(test.ExtraBefores) > 0 { // this test needs extra files, so we will make a special "root" archive for it
					extrasArchive := testlist.Test2Filename("before", tname, "-test.txtar")
					newBase, err := txtar.ParseFile(testlist.Test2Filename("before", "plain-subroot.txtar", ""))
					if err != nil {
						fmt.Printf(flames+" loading plain subroot for test %s: %s\n", tname, err)
						panic(fmt.Sprintf("Error making test %s: %s", tname, err))
					}
					for _, extra := range test.ExtraBefores {
						contents, err := os.ReadFile(testlist.Test2Filename("before", extra, ""))
						if err != nil {
							fmt.Printf(flames+" loading extra before %s for test %s: %s\n", extra, tname, err)
							panic(fmt.Sprintf("Error making test %s: %s", tname, err))
						}
						newFile := txtar.File{
							Name: extra,
							Data: contents,
						}
						newBase.Files = append(newBase.Files, newFile)
						fmt.Printf("Adding ⊕ file %s for skel %s, args:%s in 📁 %s\n", extra, cmd, args, relname)
					}
					err = os.WriteFile(extrasArchive, txtar.Format(newBase), 0644)
					if err != nil {
						fmt.Printf(flames+" writing archive %s for test %s: %s\n", extrasArchive, tname, err)
						panic(fmt.Sprintf("Error making test %s: %s", tname, err))
					}
					rootFiles = extrasArchive
				}

				allArgs := []string{"--allow-dirty", cmd}
				allArgs = append(allArgs, strings.Fields(args)...) // split on spaces to make the array exec *needs*
				fmt.Printf("Making test 🚧 skel %s, args:%s in 📁 %s\n", cmd, allArgs, relname)

				makeIt = pipe.Script(
					pipe.Print(license),
					pipe.ChDir(tmpdir),
					pipe.Exec("txtarunwrap", rootFiles, "."),
					pipe.Exec("skel", allArgs...),
					pipe.SetEnvVar(envFIXT, testlist.FIXT),
					pipe.SetEnvVar(clienv.EnvIgnore, ignore),
					pipe.Line(
						pipe.Exec("txtarwrap", "."),
						pipe.WriteFile(outputFName, 0644),
					),
				)

			}
		}

		barf, err := pipe.CombinedOutput(makeIt)
		if err != nil {
			fmt.Printf(flames+" making test %s: %s\n%s\n", tname, err, barf)
			panic(fmt.Sprintf("Error making test %s: %s", tname, err))
		}

	}

}
