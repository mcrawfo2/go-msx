// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package populate

import "context"

type Populator interface {
	Populate(ctx context.Context) error
}

type PopulatorFactory func(context.Context) (Populator, error)

type PopulatorTask struct {
	factory     PopulatorFactory
	description string
	during      []string
	order       int
}

func (p PopulatorTask) Description() string {
	return p.description
}

func (p PopulatorTask) During() []string {
	return p.during
}

func (p PopulatorTask) Order() int {
	return p.order
}

func (p PopulatorTask) Populate(ctx context.Context) error {
	populator, err := p.factory(ctx)
	if err != nil {
		return err
	}
	return populator.Populate(ctx)
}

func NewPopulatorTask(description string, order int, during []string, factory PopulatorFactory) *PopulatorTask {
	return &PopulatorTask{
		factory:     factory,
		description: description,
		during:      during,
		order:       order,
	}
}
