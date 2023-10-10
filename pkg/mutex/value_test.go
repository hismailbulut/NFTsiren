package mutex

import (
	"sync/atomic"
	"testing"
)

func Benchmark_AtomicValue(b *testing.B) {
	var v atomic.Value
	v.Store(int(0)) // This is necessary for this
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			i := v.Load().(int) + 1
			v.Store(i)
		}
	})
}

func Benchmark_AtomicPointer(b *testing.B) {
	var v atomic.Pointer[int]
	zero := int(0)
	v.Store(&zero) // This is necessary for this
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			i := *v.Load() + 1
			v.Store(&i)
		}
	})
}

func Benchmark_Value(b *testing.B) {
	var v Value[int]
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			v.Store(v.Load() + 1)
		}
	})
}

// Benchmark_AtomicValue-12        33911619          32.92 ns/op          8 B/op          0 allocs/op
// Benchmark_AtomicPointer-12      35873688          30.51 ns/op          8 B/op          1 allocs/op
// Benchmark_Value-12              15039704          81.37 ns/op          0 B/op          0 allocs/op    ->   with defer statement
// Benchmark_Value-12              15808214          75.19 ns/op          0 B/op          0 allocs/op    ->   without deferring unlocks
