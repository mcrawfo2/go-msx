// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// Run-tests runs the final tests listed in the testlist package

package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/fixtures/build/testlist"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"fmt"
	"gopkg.in/pipe.v2"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {

	var keys []string

	if len(os.Args) > 1 {
		keys = os.Args[1:]
	} else {
		keys = make([]string, 0)
		for k := range testlist.Tests {
			keys = append(keys, k)
		}
	}

	var wg sync.WaitGroup
	allPassed := true

	for _, name := range keys { // note we are running in random order and parallel if no args

		test, ok := testlist.Tests[name]
		if !ok {
			fmt.Printf("Skipping: No such test: %s of command: %s %s ⏭\n", name, test.Command, test.Args)
			continue
		}

		if test.Disabled {
			fmt.Printf("Skipping: Disabled: %s of command: %s %s ⏭\n", name, test.Command, test.Args)
			continue
		}

		fname := testlist.Test2Filename("golden", name, "-test.txtar")

		if test.SpecialRun == testlist.OrdinaryRun { // it's a normal test, so we need a golden file
			if _, err := os.Stat(fname); os.IsNotExist(err) {
				fmt.Printf("Skipping: Not made: %s of command: %s %s ⏭\n", name, test.Command, test.Args)
				continue
			}
		}

		fmt.Printf("Running: test %s of command: %s %s\n", name, test.Command, test.Args)

		wg.Add(1)

		innername := name // capture for closure
		go func() {

			defer wg.Done()

			dir, err := os.MkdirTemp("", "msx_skel_test_")
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				err := os.RemoveAll(dir)
				if err != nil {
					log.Fatal(err)
				}
			}()

			golden := testlist.Test2Filename("golden", innername, "-test.txtar")
			test := testlist.Tests[innername]

			before := testlist.Test2Filename("before", "plain-subroot.txtar", "")
			if len(test.ExtraBefores) > 0 {
				before = testlist.Test2Filename("before", innername, "-test.txtar")
			}

			genglobs := test.CmpGlobs

			allArgs := []string{"--allow-dirty", innername}
			allArgs = append(allArgs, strings.Fields(test.Args)...) // split on spaces to make the array exec *needs*

			var runIt pipe.Pipe
			if test.SpecialRun == testlist.OrdinaryRun { // it's a normal test
				runIt = pipe.Script(
					pipe.ChDir(dir),
					pipe.Exec("txtarunwrap", before, "."),
					pipe.Exec("skel", allArgs...),
					pipe.SetEnvVar(clienv.EnvCmp, genglobs),
					pipe.Exec("txtarcmp", "-debug", golden, "."),
				)
			} else { // it's a special test
				if test.SpecialRun == testlist.SpecRunPipe {
					runIt = pipe.Script(
						pipe.ChDir(dir),
						pipe.SetEnvVar(clienv.EnvCmp, genglobs),
						test.SRPipe,
					)
				} else {
					ok := test.SRFunction()
					if !ok {
						fmt.Printf("🔥🔥🔥🔥 Error 🔥 running test %s\n", innername)
						allPassed = false
					}
					fmt.Printf("%s: ✅ OK\n", innername)
					return
				}
			}

			output, err := pipe.CombinedOutput(runIt)
			if err != nil {
				fmt.Printf("🔥🔥🔥🔥 Error 🔥 running test %s: %s\n%s\n", innername, err, output)
				allPassed = false
				return
			}

			fmt.Printf("%s: ✅ OK\n", innername)
		}()

	}

	wg.Wait()

	if !allPassed {
		fmt.Printf("Some tests failed ❌\n")
		os.Exit(1)
	}

	fmt.Printf("All tests passed ✅\n")
	os.Exit(0)

}
