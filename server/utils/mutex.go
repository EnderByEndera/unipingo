package utils

import (
	"errors"
	"sync"
)

type KeyedMutex struct {
	mutexes sync.Map
}

func (m *KeyedMutex) Lock(key interface{}) {
	value, _ := m.mutexes.LoadOrStore(key, &sync.Mutex{}) // 如果没有值，则存储一个新的mutex
	mtx := value.(*sync.Mutex)
	mtx.Lock()
}

func (m *KeyedMutex) Unlock(key interface{}) (err error) {
	mtxInterface, ok := m.mutexes.Load(key)
	mtx := mtxInterface.(*sync.Mutex)
	if !ok {
		err = errors.New("key not found in mutex map")
		return
	}
	mtx.Unlock()
	return
}
