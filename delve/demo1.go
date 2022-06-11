package main

import "fmt"

func main() {
	nums := make([]int, 5)
	for i := 0; i < len(nums); i++ {
		nums[i] = i * i
	}
	fmt.Println(nums)
}

// dlv debug demo1.go

//# goroutines
// * Goroutine 1 - User: ./demo1.go:8 main.main (0x4966e2) (thread 112328)
//   Goroutine 2 - User: /usr/local/go/src/runtime/proc.go:362 runtime.gopark (0x4377d2) [force gc (idle)]
//   Goroutine 3 - User: /usr/local/go/src/runtime/proc.go:362 runtime.gopark (0x4377d2) [GC sweep wait]
//   Goroutine 4 - User: /usr/local/go/src/runtime/proc.go:362 runtime.gopark (0x4377d2) [GC scavenge wait]
//   Goroutine 5 - User: /usr/local/go/src/runtime/proc.go:362 runtime.gopark (0x4377d2) [finalizer wait]
// [5 goroutines]
