package types

type StringSet map[string]struct{}

func (s StringSet) Contains(value string) bool {
	_, ok := s[value]
	return ok
}

func (s StringSet) Add(value string) {
	s[value] = struct{}{}
}

func (s StringSet) AddAll(values ...string) {
	for _, value := range values {
		s[value] = struct{}{}
	}
}

func (s StringSet) Values() []string {
	var result []string
	for k := range s {
		result = append(result, k)
	}
	return result
}
