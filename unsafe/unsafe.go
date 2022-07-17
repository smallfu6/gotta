package main

/*
	unsafe 非常简洁, 定义了一个类型和三个函数, 其中 ArbitraryType 并不真正属于
	unsafe 包, 其表示一个任意表达式的类型, 仅用于文档目的, go 编译器会对其做
	特殊处理(TODO)
	- ./unsafe_sizeof.go
	- ./unsafe_alignof.go
	- ./unsafe_offsetof.go

	unsafe 包中定义了 unsafe.Pointer 类型, unsafe.Pointer 类型可用于表示任意
	类型的指针, 并且具备以下其他指针类型不具备的性质:
	- 任意类型的指针都可以被转换为 unsafe.Pointer
	- unsafe.Pointer 也可以被转换为任意类型的指针值
	- uintptr 类型值可以被转换为一个 unsafe.Pointer
	- unsafe.Pointer 也可以被转换为一个 uintptr 类型值
*/

// 可以通过 unsafe.Pointer 很容易的穿透go的类型安全保护, 对比以下代码:
// ./c_not_type_safe.c
// ./go_type_safe.go
//
