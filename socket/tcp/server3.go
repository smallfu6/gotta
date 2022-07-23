package main

import (
	"log"
	"net"
	"time"
)

/* 客户端在有数据关闭时服务端的Read操作 */

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// 从连接上读取数据
		time.Sleep(10 * time.Second)
		var buf = make([]byte, 10)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Println("conn read error: ", err)
			return
		}
		log.Printf("read %d bytes content is %s\n", n, string(buf[:n]))
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

// 2022/07/23 18:07:33 listen ok
// 2022/07/23 18:07:45 start to read from conn
// 2022/07/23 18:07:45 read 10 bytes content is test 0test
// 2022/07/23 18:07:55 start to read from conn
// 2022/07/23 18:07:55 read 10 bytes content is  1test 2te
// 2022/07/23 18:08:05 start to read from conn
// 2022/07/23 18:08:05 read 10 bytes content is st 3test 4
// 2022/07/23 18:08:15 start to read from conn
// 2022/07/23 18:08:15 read 10 bytes content is test 5test
// 2022/07/23 18:08:25 start to read from conn
// 2022/07/23 18:08:25 read 10 bytes content is  6test 7te
// 2022/07/23 18:08:35 start to read from conn
// 2022/07/23 18:08:35 read 10 bytes content is st 8test 9
// 2022/07/23 18:08:45 start to read from conn
// 2022/07/23 18:08:45 read 7 bytes content is test 10
// 2022/07/23 18:08:55 start to read from conn
// 2022/07/23 18:08:55 conn read error:  EOF

/*
	从输出结果来看, 在客户端关闭Socket并退出后, server3依旧没有开始执行Read
	操作; 10秒后的第一次Read操作成功读出了10字节的数据, 当执行第8次Read操作时,
	由于此时客户端已经关闭Socket并退出了, Read返回了错误EOF(代表连接断开);
	如果连接未断开, 但是客户端没有发送数据, Read 操作会阻塞

*/
