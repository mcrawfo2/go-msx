package topics

import (
	"encoding/json"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

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
	return time.Time(t).In(time.UTC).Format(TimeLayout), nil
}

func (t *Time) UnmarshalText(data string) (err error) {
	var v time.Time
	v, err = time.Parse(TimeLayout, data)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

func (t Time) ToTimeTime() time.Time {
	return time.Time(t)
}

func (t Time) String() string {
	return time.Time(t).Format(TimeLayout)
}
