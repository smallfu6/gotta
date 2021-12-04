package main

import "testing"

var Result int

// 禁用函数内联基准测试
func BenchmarkMax(b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = maxN(-1, i)
	}
	Result = r
}

// go test func_noinline_test.go  func.go -bench=.
// goos: linux
// goarch: amd64
// BenchmarkMax-8          770861887                1.45 ns/op
// PASS
// ok      command-line-arguments  1.276s

// line 18: GOMAXPROCS=8  770861887 次调用 每次调用用时 1.45ns
// 测试时间默认是1秒, 即1秒调用 7 7086 1887 次, 每次调用花费 1.45 纳秒.
