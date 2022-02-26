package main

import (
	"fmt"
	"time"
)

/*
	panic 能改变程序的控制流, 调用 panic 后会立刻停止执行当前函数的剩余代码,
	并在当前 goroutine 中递归执行调用方的 defer;


	panic 只会触发当前 goroutine 的 defer, panic 允许在 defer 嵌套多次调用;
	recover 只有在 defer 中调用才会生效, recover 只有在发生 panic 之后调用
	才会生效;

*/

// 跨协程失效
func MoreGoroutine() {
	defer fmt.Println("in main")
	go func() {
		defer fmt.Println("in goroutine")
		panic("")
	}()

	time.Sleep(1 * time.Second)
}

// in goroutine
// panic:
// main 函数中的 defer 语句未执行
/*
	defer 关键字对应的 runtime.deferproc 会将延迟调用函数与调用方所在 goroutine
	进行关联; 所以当程序发生崩溃, 只会调用当前 goroutine 的延迟调用函数;

	多个 goroutine 之间没有太多关联, 一个 goroutine 在触发 panic 时不应该只需
	其他 goroutine 的延迟函数;
*/

// 嵌套崩溃(panic 可以多次嵌套使用)
func main() {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			panic("panic again and again")
		}()
		panic("panic again")
	}()
	panic("panic once")
}

// in main
// panic: panic once
//         panic: panic again
//         panic: panic again and again
// 程序多次调用 panic 也不会影响 defer 函数的正常执行, 所以使用 defer 进行收尾
// 工作是安全的;

/*
	$GOROOT/src/runtime/runtime2.go, line912

	// A _panic holds information about an active panic.
	//
	// This is marked go:notinheap because _panic values must only ever
	// live on the stack.
	//
	// The argp and link fields are stack pointers, but don't need special
	// handling during stack growth: because they are pointer-typed and
	// _panic values only live on the stack, regular stack pointer
	// adjustment takes care of them.
	type _panic struct {
		argp      unsafe.Pointer // pointer to arguments of deferred call run during panic; cannot move - known to liblink
		arg       interface{}    // argument to panic
		link      *_panic        // link to earlier panic
		pc        uintptr        // where to return to in runtime if this panic is bypassed
		sp        unsafe.Pointer // where to return to in runtime if this panic is bypassed
		recovered bool           // whether this panic is over
		aborted   bool           // the panic was aborted
		goexit    bool
	}

	由 link 字段可以推测: panic 函数可以被连续多次调用, 它们之间通过 link 可以
	组成链表;
	结构体中的 pc, sp, goexit 都是为了修复 runtime.Goexit 带来的问题而引入的,
	runtime.Goexit 能够只结束调用该函数的 Goroutine 而不影响其他 Goroutine,
	但是该函数会被 defer 中的 panic 和 recover 取消, 引入这3个字段为了保证该
	函数一定会生效;
*/
