// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package tests

import (
	"bufio"
	"bytes"
	"os"
)

type DiffNode struct {
	Filename string
	Data     []byte
}

func (n *DiffNode) Lines() ([]string, error) {
	if n.Data == nil {
		data, err := os.ReadFile(n.Filename)
		if err != nil {
			return nil, err
		}

		n.Data = data
	}

	return n.ByteLines()
}

func (n *DiffNode) ByteLines() ([]string, error) {
	var matching []string
	s := bufio.NewScanner(bytes.NewReader(n.Data))
	s.Split(bufio.ScanLines)

	for s.Scan() {
		t := s.Text()
		matching = append(matching, t+"\n")
	}
	if s.Err() != nil {
		return nil, s.Err()
	}

	return matching, nil
}
