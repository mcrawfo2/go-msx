package types

type Slice[I any] []I

func (s Slice[I]) AnySlice() (results []any) {
	results = make([]any, 0, len(s))
	for _, v := range s {
		results = append(results, v)
	}
	return results
}
