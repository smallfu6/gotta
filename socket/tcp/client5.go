package main

import (
	"log"
	"net"
	"time"
)

/* 写入阻塞 */

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

// Write操作存在写入部分数据的情况, 比如当写入阻塞时, 杀掉 server5, 这时的输出
// 如下:
// ......
// 2022/07/23 18:55:09 write 65536 bytes this time, 2555904 bytes in total
// 2022/07/23 18:55:11 write 30689 bytes, error:write tcp 127.0.0.1:35790->127.0.0.1:8080: write: connection reset by peer
// 2022/07/23 18:55:11 write 2586593 bytes in total
// 可以看出Write在写入30689 bytes时发生了阻塞, 服务端关闭后, 客户端又写入了
// 30689 字节后才返回 broken pipe 错误, 因此程序需要对这 30689 字节数据进行
// 特殊处理;
