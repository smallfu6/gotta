package main

/*
	TCP 连接的建立

	建立 TCP Socket 连接需要经历客户端和服务端的三次握手过程, 在连接的建立过程
	中, 服务端是一个标准的 Listen+Accept 的结构, 而在客户端go语言使用 Dial
	或 DialTimeout 函数发起建立连接的请求;


	Dial 在调用后会一直阻塞, 直到连接建立成功或失败
	conn, err := net.Dial("tcp", "taobao.com:80")
	if err != nil {
		// 处理错误
	}
	// 连接建立成功, 可以进行读写操作

	DialTimeout 是带有超时机制的Dial

	对于客户端, 建立连接时可能会遇到以下几种情况:
	1. 网络不可达或对方服务未启动
	2. 对方服务的 listen backlog 队列满了(TODO)
		服务端忙, 瞬间有大量客户端尝试与服务端建立连接, 服务端可能会出现
		listen backlog队列满了, 接收连接(accept)不及时的情况, 这将导致客
		户端的Dial调用阻塞; 通常, 即便服务端不调用accept接收客户端连接,
		在backlog数量范围之内, 客户端的连接操作也都是会成功的, 因为新的
		连接已经加入服务端的内核listen队列中了, accept操作只是从这个队列
		中取出一个连接而已;
		./server2.go, ./client2.go TODO: 使用此代码验证, 结果与教程不符

		TODO: tcp 相关内核参数
		客户端初始可以成功建立的最大连接数与系统中net.ipv4.tcp_max_syn_backlog
		的设置有关, 可以通过 sysctl net.ipv4.tcp_max_syn_backlog 查看值;


	3. 若网络延迟较大, Dial 将阻塞并超时
		如果网络延迟较大, TCP握手过程将更加艰难坎坷(经历各种丢包), 时间消耗
		自然也会更长, Dial此时会阻塞; 如果经过长时间阻塞后依旧无法建立连接,
		Dial会返回类似"getsockopt: operation timed out"的错误;

		在连接建立阶段, 多数情况下Dial是可以满足需求的, 即便是阻塞一小会儿,
		但对于那些有严格的连接时间限定的Go应用, 如果一定时间内没能成功建立连接,
		程序可能需要执行一段异常处理逻辑, 为此就需要使用DialTimeout函数





*/
