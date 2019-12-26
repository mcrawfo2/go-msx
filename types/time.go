package types

import (
	"encoding/json"
	"time"
)

const timeLayout = "2006-01-02T15:04:05.999999999Z"

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).In(time.UTC).Format(timeLayout))
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var str *string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	var v time.Time
	v, err = time.Parse(timeLayout, *str)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}
