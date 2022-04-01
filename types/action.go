// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"context"
	"github.com/pkg/errors"
	"sort"
)

const filterAutoSkip = -100

type ActionFunc func(ctx context.Context) error

type ActionFuncDecorator func(action ActionFunc) ActionFunc

// ActionFilter is an ordered Decorator
type ActionFilter interface {
	Order() int
	Decorator() ActionFuncDecorator
}

// ActionFilters is an ordered sequence of ActionFilter
type ActionFilters []ActionFilter

func (a ActionFilters) Len() int {
	return len(a)
}

func (a ActionFilters) Less(i, j int) bool {
	return a[i].Order() > a[j].Order()
}

func (a ActionFilters) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ActionFilters) NextCustomOrder() int {
	next := 0
	if len(a) > 0 {
		last := a[len(a)-1].Order()
		if last <= next {
			next = last + filterAutoSkip
		}
	}
	return next
}

type OrderedDecorator struct {
	order     int
	decorator ActionFuncDecorator
}

func (d OrderedDecorator) Order() int {
	return d.order
}

func (d OrderedDecorator) Decorator() ActionFuncDecorator {
	return d.decorator
}

func NewOrderedDecorator(order int, deco ActionFuncDecorator) OrderedDecorator {
	return OrderedDecorator{
		order:     order,
		decorator: deco,
	}
}

type Operation struct {
	action  ActionFunc
	filters ActionFilters
}

func (o Operation) decoratedAction() ActionFunc {
	action := o.action

	for i := len(o.filters) - 1; i >= 0; i-- {
		deco := o.filters[i].Decorator()
		action = deco(action)
	}

	return action
}

func (o Operation) WithDecorator(deco ActionFuncDecorator) Operation {
	return o.WithFilter(NewOrderedDecorator(o.filters.NextCustomOrder(), deco))
}

func (o Operation) WithFilter(filter ActionFilter) Operation {
	o.filters = append(o.filters[:], filter)
	sort.Sort(o.filters)
	return o
}

func (o Operation) Run(ctx context.Context) error {
	action := o.decoratedAction()
	return action(ctx)
}

func NewOperation(fn ActionFunc) Operation {
	return Operation{
		action: fn,
	}
}

func RecoverErrorDecorator(action ActionFunc) ActionFunc {
	return func(ctx context.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				var e error
				if err, ok := r.(error); ok {
					e = err
				} else {
					e = errors.Errorf("Exception: %v", r)
				}

				// TODO: decorate error with backtrace
				//bt := BackTraceFromDebugStackTrace(debug.Stack())
				err = e
			}
		}()

		err = action(ctx)
		return err
	}
}
