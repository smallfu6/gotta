package main

/*
	数据争用(data race)在 go 语言中指两个协程同时访问相同的内存空间, 并且至少
	有一个写操作的情况, 这种情况通常是并发错误的根源, 也是最难调试的并发错误
	之一; 原因在于其结果是不明确的, 而且出错是在特定的条件下; 这导致很难复现
	相同的错误, 在测试阶段也不一定能测试出问题;

	检查工具 race 可以排查数据争用问题, race 可以使用在多个 go 指令中, 可以检
	测程序中的数据争用, 将报告包含发生 race 冲突的协程栈, 以及此时正在运行的
	协程栈;
	如下的命令都可以使用 race:
	- go test -race xxx
	- go run -race xxx
	- go build -race xxx
	- go install -race xxx

	竞争检查的成本因程序而异, 对于典型的程序, 内存使用量可能增加 5-10 倍, 执行
	时间会增加 2-20 倍; 同时竞争检测器为当前每个 defer 和 recover 语句额外分配
	8 字节, 在 goroutine 退出前, 这些额外分配的字节不会被回收; 这意味着如果有
	一个长期运行的 goroutine 并定期有 defer 和 recover 调用, 则程序内存的使用
	量可能无限增长; 这些内存分配不会显示到 runtime.ReadMemStats 或
	runtime/pprof 的输出中; TODO: 验证


	race 工具原理: [go 语言底层原理剖析, Page282] TODO: 2刷做深入理解
	race 工具借助了 ThreadSanitizer 工具(TODO: 了解), 被 go 语言内部通过 CGO
	的形式调用;

	对于一个变量 count, 协程 A 和协程 B 对 count 的安全访问会有两种情况:
	- 协程 A 访问结束后, 协程 B 继续执行
	- 协程 B 访问结束后, 协程 A 继续执行
	但是A和B不能同时访问count变量, 此时A和B之间的关系叫做 happened-before,
	可以用符号 -> 表示, 如果 A 先发生, B后发生, 写作 A -> B

	矢量时钟(Vector Clock)技术用来观察事件之间 happened-before 的顺序,
	该技术在分布式系统中使用广泛, 用于检测和确定分布式系统中事件的因果
	关系, 也可以用于数据争用的探测; (TODO)
	在go语言中, 每个协程在创建之初都会初始化矢量时钟, 并在读取或写入事件
	时修改自己的逻辑时钟;

*/
