package main

/*
	在Go中使用C语言的类型

	指针类型
	原生数值类型的指针类型可按Go语法在类型前面加上星号*, 比如var p *C.int;
	但void*比较特殊, 在Go中用unsafe.Pointer表示, 因为任何类型的指针值都可
	以转换为unsafe.Pointer类型, 而unsafe.Pointer类型也可以转换回任意类型
	的指针类型; TODO

	字符串类型
	C语言中并不存在原生的字符串类型, 在C中用带结尾'\0'的字符数组来表示字符串;
	而在Go中string类型是语言的原生类型, 因此两种语言的互操作势必要进行字符
	串类型的转换; 通过C.CString函数, 将Go的string类型转换为C的"字符串"类型
	后再传给C函数使用;
		s := "Hello, Cgo\n"
		cs := C.CString(s)
		C.print(cs)
	此转换相当于在C语言世界的堆上分配一块新内存空间, 转型后所得到的C字符串
	cs并不能由Go的GC管理, 必须在使用后手动释放cs所占用的内存, 如 ./cgowork.go
	中的 defer c.free(unsafe.Pointer(cs));
	通过 C.GoString 可将C的字符串(*C.char)转换为Go的string类型, ./cgowork.go,
	相当于在go世界重新分配一块内存对象, 并复制了C的字符串(foo)的信息, 后续
	这个位于go世界的新的string类型对象将和其他go对象一样接受GC的管理;

*/
