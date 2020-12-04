//go:generate mockery --inpackage --name=Publisher --structname=MockPublisher
//go:generate mockery --inpackage --name=PublisherService --structname=MockPublisherService

package stream

import "testing"

func TestDummy(t *testing.T) {
	t.Skipped()
}
