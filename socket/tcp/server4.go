package main

import (
	"log"
	"net"
	"time"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// 从连接上读取数据
		time.Sleep(10 * time.Second)
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

// 2022/07/23 18:31:51 listen ok
// 2022/07/23 18:32:13 start to read from conn
// 2022/07/23 18:32:13 read 65536 bytes, content is
// 2022/07/23 18:32:23 start to read from conn
// 2022/07/23 18:32:23 conn read 0 bytes, error: read tcp 127.0.0.1:8080->127.0.0.1:35776: i/o timeout
// 2022/07/23 18:32:33 start to read from conn
// 2022/07/23 18:32:33 conn read 0 bytes, error: read tcp 127.0.0.1:8080->127.0.0.1:35776: i/o timeout
// 2022/07/23 18:32:43 start to read from conn
// 2022/07/23 18:32:43 conn read 0 bytes, error: read tcp 127.0.0.1:8080->127.0.0.1:35776: i/o timeout
// 2022/07/23 18:32:53 start to read from conn
// 2022/07/23 18:32:53 conn read 0 bytes, error: read tcp 127.0.0.1:8080->127.0.0.1:35776: i/o timeout
// 2022/07/23 18:33:03 start to read from conn
// 2022/07/23 18:33:03 conn read 0 bytes, error: read tcp 127.0.0.1:8080->127.0.0.1:35776: i/o timeout

/*
	从执行结果看, 第一次成功读取所有数据, 反复执行了多次, 没有出现读出部分数据且返回超时错误的情况
*/
