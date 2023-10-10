package mutex

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// Do not copy Map
type Map[K comparable, V any] struct {
	guard  sync.RWMutex
	handle map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		handle: make(map[K]V),
	}
}

func (m *Map[K, V]) Store(key K, value V) {
	m.guard.Lock()
	defer m.guard.Unlock()
	m.handle[key] = value
}

func (m *Map[K, V]) Load(key K) (V, bool) {
	m.guard.RLock()
	defer m.guard.RUnlock()
	value, ok := m.handle[key]
	return value, ok
}

func (m *Map[K, V]) Has(key K) bool {
	m.guard.RLock()
	defer m.guard.RUnlock()
	_, ok := m.handle[key]
	return ok
}

func (m *Map[K, V]) Delete(key K) {
	m.guard.Lock()
	defer m.guard.Unlock()
	delete(m.handle, key)
}

func (m *Map[K, V]) Clear() {
	m.guard.Lock()
	defer m.guard.Unlock()
	m.handle = make(map[K]V)
}

func (m *Map[K, V]) Length() int {
	m.guard.RLock()
	defer m.guard.RUnlock()
	return len(m.handle)
}

// This is only for readonly access
// Do not modify map and don't call any map function in range function
func (m *Map[K, V]) Range(fn func(index int, key K, value V)) {
	m.guard.RLock()
	defer m.guard.RUnlock()
	i := 0
	for k, v := range m.handle {
		fn(i, k, v)
		i++
	}
}

func (m *Map[K, V]) String() string {
	s := strings.Builder{}
	l := m.Length()
	if l > 10 {
		s.WriteString(fmt.Sprintf("%d pairs", l))
	} else {
		m.Range(func(index int, key K, value V) {
			s.WriteString(fmt.Sprintf("%v=%v", key, value))
			if index != l-1 {
				s.WriteByte(' ')
			}
		})
	}
	return fmt.Sprintf("mutex.Map[%T, %T](%s)", *new(K), *new(V), s.String())
}

func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	m.guard.Lock()
	defer m.guard.Unlock()
	return json.Marshal(m.handle)
}

func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	m.guard.Lock()
	defer m.guard.Unlock()
	return json.Unmarshal(data, &m.handle)
}
