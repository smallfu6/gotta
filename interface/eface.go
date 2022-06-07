package main

/*
	相比于有方法的接口(iface), 空接口不需要 interfacetype 表示接口的内在类型,
	也不需要 fun 方法列表; 对于空接口, go 语言在运行时使用了特殊的 eface 类型,
	其在64位系统中占据16字节;

	type eface struct {
		_type *_type
		data unsafe.Pointer
	}

	当类型转换为 eface 时, 空接口与一般接口的处理方式是相似的, 同样面临内存逃逸,
	寻址等问题
*/
