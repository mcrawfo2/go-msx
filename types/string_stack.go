package types

type StringStack []string

func (s StringStack) Contains(value string) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

func (s StringStack) Push(value string) StringStack {
	return append(s, value)
}

func (s StringStack) Pop() StringStack {
	return s[:len(s)-1]
}

func (s StringStack) Peek() string {
	return s[len(s)-1]
}
