package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func main() {
	file := "hello_gopher.gz"
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("open file error: ", err)
		return
	}

	// 通过包裹函数 NewReader 对 io.File 实例进行包裹, 得到包裹类型
	// gzip.Reader 类型的实例 zw
	zw, _ := gzip.NewReader(f)
	defer zw.Close()

	i := 1
	for {
		buf := make([]byte, 32)
		_, err = zw.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("第%d次读取的压缩数据为: %q\n", i, buf)
				fmt.Println("读取到文件末尾, 程序退出!")
			} else {
				fmt.Printf("第%d次读取压缩数据失败: %v", i, err)
				return
			}
			return
		}

		fmt.Printf("第%d次读取的压缩数据为: %q\n", i, buf)
	}
}
