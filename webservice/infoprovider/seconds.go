package infoprovider

import (
	"fmt"
	"time"
)

type epochSeconds float64

func (e epochSeconds) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.9f", float64(e))), nil
}

func newEpochSeconds(when time.Time) epochSeconds {
	return epochSeconds(float64(when.Unix()) + (float64(when.Nanosecond()) * 1e-9))
}
