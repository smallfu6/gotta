package main

import "fmt"

/*
 * 在完成逃逸分析后, 接下来优化的阶段是闭包重写,其核心逻辑位于:
 * $GOROOT/src/cmd/compile/internal/gc/closure.go
 *
 * 闭包重写分为闭包定义后被立即调用和闭包定义后不被立即调用;
 * 在闭包被立即调用的情况下, 闭包只能被调用一次, 这时可以将闭包转换为
 * 普通函数的调用形式.
 */

// func do() {
// 	a := 1
// 	func() {
// 		fmt.Println(a)
// 		a = 2
// 	}()
// }

// 上面的闭包最终会被转换为类似正常函数调用的形式, 由于变量 a 为引用传递,
// 因此构造的新的函数的参数为 int 指针类型; 如果变量是值引用的, 那么
// 构造的新的函数的参数为 int 类型.
func do() {
	a := 1
	func1(&a)
}

func func1(a *int) {
	fmt.Println(*a)
	*a = 2
}

/*
 * TODO
 * 如果闭包定义后不被立即调用, 而是后续调用, 那么同一个闭包可能被调用多次,
 * 这时需要创建闭包对象.
 */

// TODO
// 如果变量是按值引用的, 并且该变量占用的存储空间小于 2 X sizeof(int), 那么
// 通过在函数体内创建局部变量的形式来产生该变量, 如果该变量通过指针或值引用,
// 但是占用存储空间较大, 那么捕获的变量(var)转换成指针类型的 "&var", 这两种
// 方式都需要在函数序言阶段将变量初始化为捕获变量的值.
