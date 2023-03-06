// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package testhelpers

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/mohae/deepcopy"
	"testing"
)

type Case[I Testable] struct {
	*types.Assembler[Case[I], *Case[I]]
	T        *testing.T
	Testable I
}

func (c *Case[I]) Clone() *Case[I] {
	return deepcopy.Copy(c).(*Case[I])
}

func (c *Case[I]) Test(t *testing.T) {
	c.T = t
	c.Assemble()
	c.Testable.Test(t)
}

func (c *Case[I]) DeepCopy() any {
	result := &Case[I]{
		Assembler: &types.Assembler[Case[I], *Case[I]]{
			Setups: deepcopy.Copy(c.Assembler.Setups).([]types.NamedSetupFunc[Case[I], *Case[I]]),
		},
		Testable: deepcopy.Copy(c.Testable).(I),
	}
	result.Assembler.Target = result
	return result
}

func NewCase[I Testable](testable I) *Case[I] {
	c := new(Case[I])
	c.Testable = testable
	c.Assembler = types.NewAssembler(c)
	return c
}

type FixtureCase[I Testable, F any] struct {
	*types.Assembler[FixtureCase[I, F], *FixtureCase[I, F]]
	T        *testing.T
	Testable I
	Fixture  F
}

func (c *FixtureCase[I, F]) Clone() *FixtureCase[I, F] {
	return deepcopy.Copy(c).(*FixtureCase[I, F])
}

func (c *FixtureCase[I, F]) DeepCopy() any {
	result := &FixtureCase[I, F]{
		Assembler: &types.Assembler[FixtureCase[I, F], *FixtureCase[I, F]]{
			Setups: deepcopy.Copy(c.Assembler.Setups).([]types.NamedSetupFunc[FixtureCase[I, F], *FixtureCase[I, F]]),
		},
		Testable: deepcopy.Copy(c.Testable).(I),
		Fixture:  deepcopy.Copy(c.Fixture).(F),
	}
	result.Assembler.Target = result
	return result
}

func (c *FixtureCase[I, F]) Test(t *testing.T) {
	c.T = t
	c.Assemble()
	c.Testable.Test(t)
}

func NewFixtureCase[I Testable, F any](testable I, fixture F) *FixtureCase[I, F] {
	c := new(FixtureCase[I, F])
	c.Testable = testable
	c.Fixture = fixture
	c.Assembler = types.NewAssembler(c)
	return c
}

type ServiceFixtureCaseFunc[I ServiceTestable, F any] func(c *ServiceFixtureCase[I, F], ctx context.Context)

type ServiceFixtureCase[I ServiceTestable, F any] struct {
	*types.Assembler[ServiceFixtureCase[I, F], *ServiceFixtureCase[I, F]]
	T        *testing.T
	Testable I
	Fixture  F
	Func     ServiceFixtureCaseFunc[I, F]
}

func (c *ServiceFixtureCase[I, F]) WithFunc(fn ServiceFixtureCaseFunc[I, F]) *ServiceFixtureCase[I, F] {
	c.Func = fn
	return c
}

func (c *ServiceFixtureCase[I, F]) Clone() *ServiceFixtureCase[I, F] {
	return deepcopy.Copy(c).(*ServiceFixtureCase[I, F])
}

func (c *ServiceFixtureCase[I, F]) DeepCopy() any {
	result := &ServiceFixtureCase[I, F]{
		Assembler: &types.Assembler[ServiceFixtureCase[I, F], *ServiceFixtureCase[I, F]]{
			Setups: deepcopy.Copy(c.Assembler.Setups).([]types.NamedSetupFunc[ServiceFixtureCase[I, F], *ServiceFixtureCase[I, F]]),
		},
		Testable: deepcopy.Copy(c.Testable).(I),
		Fixture:  deepcopy.Copy(c.Fixture).(F),
		Func:     c.Func,
	}
	result.Assembler.Target = result
	return result
}

func (c *ServiceFixtureCase[I, F]) Test(t *testing.T) {
	c.T = t
	c.Assemble()
	c.Testable.Test(t, func(t *testing.T, ctx context.Context) {
		c.Func(c, ctx)
	})
}

func NewServiceFixtureCase[I ServiceTestable, F any](testable I, fixture F) *ServiceFixtureCase[I, F] {
	c := new(ServiceFixtureCase[I, F])
	c.Testable = testable
	c.Fixture = fixture
	c.Assembler = types.NewAssembler(c)
	return c
}
