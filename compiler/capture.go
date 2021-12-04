package main

import "fmt"

/*
 * 变量捕获:
 * 类型检查阶段完成后, go 语言编译器会对抽象语法树(AST)进行分析和重构,
 * 从而完成一系列优化.
 * 变量捕获主要针对闭包场景, 由于闭包函数中可能引用闭包外的变量, 因此
 * 变量捕获需要明确在闭包中通过值引用或地址引用的方式来捕获变量.
 */

func main() {
	a := 1
	b := 2
	// 由于变量a在闭包之后进行了其他赋值操作, 因此在闭包中, a,b 变量的
	// 引用方式会有所不同; 在闭包中, 必须用取地址引用的方式对变量 a 进行
	// 操作, 而对变量 b 的引用将通过直接值传递的方式进行
	go func() {
		fmt.Println(a, b)
	}()
	a = 99
}

// 在编译过程中, 可通过下列命令查看当前程序闭包变量捕获的情况, 从输出中:
// - a 采取 ref 引用传递方式, assign=true 代表变量a在闭包完成后进行了赋值操作
// - b 采取了值传递的方式
// go tool compile -m=2 main.go | grep capturing
// capture.go:17:15: main.func1 capturing by ref: a (addr=true assign=true width=8)
// capture.go:17:18: main.func1 capturing by value: b (addr=false assign=false width=8)
