package types

import (
	"context"
	"sync"
)

type PartitionedMutex struct {
	mtx        sync.Mutex
	partitions map[interface{}]*sync.Mutex
}

func (m *PartitionedMutex) Locker(partition interface{}) *sync.Mutex {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var mtx *sync.Mutex
	var ok bool
	if mtx, ok = m.partitions[partition]; !ok {
		mtx = new(sync.Mutex)
		m.partitions[partition] = mtx
	}

	return mtx
}

func (m *PartitionedMutex) Lock(partition interface{}) {
	m.Locker(partition).Lock()
}

func (m *PartitionedMutex) Unlock(partition interface{}) {
	m.Locker(partition).Unlock()
}

func (m *PartitionedMutex) WithPartitionLock(partition interface{}, ctx context.Context, action ActionFunc) error {
	m.Lock(partition)
	defer m.Unlock(partition)

	return action(ctx)
}

func NewPartitionedMutex() *PartitionedMutex {
	return &PartitionedMutex{
		partitions: make(map[interface{}]*sync.Mutex),
	}
}
