package mutex

import "sync"

type Counter struct {
	m sync.Mutex
	v int64
}

// Adds 1 to the value and returns the previous value
func (counter *Counter) Increment() {
	counter.m.Lock()
	counter.v++
	counter.m.Unlock()
}

// Subtracts 1 from the value and returns the previous value
// Does nothing when value is zero
func (counter *Counter) Decrement() {
	counter.m.Lock()
	counter.v--
	counter.m.Unlock()
}

func (counter *Counter) Value() int64 {
	counter.m.Lock()
	v := counter.v
	counter.m.Unlock()
	return v
}
