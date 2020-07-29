package populate

import "context"

type Populator interface {
	Populate(ctx context.Context) error
}

type PopulatorFactory func(context.Context) Populator

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
	return p.factory(ctx).Populate(ctx)
}

func NewPopulatorTask(description string, order int, during []string, factory PopulatorFactory) *PopulatorTask {
	return &PopulatorTask{
		factory:     factory,
		description: description,
		during:      during,
		order:       order,
	}
}