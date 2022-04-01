// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cassandra

import (
	"context"
	"github.com/gocql/gocql"
)

type CompositeQueryObserver struct {
	queryObservers []gocql.QueryObserver
}

func (c CompositeQueryObserver) ObserveQuery(ctx context.Context, query gocql.ObservedQuery) {
	for _, observer := range c.queryObservers {
		observer.ObserveQuery(ctx, query)
	}
}

func NewCompositeQueryObserver(observers ...gocql.QueryObserver) *CompositeQueryObserver {
	return &CompositeQueryObserver{queryObservers: observers}
}

type CompositeBatchObserver struct {
	batchObservers []gocql.BatchObserver
}

func (c CompositeBatchObserver) ObserveBatch(ctx context.Context, batch gocql.ObservedBatch) {
	for _, observer := range c.batchObservers {
		observer.ObserveBatch(ctx, batch)
	}
}

func NewCompositeBatchObserver(observers ...gocql.BatchObserver) *CompositeBatchObserver {
	return &CompositeBatchObserver{batchObservers: observers}
}
