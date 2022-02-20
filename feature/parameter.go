package main

import (
	"fmt"
	"unsafe"
)

/*
  参数传递:
  除函数的调用惯例外, 在传递参数时是传值还是传引用也很重要, 不同的选择会影响
  在函数中修改入参时是否会影响调用方看到的数据; 传值和传引用的区别:
  - 传值: 函数调用时会复制参数, 被调用方和调用方持有不相关的两份数据;
  - 传引用: 函数调用时会传递参数的指针, 被调用方和调用方持有相同的数据, 任意一方
	做出的修改都会影响另一方;

  不同语言会选择不同的方式传递参数, go 语言选择传值的方式, 即无论是传递基本类型,
  结构体还是指针, 都会对传递的参数进行复制

*/

// 整形和数组
func IntAndArray(i int, arr [2]int) {
	i = 29
	arr[1] = 88
	fmt.Printf("in IntAndArray - i=(%d, %p) arr=(%v, %p)\n", i, &i, arr, &arr)
}

func mainForIntAndArray() {
	i := 30
	arr := [2]int{66, 77}
	fmt.Printf("before calling - i=(%d, %p) arr=(%v, %p)\n", i, &i, arr, &arr)
	IntAndArray(i, arr)
	fmt.Printf("after calling - i=(%d, %p) arr=(%v, %p)\n", i, &i, arr, &arr)
}

// before calling - i=(30, 0xc0000140f0) arr=([66 77], 0xc000014100)
// in IntAndArray - i=(29, 0xc0000140f8) arr=([66 88], 0xc000014120)
// after calling - i=(30, 0xc0000140f0) arr=([66 77], 0xc000014100)
// 可以得出 go 语言的整形和数组类型都是值传递的, 在调用时会复制内容; 但要注意
// 如果传递的数组大小非常大, 这种传值的方式会对性能造成比较大的影响

// 结构体和指针
func StructAndPointer(a MyStruct, b *MyStruct) {
	a.i = 31
	b.i = 41
	fmt.Printf("in StructAndPointer - i=(%d, %p) arr=(%v, %p)\n", a, &a, b, &b)
}

type MyStruct struct {
	i int
}

func mainForStructAndPointer() {
	a := MyStruct{i: 30}
	b := &MyStruct{i: 40}
	fmt.Printf("before calling - i=(%d, %p) arr=(%v, %p)\n", a, &a, b, &b)
	StructAndPointer(a, b)
	fmt.Printf("after calling - i=(%d, %p) arr=(%v, %p)\n", a, &a, b, &b)
}

// before calling - i=({30}, 0xc0000140f0) arr=(&{40}, 0xc00000e028)
// in StructAndPointer - i=({31}, 0xc000014108) arr=(&{41}, 0xc00000e038)
// after calling - i=({30}, 0xc0000140f0) arr=(&{41}, 0xc00000e028)
// 可以得出:
// - 传递结构体时会复制结构体中的全部内容
// - 传递结构体指针时会复制结构体指针(函数局部变量 b 存储结构体指针)

// go 结构体在内存中的布局
type My2Struct struct {
	i int
	j int
}

func CacheLayout(ms *My2Struct) {
	ptr := unsafe.Pointer(ms)
	for i := 0; i < 2; i++ {
		c := (*int)(unsafe.Pointer((uintptr(ptr) + uintptr(8*i))))
		*c += i + 1
		fmt.Printf("[%p} %d\n", c, *c)
	}
}

// 结构体在内存中是一块连续的空间, 指向结构体的指针也是指向该结构体的首地址;
// 将 My2Struct 指针修改为 int 类型, 访问新指针就会返回整形变量i, 将指针移动
// 8 字节(int占8字节) 后就能获取下一个成员变量 j

func mainForCacheLayout() {
	a := &My2Struct{i: 40, j: 50}
	CacheLayout(a)
	fmt.Printf("[%p] %v\n", a, a)
}

/*
    TODO:
	在Golang支持的数据类型中是包含指针的, 但是Golang中的指针与C/C++的指针却又
	不同, 主要表现在下面的两个方面:
	- 弱化了指针的操作, 在Golang中指针的作用仅是操作其指向的对象，不能进行类
		似于C/C++的指针运算, 例如指针相减, 指针移动等
	- 指针类型不能进行转换, 如int不能转换为int32
	上述的两个限定主要是为了简化指针的使用, 减少指针使用过程中出错的机率，提高
	代码的鲁棒性; 但是在开发过程中, 有时需要打破这些限制, 对内存进行任意的读写,
	这就需要unsafe.Pointer了

	任意类型的指针值都可以转换为unsafe.Pointer, unsafe.Pointer也可以转换为
	任意类型的指针值;
	unsafe.Pointer与uintptr可以实现相互转换;
	可以通过uintptr进行加减操作, 从而实现指针的运算;

*/

func MyFunction(ms *My2Struct) *My2Struct {
	return ms
}

// go tool compile -N -S -l parameter.go // 关于 MyFunction 函数的汇编代码
// "".MyFunction STEXT nosplit size=20 args=0x10 locals=0x0
// 	0x0000 00000 (parameter.go:111)	TEXT	"".MyFunction(SB), NOSPLIT|ABIInternal, $0-16
// 	0x0000 00000 (parameter.go:111)	FUNCDATA	$0, gclocals·524d71b8d4b4126db12e7a6de3370d94(SB)
// 	0x0000 00000 (parameter.go:111)	FUNCDATA	$1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
// 	0x0000 00000 (parameter.go:111)	MOVQ	$0, "".~r1+16(SP)  // 初始化返回值
// 	0x0009 00009 (parameter.go:112)	MOVQ	"".ms+8(SP), AX    // 复制引用
// 	0x000e 00014 (parameter.go:112)	MOVQ	AX, "".~r1+16(SP)  // 返回引用
// 	0x0013 00019 (parameter.go:112)	RET
// 	0x0000 48 c7 44 24 10 00 00 00 00 48 8b 44 24 08 48 89  H.D$.....H.D$.H.
// 	0x0010 44 24 10 c3                                      D$..

// 在这段汇编中, 当参数是指针时, 会使用 MOVQ "".ms + 8(SP),	AX 指令复制引用,
// 然后将复制的指针作为返回值传递回调用方.
// 所以将指针作为参数传入某个函数时, 函数内部会复制指针, 即同时出现两个指针
// 指向原有内存空间, 因此go语言中传指针也是使用传值的方式.

/*
 go 语言在传递参数时使用了传值的方式, 接收方在收到参数时会复制它们; 在传递数组
 或者内存占用非常大的结构体时, 应该尽量使用指针作为参数类型来避免发生数据复制
 进而影响性能.
*/
