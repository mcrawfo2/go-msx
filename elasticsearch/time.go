package elasticsearch

import (
	"encoding/json"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.000Z"

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).UTC().Format(timeFormat))
}
