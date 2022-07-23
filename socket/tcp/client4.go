package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Println("dial error: ", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	data := make([]byte, 65536) // 65536 为2的16次方
	conn.Write(data)
	log.Println("send data ok")
	time.Sleep(time.Second * 10000)
}
