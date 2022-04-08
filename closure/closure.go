package main

/* 闭包 */

import (
	"fmt"
	"time"
)

// https://zhuanlan.zhihu.com/p/92634505

func foo1(x *int) func() {
	return func() {
		*x = *x + 1
		fmt.Printf("foo1 val = %d\n", *x)
	}
}

func foo2(x int) func() {
	return func() {
		x = x + 1
		fmt.Printf("foo2 val = %d\n", x)
	}
}

// 闭包的延迟绑定
func foo0() func() {
	x := 1
	f := func() {
		fmt.Printf("foo0 val = %d\n", x)
	}
	x = 11
	return f
}

func foo3() {
	values := []int{1, 2, 3, 5}
	for _, val := range values {
		fmt.Printf("foo3 val = %d\n", val)
	}
}

func show(v interface{}) {
	fmt.Printf("foo4 val = %v\n", v)
}
func foo4() {
	values := []int{1, 2, 3, 5}
	for _, val := range values {
		go show(val)
	}
}

// goroutine 的延迟绑定
func foo5() {
	values := []int{1, 2, 3, 5}
	for _, val := range values {
		go func() {
			fmt.Printf("foo5 val = %v\n", val)
		}()
	}

}

var foo6Chan = make(chan int, 10)

func foo6() {
	for val := range foo6Chan {
		go func() {
			fmt.Printf("foo6 val = %d\n", val)
		}()
	}
}

func foo7(x int) []func() {
	var fs []func()
	values := []int{1, 2, 3, 5}
	for _, val := range values {
		fs = append(fs, func() {
			fmt.Printf("foo7 val = %d\n", x+val)
		})
	}
	return fs
}

func foo8() {
	for i := 1; i < 10; i++ {
		curTime := time.Now().UnixNano()
		go func(t1 int64) {
			t2 := time.Now().UnixNano()
			fmt.Printf("foo8 ts = %d us \n", t2-t1)
		}(curTime)
	}
}

func main() {
	foo5()
}
