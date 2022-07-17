package main

import (
	"fmt"
	"unsafe"
)

/*
	Offsetof 用于获取结构体中某字段的地址偏移量(相对于结构体变量的地址),
	Offsetof 函数应用面较窄, 仅用于求结构体中某字段的偏移值;
*/

type Foo struct {
	a int
	b string
	c [10]byte
	d float64
}

func main() {
	var f Foo
	fmt.Println(unsafe.Offsetof(f.b)) // 8
	fmt.Println(unsafe.Offsetof(f.d)) // 40
}
