package main

import (
	"bytes"
	"fmt"
	"io"
)

/*
	interface 类型不是任意类型, 同时对于一个类型的方法来说, nil 是一个合法
	的接收者, 与一些函数允许nil指针作为实参, 方法的接收者也一样, 尤其是当
	nil 是类型中有意义的零值(如map和slice类型时);

	// 在以下简单的整形链表中, nil 代表空链表
	// IntList 是整形链表
	// *IntList 的类型 nil 代表空列表
	type IntList struct {
		Value int
		Tail *IntList
	}

	// Sum 返回列表元素的总和
	func (list *IntList) Sum() int {
		if list = nil {
			return 0
		}
		return list.Value + list.Tail.Sum() // TODO: 巧妙的链式调用
	}
	当定义一个类型允许nil作为接收者时, 应当在文档中显式的标明;

*/

type TestStruct struct{}

func NilOrNot(v interface{}) bool {
	// TODO: 深入理解类型转换
	// 调用 NilOrNot 函数时会发生隐式类型转换, 除向方法传入参数外, 变量的赋值
	// 也会触发隐式类型转换; 在进行类型转换时, *TestStruct 类型会转换为
	// interface{} 类型, 其包括变量的类型(动态类型)信息 TestStruct,
	// 只是其动态值为 nil, 所以转换后的变量与 nil 不相等.
	return v == nil
}

func main1() {
	var s *TestStruct
	fmt.Println(s == nil)    // true
	fmt.Println(NilOrNot(s)) // false
	// 传入后带了类型信息, 即接口存储 *TestStruct 动态类型, 其动态值是 nil
}

//  含有空指针的非空接口
const debug = false

func main() {
	var buf *bytes.Buffer
	// var buf io.Writer
	if debug {
		buf = new(bytes.Buffer)
	}
	f(buf) // 当 debug=false, 调用f时, 把一个类型为 *bytes.Buffer 的空指针
	// 赋给了 out 参数, 所以 out 的动态值为空, 但其动态类型为 *bytes.Buffer;
	// 这表示 out 是一个包含空指针的非空接口, 所以在f中的防御性检查 out!=nil
	// 仍然为 true;
	if debug {
		// ...
	}
}

func f(out io.Writer) {
	// ...
	if out != nil {
		out.Write([]byte("done!\n"))
		// panic: runtime error: invalid memory address or nil pointer dereference
		// [signal SIGSEGV: segmentation violation code=0x1 addr=0x20 pc=0x45f73d]

		// TODO: bytes.Buffer 包源码
		// 对于某些类型, 比如 *os.File(TODO), 其方法接收者为nil(空接收值)是
		// 合法的, 但是 *bytes.Buffer 不行, 在尝试访问缓冲区时崩溃了
	}
}
