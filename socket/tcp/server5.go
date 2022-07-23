package main

import (
	"log"
	"net"
	"time"
)

/* 写阻塞 */

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// 从连接上读取数据
		time.Sleep(5 * time.Second)
		// 前5秒不读取数据, 因此当 ./client5.go
		// 一直调用 Write 尝试写入数据时, 写到一定量后就会发生阻塞
		var buf = make([]byte, 65536)
		log.Println("start to read from conn")
		c.SetReadDeadline(time.Now().Add(time.Microsecond * 10)) // 10 微妙超时时间
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes, error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			return
		}

		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("error listen:", err)
		return
	}
	defer l.Close()
	log.Println("listen ok")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}
		go handleConn(conn)
	}
}
