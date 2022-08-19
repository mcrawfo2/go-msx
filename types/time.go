// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"time"
)

const timeLayout = "2006-01-02T15:04:05.999999999Z"
const timeLayout2 = time.RFC3339Nano

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	str, err := t.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(str)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var str *string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	return t.UnmarshalText(*str)
}

func (t Time) MarshalText() (string, error) {
	return time.Time(t).In(time.UTC).Format(timeLayout), nil
}

func (t *Time) UnmarshalText(data string) (err error) {
	var v time.Time
	v, err = time.Parse(timeLayout, data)
	if err != nil {
		// Be slightly more liberal in accepting TZ as offset
		v, err = time.Parse(timeLayout2, data)
		if err != nil {
			return err
		}
	}
	*t = Time(v)
	return nil
}

func (t *Time) String() string {
	if t == nil {
		return "<nil>"
	}
	return time.Time(*t).Format(timeLayout)
}

// DEPRECATED
func (t Time) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return gocql.Marshal(info, time.Time(t))
}

// DEPRECATED
func (t *Time) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	var v time.Time
	err := gocql.Unmarshal(info, data, &v)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

func (t Time) ToTimeTime() time.Time {
	return time.Time(t)
}

func (t Time) SwaggerSchemaJson() string {
	return `{
		"type": "string",
		"format": "date-time"
	}`
}

func (t Time) DeepCopy() interface{} {
	return t
}

func ParseTime(v string) (t Time, err error) {
	err = t.UnmarshalText(v)
	return
}

func NewTime(t time.Time) Time {
	return Time(t)
}

func MaxTime() Time {
	return Time(time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC))
	// return Time(time.Unix(1<<63-62135596801, 999999999)) // actual golang max, cassandra doesn't like so much
}

func MustParseTime(ts string) Time {
	result, err := ParseTime(ts)
	if err != nil {
		panic(err)
	}
	return result
}
