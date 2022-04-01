// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"strconv"
)

var BuildNumber = "0"

func main() {
	buildNumber, _ := strconv.ParseInt(BuildNumber, 10, 64)
	skel.Run(int(buildNumber))
}
