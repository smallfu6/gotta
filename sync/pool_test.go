package main

import (
	"bytes"
	"sync"
	"testing"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func writeBufFromPool(data string) {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteString(data)
	bufPool.Put(b)
}

func writeBufFromNew(data string) *bytes.Buffer {
	b := new(bytes.Buffer)
	b.WriteString(data)
	return b
}

func BenchmarkWithoutPool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		writeBufFromNew("hello")
	}
}

func BenchmarkWithPool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		writeBufFromPool("hello")
	}
}

// go test -bench . pool_test.go
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkWithoutPool-8   	37755252	        32.32 ns/op	      64 B/op	       1 allocs/op
// BenchmarkWithPool-8      	65016584	        16.15 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	command-line-arguments	2.326s
/*
	可以看到通过sync.Pool来复用数据对象的方式可以有效降低内存分配频率,
	减轻垃圾回收压力, 从而提高处理性能
*/
