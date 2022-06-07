package main

/*
	<go语言底层原理剖析> Page173
	TODO: 基准测试, 通过基准测试分析程序的内存操作, 栈堆
*/

import "testing"

func BenchmarkDirect(b *testing.B) {
	adder := Sumer{id: 6754}
	for i := 0; i < b.N; i++ {
		adder.Add(10, 32)
	}
}

func BenchmarkInterface(b *testing.B) {
	adder := Sumer{id: 6754}
	for i := 0; i < b.N; i++ {
		Sumifier(adder).Add(10, 32)
	}
}

// go test escape.go escape_test.go  -bench=. -benchmem

// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkDirect-8      	1000000000	         0.2302 ns/op	       0 B/op	       0 allocs/op
// BenchmarkInterface-8   	1000000000	         1.151 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	command-line-arguments	1.533s
