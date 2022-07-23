package main

/*
	本节是一个典型的go网络服务端程序
	在go程序的用户层看来, goroutine 采用了"阻塞I/O模型"进行网络I/O操作,
	Socket 都是"阻塞"的, 实际上, 这种假象是go运行时中的 netpoller(网络
	轮询器, TODO)通过I/O多路复用机制模拟的, 对应底层操作系统 Socket 实际上
	是非阻塞的: $GOROOT/src/net/sock_cloexec.go; 只是运行时拦截了针对
	底层Socket的系统调用返回的错误码, 并通过 netpoller 和 goroutine 调度
	让 goroutine "阻塞"在用户层所看到的 Socket 描述符上;
	例如:
		当用户层针对某个Socket描述符发起read操作时, 如果该Socket对应的连接
		上尚无数据, 那么Go运行时会将该Socket描述符加入netpoller中监听, 直到
		Go运行时收到该Socket数据可读的通知, Go运行时才会重新唤醒等待在该
		Socket上准备读数据的那个goroutine; 而这个过程从goroutine的视角来看,
		就像是read操作一直阻塞在那个Socket描述符上;

*/

import (
	"fmt"
	"net"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// 从连接上读数据
		// ...
		// 从连接上写数据
		// ...
	}
}

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		// 启动一个新的 goroutine 处理连接
		go handleConn(c)
	}
}
