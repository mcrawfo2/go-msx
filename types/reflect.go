package types

import "reflect"

func Instantiate(t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.Slice:
		return reflect.MakeSlice(t, 1, 1).Interface()
	case reflect.Array:
		return reflect.Zero(t).Interface()
	case reflect.Map:
		return reflect.MakeMap(t).Interface()
	case reflect.Chan:
		return reflect.MakeChan(t, 0).Interface()
	default:
		return reflect.Zero(t).Interface()
	}
}
