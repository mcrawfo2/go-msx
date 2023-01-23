// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package payloads

import (
	"github.com/swaggest/jsonschema-go"
	"strings"
)

// JsonSchemaVisitor is a visitor for schemas.  It should return false to stop traversal.
type JsonSchemaVisitor func(schema *jsonschema.Schema) bool

// JsonSchemaPredicate is a predicate for schemas.  It should return true if the schema matches.
type JsonSchemaPredicate func(schema *jsonschema.Schema) bool

// WalkJsonSchema visits all schema nodes within a root jsonschema.Schema
func WalkJsonSchema(root *jsonschema.Schema, visitor JsonSchemaVisitor) {
	_ = walkJsonSchema(root, visitor)
}

func walkJsonSchema(root *jsonschema.Schema, visitor JsonSchemaVisitor) bool {
	if root == nil {
		return false
	}

	if visitor(root) {
		return true
	}

	if walkJsonSchemaOrBool(root.AdditionalItems, visitor) {
		return true
	}

	if walkItems(root.Items, visitor) {
		return true
	}

	if walkJsonSchemaOrBool(root.Contains, visitor) {
		return true
	}

	if walkJsonSchemaOrBool(root.AdditionalProperties, visitor) {
		return true
	}

	if walkMapJsonSchemaOrBool(root.Definitions, visitor) {
		return true
	}

	if walkMapJsonSchemaOrBool(root.Properties, visitor) {
		return true
	}

	if walkMapJsonSchemaOrBool(root.PatternProperties, visitor) {
		return true
	}

	if walkMapDependencies(root.Dependencies, visitor) {
		return true
	}

	if walkJsonSchemaOrBool(root.PropertyNames, visitor) {
		return true
	}

	if walkJsonSchemaOrBool(root.If, visitor) {
		return true
	}

	if walkJsonSchemaOrBool(root.Then, visitor) {
		return true
	}

	if walkJsonSchemaOrBool(root.Else, visitor) {
		return true
	}

	if walkSliceJsonSchemaOrBool(root.AllOf, visitor) {
		return true
	}

	if walkSliceJsonSchemaOrBool(root.AnyOf, visitor) {
		return true
	}

	if walkSliceJsonSchemaOrBool(root.OneOf, visitor) {
		return true
	}

	if walkJsonSchemaOrBool(root.Not, visitor) {
		return true
	}

	return false
}

func walkItems(items *jsonschema.Items, visitor JsonSchemaVisitor) bool {
	if items == nil {
		return false
	}

	if items.SchemaOrBool != nil {
		return walkJsonSchemaOrBool(items.SchemaOrBool, visitor)
	}

	return walkSliceJsonSchemaOrBool(items.SchemaArray, visitor)
}

func walkMapDependencies(dependencies map[string]jsonschema.DependenciesAdditionalProperties, visitor JsonSchemaVisitor) bool {
	for _, v := range dependencies {
		if walkJsonSchemaOrBool(v.SchemaOrBool, visitor) {
			return true
		}
	}
	return false
}

func walkMapJsonSchemaOrBool(iterable map[string]jsonschema.SchemaOrBool, visitor JsonSchemaVisitor) bool {
	var results = map[string]*jsonschema.SchemaOrBool{}
	var done = false
	for k, v := range iterable {
		tv := v
		pv := &tv
		done = walkJsonSchemaOrBool(pv, visitor)
		results[k] = pv
		if done {
			break
		}
	}
	for k, v := range results {
		iterable[k] = *v
	}
	return done
}

func walkSliceJsonSchemaOrBool(iterable []jsonschema.SchemaOrBool, visitor JsonSchemaVisitor) bool {
	var results = make([]*jsonschema.SchemaOrBool, len(iterable))
	var done = false
	for i, v := range iterable {
		tv := v
		pv := &tv
		done = walkJsonSchemaOrBool(pv, visitor)
		results[i] = pv
		if done {
			break
		}
	}
	for i, v := range results {
		if v == nil {
			break
		}
		iterable[i] = *v
	}
	return done
}

func walkJsonSchemaOrBool(root *jsonschema.SchemaOrBool, visitor JsonSchemaVisitor) bool {
	if root == nil {
		return false
	}

	return walkJsonSchema(root.TypeObject, visitor)
}

func VisitJsonSchemaWhen(p JsonSchemaPredicate, v JsonSchemaVisitor) JsonSchemaVisitor {
	return func(schema *jsonschema.Schema) bool {
		if p(schema) {
			return v(schema)
		} else {
			return false
		}
	}
}

func JsonSchemaHasRef(ref string) JsonSchemaPredicate {
	return func(schema *jsonschema.Schema) bool {
		return schema.Ref != nil && ref == *schema.Ref
	}
}

func JsonSchemaHasRefPrefix(refPrefix string) JsonSchemaPredicate {
	return func(schema *jsonschema.Schema) bool {
		return schema.Ref != nil && strings.HasPrefix(*schema.Ref, refPrefix)
	}
}
