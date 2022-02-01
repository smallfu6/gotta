package main

/* 同步原语与锁 */

/*
 * go 原生支持用户态进程(goroutine); 锁是并发编程和多线程编程里的关键概念; 锁
   并发编程中的一种同步原语(synchronization primitive), 能保证多个 goroutine
   在访问同一块内存时不会出现竞争条件(race condition)等问题;

 * sync 包中提供了基本的同步原语

	同步原语:	Cond         Once		WaitGroup
	容器:		Map			 Pool
	互斥锁:		Mutex		 RWMutex

	尽管提供了较为基础的同步功能, 但它们是一种相对原始的同步机制, 多数情况下
	应该使用抽象层级更高的 Channel 实现同步.

	锁 ./mutex.go
*/
