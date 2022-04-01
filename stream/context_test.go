// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContextWithPublisherService(t *testing.T) {
	ctx := context.Background()
	o := PublisherServiceFromContext(ctx)
	assert.Nil(t, o)
	service := new(MockPublisherService)
	ctx = ContextWithPublisherService(ctx, service)
	p := PublisherServiceFromContext(ctx)
	assert.Equal(t, service, p)
}
