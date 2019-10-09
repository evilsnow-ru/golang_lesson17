package dao

import (
	"golang_lesson17/internal/domain/model"
	"sync"
)

type mapStorage struct {
	lock       sync.RWMutex
	storageMap map[uint64]*model.Event
	index      uint64
}

func NewMapStorage() *mapStorage {
	return &mapStorage{
		lock:       sync.RWMutex{},
		storageMap: make(map[uint64]*model.Event),
		index:      0,
	}
}

func (ms *mapStorage) Store(event *model.Event) (uint64, error) {
	var id uint64
	ms.lock.Lock()
	defer ms.lock.Unlock()
	id = ms.index
	ms.index++
	ms.storageMap[id] = event
	return id, nil
}

func (ms *mapStorage) Get(id uint64) (*model.Event, error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()
	return ms.storageMap[id], nil
}

func (ms *mapStorage) Delete(id uint64) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	delete(ms.storageMap, id)
	return nil
}
