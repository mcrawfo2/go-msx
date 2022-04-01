// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import "container/list"

type StringQueue struct {
	*list.List
}

func (s StringQueue) Contains(value string) bool {
	for e := s.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == value {
			return true
		}
	}
	return false
}

func (s StringQueue) Push(value string) {
	s.PushBack(value)
}

func (s StringQueue) Pop() string {
	return s.Remove(s.Front()).(string)
}

func (s StringQueue) Peek() string {
	return s.Front().Value.(string)
}

func NewStringQueue() StringQueue {
	return StringQueue{
		list.New(),
	}
}
