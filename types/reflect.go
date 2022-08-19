// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
