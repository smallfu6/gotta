package main

/*
	TODO: GMP, goroutine 调度时间片
	TODO: slice 底层原理, copy; 引用类型在 goroutine 之间的使用
*/

import (
	"fmt"
	"runtime"
)

func gen() <-chan []int {
	c := make(chan []int)

	go func(c chan []int) { // goroutine1
		defer close(c)

		s := []int{0, 1, 2, 3}
		for i := 0; i < len(s); i++ {
			s[i] = -1

			// newSlice := make([]int, len(s))
			// copy(newSlice, s)
			// c <- newSlice
			c <- s
			runtime.Gosched()
		}
	}(c)

	return c

}

func main() { // goroutine2
	for s := range gen() {
		fmt.Println("chan", s)
	}
}
