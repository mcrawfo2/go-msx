package auditlog

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"os"
	"testing"
)

var recording *logtest.Recording
var logger *log.Logger

func TestMain(m *testing.M) {
	recording = logtest.RecordLogging()
	logger = log.NewLogger("msx.audit.auditlog.test")
	os.Exit(m.Run())
}
