package types

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

// https://stackoverflow.com/questions/48050945/how-to-unmarshal-json-into-durations
type Duration time.Duration

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if len(b) == 4 && string(b) == "null" {
		return nil
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		var err error
		dur, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(dur)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func NewDuration(duration time.Duration) Duration {
	return Duration(duration)
}
