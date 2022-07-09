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


	TODO: 函数栈帧, goroutine 结构中对 panic 和 defer 的存储
	recover() 函数调用有严格的要求, 必须在 defer 函数中直接调用, 如果调用的
	是 recover 的包装函数或在嵌套的 defer 里调用都不能捕获异常; 即 recover
	必须要和有异常的栈帧只隔一个栈帧, recover 函数才能捕获异常; 换言之 recover
	捕获的是祖父一级调用函数栈帧的异常(刚好可以跨越一层 defer 函数)

	TODO: 实践
	关于发生panic后输出的栈跟踪信息(stack trace)的识别, 可遵循以下几个要点:
	- 栈跟踪信息中每个函数/方法后面的"参数数值"个数与函数/方法原型的参数
		个数不是一一对应的;
	- 栈跟踪信息中每个函数/方法后面的"参数数值"是按照函数/方法原型参数列表中
		从左到右的参数类型的内存布局逐一展开的, 每个数值占用一个字
		(word, 64位平台下为8字节);
	- 如果是方法, 则第一个参数是receiver自身, 如果receiver是指针类型, 则
		第一个参数数值就是一个指针地址; TODO
	- 如果是非指针的实例, 则栈跟踪信息会按照其内存布局输出;
	- 函数/方法返回值放在栈跟踪信息的"参数数值"列表的后面;
	- 如果有多个返回值, 则同样按从左到右的顺序, 按照返回值类型的内存布局输出;
	- 指针类型参数: 占用栈跟踪信息的"参数数值"列表的一个位置; 数值表示指针值,
		也是指针指向的对象的地址
	- string类型参数: 由于string在内存中由两个字表示(第一个字是数据指针, 第
		二个字是string的长度), 因此在栈跟踪信息的"参数数值"列表中将占用两个
		位置;
	- slice类型参数: 由于slice类型在内存中由三个字表示(第一个字是数据指针, 第
		二个字是len, 第三个字是cap), 因此在栈跟踪信息的"参数数值"列表中将占用
		三个位置;
	- 内建整型(int、rune、byte): 由于按字逐个输出, 对于类型长度不足一个字的
		参数, 会进行合并处理;
		比如, 一个函数有5个int16类型的参数, 那么在栈跟踪信息中这5个参数将
		占用"参数数值"列表中的两个位置:
		- 第一个位置是前4个参数的"合体"
		- 第二个位置则是最后那个int16类型的参数值;
	- struct类型参数: 会按照struct中字段的内存布局顺序在栈跟踪信息中展开;
	- interface类型参数: 由于interface类型在内存中由两部分组成(一部分是
		接口类型的参数指针, 另一部分是接口值的参数指针), 因此interface类型
		参数将使用"参数数值"列表中的两个位置;

	栈跟踪输出的信息是在函数调用过程中的"快照"信息, 因此一些输出数值
	虽然看似不合理, 但由于其并不是最终值, 问题不一定由其引起;

*/

// 跨协程失效
func MoreGoroutine() {
	defer fmt.Println("in main")
	go func() {
		defer fmt.Println("in goroutine")
		panic("") // panic 只会触发当前 goroutine 的 defer
	}()

	time.Sleep(1 * time.Second)
}

// in goroutine
// panic:

// main 函数中的 defer 语句未执行
/*
	defer 关键字对应的 runtime.deferproc 会将延迟调用函数与调用方所在 goroutine
	进行关联; 所以当程序发生崩溃, 只会调用当前 goroutine 的延迟调用函数;

	多个 goroutine 之间没有太多关联, 一个 goroutine 在触发 panic 时也不应该执行
	其他 goroutine 的延迟函数;
*/

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

func main1() {
	call()
	fmt.Println("333 Helloworld") // (6)
}

func call() {
	/*

	 TODO: panic源码 $GOROOT/src/rumtime/panic.go:425

	 如果遇见panic关键字, 执行流程就会进入代码gopanic函数中, 进入之后
	 会拿到表示当前协程g的指针, 然后通过该指针拿到当前协程的defer链表,
	 通过for循环来进行执行defer, 如果在defer中又遇见了panic的话, 则会
	 释放这个defer, 通过continue去执行下一个defer, 然后就是一个一个的
	 执行defer了, 如果在defer中遇见recover, 那么将会通过mcall(recovery)
	 去执行panic

	*/

	defer func() {
		fmt.Println("11111") // (5)
	}()
	defer func() {
		fmt.Println("22222") // (4)
	}()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recover from r : ", r) // (3)
		}
	}()
	defer func() {
		fmt.Println("33333") // (2)
	}()

	fmt.Println("111 Hello, world") // (1) 执行顺序
	panic("Panic 1!")
	panic("Panic 2!")
	fmt.Println("222 Hello, world")
}

// 111 Hello, world
// 33333
// Recover from r :  Panic 1!
// 22222
// 11111
// 333 Helloworld

/*
	程序多次调用 panic 也不会影响 defer 函数的正常执行, 所以使用 defer 进行收尾
	工作是安全的;
*/
// 嵌套崩溃(panic 可以多次嵌套使用)
func main() {
	defer fmt.Println("in main")
	defer func() {
		// 因为此 defer 中存在 panic 所以会被释放, 然后执行该defer
		// 的下一个 defer, 执行完 defer 后开始执行同级 panic,
		// 然后再继续执行该 defer
		defer func() {
			// 同父 defer
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
