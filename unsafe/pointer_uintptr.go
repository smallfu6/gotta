package main

import (
	"fmt"
	"unsafe"
)

/*
	使用uintptr类型变量保存栈上变量的地址同样是有风险的, 因为Go使用的是连续栈
	的栈管理方案, 每个goroutine的默认栈大小为2KB(_StackMin = 2048); 当
	goroutine当前剩余栈空间无法满足函数/方法调用对栈空间的需求时, Go运行时
	就会新分配一块更大的内存空间作为该goroutine的新栈空间, 并将该goroutine
	的原有栈整体复制过来, 这样原栈上分配的变量的地址就会发生变化;
	TODO: 栈扩容原理
*/

func main() {
	var x = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	fmt.Printf("变量x的值=%d\n", x)
	println("变量x的地址=", &x)

	var p = uintptr(unsafe.Pointer(&x))
	var q = unsafe.Pointer(&x)

	a(x) // 执行一系列函数调用
	// 调用函数a之后, goroutine栈发生了扩容; 变更了数组x中的元素值以用
	// 于栈扩容前后的对比

	// 变更数组x中元素的值
	for i := 0; i < 10; i++ {
		x[i] += 10
	}

	println("栈扩容后, 变量x的地址=", &x)
	fmt.Printf("栈扩容后, 变量x的值=%d\n", x)

	fmt.Printf("变量p(uintptr)存储的地址上的值=%d\n", *(*[10]int)(unsafe.Pointer(p)))
	fmt.Printf("变量q(unsafe.Pointer)引用的的地址上的值=%d\n", *(*[10]int)(q))
}

func a(x [10]int) {
	var y [100]int
	b(y)
}

func b(x [100]int) {
	var y [1000]int
	c(y)
}

func c(x [1000]int) {
}

// go run -gcflags="-l" pointer_uintptr.go  (TODO:为何禁止内联优化后才能引起栈的扩容)
// 变量x的地址= 0xc000078ec0
// 变量x的值=[1 2 3 4 5 6 7 8 9 0]
// 栈扩容后, 变量x的地址= 0xc000093ec0
// 栈扩容后, 变量x的值=[11 12 13 14 15 16 17 18 19 10]
// 变量p(uintptr)存储的地址上的值=[1 2 3 4 5 6 7 8 9 0]
// 变量q(unsafe.Pointer)引用的的地址上的值=[11 12 13 14 15 16 17 18 19 10]

// 栈扩容后, 变量x的地址发生了变化, unsafe.Pointer类型变量q的值被Go运行时
// 做了同步变更; 但uintptr类型变量p只是一个整型值, 它的值是不变的,
// 因此输出uintptr类型变量p存储的地址上的值时, 得到的仍是变量x变更前的值;
