package main

/*
	在 go 中, 上下文 context.Context 设置截止日期, 同步信号, 传递请求相关值的
	结构体, 上下文与 goroutine 的关系较密切, 是 go 中独特的设计, 在其他编程语言
	中很少见到类似概念.

	context.Context 是 go1.7 引入标准库的接口(TODO: 源码)
	type Context interface {
		Deadline() (deadline time.Time, ok bool)
		返回 context.Context 被取消的时间, 即完成工作的截止日期

		Done() <-chan struct{}
		返回一个 Channel, 在当前工作完成或上下文被取消后关闭, 多次调用 Done 方
		法会返回同一个 Channel

		Err() error
		返回 context.Context 结束的原因, 只会在 Done 方法对应的 Channel 关闭时
		返回非空值:
		- 如果 context.Context 被取消, 会返回 Canceled 错误;
		- 如果 context.Context 超时, 会返回 DeadlineExceeded 错误;

		Value(key interface{}) interface{}
		从 context.Context 中获取键对应的值, 对于同一个上下文来说, 多次调用
		Value 并传入相同的 Key 会返回相同的结果, 该方法用来传递特定的数据
	}


	context 包中提供的 context.Background, context.TODO, context.WithDeadline
	和 context.WithValue 函数会返回实现该接口的私有结构体(TODO)


	设计原理:
	context.Context 的最大作用是在 Goroutine 构成的树形结构中同步信号以减少计算
	资源的浪费.

	可能会创建多个 goroutine 来处理一次请求, 而 context.Context 的作用是在不同
	的 goroutine 之间同步请求特定数据, 取消信号以及处理请求的截止日期:

                             -------> Goroutine ----> Goroutine
                            /
	Context ---> Goroutine ---------> Goroutine
	                       \
						    \-------> Goroutine

	每个 context.Context 都会从最顶层的 Goroutine 逐层传递到最底层(TODO: 结合
	gin 框架理解), context.Context 可以在上层 Goroutine 执行出现错误时将信号
	及时同步给下层.

*/
