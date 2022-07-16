package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	file := "bufio.txt"

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}

	// 通过包裹函数 NewReaderSize 对 io.File 实例进行包裹, 得到包裹类型
	// bufio.Reader 类型的实例 bio
	bio := bufio.NewReaderSize(f, 64)
	fmt.Printf("初始状态下缓冲区缓存数据数量=%d字节\n\n", bio.Buffered())

	var i int = 1
	for {
		data := make([]byte, 15)
		n, err := bio.Read(data)
		if err == io.EOF {
			fmt.Printf("第%d次读取数据, 读到文件末尾, 程序退出\n", i)
			return
		}

		if err != nil {
			fmt.Println("读取数据出错: ", err)
			return
		}

		fmt.Printf("第%d次读出数据: %q, 长度=%d字节\n\n", i, data, n)
		fmt.Printf("当前缓冲区缓存数据数量=%d字节\n\n", bio.Buffered())
		i++

		/*
			从执行结果看, 第一次读出15字节数据后, 当前缓冲区数据数是30字节,
			即第一次 bufio.Reader.Read 操作实际上从文件中读取了45字节, 其中
			15字节数据通过字节切片传递出来, 剩余的30字节则缓存在 bio 维护的
			内部缓冲区中, 第二, 三次读操作均为从该缓冲区中读取数据, 不会触发
			文件I/O操作;
		*/
	}

}

// 初始状态下缓冲区缓存数据数量=0字节
// 第1次读出数据: "I love golang!\n", 长度=15字节
// 当前缓冲区缓存数据数量=30字节
// 第2次读出数据: "I love golang!\n", 长度=15字节
// 当前缓冲区缓存数据数量=15字节
// 第3次读出数据: "I love golang!\n", 长度=15字节
// 当前缓冲区缓存数据数量=0字节
// 第4次读取数据, 读到文件末尾, 程序退出
