package main

/*
 * SSA 生成
 * 在 ./walk.go 遍历函数后, 编译器会将抽象语法树转换为 SSA(Static Single
 * Assignment, 静态单赋值); SSA 被大多数现代编译器(包括GCC和LLVM)使用,
 * 在 go1.7 中被正式引入并替换了之前的编译器后端, 用于最终生成更有效的机器码;
 * 在 SSA 生成阶段, 每个变量在声明之前都需要被定义, 并且每个变量只能被赋值
 * 一次.
 */

// y := 1
// y := 2
// x = y
// 在上面的代码中, 变量 y 被赋值了两次, 不符合 SSA 的规则, y := 1 这条语句是
// 无效的, 可以转化为如下形式:
// y1 := 1
// y2 := 2
// x1 = y2
// 通过 SSA 很容易识别出 y1 是无效代码并将其清除.

// SSA 生成阶段是编译器进行后续优化的保证, 如常量传播(Constant Propagation),
// 无效代码清除, 消除冗余, 强度降低(Strength Reduction)等. TODO

// 大部分与 SSA 相关的代码位于 $GOROOT/src/cmd/compile/internal/ssa 文件夹中,
// 但是将抽象语法树转换为 SSA 的逻辑位于:
// $GOROOT/src/cmd/compile/internal/gc/ssa.go 文件中; 在 ssa/README.md 文件
// 中, 有对 SSA 生成阶段比较详细的描述.

// 可以在编译时指定 GOSSAFUNC=main 查看 SSA 初始及其后续优化阶段生成的代码片段.
var d uint8

func main() {
	var a uint8 = 1
	a = 2
	if true {
		a = 3
	}
	d = a
}

// 生成 ssa.html 文件
// GOSSAFUNC=main GOOS=linux GOARCH=amd64 go tool compile ssa.go

/*
 * TODO
 * 初始阶段结束后, 编译器将根据生成的 SSA 进行一系列重写和优化, SSA 最终的阶段
 * 叫做 genssa; 在上例的 genssa 阶段中, 编译器清除了无效的代码及不会进入的 if
 * 分支, 并且将常量 Op 操作变为了 amd64 下特定的 MOVBstoreconst 操作.
 */
