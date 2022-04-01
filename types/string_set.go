// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

type StringSet map[string]struct{}

func (s StringSet) Contains(value string) bool {
	_, ok := s[value]
	return ok
}

func (s StringSet) Add(value string) {
	s[value] = struct{}{}
}

func (s StringSet) AddAll(values ...string) {
	for _, value := range values {
		s[value] = struct{}{}
	}
}

func (s StringSet) Values() []string {
	var result []string
	for k := range s {
		result = append(result, k)
	}
	return result
}

func (s StringSet) Sub(other StringSet) []string {
	var results []string
	for k := range s {
		if !other.Contains(k) {
			results = append(results, k)
		}
	}
	return results
}
