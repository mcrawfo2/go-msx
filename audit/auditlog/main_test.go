package auditlog

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"os"
	"testing"
)

var recording *log.Recording
var logger *log.Logger

func TestMain(m *testing.M) {
	recording = log.RecordLogging()
	logger = log.NewLogger("msx.audit.auditlog.test")
	os.Exit(m.Run())
}
