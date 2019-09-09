//go:generate protoc -I/usr/local/include/ -I. --go_out=. messages.proto
package lesson17

import (
	"sync"
	"sync/atomic"
)

var storage map[uint64]*Event = make(map[uint64]*Event)
var lock sync.RWMutex
var index uint64

//AddEvent store new event
func AddEvent(event *Event) uint64 {
	id := atomic.AddUint64(&index, 0)
	lock.Lock()
	storage[id] = event
	lock.Unlock()
	return id
}

//DelEvent delete event from storage
func DelEvent(id uint64) {
	lock.Lock()
	delete(storage, id)
	lock.Unlock()
}

//GetEvent returns event by id
func GetEvent(id uint64) *Event {
	lock.RLock()
	event, ok := storage[id]
	lock.RUnlock()

	if ok {
		return event
	}
	return nil
}
