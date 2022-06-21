package main

/*
	go语言精通之路[go语言的字符串类型]
	string 的底层是只读的字节数组, 对 string 的底层的数据存储区仅能进行只读操作,
	一旦试图修改那块区域的数据, 会得到 SIGBUS 的运行时错误;
*/

import (
	"fmt"
	"unsafe"
)

func main() {
	var s string = "hello"
	fmt.Println("original strinbg: ", s)

	// 试图通过 unsafe 指针改变原始 string
	modifyString(&s)
	fmt.Println(s)
}

// TODO: 熟悉 unsafe.Pointer 的用法
func modifyString(s *string) {
	// 取出第一个8字节的值
	p := (*uintptr)(unsafe.Pointer(s))

	// 获取底层数组的地址
	var array *[5]byte = (*[5]byte)(unsafe.Pointer(p))

	var len *int = (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + unsafe.Sizeof((*uintptr)(nil))))
	for i := 0; i < (*len); i++ {
		fmt.Printf("%p => %c\n", &((*array)[i]), (*array)[i])
		p1 := &((*array)[i])
		v := (*p1)
		(*p1) = v + 1
	}
}

/*
	original strinbg:  hello
	0xc000010250 => °
	0xc000010251 => ]
	0xc000010252 => I
	0xc000010253 =>
	0xc000010254 =>
	unexpected fault address 0x1014a5eb1
	fatal error: fault
	[signal SIGSEGV: segmentation violation code=0x1 addr=0x1014a5eb1 pc=0x45b83c]

*/
