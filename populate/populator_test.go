// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package populate

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewPopulatorTask(t *testing.T) {
	const description = "description"
	const order = 1000
	var during = []string{"mach1"}
	var populator = new(MockTask)
	var factory PopulatorFactory = func(ctx context.Context) (Populator, error) {
		return populator, nil
	}
	var populatorRan = false
	populator.On("Populate", mock.AnythingOfType("*context.emptyCtx")).
		Run(func(args mock.Arguments) {
			populatorRan = true
		}).
		Return(nil)

	populatorTask := NewPopulatorTask(
		description,
		order,
		during,
		factory)

	assert.Equal(t, description, populatorTask.Description())
	assert.Equal(t, order, populatorTask.Order())
	assert.True(t, reflect.DeepEqual(during, populatorTask.During()))

	err := populatorTask.Populate(context.Background())
	assert.NoError(t, err)
	assert.True(t, populatorRan)
}
