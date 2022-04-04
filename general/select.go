package main

/*
TODO: 操作系统中的 i/o  多路复用, select, epoll
	select 是操作系统中的系统调用, 我们经常使用 select, poll, epoll 等函数
	构建 i/o 多路复用模型提升程序性能(TODO: i/o 多路复用);
	c 语言的 select 系统调用可以同时监听多个文件描述符的可读可写状态, go
	语言中的 select 也能够让 goroutine 同时等待多个 channel 可读或可写,
	在多个文件或者 channel 状态改变之前, select 会一直阻塞当前线程或 goroutine;

	select 中虽然有多个 case, 但是这些 case 中的表达式必须都是 channel 的
	收发操作;

	select 控制结构中包含 default 语句, 可以实现非阻塞收发;
	select 在遇到多个 channel 同时响应时, 会随机执行一种情况;
	TODO: select 底层原理, go 汇编

	有多个 case 同时满足执行条件时, 如果按顺序依次判断, 那么后面的条件永远
	得不到执行, 引入随机性是为了避免饥饿问题发生;
*/
