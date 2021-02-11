//go:generate mockery --inpackage --name=Publisher --structname=MockPublisher
//go:generate mockery --inpackage --name=PublisherService --structname=MockPublisherService
//go:generate mockery --inpackage --name=Dispatcher --structname=MockDispatcher

package stream

import "testing"

func TestDummy(t *testing.T) {
	t.Skipped()
}
