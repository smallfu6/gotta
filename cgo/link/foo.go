package main

/*
	Go提供了#cgo指示符，可以用它指定Go源码在编译后与哪些共享库进行链接
*/

// #cgo CFLAGS: -I${SRCDIR}
// #cgo LDFLAGS: -L${SRCDIR} -lfoo
// #include <stdio.h>
// #include <stdlib.h>
// #include "foo.h"
import "C"
import "fmt"

func main() {
	fmt.Println(C.count)
	C.foo()
}

/*
	链接静态共享库
	通过#cgo指示符告诉Go编译器在当前源码目录(${SRCDIR}会在编译过程中自动转换
	为当前源码所在目录的绝对路径)下查找头文件foo.h, 并链接当前源码目录下的
	libfoo共享库;
	C.count变量和C.foo函数的定义都在libfoo共享库中

	生成静态共享库文件 ./libfoo.so
	gcc -c foo.c 生成 foo.o
	ar rv libfoo.a foo.o 生成 libfoo.a (TODO: ar工具)
	go build foo.go


	go 同样支持链接动态共享库, 使用下面的命令将 ./foo.c 编译为一个动态共享库
	gcc -c foo.c  生成 foo.o
	gcc -shared -Wl,-soname,libfoo.so -o libfoo.so foo.o  生成 libfoo.so
	然后将 libfoo.so 拷贝到 /lib/x86_64-linux-gnu/libfoo.so
	go build foo.go
	使用 ldd 查看编译foo.go后生成的二进制文件foo的动态共享库依赖情况
	ldd foo
	    linux-vdso.so.1 (0x00007ffddf1ca000)
        libfoo.so => /lib/x86_64-linux-gnu/libfoo.so (0x00007f55a1003000)
        libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f55a0fe0000)
        libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f55a0dee000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f55a1026000)

	go build foo.go

*/
