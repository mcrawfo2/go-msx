package main

import (
	_ "cto-github.cisco.com/NFV-BU/administrationservice/pkg/populate"
	_ "cto-github.cisco.com/NFV-BU/catalogservice/pkg/populate"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	_ "cto-github.cisco.com/NFV-BU/go-msx/integration/serviceconfigmanager/populate"
	_ "cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement/populate"
	_ "cto-github.cisco.com/NFV-BU/go-msx-populator/populate"
)

const (
	appName = "${app.name}"
)

func main() {
	app.Run(appName)
}
