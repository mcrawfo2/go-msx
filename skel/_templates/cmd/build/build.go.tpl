package main

import (
//#if EXTERNAL
	build "cto-github.cisco.com/NFV-BU/go-msx/build"
//#else EXTERNAL
	build "cto-github.cisco.com/NFV-BU/go-msx-build/pkg"
//#endif EXTERNAL
)

func main() {
	build.Run()
}
