package main

import (
	"compress/gzip"
	"fmt"
	"os"
)

/*
	通过包裹函数返回的包裹类型还可以实现对读出或写入数据的变换, 比如压缩等;
	go 标准库中的 compress/gzip 就提供了这种包裹函数与包裹类型;
	TODO: gzip 包的常规使用

*/

func main() {
	file := "hello_gopher.gz"

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}
	defer f.Close()

	// 通过包裹函数 NewWriter 对 io.File 实例进行包裹, 得到包裹类型
	// gzip.Writer 类型的实例zw
	zw := gzip.NewWriter(f)
	defer zw.Close() // zw.Close 方法调用会将压缩变换后的数据刷新到文件实例中

	_, err = zw.Write([]byte("hello, gopher! I love golang!!"))
	if err != nil {
		fmt.Println("write compressed data error:", err)
		return
	}

	fmt.Println("write compressed data ok")
}
