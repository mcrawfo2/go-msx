// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"context"
	"github.com/pkg/errors"
	"reflect"
)

var ErrUnknownValueType = errors.New("Unknown value type")

type HandlerArguments interface {
	GenerateArgument(ctx context.Context, t HandlerValueType) (reflect.Value, error)
}

type HandlerResults interface {
	HandleResult(t HandlerValueType, v reflect.Value) error
}

type HandlerContext interface {
	HandlerArguments
	HandlerResults
}

type contextKeyString string

const contextKeyHandlerContext = contextKeyString("HandlerContext")

func ContextWithHandlerContext(ctx context.Context, handlerContext HandlerContext) context.Context {
	return context.WithValue(ctx, contextKeyHandlerContext, handlerContext)
}

func handlerContextFromContext(ctx context.Context) HandlerContext {
	value, _ := ctx.Value(contextKeyHandlerContext).(HandlerContext)
	return value
}

// Handler represents an action which has (optional) arguments and return values
type Handler struct {
	fn  reflect.Value
	in  []HandlerValueType
	out []HandlerValueType
}

func (h *Handler) reflectIn(inr HandlerValueTypeReflector) (err error) {
	if inr == nil {
		// Use the default argument type reflector when unspecified
		inr = NewHandlerValueTypeReflector(
			DefaultHandlerArgumentValueTypeSet,
		)
	}

	funcType := h.fn.Type()

	var results []HandlerValueType
	for i := 0; i < funcType.NumIn(); i++ {
		var argType = funcType.In(i)

		var valueType HandlerValueType
		valueType = inr.ValueType(argType)
		if valueType.IsEmpty() {
			err = errors.Wrapf(ErrUnknownValueType, "Failed to identify type for handler argument %d: %v", i, argType)
			return
		}

		results = append(results, valueType)
	}

	h.in = results

	return
}

func (h *Handler) reflectOut(outr HandlerValueTypeReflector) (err error) {
	if outr == nil {
		// Use the default result type reflector when unspecified
		outr = NewHandlerValueTypeReflector(
			DefaultHandlerResultValueTypeSet,
		)
	}

	funcType := h.fn.Type()

	var results []HandlerValueType
	for i := 0; i < funcType.NumOut(); i++ {
		var argType = funcType.Out(i)

		var valueType HandlerValueType
		valueType = outr.ValueType(argType)
		if valueType.IsEmpty() {
			err = errors.Wrapf(ErrUnknownValueType, "Failed to identify type for handler result %d: %v", i, argType)
			return
		}

		results = append(results, valueType)
	}

	h.out = results

	return
}

func (h Handler) args(ctx context.Context) (results []reflect.Value, err error) {
	in := handlerContextFromContext(ctx).(HandlerArguments)
	if in == nil {
		in = ContextArguments{}
	}

	// TODO: Indirection

	results = make([]reflect.Value, len(h.in))
	for i, inValueType := range h.in {
		var v reflect.Value
		v, err = in.GenerateArgument(ctx, inValueType)
		if err != nil {
			return nil, err
		}

		results[i] = v
	}

	return
}

func (h Handler) results(ctx context.Context, values []reflect.Value) (err error) {
	out := handlerContextFromContext(ctx).(HandlerResults)
	if out == nil {
		// Make sure we handle error return value even when no result sink is supplied
		out = ErrorResults{}
	}

	// TODO: Indirection

	for i, outValueType := range h.out {
		err = out.HandleResult(outValueType, values[i])
		if err != nil {
			break
		}
	}

	return
}

func (h Handler) Call(ctx context.Context) error {
	argValues, err := h.args(ctx)
	if err != nil {
		return err
	}

	resultValues := h.fn.Call(argValues)

	err = h.results(ctx, resultValues)
	if err != nil {
		return err
	}

	return nil
}

