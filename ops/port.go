// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"github.com/pkg/errors"
	"reflect"
)

type Port struct {
	Type       string // Tag name
	StructType reflect.Type
	Fields     PortFields
}

func (p *Port) WithField(f *PortField) *Port {
	p.Fields = append(p.Fields, f)
	return p
}

func (p Port) NewStruct() *interface{} {
	if p.StructType == nil {
		return nil
	}

	// heap allocate a struct
	var result interface{}
	result = reflect.New(p.StructType).Interface()
	return &result
}

func NewPort(typ string, structType reflect.Type) (*Port, error) {
	if typ == "" {
		return nil, errors.Errorf("Invalid port type %q", typ)
	}

	if structType == nil {
		return nil, errors.Errorf("Null struct type")
	} else if structType.Kind() != reflect.Struct {
		return nil, errors.Errorf("Not a struct type: %v", structType)
	}

	return &Port{Type: typ, StructType: structType}, nil
}
