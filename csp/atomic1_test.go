package main

/* 针对共享整形变量的无锁读写 */

import (
	"sync"
	"sync/atomic"
	"testing"
)

var n1 int64

func addSyncByAtomic(delta int64) int64 {
	return atomic.AddInt64(&n1, delta)

}
func readSyncByAtomic() int64 {
	return atomic.LoadInt64(&n1)
}

var n2 int64
var rwmu sync.RWMutex

func addSyncByRWMutex(delta int64) {
	rwmu.Lock()
	n2 += delta
	rwmu.Unlock()
}

func readSyncByRWMutex() int64 {
	var n int64
	rwmu.RLock()
	n = n2
	rwmu.RUnlock()
	return n
}

func BenchmarkAddSyncByAtomic(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			addSyncByAtomic(1)
		}
	})
}

func BenchmarkReadSyncByAtomic(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			readSyncByAtomic()
		}
	})
}

func BenchmarkAddSyncByRWMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			addSyncByRWMutex(1)
		}
	})
}

func BenchmarkReadSyncByRWMutex(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			readSyncByRWMutex()
		}
	})
}

// go test -bench . atomic1_test.go -cpu 2 >> atomic1_test.go
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkAddSyncByAtomic-2     	66345073	        18.49 ns/op
// BenchmarkReadSyncByAtomic-2    	1000000000	         0.6516 ns/op
// BenchmarkAddSyncByRWMutex-2    	37699429	        29.84 ns/op
// BenchmarkReadSyncByRWMutex-2   	50949586	        25.77 ns/op
// PASS
// ok  	command-line-arguments	4.462s

// go test -bench . atomic1_test.go -cpu 8 >> atomic1_test.go
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkAddSyncByAtomic-8     	82798820	        14.28 ns/op
// BenchmarkReadSyncByAtomic-8    	1000000000	         0.2071 ns/op
// BenchmarkAddSyncByRWMutex-8    	21368559	        56.39 ns/op
// BenchmarkReadSyncByRWMutex-8   	32455552	        36.64 ns/op
// PASS
// ok  	command-line-arguments	3.921s

// go test -bench . atomic1_test.go -cpu 16 >> atomic1_test.go
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkAddSyncByAtomic-16      	75321930	        15.86 ns/op
// BenchmarkReadSyncByAtomic-16     	1000000000	         0.2082 ns/op
// BenchmarkAddSyncByRWMutex-16     	19459443	        62.12 ns/op
// BenchmarkReadSyncByRWMutex-16    	31138516	        36.73 ns/op
// PASS
// ok  	command-line-arguments	3.904s

// 利用原子操作的无锁并发写的性能随着并发量增大几乎保持恒定, 利用原子操作
// 的无锁并发读的性能随着并发量增大有持续提升的趋势, 并且以 cpu 8 为例,
// 原子锁的无锁并发读的性能约为读锁的 180 倍;