func NewHandler(fn interface{}, inr, outr HandlerValueTypeReflector) (h *Handler, err error) {
	funcValue := reflect.ValueOf(fn)
	if funcValue.Kind() != reflect.Func {
		return nil, errors.Errorf("Expected handler function type, got %T", fn)
	}

	h = &Handler{
		fn: funcValue,
	}

	if err = h.reflectIn(inr); err != nil {
		return nil, err
	}

	if err = h.reflectOut(outr); err != nil {
		return nil, err
	}

	return
}

type HandlerValueTypeReflector interface {
	ValueType(t reflect.Type) HandlerValueType
}

type HandlerValueType struct {
	Indirections int
	ValueType    reflect.Type
}

func (h HandlerValueType) IsEmpty() bool {
	return h.ValueType == nil
}

func (h *HandlerValueType) IncIndirections() {
	h.Indirections++
}

func NewHandlerValueType(t reflect.Type) HandlerValueType {
	return HandlerValueType{
		Indirections: 0,
		ValueType:    t,
	}
}

// IndirectionHandlerValueTypeReflector matches both direct and indirect values of the nested reflector
type IndirectionHandlerValueTypeReflector struct {
	next HandlerValueTypeReflector
}

func (d IndirectionHandlerValueTypeReflector) ValueType(t reflect.Type) (result HandlerValueType) {
	// Direct
	result = d.next.ValueType(t)
	if !result.IsEmpty() {
		return
	} else if t.Kind() != reflect.Ptr {
		return
	}

	// Indirect
	result = d.ValueType(t.Elem())
	if result.IsEmpty() {
		return
	}

	result.IncIndirections()
	return
}

// ValueTypeReflector matches only direct values of the specified types
type ValueTypeReflector struct {
	types TypeSet
}

func (v ValueTypeReflector) ValueType(t reflect.Type) HandlerValueType {
	_, ok := v.types[t]
	if ok {
		return NewHandlerValueType(t)
	}

	return HandlerValueType{}
}

// NewValueTypeReflector creates a new ValueTypeReflector, matching only direct values of the specified types
func NewValueTypeReflector(typeSets ...TypeSet) HandlerValueTypeReflector {
	ts := make(TypeSet)
	for _, typeSet := range typeSets {
		ts = ts.With(typeSet)
	}

	return ValueTypeReflector{
		types: ts,
	}
}

// NewHandlerValueTypeReflector creates an IndirectionHandlerValueTypeReflector,
// matching both direct and indirect values of the specified types
func NewHandlerValueTypeReflector(typeSets ...TypeSet) HandlerValueTypeReflector {
	return IndirectionHandlerValueTypeReflector{
		next: NewValueTypeReflector(typeSets...),
	}
}

var errorValue error
var errorType = reflect.TypeOf(&errorValue).Elem()
var contextContextValue context.Context
var contextContextType = reflect.TypeOf(&contextContextValue).Elem()

var DefaultHandlerResultValueTypeSet = NewTypeSet(errorType)
var DefaultHandlerArgumentValueTypeSet = NewTypeSet(contextContextType)

type ContextArguments struct{}

func (c ContextArguments) GenerateArgument(ctx context.Context, t HandlerValueType) (reflect.Value, error) {
	switch t.ValueType {
	case contextContextType:
		return reflect.ValueOf(ctx), nil
	}

	return reflect.Value{}, nil
}

type ErrorResults struct{}

func (e ErrorResults) HandleResult(t HandlerValueType, v reflect.Value) error {
	if t.ValueType == errorType {
		err, _ := v.Interface().(error)
		return err
	}
	return nil
}

func (e ErrorResults) Finish() error {
	return nil
}

type TypeSet map[reflect.Type]struct{}

func (t TypeSet) With(o TypeSet) TypeSet {
	for k := range o {
		t[k] = struct{}{}
	}
	return t
}

func (t TypeSet) WithType(o reflect.Type) TypeSet {
	t[o] = struct{}{}
	return t
}

func (t TypeSet) WithTypes(types ...reflect.Type) TypeSet {
	for _, typ := range types {
		t.WithType(typ)
	}
	return t
}

func NewTypeSet(types ...reflect.Type) TypeSet {
	return make(TypeSet).WithTypes(types...)
}
