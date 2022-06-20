package main

import "testing"

const slcieSize = 10000

func BenchmarkSliceInitWithoutCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sl := make([]int, 0)
		for i := 0; i < slcieSize; i++ {
			sl = append(sl, i)
		}
	}
}

func BenchmarkSliceInitWithCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sl := make([]int, slcieSize)
		for i := 0; i < slcieSize; i++ {
			sl = append(sl, i)
		}
	}
}

/*
	go test  -bench=. -benchmem slice_benchmark_test.go

	goos: linux
	goarch: amd64
	cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
	BenchmarkSliceInitWithoutCap-8             29730             40582 ns/op          357626 B/op         19 allocs/op
	BenchmarkSliceInitWithCap-8                24931             47890 ns/op          507905 B/op          4 allocs/op
	PASS
	ok      command-line-arguments  3.303s
*/
