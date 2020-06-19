package ringmap

import (
	"github.com/elliotchance/orderedmap"
)

// RingMap the ordered map data structure
type RingMap struct {
	orderedMap *orderedmap.OrderedMap
	capacity   int
}

// NewRingMap creates a new ordered map with a maximum size
func NewRingMap(capacity int) *RingMap {
	return &RingMap{
		orderedMap: orderedmap.NewOrderedMap(),
		capacity:   capacity,
	}
}

// Get returns the value for a key. If the key does not exist, the second return
// parameter will be false and the value will be nil.
func (m *RingMap) Get(key interface{}) (interface{}, bool) {
	return m.orderedMap.Get(key)
}

// Set will set (or replace) a value for a key. If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).  If a new key is being added and the map is
// full, then the front element will be deleted to make room for the new element.
func (m *RingMap) Set(key, value interface{}) bool {
	_, didExist := m.Get(key)

	if !didExist {
		if m.IsFull() {
			m.Delete(m.Front().Key)
		}
	}
	m.orderedMap.Set(key, value)

	return !didExist
}

// Put will set a value for a key. If the key already exists, it will be deleted
// from and a recreated at the end of the list.  If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).  If a new key is being added and the map is
// full, then the front element will be deleted to make room for the new element.
func (m *RingMap) Put(key, value interface{}) bool {
	_, didExist := m.Get(key)

	if didExist {
		m.Delete(m.Front().Key)
	}
	return m.orderedMap.Set(key, value)
}

// GetOrDefault returns the value for a key. If the key does not exist, returns
// the default value instead.
func (m *RingMap) GetOrDefault(key, defaultValue interface{}) interface{} {
	return m.orderedMap.GetOrDefault(key, defaultValue)
}

// Len returns the number of elements in the map.
func (m *RingMap) Len() int {
	return m.orderedMap.Len()
}

// Capacity returns the capacity of the map
func (m *RingMap) Capacity() int {
	return m.capacity
}

// IsFull returns true if the number of elements in the map is Capacity()
func (m *RingMap) IsFull() bool {
	return m.orderedMap.Len() == m.capacity
}

// Keys returns all of the keys in the order they were inserted. If a key was
// replaced it will retain the same position. To ensure most recently set keys
// are always at the end you must always Delete before Set.
func (m *RingMap) Keys() (keys []interface{}) {
	return m.orderedMap.Keys()
}

// Delete will remove a key from the map. It will return true if the key was
// removed (the key did exist).
func (m *RingMap) Delete(key interface{}) (didDelete bool) {
	return m.orderedMap.Delete(key)
}

// Front will return the element that is the first (oldest Set element). If
// there are no elements this will return nil.
func (m *RingMap) Front() *orderedmap.Element {
	return m.orderedMap.Front()
}

// Back will return the element that is the last (most recent Set element). If
// there are no elements this will return nil.
func (m *RingMap) Back() *orderedmap.Element {
	return m.orderedMap.Back()
}
