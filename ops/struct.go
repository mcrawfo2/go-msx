// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"github.com/fatih/structtag"
	"github.com/pkg/errors"
	"github.com/swaggest/refl"
	"reflect"
)

type StructFieldVisitor interface {
	VisitField(f reflect.StructField) error
	EnterAnonymousStructField(f reflect.StructField)
	ExitAnonymousStructField(f reflect.StructField)
}

var ErrRecursiveStruct = errors.New("Recursion found but not supported")

func WalkStruct(st reflect.Type, visitor StructFieldVisitor) error {
	return walkStruct(st, visitor, map[reflect.Type]struct{}{})
}

func walkStruct(st reflect.Type, visitor StructFieldVisitor, active map[reflect.Type]struct{}) error {
	if _, ok := active[st]; ok {
		return errors.Wrapf(ErrRecursiveStruct, "Recursive type: %+v", st)
	}
	active[st] = struct{}{}

	for i := 0; i < st.NumField(); i++ {
		sf := st.Field(i)

		deepType := refl.DeepIndirect(sf.Type)
		if sf.Anonymous && deepType.Kind() == reflect.Struct {
			visitor.EnterAnonymousStructField(sf)
			if err := walkStruct(deepType, visitor, active); err != nil {
				return err
			}
			visitor.ExitAnonymousStructField(sf)
		} else if err := visitor.VisitField(sf); err != nil {
			return err
		}
	}

	delete(active, st)
	return nil
}

func LookupTag(tags *structtag.Tags, key string) (*structtag.Tag, bool) {
	for _, tag := range tags.Tags() {
		if tag.Key == key {
			return tag, true
		}
	}
	return nil, false
}
