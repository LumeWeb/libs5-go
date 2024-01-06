package structs

import (
	"github.com/emirpasic/gods/maps/hashmap"
	"log"
	"sync"
)

type Map struct {
	*hashmap.Map
	mutex *sync.RWMutex
}

func NewMap() *Map {
	return &Map{
		Map:   hashmap.New(),
		mutex: &sync.RWMutex{},
	}
}

func (m *Map) Get(key interface{}) (value interface{}, found bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Map.Get(key)
}

func (m *Map) GetInt(key interface{}) (value *int) {
	val, found := m.Get(key)

	if !found {
		return nil
	}

	if intValue, ok := val.(int); ok {
		value = &intValue
	} else {
		log.Fatalf("value is not an int: %v", val)
	}

	return value
}

func (m *Map) GetString(key interface{}) (value *string) {
	val, found := m.Get(key)

	if !found {
		return nil
	}

	if _, ok := val.(string); ok {
		value = val.(*string)
	} else {
		log.Fatalf("value is not a string: %v", value)
	}

	return
}

func (m *Map) Put(key interface{}, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Map.Put(key, value)
}

func (m *Map) PutInt(key interface{}, value int) {
	m.Put(key, value)
}

func (m *Map) Remove(key interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Map.Remove(key)
}

func (m *Map) Keys() []interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Map.Keys()
}

func (m *Map) Values() []interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Map.Values()
}

func (m *Map) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Map.Size()
}

func (m *Map) Empty() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Map.Empty()
}

func (m *Map) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Map.Clear()
}

func (m *Map) String() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Map.String()
}

func (m *Map) GetKey(value interface{}) (key interface{}, found bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Map.Get(value)
}

func (m *Map) Contains(value interface{}) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, has := m.Map.Get(value)

	return has
}
