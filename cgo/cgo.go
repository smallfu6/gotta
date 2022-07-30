package main

/*
	cgo 使用场景:
	- 对go内存GC的延迟敏感, 需要手动进行内存管理(分配和释放)
	- 为一些c语言专有的且没有go替代品的库制作go绑定(binding)或包装;
		如: Oracle 提供了C版本OCI库(Oracle Call Interface), 但并未提供go版本
		以及连接数据库的协议细节, 因此只能通过包装C语言的OCI版本与Oracle数据
		库通信; 类似的还有一些图形化驱动程序以及图形化的窗口系统接口;
	- 与遗留的且重构难度较大的C代码进行交互

*/
