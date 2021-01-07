package manage

import "testing"

func Test_Implementations(t *testing.T) {
	// Ensure MockManage is up to date
	var _ Api = new(MockManage)
}
