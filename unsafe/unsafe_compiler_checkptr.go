package main

import "unsafe"

/*
	做完指针运算后, 转换后的unsafe.Pointer仍应指向原先的内存对象
*/

func main() {
	var n = 5
	b := make([]byte, n)
	end := unsafe.Pointer(uintptr(unsafe.Pointer(&b[0])) + uintptr(n+10))
	_ = end
}

// Go编译器检查到了越界的指针运算
// go run -race unsafe_compiler_checkptr.go
// fatal error: checkptr: pointer arithmetic result points to invalid allocation
