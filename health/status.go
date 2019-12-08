package health

const (
	StatusUp      Status = 1
	StatusDown    Status = 0
	StatusUnknown Status = -1
)

type Status int

func (s Status) String() string {
	switch s {
	case StatusUp:
		return "UP"
	case StatusDown:
		return "DOWN"
	default:
		return "UNKNOWN"
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func ParseStatus(status string) Status {
	switch status {
	case "UP":
		return StatusUp
	case "DOWN":
		return StatusDown
	default:
		return StatusUnknown
	}
}

func (s Status) Aggregate(other Status) Status {
	if (s == StatusDown && other == StatusUnknown) || s == StatusUp {
		return other
	}
	return s
}
