package main

import "testing"

var Result int

// 函数内联基准测试
func BenchmarkMax(b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = maxY(-1, i)
	}
	Result = r
}

// go test func_inline_test.go  func.go -bench=.
// goos: linux
// goarch: amd64
// BenchmarkMax-8          1000000000               0.362 ns/op
// PASS
// ok      command-line-arguments  0.406s

// line 19: GOMAXPROCS=8  1000000000次调用 每次调用用时 0.362ns
// 测试时间默认是1秒, 即1秒调用 10 0000 0000 次, 每次调用花费0.362 纳秒.
