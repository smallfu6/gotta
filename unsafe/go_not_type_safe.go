package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var a uint32 = 0x12345678
	fmt.Printf("0x%x\n", a) // 0x12345678

	p := unsafe.Pointer(&a) // 任意指针类型转换为 unsafe.Pointer
	b := (*[4]byte)(p)      // unsafe.Pointer 可以转换为任意的指针类型
	b[0] = 0x23
	b[1] = 0x45
	b[2] = 0x67
	b[3] = 0x8a

	fmt.Printf("0x%x\n", a) // 0x8a674523 小端

	// 原本被解释为 uint32 类型的一段内存(起始地址为&a, 长度为4字节), 通过
	// unsafe.Pointer 被重新解释为 [4]byte 并且通过变量 b(*[4]byte类型)对该
	// 段内存进行修改;
}
