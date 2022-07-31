package main

// #include <stdio.h>
// #include <stdlib.h>
//
// void print(char *str) {
//		printf("%s\n", str);
// }
//
// char *foo = "hellofoo";
import "C" // 与上面的C代码之间不能用空行分隔, 这里的"C"不是包名, 而是一种类似
// 名字空间的概念, 也可以理解为伪包名, C语言所有语法元素均在该伪包下面; 访问
// C语法元素时都要在其前面加上伪包C的前缀;
// 上面的C代码也可以使用 /* */ 进行注释, cgo 仍能解析

import (
	"fmt"
	"unsafe"
)

func main() {
	s := "Hello, Cgo"
	cs := C.CString(s)
	fmt.Printf("%T\n", cs) // *main._Ctype_char
	defer C.free(unsafe.Pointer(cs))
	C.print(cs)                                  // Hello, Cgo
	fmt.Printf("%T: %[1]s\n", C.GoString(C.foo)) // string
}

// 可以通过go build -x -v输出带有cgo代码的Go源文件的构建细节
// go build调用了名为cgo的工具, cgo会识别和读取Go源文件中的C代码, 并将其提取
// 后交给外部的C编译器(clang或gcc)编译, 最后与Go源码编译后的目标文件链接成一
// 个可执行程序
