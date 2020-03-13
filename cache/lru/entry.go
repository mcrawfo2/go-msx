package lru

type entry struct {
	key     string
	expires int64
	index   int
	value   interface{}
}
