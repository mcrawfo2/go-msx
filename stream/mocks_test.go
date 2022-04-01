// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --inpackage --name=Publisher --structname=MockPublisher
//go:generate mockery --inpackage --name=MessagePublisher --structname=MockMessagePublisher
//go:generate mockery --inpackage --name=Subscriber --structname=MockSubscriber
//go:generate mockery --inpackage --name=PublisherService --structname=MockPublisherService
//go:generate mockery --inpackage --name=Dispatcher --structname=MockDispatcher
//go:generate mockery --inpackage --name=Provider --structname=MockProvider

package stream

import "testing"

func TestMockImplements(t *testing.T) {
	var _ Publisher = new(MockPublisher)
	var _ MessagePublisher = new(MockMessagePublisher)
	var _ Subscriber = new(MockSubscriber)
	var _ PublisherService = new(MockPublisherService)
	var _ Dispatcher = new(MockDispatcher)
	var _ Provider = new(MockProvider)
}
