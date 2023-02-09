// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package tests

// We will run all the test scripts in ./fixtures using testscript
//
// *Very* helpful article here: https://bitfieldconsulting.com/golang/test-scripts -- *recommended*
//
// see also: https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript

// TODO
//func TestSkel(t *testing.T) {
//	testscript.Run(t, testscript.Params{
//		Dir: "fixtures/scripts",
//	})
//}
//
//func TestMain(m *testing.M) {
//	os.Exit(testscript.RunMain(m, map[string]func() int{
//		"skel": skel.Main,
//	}))
//}

// TestHelp runs through all the skel help commands and checks they return without error
//func TestHelp(t *testing.T) {
//
//	for test := range testlist.Tests {
//		cmd := exec.Command("skel", "help", test)
//		err := cmd.Run()
//		if err != nil {
//			t.Errorf("Test %s failed: %s", test, err)
//		}
//	}
//
//}
