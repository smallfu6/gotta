package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

/*
	对共享自定义类型变量的无锁读写
	atomic通过Value类型的装拆箱操作实现了对任意自定义类型的
	原子操作(Load和Store), 从而实现对共享自定义类型变量无锁读写的支持;
*/

type Config struct {
	sync.RWMutex
	data string
}

func BenchmarkRWMutexSet(b *testing.B) {
	config := Config{}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			config.Lock()
			config.data = "hello"
			config.Unlock()
		}
	})
}

func BenchmarkRWMutexGet(b *testing.B) {
	config := Config{data: "hello"}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			config.RLock()
			_ = config.data
			config.RUnlock()
		}
	})
}

func BenchmarkAtomicSet(b *testing.B) {
	var config atomic.Value
	c := Config{data: "hello"}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			config.Store(c)
			// 这里没有用到 Config.RWMutex, 不需要考虑锁类型的复制问题
		}
	})
}

func BenchmarkAtomicGet(b *testing.B) {
	var config atomic.Value
	config.Store(Config{data: "hello"})
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = config.Load().(Config)
		}
	})
}

// go test -bench . atomic2_test.go -cpu 2
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkRWMutexSet-2   	37391205	        30.33 ns/op	       0 B/op	       0 allocs/op
// BenchmarkRWMutexGet-2   	34338184	        34.71 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAtomicSet-2    	32154892	        36.52 ns/op	      48 B/op	       1 allocs/op
// BenchmarkAtomicGet-2    	1000000000	         0.6470 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	command-line-arguments	4.325s

// go test -bench . atomic2_test.go -cpu 8
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkRWMutexSet-8   	21476626	        53.85 ns/op	       0 B/op	       0 allocs/op
// BenchmarkRWMutexGet-8   	31414195	        36.34 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAtomicSet-8    	44187176	        27.41 ns/op	      48 B/op	       1 allocs/op
// BenchmarkAtomicGet-8    	1000000000	         0.3519 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	command-line-arguments	4.024s

// go test -bench . atomic2_test.go -cpu 16
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkRWMutexSet-16    	21032923	        59.59 ns/op	       0 B/op	       0 allocs/op
// BenchmarkRWMutexGet-16    	31384052	        36.37 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAtomicSet-16     	39095743	        30.70 ns/op	      48 B/op	       1 allocs/op
// BenchmarkAtomicGet-16     	1000000000	         0.3954 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	command-line-arguments	4.159s

// 利用原子操作的无锁并发写的性能随着并发量的增大而小幅下降; 利用原子操作
// 的无锁并发读的性能随着并发量增大有提升的趋势, 并且性能约为读锁的100倍;
