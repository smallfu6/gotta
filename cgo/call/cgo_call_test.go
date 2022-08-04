package main

import "testing"

func BenchmarkCGO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CallCFunc()
	}
}

func BenchmarkGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CallGoFunc()
	}
}

// 注意务必通过-gcflags '-l'关闭内联优化，这样才能得到公平的测试结果:
// go test -bench . -gcflags '-l' cgo_call_test.go cgo_call.go
// goos: linux
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkCGO-8   	24487000	        47.68 ns/op
// BenchmarkGo-8    	489380533	         2.522 ns/op
// PASS
// ok  	command-line-arguments	2.704s
