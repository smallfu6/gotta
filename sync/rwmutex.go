package main

/* sync.RWMutex (TODO: 源码) */

/*
	读写互斥锁 sync.RWMutex 是细粒度的互斥锁, 不限制资源的并发读, 但是读写,
	写写操作无法并行执行.

	常见服务的资源读写比例会非常高, 因为大多数读请求之间不会相互影响, 所有可以
	分离读写操作, 以此提高服务的性能

	type RWMutex struct {
		w           Mutex
		writerSem   uint32
		readerSem   uint32
		readerCount int32
		readerWait  int32
	}
	- w 复用互斥锁提供的能力;
	- writerSem 和 readerSem 分别用于写等待读和读等待写;
	- readerCount 存储了当前正在执行的读操作数量;
	- readerWait 表示当写操作被阻塞时等待的读操作个数;



*/
