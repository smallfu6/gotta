package main

import (
	"fmt"
	"unsafe"
)

/* Sizeof 用于获取一个表达式值所占内存的大小 */

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
	fmt.Println(unsafe.Sizeof(i))           // 8
	fmt.Println(unsafe.Sizeof(a))           // 800
	fmt.Println(unsafe.Sizeof(sl))          // 24 reflect.SliceHeader
	fmt.Println(unsafe.Sizeof(f))           // 48 8(int) + 16(reflect.StringHeader) + 10([10]byte) + 8(float64) = 42
	fmt.Println(unsafe.Sizeof(f.c))         // 10
	fmt.Println(unsafe.Sizeof((*int)(nil))) // 8
	// Sizeof 函数不支持直接传入无类型信息的 nil 值, 必须显式告知 Sizeof 传入
	// 的nil 是什么类型
}
