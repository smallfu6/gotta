package main

/*
 * 当函数可以被内联时, 函数将被纳入调用函数.
 * 函数参数与返回值在编译器内联阶段都将转换为声明语句, 并通过 goto 语义跳转到
 * 调用者函数语句中, 在后续编译器阶段还将对内联结构做进一步优化(TODO)
 */

func small() string {
	s := "hello, " + "world!"
	return s
}

func fib(index int) int {
	if index < 2 {
		return index
	}
	return fib(index-1) + fib(index-2)
}

func main() {
	small()
	fib(65)
}

// go tool compile -m=2 inline.go
// 在编译时加入 -m=2(2?)标志, 可打印出函数的内联调试信息, 可以看出
// small 函数可以被内联, 而 fib(斐波那契)函数为递归函数, 不能被内联
// inline.go:3:6: can inline small with cost 7 as: func() string { s := "hello, world!"; return s }
// inline.go:8:6: cannot inline fib: recursive(递归)
// inline.go:15:6: can inline main with cost 69 as: func() { small(); fib(65) }
// inline.go:16:7: inlining call to small func() string { s := "hello, world!"; return s }
