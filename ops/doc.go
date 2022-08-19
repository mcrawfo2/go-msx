// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import "cto-github.cisco.com/NFV-BU/go-msx/types"

type Documentor[I any] interface {
	DocType() string
	Document(*I) error
}

type Documented[I any] interface {
	Documentor(pred DocumentorPredicate[I]) Documentor[I]
}

type DocumentorPredicate[I any] func(doc Documentor[I]) bool

type Documentors[I any] []Documentor[I]

func (o Documentors[I]) WithDocumentor(doc Documentor[I]) Documentors[I] {
	return append(o, doc)
}

func (o Documentors[I]) Documentor(pred DocumentorPredicate[I]) Documentor[I] {
	for _, doc := range o {
		if pred(doc) {
			return doc
		}
	}
	return nil
}

func DocumentorWithType[I any](o Documented[I], docType string) types.Optional[Documentor[I]] {
	pred := WithDocType[I](docType)
	doc := o.Documentor(pred)
	if doc == nil {
		return types.OptionalEmpty[Documentor[I]]()
	}
	return types.OptionalOf(doc)
}

func WithDocType[I any](docType string) DocumentorPredicate[I] {
	return func(d Documentor[I]) bool {
		return d.DocType() == docType
	}
}

type DocumentElementMutator[I any] func(*I)
