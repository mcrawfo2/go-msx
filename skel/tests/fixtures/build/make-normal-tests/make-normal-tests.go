// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// This program creates all the txtar files corresponding to normal tests given in
// the testlist package -- it skips the special build tests

// Requires the FIXT env var to be set to the path of the fixtures directory

package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/fixtures/build/testlist"
	"fmt"
	"gopkg.in/pipe.v2"
	"os"
	"sort"
	"strings"
)

func main() {

	fixt := os.Getenv("FIXT")
	if len(fixt) == 0 {
		panic("FIXT env var not set")
	}

	// random order turns out to be aggravating, and we are not actually running the tests
	keys := make([]string, 0)
	for k := range testlist.Tests {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, tname := range keys {
		test := testlist.Tests[tname]
		if !test.SpecialBuild && !test.Disabled { // skip special build tests and broken ones

			cmd := strings.TrimSpace(test.Command + " " + test.Args)
			outputFName := fmt.Sprintf("%s/final/%s-test.txtar", fixt, tname)
			fmt.Printf("Making test of:\nskel %s %s\n in %s\n\n", test.Command, test.Args, outputFName)

			tsargs := []string{"-e", "FIXT=" + fixt,
				"-e", "TEST_NAME=" + tname,
				"-e", "TEST_CMD=" + cmd}

			fmt.Printf("tsargs: %s\n", tsargs)

			makeIt := pipe.Line(
				pipe.ReadFile(fixt+"/build/make-any.stub"),
				pipe.Exec("txtarwrap", "-", fixt+"/before/plain-subroot.txtar"),
				pipe.Exec("testscript", tsargs...),
			)

			output, err := pipe.CombinedOutput(makeIt)
			if err != nil {
				fmt.Printf("Output>>>>>>>>>>>: %s", output)
				panic(fmt.Sprintf("Error making test %s: %s", tname, err))
			}

		}
	}

}
