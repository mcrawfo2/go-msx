// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

type SetupFunc[I any, PI *I] func(PI)

type NamedSetupFunc[I any, PI *I] struct {
	Name string
	Func SetupFunc[I, PI]
}

func (f NamedSetupFunc[I, PI]) IsAnonymous() bool {
	return f.Name == ""
}

type Assembler[I any, PI *I] struct {
	Setups []NamedSetupFunc[I, PI]
	Target PI
}

func (s *Assembler[I, PI]) Assemble() {
	for _, setup := range s.Setups {
		setup.Func(s.Target)
	}
}

// WithSetup adds an anonymous Setup function
func (s *Assembler[I, PI]) WithSetup(setupFunc SetupFunc[I, PI]) PI {
	namedSetup := NamedSetupFunc[I, PI]{
		Func: setupFunc,
	}
	s.Setups = append(s.Setups, namedSetup)
	return s.Target
}

// WithoutSetup removes all anonymous Setup functions
func (s *Assembler[I, PI]) WithoutSetup() PI {
	var indices []int
	for i, namedSetup := range s.Setups {
		if namedSetup.IsAnonymous() {
			indices = append(indices, i)
		}
	}
	s.Setups = Slice[NamedSetupFunc[I, PI]](s.Setups).RemoveAll(indices...)
	return s.Target
}

// WithNamedSetup adds or replaces a named Setup function
func (s *Assembler[I, PI]) WithNamedSetup(name string, setupFunc SetupFunc[I, PI]) PI {
	named := NamedSetupFunc[I, PI]{
		Name: name,
		Func: setupFunc,
	}

	for i := range s.Setups {
		if s.Setups[i].Name == name {
			s.Setups[i] = named
			return s.Target
		}
	}

	s.Setups = append(s.Setups, named)
	return s.Target
}

// WithoutNamedSetup removes a named Setup function
func (s *Assembler[I, PI]) WithoutNamedSetup(name string) PI {
	for i := range s.Setups {
		if s.Setups[i].Name == name {
			s.Setups = Slice[NamedSetupFunc[I, PI]](s.Setups).RemoveAll(i)
			break
		}
	}
	return s.Target
}

// NewAssembler creates a new Assembler with the specified target
func NewAssembler[I any, PI *I](target PI) *Assembler[I, PI] {
	result := new(Assembler[I, PI])
	result.Target = target
	return result
}
