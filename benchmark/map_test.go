package main

import "testing"

const mapSize = 10000

func BenchmarkMapInitWithoutCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := make(map[int]int)
		for i := 0; i < mapSize; i++ {
			m[i] = i
		}
	}
}

func BenchmarkMapInitWithCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := make(map[int]int, mapSize)
		for i := 0; i < mapSize; i++ {
			m[i] = i
		}
	}
}

/*
	go test -bench=. -benchmem ./map_test.go

	goos: linux
	goarch: amd64
	cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
	BenchmarkMapInitWithoutCap-8   	    2128	    775235 ns/op	  687062 B/op	     275 allocs/op
	BenchmarkMapInitWithCap-8      	    2224	    501831 ns/op	  322223 B/op	      11 allocs/op
	PASS
	ok  	command-line-arguments	2.881s
*/
