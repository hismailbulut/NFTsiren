package mutex

import (
	"math/rand"
	"sync"
	"testing"
)

const randMax = 10000

func rint() int {
	return rand.Intn(randMax)
}

func Benchmark_SyncMap(b *testing.B) {
	var m sync.Map
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rint()
			i, ok := m.Load(key)
			if ok {
				m.Store(key, i.(int)+1)
			} else {
				m.Store(key, 1)
			}
		}
	})
}

func Benchmark_Map(b *testing.B) {
	m := NewMap[int, int]()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rint()
			i, ok := m.Load(key)
			if ok {
				m.Store(key, i+1)
			} else {
				m.Store(key, 1)
			}
		}
	})
}

// Benchmark_SyncMap-12     5461548          219.4 ns/op            28 B/op          2 allocs/op
// Benchmark_Map-12         6444844          184.2 ns/op             0 B/op          0 allocs/op         Mutex
// Benchmark_Map-12         7533050          154.7 ns/op             0 B/op          0 allocs/op         RWMutex
// Benchmark_Map-12         8000852          146.9 ns/op             0 B/op          0 allocs/op         RWMutex - no defer
// We have to use deferred unlock in a map because if the map not initialized with NewMap then it will panic
// and using non-deferred operations most likely cause a deadlock
