package lru

import (
	"container/heap"
)

type entryHeap []*entry

func (e entryHeap) Len() int { return len(e) }

func (e entryHeap) Less(i, j int) bool {
	return e[i].expires < e[j].expires
}

func (e entryHeap) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
	e[i].index = i
	e[j].index = j
}

func (e *entryHeap) Push(x interface{}) {
	n := len(*e)
	item := x.(*entry)
	item.index = n
	*e = append(*e, item)
}

func (e *entryHeap) Pop() interface{} {
	old := *e
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*e = old[0 : n-1]
	return item
}

func (e entryHeap) countExpired(now int64) int {
	var j = 0
	for i := len(e) - 1; i >= 0; i-- {
		if e[i].expires > now {
			break
		}
		j++
	}
	return j
}

func (e *entryHeap) expire(expires int64, keys []string) int {
	expiredCount := e.countExpired(expires)
	if expiredCount > len(keys) {
		expiredCount = len(keys)
	}

	old := *e
	n := len(old)

	j := 0
	for i := n - expiredCount; i < n; i++ {
		item := old[i]
		old[i] = nil
		keys[j] = item.key
		item.index = -1
	}

	*e = old[0 : n-expiredCount]
	return expiredCount
}

func (e *entryHeap) update(item *entry, expires int64) {
	item.expires = expires
	heap.Fix(e, item.index)
}
