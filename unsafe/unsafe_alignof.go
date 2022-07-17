package main

import (
	"fmt"
	"unsafe"
)

/*
	Alignof 用于获取一个表达式的内存地址对齐系数; TODO: 对齐系数及其应用
	对齐系数(alignment factor)是一个计算机体系架构层面的术语, 在不同的计算机
	体系结构下, 处理器对变量地址都有着对齐要求, 即变量的地址必须可被该变量
	的对齐数整除;

	var x unsafe.ArbitraryType // unsafe.ArbitraryType 表示任意类型
	b := uintptr(unsafe.Pointer(&x)) % unsafe.Alignof(x) == 0
	fmt.Println(b) // true
*/

type Foo struct {
	a int
	b string
	c [10]byte
	d float64
}

func main() {
	var i int = 5
	var a = [100]int{}
	var sl = a[:]
	var f Foo
	fmt.Println(unsafe.Alignof(i))          // 8
	fmt.Println(unsafe.Alignof(f.a))        // 8
	fmt.Println(unsafe.Alignof(a))          // 8
	fmt.Println(unsafe.Alignof(sl))         // 8
	fmt.Println(unsafe.Alignof(f))          // 8
	fmt.Println(unsafe.Alignof(f.c))        // 1
	fmt.Println(unsafe.Alignof(struct{}{})) // 1 空结构体的对齐系数为1(TODO)
	fmt.Println(unsafe.Alignof([0]int{}))   // 8 长度为0 的数组, 其对齐系数
	// 依然与其元素类型的对齐系数相同(TODO)

}
