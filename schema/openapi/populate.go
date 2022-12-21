// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/refl"
	"reflect"
	"strings"
)

func PopulateFieldsFromTags(s *openapi3.Schema, tag reflect.StructTag) (err error) {
	if err = refl.PopulateFieldsFromTags(s, tag); err != nil {
		err = errors.Wrap(err, "Failed to populate parameter fields from struct tags")
		return
	}

	if err = PopulateInterfaceFieldsFromTags(s, tag); err != nil {
		err = errors.Wrap(err, "Failed to populate parameter fields from struct tags")
		return
	}

	if err = PopulateEnumFromTags(s, tag); err != nil {
		err = errors.Wrap(err, "Failed to parse enumeration from struct tags")
		return
	}

	return
}

func PopulateInterfaceFieldsFromTags(s *openapi3.Schema, tag reflect.StructTag) error {
	pv := reflect.ValueOf(s).Elem()
	pt := pv.Type()

	var errs = make(types.ErrorMap)

	for i := 0; i < pv.NumField(); i++ {
		ptf := pt.Field(i)
		tagName := strings.ToLower(ptf.Name[0:1]) + ptf.Name[1:]

		switch tagName {
		case "const":
		case "default":
		default:
			continue
		}

		_, ok := tag.Lookup(tagName)
		if !ok {
			continue
		}

		st := s.Type
		if st == nil || *st == openapi3.SchemaTypeArray || *st == openapi3.SchemaTypeObject {
			continue
		}

		pvf := pv.Field(i)

		var err error
		var tv interface{}

		switch *st {
		case openapi3.SchemaTypeString:
			var v string
			refl.ReadStringTag(tag, tagName, &v)
			tv = v
		case openapi3.SchemaTypeNumber:
			var v float64
			err = refl.ReadFloatTag(tag, tagName, &v)
			tv = v
		case openapi3.SchemaTypeBoolean:
			var v bool
			err = refl.ReadBoolTag(tag, tagName, &v)
			tv = v
		case openapi3.SchemaTypeInteger:
			var v int64
			err = refl.ReadIntTag(tag, tagName, &v)
			tv = v
		}

		if err != nil {
			errs[ptf.Name] = err
		} else {
			pvf.Set(reflect.ValueOf(&tv))
		}
	}

	return errs.Filter()
}

// PopulateEnumFromTags loads enum from struct tag (comma-separated string).
func PopulateEnumFromTags(s *openapi3.Schema, tag reflect.StructTag) error {
	enumTag, ok := tag.Lookup("enum")
	if !ok || enumTag == "" {
		return nil
	}

	var items []interface{}

	st := s.Type
	if st == nil || *st == openapi3.SchemaTypeArray || *st == openapi3.SchemaTypeObject {
		// Only scalars supported currently
		return nil
	}

	errs := types.ErrorList{}
	es := strings.Split(enumTag, ",")
	items = make([]interface{}, 0, len(es))

	for _, e := range es {
		var err error
		var tv interface{}

		switch *st {
		case openapi3.SchemaTypeString:
			tv, err = cast.ToStringE(e)
		case openapi3.SchemaTypeNumber:
			tv, err = cast.ToFloat64E(e)
		case openapi3.SchemaTypeBoolean:
			tv, err = cast.ToBoolE(e)
		case openapi3.SchemaTypeInteger:
			tv, err = cast.ToInt64E(e)
		}

		if err != nil {
			errs = append(errs, err)
		} else {
			items = append(items, tv)
		}
	}

	if len(items) > 0 {
		s.Enum = items
	}

	return nil
}
