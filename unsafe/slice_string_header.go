package main

/*
	模式6: reflect.SliceHeader和reflect.StringHeader必须通过模式1构建
*/

import (
	"fmt"
	"reflect"
	"unsafe"
)

func newSlice() *[]byte {
	var b = []byte("hello, gopher")
	return &b
}

func main() {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(newSlice())) //  模式1
	var p = (*[]byte)(unsafe.Pointer(bh))
	fmt.Printf("%q\n", *p)
	fmt.Println(bh)

	var a = [...]byte{'I', ' ', 'l', 'o', 'v', 'e', ' ', 'G', 'o', '!', '!'}
	bh.Data = uintptr(unsafe.Pointer(&a))
	bh.Len = len(a)
	bh.Cap = len(a)
	fmt.Printf("%q\n", *p)
	/*
		通过模式1构建的reflect.SliceHeader实例bh对newSlice返回的切片对象
		具有对象引用作用(TODO), 可以保证newSlice返回的对象不会被垃圾回收掉,
		后续反向转换成*[]byte依旧有效;
	*/
}

// "hello, gopher"
// &{824633802992 13 13}
// "I love Go!!"
