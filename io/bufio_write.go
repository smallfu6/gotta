package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
	有一种接口的常见应用模式: 包裹函数(wrapper function): 接收类型参数并返回
	与其参数类型相同的返回值; 如:
		func WrapperFunc(param InterfaceType) InterfaceType

	通过包裹函数返回的包裹类型可以实现对输入数据的过滤, 装饰, 变换等操作, 并
	将结果再次返回给调用者; go 标准库的读写模型广泛运用了包裹函数模式, 并且
	基于这种模式实现了有缓冲 I/O, 数据格式变换等;


	如果对文件的读写都是无缓冲的, 即每次读都会驱动磁盘运转来读取数据, 每次写(
	并随后调用Sync)也都会对数据进行落盘处理, 这种频繁的磁盘 I/O 是无缓冲 I/O
	模式性能不高的主因; 任何软件工程遇到的问题都可以通过增加一个中间层来解决,
	于是出现了带缓冲 I/O;
	带缓冲 I/O 模式通过维护一个中间的缓存来降低数据读写时磁盘操作的频度, go
	标准库中的带缓冲 I/O 是通过包裹函数创建的包裹类型实现的;

	标准库通过包裹函数模式轻松实现了带缓冲的I/O, 充分展示了标准库读写模型的
	优势, 且 bufio 不仅可用于磁盘文件, 还可以用于包裹任何实现了 io.Writer
	和 io.Reader 接口的类型(如网络连接), 为其提供缓冲I/O的特性;

*/

func main() {
	file := "bufio.txt"

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}

	defer func() {
		f.Sync() // TODO: 内存数据落盘, 常用于哪种情况?
		f.Close()
	}()

	data := []byte("I love golang!\n")

	// TODO: bufio 包源码, 掌握其常用方法的使用
	// 通过 NewWriterSize 对 io.File 实例进行包裹, 得到包裹类型 bufio.Writer
	// 类型的实例 bio
	bio := bufio.NewWriterSize(f, 32) // 初始缓冲区大小为32字节

	// 将15字节写入bio缓冲区, 缓冲区缓存15字节, bufio.txt 中内容为空
	bio.Write(data)

	// 将15字节写入bio缓冲区, 缓冲区缓存15字节, bufio.txt 中内容为空
	bio.Write(data)

	// 将15字节写入bio缓冲区后, bufio 开始将32字节写入 bufio.txt 文件中,
	// 缓冲区空出 32 字节的位置
	bio.Write(data) // bio 缓冲区中仍然缓存(15*3-32) 字节

	/*
		以上三次写入, 实际上仅执行了一次真正的文件 I/O
	*/

	// 将缓冲区的剩余数据都写入磁盘中
	bio.Flush()

	// TODO: 日志包是不是使用了这种方式? 研究日志包
}
