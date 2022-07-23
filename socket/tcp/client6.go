package main

import (
	"log"
	"net"
	"time"
)

/* 写入超时 */

func main() {
	log.Println("begin dial")
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Println("dial error: ", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	data := make([]byte, 65536) // 65536 为2的16次方
	var total int
	for {
		conn.SetWriteDeadline(time.Now().Add(time.Microsecond * 10)) // 10毫秒写入超时
		n, err := conn.Write(data)
		if err != nil {
			total += n
			log.Printf("write %d bytes, error:%s\n", n, err)
			break
		}
		total += n
		log.Printf("write %d bytes this time, %d bytes in total\n", n, total)
	}
	log.Printf("write %d bytes in total\n", total)
	time.Sleep(time.Second * 10000)
}

// ......
// 2022/07/23 19:06:24 write 65536 bytes this time, 2424832 bytes in total
// 2022/07/23 19:06:24 write 65536 bytes this time, 2490368 bytes in total
// 2022/07/23 19:06:24 write 65536 bytes this time, 2555904 bytes in total
// 2022/07/23 19:06:24 write 30689 bytes, error:write tcp 127.0.0.1:35794->127.0.0.1:8080: i/o timeout
// 2022/07/23 19:06:24 write 2586593 bytes in total

// 可知写入超时时, 仍然存数据部分写入的情况
