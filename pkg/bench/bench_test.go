package bench

import (
	"runtime"
	"runtime/debug"
	"testing"
)

func Benchmark_BenchmarkBegin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		end := Begin()
		end("test")
	}
}

func Benchmark_BenchmarkBeginParallel(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			end := Begin()
			end("test")
		}
	})
}

func Benchmark_RuntimeReadMemStats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var m1 runtime.MemStats
		runtime.ReadMemStats(&m1)
		_ = m1
	}
}

func Benchmark_DebugFreeOsMemory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		debug.FreeOSMemory()
	}
}

func Benchmark_RuntimeGC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.GC()
	}
}
