// This package contains encapsulation of generic values with mutex lock
package mutex

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Do not copy Value
// This value should only used with basic and concrete types
// Be carefull with pointer types. This will only return the pointer safely, does not protect the underlying operations
type Value[T any] struct {
	guard  sync.RWMutex
	handle T
}

func (value *Value[T]) Store(v T) {
	value.guard.Lock()
	value.handle = v
	value.guard.Unlock()
}

func (value *Value[T]) Load() (v T) {
	value.guard.RLock()
	v = value.handle
	value.guard.RUnlock()
	return
}

func (value *Value[T]) String() string {
	v := value.Load()
	return fmt.Sprintf("mutex.Value[%T](%v)", v, v)
}

func (value *Value[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Load())
}

func (value *Value[T]) UnmarshalJSON(data []byte) error {
	zero := *new(T)
	err := json.Unmarshal(data, &zero)
	if err != nil {
		return err
	}
	value.Store(zero)
	return nil
}
