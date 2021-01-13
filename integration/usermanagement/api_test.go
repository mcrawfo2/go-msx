package usermanagement

import "testing"

func Test_Implementations(t *testing.T) {
	var _ Api = new(MockUserManagement)
	var _ Api = new(Integration)
}

