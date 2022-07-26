package main

/*
	TODO: tcp 协议
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



	Socket 读写

	连接建立起来后, 需要在连接上进行读写以完成业务逻辑, go 运行时隐藏了 I/O 多路
	复用的复杂性, go 运行时隐藏了I/O多路复用的复杂性; 开发者只需采用 goroutine
	+阻塞I/O模型即可满足大部分场景需求;
	Dial 连接成功后会返回一个 net.Conn 接口类型的变量值, 这个接口变量的底层
	类型为 *TCPConn:

	// TCPConn is an implementation of the Conn interface for TCP network
	// connections.
	type TCPConn struct {
		conn
	}

	因此 net.Conn 继承了 conn 类型的 Read 和 Write 方法, 后续通过Dial函数
	返回值调用的 Write 和 Read 方法均是 net.conn 的方法

	type conn struct {
		fd *netFD
	}

	func (c *conn) ok() bool { return c != nil && c.fd != nil }

	// Implementation of the Conn interface.

	// Read implements the Conn Read method.
	func (c *conn) Read(b []byte) (int, error) {
		if !c.ok() {
			return 0, syscall.EINVAL
		}
		n, err := c.fd.Read(b)
		if err != nil && err != io.EOF {
			err = &OpError{Op: "read", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
		}
		return n, err
	}

	// Write implements the Conn Write method.
	func (c *conn) Write(b []byte) (int, error) {
		if !c.ok() {
			return 0, syscall.EINVAL
		}
		n, err := c.fd.Write(b)
		if err != nil {
			err = &OpError{Op: "write", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
		}
		return n, err
	}

	conn.Read 有以下几种场景:
	1. Socket 中无数据
		连接建立后, 如果客户端未发送数据, 服务端会阻塞在Socket的读操作上(TODO:
		应该为服务端未发送数据, 客户端会阻塞在Socket的读操作上?), 与阻塞I/O
		模型的行为模式是一致的,  执行该读操作的goroutine也会被挂起; Go运行时
		会监视该Socket, 直到其有数据读事件才会重新调度该Socket对应的
		goroutine完成读操作;

	2. Socket 中有部分数据
		如果Socket中有部分数据就绪, 且数据数量小于一次读操作所期望读出的
		数据长度, 那么读操作将会成功读出这部分数据并返回, 而不是等待期望
		长度数据全部读取后再返回;

	3. Socket 中有足够多的数据
		如果连接上有数据, 且数据长度大于或等于一次Read操作所期望读出的数据
		长度, 那么Read将会成功读出这部分数据并返回; 服务端(TODO: 客户端?)再
		次读取时会把剩余的数据继续读出;

	4. Socket 关闭
		如果客户端主动关闭了 Socket, 服务端的读要分为有数据关闭和无数据
		关闭两种情况;

		- 有数据关闭: 有数据关闭是指在客户端关闭连接(Socket)时, Socket中
			还有服务端尚未读取的数据;
			./client3.go, ./client3.go
	5. 读操作超时
		有些场合对读操作的阻塞时间有严格限制, 在这种情况下, 读操作的行为
		到底是什么样的呢？在返回超时错误时, 是否也同时读出了一部分数据呢?
		使用 ./server4.go, ./client4.go 模拟超时情况可知不会出现读出部分
		数据且返回超时错误的情况

	6. 成功写
		"成功写"指的就是Write调用返回的n与预期要写入的数据长度相等, 且error = nil

	7. 写阻塞
		TCP通信连接两端的操作系统内核都会为该连接保留数据缓冲区, 一端调用
		Write后, 实际上数据是写入操作系统协议栈的数据缓冲区中的(TODO);
		TCP是全双工通信(TODO: tcp 协议), 因此每个方向都有独立的数据缓冲区,
		当发送方将对方的接收缓冲区及自身的发送缓冲区都写满后, Write调用
		就会阻塞;
		./server5.go, ./client5.go

	8. 写入部分数据
		./client5.go 末尾

	9. 写入超时
		由 ./server6.go, ./client6.go 验证写入超时时, 仍然存在数据部分写入的
		情况; 因此虽然Go提供了阻塞I/O的便利, 但在调用Read和Write时依旧要结合
		这两个方法返回的n和err的结果来做出正确处理;

	10. goroutine安全的并发读写
		goroutine的网络编程模型决定了存在不同goroutine间共享conn的情况(TODO),
		那么conn的读写是不是goroutine并发安全的呢?

		先从应用的角度来看看并发Read操作和Write操作的goroutine安全的必要性:
		- Read
			由于TCP是面向字节流的(TODO: 粘包), conn.Read无法正确区分数据的
			业务边界, 因此多个goroutine对同一个conn进行Read操作的意义不大,
			goroutine读到不完整的业务包反倒增加了业务处理的难度;
		- Write
			对于Write操作而言, 倒是有多个goroutine并发写的情况;

		// $GOROOT/src/net/net.go
		type conn struct {
			fd *netFD
		}

		net.conn 只是 *netFD 的外层包裹结构, 最终 Write 和 Read 都是在 fd
		字段上的操作;

		// $GOROOT/src/net/fd_unix.go
		// Network file descriptor.
		type netFD struct {
			pfd poll.FD

			// immutable until Close
			family      int
			sotype      int
			isConnected bool // handshake completed or use of association with peer
			net         string
			laddr       Addr
			raddr       Addr
		}

		// FD is a file descriptor. The net and os packages use this type as a
		// field of a larger type representing a network connection or OS file.
		type FD struct {
			// Lock sysfd and serialize access to Read and Write methods.
			fdmu fdMutex

			// System file descriptor. Immutable until Close.
			Sysfd int

			// I/O poller.
			pd pollDesc

			// Writev cache.
			iovecs *[]syscall.Iovec

			// Semaphore signaled when file is closed.
			csema uint32

			// Non-zero if this file has been set to blocking mode.
			isBlocking uint32

			// Whether this is a streaming descriptor, as opposed to a
			// packet-based descriptor like a UDP socket. Immutable.
			IsStream bool

			// Whether a zero byte read indicates EOF. This is false for a
			// message based socket connection.
			ZeroReadIsEOF bool

			// Whether this is a file rather than a network socket.
			isFile bool
		}

		poll.FD 类型中包含了一个运行时实现的 fdMutex 类型字段, 用来串行化对
		该netFD对应sysfd的Write和Read操作; 也就是说, 所有对conn的Read和Write
		操作都是由fdMutex来同步的; netFD 的Read和Write方法的实现也证实了这
		一点;(TODO: 源码)

		每次Write操作都是受锁保护的, 直到此次数据全部写完; 因此在应用层面,
		要想保证多个goroutine在一个conn上的Write操作是安全的, 需要让每一次
		Write操作完整地写入一个业务包; 一旦将业务包的写入拆分为多次Write操作,
		就无法保证某个goroutine的某业务包数据在conn上发送的连续性;

		同时可以看出, 即便是Read操作, 也是有锁保护的; 多个goroutine对同一
		conn的并发读不会出现读出内容重叠的情况, 但内容断点是依运行时调度来
		随机确定的; 所以存在一个业务包数据三分之一的内容被goroutine-1读走,
		而另三分之二被goroutine-2读走的情况; 比如一个完整数据包"world",
		当goroutine的读缓冲区长度小于5时, 存在这样一种可能: 一个goroutine
		读出"worl", 而另一个goroutine读出"d";


		关闭连接

		在己方已经关闭的Socket上再进行Read和Write操作, 会得到"use of closed
		network connection"的错误;
		在对方关闭的Socket上执行Read操作会得到EOF错误, 但Write操作依然会成功,
		因为数据会成功写入己方的内核Socket缓冲区中, 即便最终发不到对方的
		Socket缓冲区(因为己方Socket尚未关闭); 因此当发现对方Socket关闭时,
		己方应该正确处理自己的Socket, 再继续进行Write操作已经无任何意义了;

*/
