package main

/*
 * 函数内联:
 * 指将较小的函数直接组合进去调用者的函数, 这是现代编译器的一种核心技术.
 * 函数内联的优势在于可以减少函数调用带来的开销.
 * 对于 go 语言, 函数调用的成本在与参数与返回值栈复制, 较小的栈寄存器开销
 * 以及函数序言部分(?)的检查栈扩容(go语言中的栈可动态扩容), 同时, 函数内联
 * 是其他编译器优化(如无效代码删除)的基础
 *
 * go 语言编译器会计算函数内联花费的成本, 只有执行相对简单的函数时才会内联,
 * 函数内联的核心逻辑位于 $GOROOT/src/cmd/compile/internal/gc/inl.go 中.
 *
 *
 * 当函数中有 for, range, go, select 等语句时, 该函数不会被内联, 当函数执行
 * 过于复杂(如太多的语句或函数为递归函数)时, 也不会执行内联
 *
 */

// 使用 bench  对 max 函数调用进行测试, 在函数的注释前方加上 //go:noinline
// 时, 代表当前函数禁止进行函数内联优化

// 如果希望程序中的所有函数都不执行内联, 可添加编译器选项 "-l"
// go build -gcflags="-l" xxx.go
// go tool compile -l xxx.go

//go:noinline
func maxN(a, b int) int {
	// 禁用内联
	if a > b {
		return a
	}
	return b
}

func maxY(a, b int) int {
	// 内联
	if a > b {
		return a
	}
	return b
}

// 通过 ./func_inline_test.go 和 ./func_noinline_test.go 的基准测试对比, 可
// 看出内联后, max 函数的执行时间显著少于非内联函数调用花费的时间, 这里的
// 消耗主要来自函数调用增加的执行指令
